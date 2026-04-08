package realtime

import (
	"context"
	"sync"
	"time"

	optionsTypes "github.com/wisp-trading/sdk/pkg/markets/options/types"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	optionsconnector "github.com/wisp-trading/sdk/pkg/types/connector/options"
	"github.com/wisp-trading/sdk/pkg/types/logging"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
)

type optionsRealtimeIngestor struct {
	connector optionsconnector.WebSocketConnector
	watchlist optionsTypes.OptionsWatchlist
	store     optionsTypes.OptionsStore
	logger    logging.ApplicationLogger
	isActive  bool
	ctx       context.Context
	cancel    context.CancelFunc
	wg        sync.WaitGroup

	// pendingExpirations tracks expirations skipped at Start() due to absent strikes.
	// processWatchlistEvents retries them when StrikesUpdated fires.
	pendingMu          sync.Mutex
	pendingExpirations map[expirationRef]struct{}
}

type expirationRef struct {
	pair       portfolio.Pair
	expiration time.Time
}

func (o *optionsRealtimeIngestor) Start(ctx context.Context) error {
	if o.isActive {
		return nil
	}

	o.isActive = true
	o.ctx, o.cancel = context.WithCancel(ctx)
	o.pendingExpirations = make(map[expirationRef]struct{})

	exchangeName := o.connector.GetConnectorInfo().Name
	expirations := o.watchlist.GetWatchedExpirations(exchangeName)

	// Subscribe to all watched expirations. Strikes are resolved from the watchlist,
	// which is populated by the batch ingestor via SetStrikes.
	for pair, exps := range expirations {
		for _, expiration := range exps {
			o.subscribeExpiration(exchangeName, pair, expiration)
		}
	}

	// Subscribe to any already fast-watched instruments (order books)
	for _, contract := range o.watchlist.GetWatchedInstruments(exchangeName) {
		c := toConnectorContract(contract)
		if err := o.connector.SubscribeOrderBook(&c); err != nil {
			o.logger.Errorf("failed to subscribe to order book for %s: %v", contract, err)
		}
	}

	// processWatchlistEvents and processUpdates are tracked so Stop() can wait for
	// them to exit before running connector cleanup.
	o.wg.Add(2)
	go o.processWatchlistEvents()
	go o.processUpdates()

	return nil
}

func (o *optionsRealtimeIngestor) Stop() error {
	if !o.isActive {
		return nil
	}

	o.isActive = false

	// Cancel context first so both goroutines begin exiting.
	if o.cancel != nil {
		o.cancel()
	}

	// Wait for goroutines to finish before touching the connector — prevents
	// concurrent calls to Subscribe/Unsubscribe during shutdown.
	o.wg.Wait()

	exchangeName := o.connector.GetConnectorInfo().Name

	for pair, exps := range o.watchlist.GetWatchedExpirations(exchangeName) {
		for _, expiration := range exps {
			strikes := o.watchlist.GetAvailableStrikes(exchangeName, pair, expiration)
			if len(strikes) == 0 {
				continue
			}
			contracts := strikesToContracts(pair, expiration, strikes)
			if err := o.connector.UnsubscribeExpirationUpdates(pair, expiration, contracts); err != nil {
				o.logger.Errorf("failed to unsubscribe expiration updates for %s %s: %v", pair.Symbol(), expiration, err)
			}
		}
	}

	for _, contract := range o.watchlist.GetWatchedInstruments(exchangeName) {
		c := toConnectorContract(contract)
		if err := o.connector.UnsubscribeOrderBook(&c); err != nil {
			o.logger.Warnf("failed to unsubscribe order book for %s: %v", contract, err)
		}
	}

	return nil
}

func (o *optionsRealtimeIngestor) IsActive() bool {
	return o.isActive
}

func (o *optionsRealtimeIngestor) GetMarketType() connector.MarketType {
	return connector.MarketTypeOptions
}

func (o *optionsRealtimeIngestor) GetActiveConnections() map[connector.ExchangeName]interface{} {
	connMap := make(map[connector.ExchangeName]interface{})
	connMap[o.connector.GetConnectorInfo().Name] = o.connector
	return connMap
}

// subscribeExpiration subscribes to ticker updates for a given expiration.
// If strikes are not yet available it records the expiration as pending so
// processWatchlistEvents can retry when StrikesUpdated fires.
func (o *optionsRealtimeIngestor) subscribeExpiration(
	exchangeName connector.ExchangeName,
	pair portfolio.Pair,
	expiration time.Time,
) {
	strikes := o.watchlist.GetAvailableStrikes(exchangeName, pair, expiration)
	if len(strikes) == 0 {
		o.logger.Warnf("no known strikes for %s %s — will retry when batch ingestor populates them",
			pair.Symbol(), expiration)
		o.pendingMu.Lock()
		o.pendingExpirations[expirationRef{pair: pair, expiration: expiration}] = struct{}{}
		o.pendingMu.Unlock()
		return
	}

	contracts := strikesToContracts(pair, expiration, strikes)
	if err := o.connector.SubscribeExpirationUpdates(pair, expiration, contracts); err != nil {
		o.logger.Errorf("failed to subscribe to expiration updates for %s %s: %v",
			pair.Symbol(), expiration, err)
		return
	}

	// No longer pending
	o.pendingMu.Lock()
	delete(o.pendingExpirations, expirationRef{pair: pair, expiration: expiration})
	o.pendingMu.Unlock()
}

// processWatchlistEvents handles dynamic watchlist changes:
//   - InstrumentWatched / InstrumentUnwatched → order book subscribe/unsubscribe
//   - StrikesUpdated → retry any pending expirations that had no strikes at Start()
func (o *optionsRealtimeIngestor) processWatchlistEvents() {
	defer o.wg.Done()

	exchangeName := o.connector.GetConnectorInfo().Name
	events := o.watchlist.Subscribe(exchangeName)
	defer o.watchlist.Unsubscribe(exchangeName)

	for {
		select {
		case <-o.ctx.Done():
			return
		case ev, ok := <-events:
			if !ok {
				return
			}
			switch ev.Type {
			case optionsTypes.WatchEventInstrumentWatched:
				if ev.Contract == nil {
					continue
				}
				c := toConnectorContract(*ev.Contract)
				if err := o.connector.SubscribeOrderBook(&c); err != nil {
					o.logger.Errorf("failed to subscribe order book for %v: %v", ev.Contract, err)
				}

			case optionsTypes.WatchEventInstrumentUnwatched:
				if ev.Contract == nil {
					continue
				}
				c := toConnectorContract(*ev.Contract)
				if err := o.connector.UnsubscribeOrderBook(&c); err != nil {
					o.logger.Warnf("failed to unsubscribe order book for %v: %v", ev.Contract, err)
				}

			case optionsTypes.WatchEventStrikesUpdated:
				// Retry any expirations that were skipped at Start() due to absent strikes.
				o.pendingMu.Lock()
				ref := expirationRef{pair: ev.Pair, expiration: ev.Expiration}
				_, isPending := o.pendingExpirations[ref]
				o.pendingMu.Unlock()

				if isPending {
					o.subscribeExpiration(exchangeName, ev.Pair, ev.Expiration)
				}
			}
		}
	}
}

// toConnectorContract converts the domain OptionContract to the connector type.
func toConnectorContract(c optionsTypes.OptionContract) optionsconnector.OptionContract {
	return optionsconnector.OptionContract(c)
}

// strikesToContracts expands a list of strikes into CALL and PUT contracts.
func strikesToContracts(pair portfolio.Pair, expiration time.Time, strikes []float64) []optionsconnector.OptionContract {
	contracts := make([]optionsconnector.OptionContract, 0, len(strikes)*2)
	for _, strike := range strikes {
		contracts = append(contracts,
			optionsconnector.OptionContract{Pair: pair, Strike: strike, Expiration: expiration, OptionType: "CALL"},
			optionsconnector.OptionContract{Pair: pair, Strike: strike, Expiration: expiration, OptionType: "PUT"},
		)
	}
	return contracts
}

func (o *optionsRealtimeIngestor) processUpdates() {
	defer o.wg.Done()

	updateChannels := o.connector.GetOptionUpdateChannels()

	for {
		select {
		case <-o.ctx.Done():
			return
		default:
		}

		for _, updateChan := range updateChannels {
			select {
			case update, ok := <-updateChan:
				if !ok {
					return
				}

				contract := optionsTypes.OptionContract{
					Pair:       update.Contract.Pair,
					Strike:     update.Contract.Strike,
					Expiration: update.Contract.Expiration,
					OptionType: update.Contract.OptionType,
				}

				o.store.SetMarkPrice(contract, update.MarkPrice)
				o.store.SetUnderlyingPrice(contract, update.UnderlyingPrice)
				o.store.SetGreeks(contract, optionsTypes.Greeks{
					Delta: update.Greeks.Delta,
					Gamma: update.Greeks.Gamma,
					Theta: update.Greeks.Theta,
					Vega:  update.Greeks.Vega,
					Rho:   update.Greeks.Rho,
				})
				o.store.SetIV(contract, update.IV)

			case <-o.ctx.Done():
				return
			default:
			}
		}
	}
}
