package realtime

import (
	"context"

	optionsTypes "github.com/wisp-trading/sdk/pkg/markets/options/types"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	optionsconnector "github.com/wisp-trading/sdk/pkg/types/connector/options"
	"github.com/wisp-trading/sdk/pkg/types/logging"
)

type optionsRealtimeIngestor struct {
	connector optionsconnector.WebSocketConnector
	watchlist optionsTypes.OptionsWatchlist
	store     optionsTypes.OptionsStore
	logger    logging.ApplicationLogger
	isActive  bool
	ctx       context.Context
	cancel    context.CancelFunc
}

func (o *optionsRealtimeIngestor) Start(ctx context.Context) error {
	if o.isActive {
		return nil
	}

	o.isActive = true
	o.ctx, o.cancel = context.WithCancel(ctx)

	exchangeName := o.connector.GetConnectorInfo().Name
	expirations := o.watchlist.GetWatchedExpirations(exchangeName)

	// Subscribe to all watched expirations
	for pair, exps := range expirations {
		for _, expiration := range exps {
			if err := o.connector.SubscribeExpirationUpdates(pair, expiration); err != nil {
				o.logger.Errorf("failed to subscribe to expiration updates for %s %s: %v", pair.Symbol(), expiration, err)
			}
		}
	}

	// Start processing updates
	go o.processUpdates()

	return nil
}

func (o *optionsRealtimeIngestor) Stop() error {
	if !o.isActive {
		return nil
	}

	exchangeName := o.connector.GetConnectorInfo().Name
	expirations := o.watchlist.GetWatchedExpirations(exchangeName)

	// Unsubscribe from all expirations
	for pair, exps := range expirations {
		for _, expiration := range exps {
			if err := o.connector.UnsubscribeExpirationUpdates(pair, expiration); err != nil {
				o.logger.Errorf("failed to unsubscribe from expiration updates for %s %s: %v", pair.Symbol(), expiration, err)
			}
		}
	}

	o.isActive = false
	if o.cancel != nil {
		o.cancel()
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

func (o *optionsRealtimeIngestor) processUpdates() {
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
