package market

import (
	"context"
	"fmt"

	"github.com/backtesting-org/kronos-sdk/pkg/types/connector"
	storeTypes "github.com/backtesting-org/kronos-sdk/pkg/types/data/stores/market"
	"github.com/backtesting-org/kronos-sdk/pkg/types/kronos/analytics"
	"github.com/backtesting-org/kronos-sdk/pkg/types/kronos/numerical"
	"github.com/backtesting-org/kronos-sdk/pkg/types/portfolio"
)

// baseMarketService implements common market operations that work identically
// across all market types (spot, perp, futures, options, etc.)
// Market-specific services embed this to inherit common functionality.
type baseMarketService struct {
	store storeTypes.MarketStore
}

func newBaseMarketService(store storeTypes.MarketStore) baseMarketService {
	return baseMarketService{store: store}
}

// Price returns the current price for an asset
func (b *baseMarketService) Price(ctx context.Context, asset portfolio.Asset, exchange ...connector.ExchangeName) (numerical.Decimal, error) {
	priceMap := b.store.GetAssetPrices(asset)

	if len(exchange) > 0 && exchange[0] != "" {
		price, exists := priceMap[exchange[0]]
		if !exists {
			return numerical.Zero(), fmt.Errorf("no price found for %s on %s", asset.Symbol(), exchange[0])
		}
		return price.Price, nil
	}

	// Return first available
	if len(priceMap) == 0 {
		return numerical.Zero(), fmt.Errorf("no price data available for %s", asset.Symbol())
	}

	for _, price := range priceMap {
		return price.Price, nil
	}

	return numerical.Zero(), fmt.Errorf("no price found for %s", asset.Symbol())
}

// Prices returns prices across all exchanges for this market type
func (b *baseMarketService) Prices(ctx context.Context, asset portfolio.Asset) map[connector.ExchangeName]numerical.Decimal {
	result := make(map[connector.ExchangeName]numerical.Decimal)
	priceMap := b.store.GetAssetPrices(asset)
	for exchange, price := range priceMap {
		result[exchange] = price.Price
	}
	return result
}

// OrderBook returns the order book for an asset
func (b *baseMarketService) OrderBook(ctx context.Context, asset portfolio.Asset, exchange ...connector.ExchangeName) (*connector.OrderBook, error) {
	if len(exchange) > 0 && exchange[0] != "" {
		ob := b.store.GetOrderBook(asset, exchange[0])
		if ob == nil {
			return nil, fmt.Errorf("no order book found for %s on %s", asset.Symbol(), exchange[0])
		}
		return ob, nil
	}

	// Return first available
	orderBooks := b.store.GetOrderBooks(asset)
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

// GetKlines returns historical kline/candlestick data
func (b *baseMarketService) GetKlines(asset portfolio.Asset, exchange connector.ExchangeName, interval string, limit int) []connector.Kline {
	return b.store.GetKlines(asset, exchange, interval, limit)
}

// GetTradableQuantity calculates available liquidity
func (b *baseMarketService) GetTradableQuantity(ctx context.Context, asset portfolio.Asset, opts ...analytics.LiquidityOptions) numerical.Decimal {
	options := DefaultLiquidityOptions()
	if len(opts) > 0 {
		options = opts[0]
	}

	orderBook, err := b.OrderBook(ctx, asset)
	if err != nil {
		return numerical.Zero()
	}

	return getTradableQuantity(ctx, orderBook, options)
}
