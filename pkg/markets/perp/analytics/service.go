package analytics

import (
	"context"
	"fmt"

	"github.com/wisp-trading/sdk/pkg/markets/perp/types"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	perpConn "github.com/wisp-trading/sdk/pkg/types/connector/perp"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	analyticsTypes "github.com/wisp-trading/sdk/pkg/types/wisp/analytics"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

// service implements analytics.PerpMarket backed by the typed perp MarketStore.
// No interface casting — the store already satisfies OrderBook, Kline, and FundingRate extensions.
type service struct {
	store types.MarketStore
}

// New creates a new perp analytics service.
func New(store types.MarketStore) types.PerpMarket {
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

// Prices returns prices across all perp exchanges for an asset.
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

// FundingRate returns the latest funding rate for an asset on the specified exchange.
func (s *service) FundingRate(ctx context.Context, asset portfolio.Pair, exchange connector.ExchangeName) (*perpConn.FundingRate, error) {
	rate := s.store.GetFundingRate(asset, exchange)
	if rate == nil {
		return nil, fmt.Errorf("no funding rate found for %s on %s", asset.Symbol(), exchange)
	}
	return rate, nil
}

// FundingRates returns funding rates across all perp exchanges for an asset.
func (s *service) FundingRates(ctx context.Context, asset portfolio.Pair) map[connector.ExchangeName]perpConn.FundingRate {
	return s.store.GetFundingRatesForAsset(asset)
}

// GetAllAssetsWithFundingRates returns all assets that have funding rate data.
func (s *service) GetAllAssetsWithFundingRates(ctx context.Context) []portfolio.Pair {
	return s.store.GetAllAssetsWithFundingRates()
}
