package realtime

import (
	"context"

	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/connector/prediction"
	"github.com/wisp-trading/sdk/pkg/types/data/ingestors/realtime"
	predictionStore "github.com/wisp-trading/sdk/pkg/types/data/stores/market/prediction"
	"github.com/wisp-trading/sdk/pkg/types/logging"
)

// orderBookExtension handles WebSocket subscriptions for prediction order book updates.
type predictionOrderBookExtension struct {
	store  predictionStore.OrderBookStoreExtension
	logger logging.ApplicationLogger
}

func NewPredictionOrderBookExtension(
	store predictionStore.OrderBookStoreExtension,
	logger logging.ApplicationLogger,
) realtime.PredictionExtension {
	return &predictionOrderBookExtension{
		store:  store,
		logger: logger,
	}
}

// Subscribe is called by the prediction realtime ingestor for each market it wants to watch.
func (e *predictionOrderBookExtension) Subscribe(
	wsConn interface{},
	exchangeName connector.ExchangeName,
	market prediction.Market,
) error {
	wsConnector, ok := wsConn.(prediction.WebSocketConnector)
	if !ok {
		e.logger.Debug("WebSocket connector %s does not support prediction order book subscriptions", exchangeName)
		return nil
	}

	if err := wsConnector.SubscribeOrderBook(market); err != nil {
		e.logger.Error("Failed to subscribe prediction order books for market %s on %s: %v",
			market.MarketID, exchangeName, err)
		return err
	}

	e.logger.Info("Subscribed prediction order books for market %s on %s", market.MarketID, exchangeName)
	return nil
}

func (e *predictionOrderBookExtension) ProcessChannels(
	wsConn interface{},
	exchangeName connector.ExchangeName,
	ctx context.Context,
) {
	wsConnector, ok := wsConn.(prediction.WebSocketConnector)
	if !ok {
		return
	}

	orderBookChan := wsConnector.GetOrderBookUpdates()
	e.logger.Info("Starting prediction order book channel processor for %s", exchangeName)

	for {
		select {
		case <-ctx.Done():
			e.logger.Debug("Context cancelled, stopping prediction order book channel for %s", exchangeName)
			return

		case update, ok := <-orderBookChan:
			if !ok {
				e.logger.Debug("Prediction order book channel closed for %s", exchangeName)
				return
			}

			e.store.UpdateOrderBook(exchangeName, update.MarketID, update.OutcomeID, update.OrderBook)

			if len(update.Bids) > 0 && len(update.Asks) > 0 {
				e.logger.Debug(
					"WS updated prediction order book for market %s / outcome %s on %s - bid: %s, ask: %s",
					update.MarketID,
					update.OutcomeID,
					exchangeName,
					update.Bids[0].Price.StringFixed(2),
					update.Asks[0].Price.StringFixed(2),
				)
			}
		}
	}
}

func (e *predictionOrderBookExtension) Unsubscribe(
	wsConn interface{},
	exchangeName connector.ExchangeName,
	market prediction.Market,
) error {
	wsConnector, ok := wsConn.(prediction.WebSocketConnector)
	if !ok {
		return nil
	}

	if err := wsConnector.UnsubscribeMarket(market); err != nil {
		e.logger.Warn("Failed to unsubscribe prediction market %s on %s: %v",
			market.MarketID, exchangeName, err)
		return err
	}

	e.logger.Info("Unsubscribed prediction market %s on %s", market.MarketID, exchangeName)
	return nil
}
