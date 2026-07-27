// Package analytics provides shared pair-market analytics used by spot and perp.
// Domain packages embed Service and add market-specific methods (e.g. funding).
package analytics

import (
	"context"
	"fmt"

	market "github.com/wisp-trading/sdk/pkg/markets/base/types/stores/market"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	analyticsTypes "github.com/wisp-trading/sdk/pkg/types/wisp/analytics"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

// PairDataReader is the store surface needed for pair price/book/kline analytics.
// Satisfied by spot and perp MarketStore (base market + orderbook + kline extensions).
type PairDataReader interface {
	GetPairPrices(pair portfolio.Pair) market.PriceMap
	GetOrderBook(pair portfolio.Pair, exchange connector.ExchangeName) *connector.OrderBook
	GetKlines(pair portfolio.Pair, exchange connector.ExchangeName, interval string, limit int) []connector.Kline
}

// Service implements shared pair-market analytics methods.
type Service struct {
	store PairDataReader
}

// NewService creates a pair analytics service backed by store.
func NewService(store PairDataReader) *Service {
	return &Service{store: store}
}

// Price returns the current price for a pair, optionally on a specific exchange.
func (s *Service) Price(_ context.Context, asset portfolio.Pair, exchange ...connector.ExchangeName) (numerical.Decimal, error) {
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

// Prices returns prices across all exchanges that have data for the pair.
func (s *Service) Prices(_ context.Context, asset portfolio.Pair) map[connector.ExchangeName]numerical.Decimal {
	result := make(map[connector.ExchangeName]numerical.Decimal)
	for exchange, price := range s.store.GetPairPrices(asset) {
		result[exchange] = price.Price
	}
	return result
}

// OrderBook returns the order book for a pair on an exchange.
func (s *Service) OrderBook(_ context.Context, asset portfolio.Pair, exchange connector.ExchangeName) (*connector.OrderBook, error) {
	ob := s.store.GetOrderBook(asset, exchange)
	if ob == nil {
		return nil, fmt.Errorf("no order book found for %s on %s", asset.Symbol(), exchange)
	}
	return ob, nil
}

// GetKlines returns historical klines from the store.
func (s *Service) GetKlines(asset portfolio.Pair, exchange connector.ExchangeName, interval string, limit int) []connector.Kline {
	return s.store.GetKlines(asset, exchange, interval, limit)
}

// GetTradableQuantity calculates available liquidity from the order book.
func (s *Service) GetTradableQuantity(ctx context.Context, asset portfolio.Pair, opts ...analyticsTypes.LiquidityOptions) numerical.Decimal {
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
