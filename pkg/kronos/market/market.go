package market

import (
	"fmt"
	"sort"

	"github.com/backtesting-org/kronos-sdk/pkg/types/connector"
	"github.com/backtesting-org/kronos-sdk/pkg/types/portfolio"
	"github.com/backtesting-org/kronos-sdk/pkg/types/portfolio/store"
	"github.com/shopspring/decimal"
)

// MarketService provides user-friendly methods for market data access.
// All methods handle data fetching internally.
type MarketService struct {
	store store.Store
}

// NewMarketService creates a new MarketService
func NewMarketService(store store.Store) *MarketService {
	return &MarketService{
		store: store,
	}
}

// MarketOptions configures market data queries
type MarketOptions struct {
	Exchange       connector.ExchangeName // Optional: defaults to first available exchange
	InstrumentType connector.Instrument   // Optional: defaults to perpetual
}

// GetAllAssetsWithFundingRates returns all assets that have funding rate data.
func (s *MarketService) GetAllAssetsWithFundingRates() []portfolio.Asset {
	return s.store.GetAllAssetsWithFundingRates()
}

// FundingRates returns funding rates for an asset across all exchanges.
// Returns a map of exchange name to funding rate.
func (s *MarketService) FundingRates(asset portfolio.Asset) map[connector.ExchangeName]connector.FundingRate {
	fundingMap := s.store.GetFundingRatesForAsset(asset)
	if fundingMap == nil {
		return make(map[connector.ExchangeName]connector.FundingRate)
	}
	return fundingMap
}

// FundingRate returns the funding rate for an asset on a specific exchange.
func (s *MarketService) FundingRate(asset portfolio.Asset, exchange connector.ExchangeName) (*connector.FundingRate, error) {
	rate := s.store.GetFundingRate(asset, exchange)
	if rate == nil {
		return nil, fmt.Errorf("no funding rate found for %s on %s", asset.Symbol(), exchange)
	}
	return rate, nil
}

// Price returns the current price for an asset.
// If exchange is not specified in opts, returns price from first available exchange.
func (s *MarketService) Price(asset portfolio.Asset, opts ...MarketOptions) (decimal.Decimal, error) {
	options := s.parseOptions(opts...)

	if options.Exchange != "" {
		// Get price from specific exchange
		price := s.store.GetAssetPrice(asset, options.Exchange)
		if price == nil {
			return decimal.Zero, fmt.Errorf("no price found for %s on %s", asset.Symbol(), options.Exchange)
		}
		return price.Price, nil
	}

	// Get price from any available exchange
	priceMap := s.store.GetAssetPrices(asset)
	if len(priceMap) == 0 {
		return decimal.Zero, fmt.Errorf("no price data available for %s", asset.Symbol())
	}

	// Return first available price
	for _, price := range priceMap {
		return price.Price, nil
	}

	return decimal.Zero, fmt.Errorf("no price found for %s", asset.Symbol())
}

// Prices returns prices for an asset across all exchanges.
func (s *MarketService) Prices(asset portfolio.Asset) map[connector.ExchangeName]decimal.Decimal {
	priceMap := s.store.GetAssetPrices(asset)
	result := make(map[connector.ExchangeName]decimal.Decimal)

	for exchange, price := range priceMap {
		result[exchange] = price.Price
	}

	return result
}

// OrderBook returns the order book for an asset.
// If exchange is not specified, returns order book from first available exchange.
func (s *MarketService) OrderBook(asset portfolio.Asset, opts ...MarketOptions) (*connector.OrderBook, error) {
	options := s.parseOptions(opts...)

	if options.Exchange != "" {
		// Get order book from specific exchange
		ob := s.store.GetOrderBook(asset, options.Exchange, options.InstrumentType)
		if ob == nil {
			return nil, fmt.Errorf("no order book found for %s on %s", asset.Symbol(), options.Exchange)
		}
		return ob, nil
	}

	// Get order book from first available exchange
	orderBooks := s.store.GetOrderBooks(asset)
	if len(orderBooks) == 0 {
		return nil, fmt.Errorf("no order book data available for %s", asset.Symbol())
	}

	// Return first available order book for the specified instrument type
	for _, instrumentMap := range orderBooks {
		if ob, exists := instrumentMap[options.InstrumentType]; exists && ob != nil {
			return ob, nil
		}
	}

	return nil, fmt.Errorf("no order book found for %s with instrument type %s", asset.Symbol(), options.InstrumentType)
}

// ArbitrageOpportunity represents a price discrepancy across exchanges
type ArbitrageOpportunity struct {
	Asset         portfolio.Asset
	BuyExchange   connector.ExchangeName
	SellExchange  connector.ExchangeName
	BuyPrice      decimal.Decimal
	SellPrice     decimal.Decimal
	SpreadBps     decimal.Decimal // Spread in basis points
	SpreadPercent decimal.Decimal // Spread as percentage
}

// FindArbitrage finds arbitrage opportunities for an asset across exchanges.
// Returns opportunities sorted by spread (highest first).
func (s *MarketService) FindArbitrage(asset portfolio.Asset, minSpreadBps decimal.Decimal) []ArbitrageOpportunity {
	priceMap := s.store.GetAssetPrices(asset)
	if len(priceMap) < 2 {
		return nil // Need at least 2 exchanges for arbitrage
	}

	// Convert map to sorted slice for consistent comparison
	type exchangePrice struct {
		exchange connector.ExchangeName
		price    decimal.Decimal
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

	var opportunities []ArbitrageOpportunity

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
			spreadPercent := spread.Div(buyPrice).Mul(decimal.NewFromInt(100))
			spreadBps := spreadPercent.Mul(decimal.NewFromInt(100))

			// Only include if spread exceeds minimum
			if spreadBps.GreaterThanOrEqual(minSpreadBps) {
				opportunities = append(opportunities, ArbitrageOpportunity{
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
func (s *MarketService) parseOptions(opts ...MarketOptions) MarketOptions {
	if len(opts) > 0 {
		options := opts[0]
		if options.InstrumentType == "" {
			options.InstrumentType = connector.TypePerpetual
		}
		return options
	}
	return MarketOptions{
		InstrumentType: connector.TypePerpetual,
	}
}
