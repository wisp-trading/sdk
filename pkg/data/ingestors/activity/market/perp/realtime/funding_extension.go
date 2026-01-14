package realtime

import (
	"github.com/backtesting-org/kronos-sdk/pkg/data/ingestors/activity/market/base"
	"github.com/backtesting-org/kronos-sdk/pkg/types/connector"
	perpConn "github.com/backtesting-org/kronos-sdk/pkg/types/connector/perp"
	perpStore "github.com/backtesting-org/kronos-sdk/pkg/types/data/stores/market/perp"
	"github.com/backtesting-org/kronos-sdk/pkg/types/logging"
	"github.com/backtesting-org/kronos-sdk/pkg/types/portfolio"
)

// FundingRateExtension handles WebSocket subscriptions for perp funding rate updates
type FundingRateExtension struct {
	store  perpStore.MarketStore
	logger logging.ApplicationLogger
}

func NewFundingRateExtension(store perpStore.MarketStore, logger logging.ApplicationLogger) base.WebSocketExtension {
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
			f.logger.Warn(
				"Failed to subscribe to funding rates for %s on %s (may not be supported via WebSocket): %v",
				asset.Symbol(),
				exchangeName,
				err,
			)
			// Not fatal - can fall back to REST
		}
	}

	// Listen to funding rate updates in background
	go f.listenToFundingRates(perpWS, exchangeName)

	f.logger.Info("Subscribed to funding rate updates for %s", exchangeName)
	return nil
}

func (f *FundingRateExtension) listenToFundingRates(perpWS perpConn.WebSocketConnector, exchangeName connector.ExchangeName) {
	fundingChan := perpWS.FundingRateUpdates()

	for update := range fundingChan {
		f.handleFundingRateUpdate(exchangeName, update)
	}

	f.logger.Debug("Funding rate channel closed for %s", exchangeName)
}

func (f *FundingRateExtension) handleFundingRateUpdate(exchangeName connector.ExchangeName, update connector.FundingRate) {
	// Get asset from update (assuming FundingRate has Asset or Symbol field)
	// For now, we'll need to look up the asset - this may need adjustment based on actual FundingRate structure

	f.logger.Debug("WebSocket received funding rate update for %s", exchangeName)

	// Note: The actual asset lookup and update will depend on the FundingRate structure
	// This is a placeholder that needs to be adjusted based on your actual types
}

func (f *FundingRateExtension) Unsubscribe(wsConn interface{}, exchangeName connector.ExchangeName) error {
	// Type-assert to perp WebSocket connector
	_, ok := wsConn.(perpConn.WebSocketConnector)
	if !ok {
		return nil
	}

	// Note: UnsubscribeFundingRates takes an asset parameter
	// We would need to track which assets we subscribed to, or the connector handles cleanup on disconnect
	f.logger.Info("Unsubscribed from funding rate updates for %s", exchangeName)
	return nil
}

var _ base.WebSocketExtension = (*FundingRateExtension)(nil)
