package market

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/wisp-trading/sdk/pkg/monitoring/profiling"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	marketTypes "github.com/wisp-trading/sdk/pkg/types/data/stores/market"
	perpTypes "github.com/wisp-trading/sdk/pkg/types/data/stores/market/perp"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	"github.com/wisp-trading/sdk/pkg/types/wisp/analytics"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

// marketService is the concrete implementation of analytics.Market.
// It provides read-only access to market data via the registry.
type marketService struct {
	registry marketTypes.MarketRegistry
	spot     *spotMarketService
	perp     *perpMarketService
}

// NewMarketService creates a new market service using the market registry.
func NewMarketService(registry marketTypes.MarketRegistry) analytics.Market {
	spotStore := registry.Get(marketTypes.MarketTypeSpot)
	perpStore := registry.Get(marketTypes.MarketTypePerp)

	return &marketService{
		registry: registry,
		spot:     newSpotMarketService(spotStore),
		perp:     newPerpMarketService(perpStore.(perpTypes.MarketStore)),
	}
}

// Spot returns the spot market service for spot-specific operations
func (s *marketService) Spot() analytics.SpotMarket {
	return s.spot
}

// Perp returns the perp market service for perp-specific operations
func (s *marketService) Perp() analytics.PerpMarket {
	return s.perp
}

// Price returns the current price for an asset from any exchange (spot or perp).
// If exchange is specified, returns price from that exchange.
// Otherwise returns first available price.
func (s *marketService) Price(ctx context.Context, asset portfolio.Pair, exchange ...connector.ExchangeName) (numerical.Decimal, error) {
	var targetExchange connector.ExchangeName
	if len(exchange) > 0 {
		targetExchange = exchange[0]
	}

	// Iterate all registered market stores
	for _, store := range s.registry.GetAll() {
		prices := store.GetAssetPrices(asset)

		if targetExchange != "" {
			if price, exists := prices[targetExchange]; exists {
				return price.Price, nil
			}
		} else if len(prices) > 0 {
			// Return first available price
			for _, price := range prices {
				return price.Price, nil
			}
		}
	}

	if targetExchange != "" {
		return numerical.Zero(), fmt.Errorf("no price found for %s on %s", asset.Symbol(), targetExchange)
	}
	return numerical.Zero(), fmt.Errorf("no price data available for %s", asset.Symbol())
}

// Prices returns prices for an asset across all exchanges (both spot and perp).
func (s *marketService) Prices(ctx context.Context, asset portfolio.Pair) map[connector.ExchangeName]numerical.Decimal {
	result := make(map[connector.ExchangeName]numerical.Decimal)

	// Iterate all registered market stores
	for _, store := range s.registry.GetAll() {
		prices := store.GetAssetPrices(asset)
		for exchange, price := range prices {
			result[exchange] = price.Price
		}
	}

	return result
}

// GetKlines returns historical kline data for an asset on the specified exchange.
// Automatically searches all registered market stores to find which one has this exchange.
// The user doesn't need to know whether the exchange is spot, perp, futures, etc.
func (s *marketService) Klines(asset portfolio.Pair, exchange connector.ExchangeName, interval string, limit int) []connector.Kline {
	// Iterate all registered stores to find which one has data for this exchange
	for _, store := range s.registry.GetAll() {
		klines := store.GetKlines(asset, exchange, interval, limit)
		if len(klines) > 0 {
			return klines
		}
	}

	// No data found in any store
	return nil
}

// GetOrderBook returns the order book for an asset on the specified exchange.
// Automatically searches all registered market stores to find which one has this exchange.
// If no exchange is specified, returns the first available orderbook.
func (s *marketService) OrderBook(ctx context.Context, asset portfolio.Pair, exchange ...connector.ExchangeName) (*connector.OrderBook, error) {
	var targetExchange connector.ExchangeName
	if len(exchange) > 0 {
		targetExchange = exchange[0]
	}

	// Iterate all registered stores to find which one has data for this exchange/asset
	for _, store := range s.registry.GetAll() {
		if targetExchange != "" {
			// Looking for specific exchange
			ob := store.GetOrderBook(asset, targetExchange)
			if ob != nil {
				return ob, nil
			}
		} else {
			// Get first available orderbook
			orderBooks := store.GetOrderBooks(asset)
			for _, ob := range orderBooks {
				if ob != nil {
					return ob, nil
				}
			}
		}
	}

	if targetExchange != "" {
		return nil, fmt.Errorf("no order book found for %s on %s", asset.Symbol(), targetExchange)
	}
	return nil, fmt.Errorf("no order book data available for %s", asset.Symbol())
}

// FindArbitrage finds arbitrage opportunities for an asset across exchanges.
// Returns opportunities sorted by spread (highest first).
// Searches all registered market types for arbitrage opportunities.
func (s *marketService) FindArbitrage(ctx context.Context, asset portfolio.Pair, minSpreadBps numerical.Decimal) []analytics.ArbitrageOpportunity {
	start := time.Now()
	defer func() {
		if profCtx := profiling.FromContext(ctx); profCtx != nil {
			profCtx.RecordIndicator("FindArbitrage", time.Since(start))
		}
	}()

	// Combine prices from all registered stores
	priceMap := make(map[connector.ExchangeName]connector.Price)

	for _, store := range s.registry.GetAll() {
		prices := store.GetAssetPrices(asset)
		for exchange, price := range prices {
			priceMap[exchange] = price
		}
	}

	if len(priceMap) < 2 {
		return nil // Need at least 2 exchanges for arbitrage
	}

	// ...existing code for arbitrage calculation...
	type exchangePrice struct {
		exchange connector.ExchangeName
		price    numerical.Decimal
	}

	var prices []exchangePrice
	for exchange, priceData := range priceMap {
		prices = append(prices, exchangePrice{
			exchange: exchange,
			price:    priceData.Price,
		})
	}

	sort.Slice(prices, func(i, j int) bool {
		return prices[i].price.LessThan(prices[j].price)
	})

	var opportunities []analytics.ArbitrageOpportunity

	for i := 0; i < len(prices); i++ {
		for j := i + 1; j < len(prices); j++ {
			buyPrice := prices[i].price
			sellPrice := prices[j].price

			if buyPrice.IsZero() {
				continue
			}

			spread := sellPrice.Sub(buyPrice)
			spreadPercent := spread.Div(buyPrice).Mul(numerical.NewFromInt(100))
			spreadBps := spreadPercent.Mul(numerical.NewFromInt(100))

			if spreadBps.GreaterThanOrEqual(minSpreadBps) {
				opportunities = append(opportunities, analytics.ArbitrageOpportunity{
					Asset:         asset,
					BuyExchange:   prices[i].exchange,
					SellExchange:  prices[j].exchange,
					BuyPrice:      buyPrice,
					SellPrice:     sellPrice,
					SpreadBps:     spreadBps,
					SpreadPercent: spreadPercent,
				})
			}
		}
	}

	sort.Slice(opportunities, func(i, j int) bool {
		return opportunities[i].SpreadBps.GreaterThan(opportunities[j].SpreadBps)
	})

	return opportunities
}
