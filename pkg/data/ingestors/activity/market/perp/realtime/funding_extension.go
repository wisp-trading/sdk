package realtime

import (
	"context"

	"github.com/backtesting-org/kronos-sdk/pkg/types/connector"
	perpConn "github.com/backtesting-org/kronos-sdk/pkg/types/connector/perp"
	"github.com/backtesting-org/kronos-sdk/pkg/types/data/ingestors/realtime"
	perpStore "github.com/backtesting-org/kronos-sdk/pkg/types/data/stores/market/perp"
	"github.com/backtesting-org/kronos-sdk/pkg/types/logging"
	"github.com/backtesting-org/kronos-sdk/pkg/types/portfolio"
)

// FundingRateExtension handles WebSocket subscriptions for perp funding rate updates
type FundingRateExtension struct {
	store  perpStore.MarketStore
	logger logging.ApplicationLogger
}

func NewFundingRateExtension(store perpStore.MarketStore, logger logging.ApplicationLogger) realtime.WebSocketExtension {
	return &FundingRateExtension{
		store:  store,
		logger: logger,
	}
}

func (f *FundingRateExtension) Subscribe(wsConn interface{}, exchangeName connector.ExchangeName, assets []portfolio.Asset) error {
	// Type-assert to perp WebSocket connector
	perpWS, ok := wsConn.(perpConn.WebSocketConnector)
	if !ok {
		f.logger.Debug("WebSocket connector %s does not support perp operations", exchangeName)
		return nil
	}

	// Subscribe to funding rates for each asset
	// Some exchanges provide funding rates via WebSocket, others only via REST
	for _, asset := range assets {
		if err := perpWS.SubscribeFundingRates(asset); err != nil {
			f.logger.Warn("Failed to subscribe to funding rates for %s on %s (may not be supported via WebSocket): %v",
				asset.Symbol(), exchangeName, err)
			// Not fatal - can fall back to REST
		} else {
			f.logger.Info("Subscribed to funding rates for %s on %s", asset.Symbol(), exchangeName)
		}
	}

	return nil
}

func (f *FundingRateExtension) ProcessChannels(wsConn interface{}, exchangeName connector.ExchangeName, ctx context.Context) {
	// Type-assert to perp WebSocket connector
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
func (f *FundingRateExtension) handleFundingRateUpdate(exchangeName connector.ExchangeName, update perpConn.FundingRate) {
	// Get asset from the funding rate update
	asset := update.Asset

	// Store the funding rate update
	f.store.UpdateFundingRate(asset, exchangeName, update)

	// Update last updated timestamp
	f.store.UpdateLastUpdated(marketTypes.UpdateKey{
		DataType: "funding_rates", // Use string literal for perp-specific data type
		Asset:    asset,
		Exchange: exchangeName,
	})

	f.logger.Debug("WebSocket updated funding rate for %s on %s = %s",
		asset.Symbol(), exchangeName, update.Rate.String())
}

func (f *FundingRateExtension) Unsubscribe(wsConn interface{}, exchangeName connector.ExchangeName) error {
	// Type-assert to perp WebSocket connector
	perpWS, ok := wsConn.(perpConn.WebSocketConnector)
	if !ok {
		return nil
	}

	// Note: The connector should handle cleanup of subscriptions on disconnect
	// Individual UnsubscribeFundingRates calls per asset could be added if needed
	f.logger.Info("Unsubscribing from funding rate updates for %s", exchangeName)

	_ = perpWS // Placeholder to avoid unused variable warning

	return nil
}

var _ realtime.WebSocketExtension = (*FundingRateExtension)(nil)
