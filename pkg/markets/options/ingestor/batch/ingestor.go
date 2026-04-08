package batch

import (
	"time"

	optionsTypes "github.com/wisp-trading/sdk/pkg/markets/options/types"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	optionsconnector "github.com/wisp-trading/sdk/pkg/types/connector/options"
	"github.com/wisp-trading/sdk/pkg/types/logging"
)

type optionsIngestor struct {
	connector optionsconnector.Connector
	watchlist optionsTypes.OptionsWatchlist
	store     optionsTypes.OptionsStore
	logger    logging.ApplicationLogger
	isActive  bool
	ticker    *time.Ticker
	stopChan  chan struct{}
}

func (o *optionsIngestor) Start(interval time.Duration) error {
	if o.isActive {
		return nil
	}

	o.isActive = true
	o.ticker = time.NewTicker(interval)

	go func() {
		for {
			select {
			case <-o.ticker.C:
				o.CollectNow()
			case <-o.stopChan:
				o.ticker.Stop()
				return
			}
		}
	}()

	return nil
}

func (o *optionsIngestor) Stop() error {
	if !o.isActive {
		return nil
	}

	o.isActive = false
	close(o.stopChan)
	return nil
}

func (o *optionsIngestor) IsActive() bool {
	return o.isActive
}

func (o *optionsIngestor) CollectNow() {
	exchangeName := o.connector.GetConnectorInfo().Name
	expirations := o.watchlist.GetWatchedExpirations(exchangeName)

	for pair, exps := range expirations {
		for _, expiration := range exps {
			expirationData, err := o.connector.GetExpirationData(pair, expiration)
			if err != nil {
				o.logger.Errorf("failed to get expiration data for %s %s: %v", pair.Symbol(), expiration, err)
				continue
			}

			// Store mark prices, Greeks, and IV for each strike and option type
			for strike, optionTypes := range expirationData {
				for optionType, optionData := range optionTypes {
					contract := optionsTypes.OptionContract{
						Pair:       pair,
						Strike:     strike,
						Expiration: expiration,
						OptionType: optionType,
					}

					o.store.SetMarkPrice(contract, optionData.MarkPrice)
					o.store.SetUnderlyingPrice(contract, optionData.UnderlyingPrice)
					o.store.SetGreeks(contract, optionsTypes.Greeks{
						Delta: optionData.Greeks.Delta,
						Gamma: optionData.Greeks.Gamma,
						Theta: optionData.Greeks.Theta,
						Vega:  optionData.Greeks.Vega,
						Rho:   optionData.Greeks.Rho,
					})
					o.store.SetIV(contract, optionData.IV)
				}
			}
		}
	}
}

func (o *optionsIngestor) GetMarketType() connector.MarketType {
	return connector.MarketTypeOptions
}
