package market

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/backtesting-org/kronos-sdk/pkg/monitoring/profiling"
	"github.com/backtesting-org/kronos-sdk/pkg/types/connector"
	"github.com/backtesting-org/kronos-sdk/pkg/types/data/stores/market"
	"github.com/backtesting-org/kronos-sdk/pkg/types/kronos/analytics"
	"github.com/backtesting-org/kronos-sdk/pkg/types/kronos/numerical"
	"github.com/backtesting-org/kronos-sdk/pkg/types/portfolio"
)

// marketService is the concrete implementation of analytics.Market.
type marketService struct {
	store market.MarketStore
}

// NewMarketService creates a new market service.
func NewMarketService(store market.MarketStore) analytics.Market {
	return &marketService{
		store: store,
	}
}

// GetAllAssetsWithFundingRates returns all assets that have funding rate data.
func (s *marketService) GetAllAssetsWithFundingRates(ctx context.Context) []portfolio.Asset {
	return nil

	//return s.store.GetAllAssetsWithFundingRates()
}

// FundingRates returns funding rates for an asset across all exchanges.
// Returns a map of exchange name to funding rate.
func (s *marketService) FundingRates(ctx context.Context, asset portfolio.Asset) map[connector.ExchangeName]connector.FundingRate {
	return nil

	//fundingMap := s.store.GetFundingRatesForAsset(asset)
	//if fundingMap == nil {
	//	return make(map[connector.ExchangeName]connector.FundingRate)
	//}
	//return fundingMap
}

// FundingRate returns the funding rate for an asset on a specific exchange.
func (s *marketService) FundingRate(ctx context.Context, asset portfolio.Asset, exchange connector.ExchangeName) (*connector.FundingRate, error) {

	return nil, nil
	//rate := s.store.GetFundingRate(asset, exchange)
	//if rate == nil {
	//	return nil, fmt.Errorf("no funding rate found for %s on %s", asset.Symbol(), exchange)
	//}
	//return rate, nil
}

// Price returns the current price for an asset.
// If exchange is not specified in opts, returns price from first available exchange.
func (s *marketService) Price(ctx context.Context, asset portfolio.Asset, opts ...analytics.MarketOptions) (numerical.Decimal, error) {
	options := s.parseOptions(opts...)

	if options.Exchange != "" {
		// Get price from specific exchange
		price := s.store.GetAssetPrice(asset, options.Exchange)
		if price == nil {
			return numerical.Zero(), fmt.Errorf("no price found for %s on %s", asset.Symbol(), options.Exchange)
		}
		return price.Price, nil
	}

	// Get price from any available exchange
	priceMap := s.store.GetAssetPrices(asset)
	if len(priceMap) == 0 {
		return numerical.Zero(), fmt.Errorf("no price data available for %s", asset.Symbol())
	}

	// Return first available price
	for _, price := range priceMap {
		return price.Price, nil
	}

	return numerical.Zero(), fmt.Errorf("no price found for %s", asset.Symbol())
}

// Prices returns prices for an asset across all exchanges.
func (s *marketService) Prices(ctx context.Context, asset portfolio.Asset) map[connector.ExchangeName]numerical.Decimal {
	priceMap := s.store.GetAssetPrices(asset)
	result := make(map[connector.ExchangeName]numerical.Decimal)

	for exchange, price := range priceMap {
		result[exchange] = price.Price
	}

	return result
}

// OrderBook returns the order book for an asset.
// If exchange is not specified, returns order book from first available exchange.
func (s *marketService) OrderBook(ctx context.Context, asset portfolio.Asset, opts ...analytics.MarketOptions) (*connector.OrderBook, error) {
	//options := s.parseOptions(opts...)
	//
	//if options.Exchange != "" {
	//	// Get order book from specific exchange
	//	ob := s.store.GetOrderBook(asset, options.Exchange, options.InstrumentType)
	//	if ob == nil {
	//		return nil, fmt.Errorf("no order book found for %s on %s", asset.Symbol(), options.Exchange)
	//	}
	//	return ob, nil
	//}
	//
	//// Get order book from first available exchange
	//orderBooks := s.store.GetOrderBooks(asset)
	//if len(orderBooks) == 0 {
	//	return nil, fmt.Errorf("no order book data available for %s", asset.Symbol())
	//}
	//
	//// Return first available order book for the specified instrument type
	//for _, instrumentMap := range orderBooks {
	//	if ob, exists := instrumentMap[options.InstrumentType]; exists && ob != nil {
	//		return ob, nil
	//	}
	//}

	return nil, nil
}

// FindArbitrage finds arbitrage opportunities for an asset across exchanges.
// Returns opportunities sorted by spread (highest first).
func (s *marketService) FindArbitrage(ctx context.Context, asset portfolio.Asset, minSpreadBps numerical.Decimal) []analytics.ArbitrageOpportunity {
	start := time.Now()
	defer func() {
		if profCtx := profiling.FromContext(ctx); profCtx != nil {
			profCtx.RecordIndicator("FindArbitrage", time.Since(start))
		}
	}()

	priceMap := s.store.GetAssetPrices(asset)
	if len(priceMap) < 2 {
		return nil // Need at least 2 exchanges for arbitrage
	}

	// Convert map to sorted slice for consistent comparison
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

	// Sort by price
	sort.Slice(prices, func(i, j int) bool {
		return prices[i].price.LessThan(prices[j].price)
	})

	var opportunities []analytics.ArbitrageOpportunity

	// Compare lowest price exchanges with highest price exchanges
	for i := 0; i < len(prices); i++ {
		for j := i + 1; j < len(prices); j++ {
			buyPrice := prices[i].price
			sellPrice := prices[j].price

			if buyPrice.IsZero() {
				continue
			}

			// Calculate spread
			spread := sellPrice.Sub(buyPrice)
			spreadPercent := spread.Div(buyPrice).Mul(numerical.NewFromInt(100))
			spreadBps := spreadPercent.Mul(numerical.NewFromInt(100))

			// Only include if spread exceeds minimum
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

	// Sort by spread descending (best opportunities first)
	sort.Slice(opportunities, func(i, j int) bool {
		return opportunities[i].SpreadBps.GreaterThan(opportunities[j].SpreadBps)
	})

	return opportunities
}

// parseOptions extracts options with defaults
func (s *marketService) parseOptions(opts ...analytics.MarketOptions) analytics.MarketOptions {
	if len(opts) > 0 {
		options := opts[0]
		if options.InstrumentType == "" {
			options.InstrumentType = connector.TypePerpetual
		}
		return options
	}
	return analytics.MarketOptions{
		InstrumentType: connector.TypePerpetual,
	}
}
