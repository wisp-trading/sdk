package realtime

import (
	"context"

	"github.com/wisp-trading/sdk/pkg/markets/base/types/ingestors/realtime"
	"github.com/wisp-trading/sdk/pkg/markets/base/types/stores/market"
	perpTypes "github.com/wisp-trading/sdk/pkg/markets/perp/types"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	perpConn "github.com/wisp-trading/sdk/pkg/types/connector/perp"
	"github.com/wisp-trading/sdk/pkg/types/logging"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
)

// fundingRateExtension handles WebSocket subscriptions for perp funding rate updates
type fundingRateExtension struct {
	store  perpTypes.MarketStore
	logger logging.ApplicationLogger
}

func NewFundingRateExtension(store perpTypes.MarketStore, logger logging.ApplicationLogger) realtime.WebSocketExtension {
	return &fundingRateExtension{
		store:  store,
		logger: logger,
	}
}

func (f *fundingRateExtension) Subscribe(wsConn interface{}, exchangeName connector.ExchangeName, assets []portfolio.Pair) error {
	perpWS, ok := wsConn.(perpConn.WebSocketConnector)
	if !ok {
		f.logger.Debug("WebSocket connector %s does not support perp operations", exchangeName)
		return nil
	}

	for _, asset := range assets {
		if err := perpWS.SubscribeFundingRates(asset); err != nil {
			f.logger.Warn("Failed to subscribe to funding rates for %s on %s (may not be supported via WebSocket): %v",
				asset.Symbol(), exchangeName, err)
			// Not fatal — can fall back to REST
		} else {
			f.logger.Info("Subscribed to funding rates for %s on %s", asset.Symbol(), exchangeName)
		}
	}

	return nil
}

func (f *fundingRateExtension) ProcessChannels(wsConn interface{}, exchangeName connector.ExchangeName, ctx context.Context) {
	perpWS, ok := wsConn.(perpConn.WebSocketConnector)
	if !ok {
		return
	}

	fundingChan := perpWS.FundingRateUpdates()
	f.logger.Info("Starting funding rate channel processor for %s", exchangeName)

	for {
		select {
		case <-ctx.Done():
			f.logger.Debug("Context cancelled, stopping funding rate channel for %s", exchangeName)
			return
		case update, ok := <-fundingChan:
			if !ok {
				f.logger.Debug("Funding rate channel closed for %s", exchangeName)
				return
			}
			f.handleFundingRateUpdate(exchangeName, update)
		}
	}
}

func (f *fundingRateExtension) handleFundingRateUpdate(exchangeName connector.ExchangeName, update perpConn.FundingRate) {
	asset := update.Pair

	f.store.UpdateFundingRate(asset, exchangeName, update)

	f.store.UpdateLastUpdated(market.UpdateKey{
		DataType: "funding_rates",
		Pair:     asset,
		Exchange: exchangeName,
	})

	f.logger.Debug("WebSocket updated funding rate for %s on %s = %s",
		asset.Symbol(), exchangeName, update.CurrentRate.String())
}

func (f *fundingRateExtension) Unsubscribe(wsConn interface{}, exchangeName connector.ExchangeName) error {
	perpWS, ok := wsConn.(perpConn.WebSocketConnector)
	if !ok {
		return nil
	}

	f.logger.Info("Unsubscribing from funding rate updates for %s", exchangeName)
	_ = perpWS
	return nil
}

var _ realtime.WebSocketExtension = (*fundingRateExtension)(nil)
