package market

import (
	"context"
	"fmt"

	"github.com/wisp-trading/sdk/pkg/types/connector"
	marketTypes "github.com/wisp-trading/sdk/pkg/types/data/stores/market"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	"github.com/wisp-trading/sdk/pkg/types/wisp/analytics"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

// spotMarketService implements analytics.SpotMarket
// It embeds baseMarketService for price operations and implements extension-dependent methods directly.
type spotMarketService struct {
	baseMarketService
	store marketTypes.MarketStore // Store that should have OrderBook and Kline extensions
}

// newSpotMarketService creates a spot market service
func newSpotMarketService(store marketTypes.MarketStore) *spotMarketService {
	return &spotMarketService{
		baseMarketService: newBaseMarketService(store),
		store:             store,
	}
}

// OrderBook returns the order book for an asset (uses OrderBookStoreExtension)
func (s *spotMarketService) OrderBook(ctx context.Context, asset portfolio.Pair, exchange ...connector.ExchangeName) (*connector.OrderBook, error) {
	// Cast to extension interface
	orderBookStore, ok := s.store.(marketTypes.OrderBookStoreExtension)
	if !ok {
		return nil, fmt.Errorf("spot store does not support order books")
	}

	if len(exchange) > 0 && exchange[0] != "" {
		ob := orderBookStore.GetOrderBook(asset, exchange[0])
		if ob == nil {
			return nil, fmt.Errorf("no order book found for %s on %s", asset.Symbol(), exchange[0])
		}
		return ob, nil
	}

	// Return first available
	orderBooks := orderBookStore.GetOrderBooks(asset)
	if len(orderBooks) == 0 {
		return nil, fmt.Errorf("no order book data available for %s", asset.Symbol())
	}

	for _, ob := range orderBooks {
		if ob != nil {
			return ob, nil
		}
	}

	return nil, fmt.Errorf("no order book found for %s", asset.Symbol())
}

// GetKlines returns historical kline/candlestick data (uses KlineStoreExtension)
func (s *spotMarketService) GetKlines(asset portfolio.Pair, exchange connector.ExchangeName, interval string, limit int) []connector.Kline {
	// Cast to extension interface
	if klineStore, ok := s.store.(marketTypes.KlineStoreExtension); ok {
		return klineStore.GetKlines(asset, exchange, interval, limit)
	}
	return nil
}

// GetTradableQuantity calculates available liquidity (uses OrderBookStoreExtension)
func (s *spotMarketService) GetTradableQuantity(ctx context.Context, asset portfolio.Pair, opts ...analytics.LiquidityOptions) numerical.Decimal {
	options := DefaultLiquidityOptions()
	if len(opts) > 0 {
		options = opts[0]
	}

	orderBook, err := s.OrderBook(ctx, asset)
	if err != nil {
		return numerical.Zero()
	}

	return getTradableQuantity(ctx, orderBook, options)
}
