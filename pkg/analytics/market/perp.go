package market

import (
	"context"
	"fmt"

	perpTypes "github.com/wisp-trading/sdk/pkg/markets/perp/types"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/connector/perp"
	marketTypes "github.com/wisp-trading/sdk/pkg/types/data/stores/market"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	"github.com/wisp-trading/sdk/pkg/types/wisp/analytics"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

// perpMarketService implements analytics.PerpMarket
// It embeds baseMarketService for price operations and implements extension-dependent methods directly.
type perpMarketService struct {
	baseMarketService
	store marketTypes.MarketStore // Store that should have OrderBook, Kline, and Funding extensions
}

// newPerpMarketService creates a perp market service
func newPerpMarketService(store marketTypes.MarketStore) *perpMarketService {
	return &perpMarketService{
		baseMarketService: newBaseMarketService(store),
		store:             store,
	}
}

// ========== Extension-dependent methods (OrderBook, Klines) ==========

// OrderBook returns the order book for an asset (uses OrderBookStoreExtension)
func (s *perpMarketService) OrderBook(ctx context.Context, asset portfolio.Pair, exchange ...connector.ExchangeName) (*connector.OrderBook, error) {
	// Cast to extension interface
	orderBookStore, ok := s.store.(marketTypes.OrderBookStoreExtension)
	if !ok {
		return nil, fmt.Errorf("perp store does not support order books")
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
func (s *perpMarketService) GetKlines(asset portfolio.Pair, exchange connector.ExchangeName, interval string, limit int) []connector.Kline {
	// Cast to extension interface
	if klineStore, ok := s.store.(marketTypes.KlineStoreExtension); ok {
		return klineStore.GetKlines(asset, exchange, interval, limit)
	}
	return nil
}

// GetTradableQuantity calculates available liquidity (uses OrderBookStoreExtension)
func (s *perpMarketService) GetTradableQuantity(ctx context.Context, asset portfolio.Pair, opts ...analytics.LiquidityOptions) numerical.Decimal {
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

// ========== Perp-specific methods (FundingRateStoreExtension) ==========

func (s *perpMarketService) FundingRate(ctx context.Context, asset portfolio.Pair, exchange connector.ExchangeName) (*perp.FundingRate, error) {
	// Cast to extension interface
	fundingStore, ok := s.store.(perpTypes.FundingRateStoreExtension)
	if !ok {
		return nil, fmt.Errorf("perp store does not support funding rates")
	}

	rate := fundingStore.GetFundingRate(asset, exchange)
	if rate == nil {
		return nil, fmt.Errorf("no funding rate found for %s on %s", asset.Symbol(), exchange)
	}
	return rate, nil
}

// FundingRates returns funding rates across all perp exchanges
func (s *perpMarketService) FundingRates(ctx context.Context, asset portfolio.Pair) map[connector.ExchangeName]perp.FundingRate {
	// Cast to extension interface
	if fundingStore, ok := s.store.(perpTypes.FundingRateStoreExtension); ok {
		return fundingStore.GetFundingRatesForAsset(asset)
	}
	return make(map[connector.ExchangeName]perp.FundingRate)
}

// GetAllAssetsWithFundingRates returns all assets that have funding rate data
func (s *perpMarketService) GetAllAssetsWithFundingRates(ctx context.Context) []portfolio.Pair {
	// Cast to extension interface
	if fundingStore, ok := s.store.(perpTypes.FundingRateStoreExtension); ok {
		return fundingStore.GetAllAssetsWithFundingRates()
	}
	return nil
}
