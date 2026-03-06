package analytics

import (
	"context"
	"fmt"

	spotTypes "github.com/wisp-trading/sdk/pkg/markets/spot/types"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	analyticsTypes "github.com/wisp-trading/sdk/pkg/types/wisp/analytics"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

// service implements analytics.SpotMarket backed by the typed spot MarketStore.
// No interface casting — the store already satisfies OrderBook and Kline extensions.
type service struct {
	store spotTypes.MarketStore
}

// New creates a new spot analytics service.
func New(store spotTypes.MarketStore) spotTypes.SpotMarket {
	return &service{store: store}
}

// Price returns the current price for an asset on the specified exchange.
func (s *service) Price(ctx context.Context, asset portfolio.Pair, exchange ...connector.ExchangeName) (numerical.Decimal, error) {
	priceMap := s.store.GetPairPrices(asset)

	if len(exchange) > 0 && exchange[0] != "" {
		price, exists := priceMap[exchange[0]]
		if !exists {
			return numerical.Zero(), fmt.Errorf("no price found for %s on %s", asset.Symbol(), exchange[0])
		}
		return price.Price, nil
	}

	if len(priceMap) == 0 {
		return numerical.Zero(), fmt.Errorf("no price data available for %s", asset.Symbol())
	}

	for _, price := range priceMap {
		return price.Price, nil
	}

	return numerical.Zero(), fmt.Errorf("no price data available for %s", asset.Symbol())
}

// Prices returns prices across all spot exchanges for an asset.
func (s *service) Prices(ctx context.Context, asset portfolio.Pair) map[connector.ExchangeName]numerical.Decimal {
	result := make(map[connector.ExchangeName]numerical.Decimal)
	for exchange, price := range s.store.GetPairPrices(asset) {
		result[exchange] = price.Price
	}
	return result
}

// OrderBook returns the order book for an asset on the specified exchange.
func (s *service) OrderBook(ctx context.Context, asset portfolio.Pair, exchange connector.ExchangeName) (*connector.OrderBook, error) {
	ob := s.store.GetOrderBook(asset, exchange)
	if ob == nil {
		return nil, fmt.Errorf("no order book found for %s on %s", asset.Symbol(), exchange)
	}
	return ob, nil
}

// GetKlines returns historical kline data for an asset on the specified exchange.
func (s *service) GetKlines(asset portfolio.Pair, exchange connector.ExchangeName, interval string, limit int) []connector.Kline {
	return s.store.GetKlines(asset, exchange, interval, limit)
}

// GetTradableQuantity calculates available liquidity from the order book.
func (s *service) GetTradableQuantity(ctx context.Context, asset portfolio.Pair, opts ...analyticsTypes.LiquidityOptions) numerical.Decimal {
	options := analyticsTypes.DefaultLiquidityOptions()
	if len(opts) > 0 {
		options = opts[0]
	}

	ob, err := s.OrderBook(ctx, asset, options.Exchange)
	if err != nil {
		return numerical.Zero()
	}

	return analyticsTypes.CalculateTradableQuantity(ob, options)
}
