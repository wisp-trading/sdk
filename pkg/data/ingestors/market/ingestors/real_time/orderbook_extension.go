package real_time

import (
	"context"
	"sync"

	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/data/ingestors/realtime"
	"github.com/wisp-trading/sdk/pkg/types/data/stores/market"
	"github.com/wisp-trading/sdk/pkg/types/logging"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
)

// orderBookExtension handles WebSocket subscriptions for order book updates
type orderBookExtension struct {
	store  market.OrderBookStoreExtension
	logger logging.ApplicationLogger
}

func NewOrderBookExtension(store market.OrderBookStoreExtension, logger logging.ApplicationLogger) realtime.WebSocketExtension {
	return &orderBookExtension{
		store:  store,
		logger: logger,
	}
}

func (o *orderBookExtension) Subscribe(wsConn interface{}, exchangeName connector.ExchangeName, assets []portfolio.Pair) error {
	// Type-assert to WebSocketSubscriber
	wsSubscriber, ok := wsConn.(realtime.WebSocketSubscriber)
	if !ok {
		o.logger.Debug("WebSocket connector %s does not support order book subscriptions", exchangeName)
		return nil
	}

	// Subscribe to order books for all pairs
	for _, pair := range assets {
		if err := wsSubscriber.SubscribeOrderBook(pair); err != nil {
			o.logger.Error("Failed to subscribe to order book for %s on %s: %v",
				pair.Symbol(), exchangeName, err)
		} else {
			o.logger.Info("Subscribed to order book for %s on %s", pair.Symbol(), exchangeName)
		}
	}

	return nil
}

func (o *orderBookExtension) ProcessChannels(wsConn interface{}, exchangeName connector.ExchangeName, ctx context.Context) {
	// Type-assert to WebSocketSubscriber
	wsSubscriber, ok := wsConn.(realtime.WebSocketSubscriber)
	if !ok {
		return
	}

	channels := wsSubscriber.GetOrderBookChannels()
	o.logger.Info("Processing %d order book channels for %s", len(channels), exchangeName)

	var wg sync.WaitGroup
	for channelKey, orderBookChan := range channels {
		wg.Add(1)
		go func(key string, ch <-chan connector.OrderBook) {
			defer wg.Done()
			o.processOrderBookChannel(ctx, exchangeName, key, ch)
		}(channelKey, orderBookChan)
	}

	wg.Wait()
	o.logger.Info("All order book channels closed for %s", exchangeName)
}

func (o *orderBookExtension) processOrderBookChannel(ctx context.Context, exchangeName connector.ExchangeName, channelKey string, orderBookChan <-chan connector.OrderBook) {
	o.logger.Debug("Starting order book channel processor for %s on %s", channelKey, exchangeName)

	for {
		select {
		case <-ctx.Done():
			o.logger.Debug("Context cancelled, stopping order book channel %s", channelKey)
			return

		case update, ok := <-orderBookChan:
			if !ok {
				o.logger.Debug("Order book channel %s closed", channelKey)
				return
			}

			// Write to store
			o.store.UpdateOrderBook(update.Pair, exchangeName, update)

			if len(update.Bids) > 0 && len(update.Asks) > 0 {
				o.logger.Debug(
					"WebSocket updated order book for %s on %s - bid: %s, ask: %s",
					update.Pair.Symbol(),
					exchangeName,
					update.Bids[0].Price.StringFixed(2),
					update.Asks[0].Price.StringFixed(2),
				)
			}
		}
	}
}

func (o *orderBookExtension) Unsubscribe(wsConn interface{}, exchangeName connector.ExchangeName) error {
	// Note: The connector should handle cleanup of subscriptions on disconnect
	o.logger.Info("Unsubscribing from order book updates for %s", exchangeName)
	return nil
}

var _ realtime.WebSocketExtension = (*orderBookExtension)(nil)
