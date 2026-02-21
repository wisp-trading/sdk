package market

import (
	"context"
	"fmt"

	"github.com/wisp-trading/sdk/pkg/types/connector"
	storeTypes "github.com/wisp-trading/sdk/pkg/types/data/stores/market"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

// baseMarketService implements ONLY operations that use the core MarketStore interface
// It does NOT use extensions (OrderBook, Klines, etc.) to maintain true separation.
// Market-specific services embed this to inherit price functionality only.
type baseMarketService struct {
	store storeTypes.MarketStore
}

func newBaseMarketService(store storeTypes.MarketStore) baseMarketService {
	return baseMarketService{store: store}
}

// Price returns the current price for an asset (uses only core MarketStore.GetPairPrices)
func (b *baseMarketService) Price(ctx context.Context, asset portfolio.Pair, exchange ...connector.ExchangeName) (numerical.Decimal, error) {
	priceMap := b.store.GetPairPrices(asset)

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

// Prices returns prices across all exchanges (uses only core MarketStore.GetPairPrices)
func (b *baseMarketService) Prices(ctx context.Context, asset portfolio.Pair) map[connector.ExchangeName]numerical.Decimal {
	result := make(map[connector.ExchangeName]numerical.Decimal)
	priceMap := b.store.GetPairPrices(asset)
	for exchange, price := range priceMap {
		result[exchange] = price.Price
	}
	return result
}
