package real_time

import (
	"context"
	"sync"

	"github.com/wisp-trading/sdk/pkg/markets/base/types/ingestors/realtime"
	"github.com/wisp-trading/sdk/pkg/markets/base/types/stores/market"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/logging"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
)

// klineExtension handles WebSocket subscriptions for kline (candlestick) updates
type klineExtension struct {
	store     market.KlineStoreExtension
	logger    logging.ApplicationLogger
	intervals []string
}

func NewKlineExtension(store market.KlineStoreExtension, logger logging.ApplicationLogger, intervals []string) realtime.WebSocketExtension {
	return &klineExtension{
		store:     store,
		logger:    logger,
		intervals: intervals,
	}
}

func (k *klineExtension) Subscribe(wsConn interface{}, exchangeName connector.ExchangeName, assets []portfolio.Pair) error {
	// Type-assert to WebSocketSubscriber
	wsSubscriber, ok := wsConn.(realtime.WebSocketSubscriber)
	if !ok {
		k.logger.Debug("WebSocket connector %s does not support kline subscriptions", exchangeName)
		return nil
	}

	// Subscribe to klines for all pairs and intervals
	for _, pair := range assets {
		for _, interval := range k.intervals {
			if err := wsSubscriber.SubscribeKlines(pair, interval); err != nil {
				k.logger.Debug("Failed to subscribe to %s klines for %s on %s: %v",
					interval, pair.Symbol(), exchangeName, err)
			}
		}
	}

	return nil
}

func (k *klineExtension) ProcessChannels(wsConn interface{}, exchangeName connector.ExchangeName, ctx context.Context) {
	// Type-assert to WebSocketSubscriber
	wsSubscriber, ok := wsConn.(realtime.WebSocketSubscriber)
	if !ok {
		return
	}

	channels := wsSubscriber.GetKlineChannels()
	k.logger.Info("Processing %d kline channels for %s", len(channels), exchangeName)

	var wg sync.WaitGroup
	for channelKey, klineChan := range channels {
		wg.Add(1)
		go func(key string, ch <-chan connector.Kline) {
			defer wg.Done()
			k.processKlineChannel(ctx, exchangeName, key, ch)
		}(channelKey, klineChan)
	}

	wg.Wait()
	k.logger.Info("All kline channels closed for %s", exchangeName)
}

func (k *klineExtension) processKlineChannel(ctx context.Context, exchangeName connector.ExchangeName, channelKey string, klineChan <-chan connector.Kline) {
	k.logger.Debug("Starting kline channel processor for %s on %s", channelKey, exchangeName)

	for {
		select {
		case <-ctx.Done():
			k.logger.Debug("Context cancelled, stopping kline channel %s", channelKey)
			return

		case kline, ok := <-klineChan:
			if !ok {
				k.logger.Debug("Kline channel %s closed", channelKey)
				return
			}

			pair := kline.Pair

			// Write to store
			k.store.UpdateKline(pair, exchangeName, kline)

			k.logger.Debug("WebSocket updated %s kline for %s on %s",
				kline.Interval, pair.Symbol(), exchangeName)
		}
	}
}

func (k *klineExtension) Unsubscribe(wsConn interface{}, exchangeName connector.ExchangeName) error {
	// Note: The connector should handle cleanup of subscriptions on disconnect
	k.logger.Info("Unsubscribing from kline updates for %s", exchangeName)
	return nil
}

var _ realtime.WebSocketExtension = (*klineExtension)(nil)
