package kronos

import (
	"fmt"
	"sort"

	"github.com/backtesting-org/kronos-sdk/pkg/types/connector"
	"github.com/backtesting-org/kronos-sdk/pkg/types/logging"
	"github.com/backtesting-org/kronos-sdk/pkg/types/portfolio"
	"github.com/backtesting-org/kronos-sdk/pkg/types/portfolio/store"
	"github.com/shopspring/decimal"
)

// MarketService provides user-friendly methods for market data access.
// All methods handle data fetching internally.
type MarketService struct {
	store  store.Store
	logger logging.ApplicationLogger
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
		s.logger.Debug("No funding rates found for asset", "asset", asset.Symbol())
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

	// Get order book from any available exchange
	obMap := s.store.GetOrderBooks(asset)
	if len(obMap) == 0 {
		return nil, fmt.Errorf("no order book data available for %s", asset.Symbol())
	}

	// Return first available order book
	for _, instruments := range obMap {
		for _, ob := range instruments {
			return ob, nil
		}
	}

	return nil, fmt.Errorf("no order book found for %s", asset.Symbol())
}

// OrderBooks returns order books for an asset across all exchanges.
func (s *MarketService) OrderBooks(asset portfolio.Asset) map[connector.ExchangeName]*connector.OrderBook {
	obMap := s.store.GetOrderBooks(asset)
	result := make(map[connector.ExchangeName]*connector.OrderBook)

	for exchange, instruments := range obMap {
		// Get the first instrument type available for each exchange
		for _, ob := range instruments {
			result[exchange] = ob
			break
		}
	}

	return result
}

// ArbitrageOpportunity represents a potential arbitrage opportunity
type ArbitrageOpportunity struct {
	Asset           portfolio.Asset
	BuyExchange     connector.ExchangeName
	SellExchange    connector.ExchangeName
	BuyPrice        decimal.Decimal
	SellPrice       decimal.Decimal
	SpreadBps       decimal.Decimal // Spread in basis points
	SpreadPercent   decimal.Decimal // Spread as percentage
	EstimatedProfit decimal.Decimal // Spread minus estimated fees
}

// FindArbitrage compares prices across all exchanges for an asset and returns arbitrage opportunities.
// Returns slice of opportunities sorted by spread (highest first).
func (s *MarketService) FindArbitrage(asset portfolio.Asset, estimatedFeeBps ...decimal.Decimal) []ArbitrageOpportunity {
	priceMap := s.store.GetAssetPrices(asset)
	if len(priceMap) < 2 {
		// Need at least 2 exchanges for arbitrage
		return nil
	}

	// Default fee of 10 bps (0.1%) if not specified
	feeBps := decimal.NewFromInt(10)
	if len(estimatedFeeBps) > 0 {
		feeBps = estimatedFeeBps[0]
	}

	var opportunities []ArbitrageOpportunity

	// Compare all exchange pairs
	exchanges := make([]connector.ExchangeName, 0, len(priceMap))
	for exchange := range priceMap {
		exchanges = append(exchanges, exchange)
	}

	for i := 0; i < len(exchanges); i++ {
		for j := i + 1; j < len(exchanges); j++ {
			ex1 := exchanges[i]
			ex2 := exchanges[j]
			price1 := priceMap[ex1].Price
			price2 := priceMap[ex2].Price

			// Skip if either price is zero
			if price1.IsZero() || price2.IsZero() {
				continue
			}

			var buyEx, sellEx connector.ExchangeName
			var buyPrice, sellPrice decimal.Decimal

			if price1.LessThan(price2) {
				buyEx = ex1
				sellEx = ex2
				buyPrice = price1
				sellPrice = price2
			} else {
				buyEx = ex2
				sellEx = ex1
				buyPrice = price2
				sellPrice = price1
			}

			// Calculate spread
			spread := sellPrice.Sub(buyPrice)
			spreadPercent := spread.Div(buyPrice).Mul(decimal.NewFromInt(100))
			spreadBps := spreadPercent.Mul(decimal.NewFromInt(100))

			// Calculate estimated profit (spread minus fees)
			// Assuming fee on both buy and sell
			totalFeeBps := feeBps.Mul(decimal.NewFromInt(2))
			profitBps := spreadBps.Sub(totalFeeBps)

			// Only include if there's potential profit after fees
			if profitBps.GreaterThan(decimal.Zero) {
				opportunities = append(opportunities, ArbitrageOpportunity{
					Asset:           asset,
					BuyExchange:     buyEx,
					SellExchange:    sellEx,
					BuyPrice:        buyPrice,
					SellPrice:       sellPrice,
					SpreadBps:       spreadBps,
					SpreadPercent:   spreadPercent,
					EstimatedProfit: profitBps,
				})
			}
		}
	}

	// Sort by estimated profit (highest first)
	sort.Slice(opportunities, func(i, j int) bool {
		return opportunities[i].EstimatedProfit.GreaterThan(opportunities[j].EstimatedProfit)
	})

	return opportunities
}

// GetAllAssetsWithOrderBooks returns all assets that have order book data.
func (s *MarketService) GetAllAssetsWithOrderBooks() []portfolio.Asset {
	return s.store.GetAllAssetsWithOrderBooks()
}

// GetAllAssetsWithPrices returns all unique assets that have price data.
func (s *MarketService) GetAllAssetsWithPrices() []portfolio.Asset {
	// Get all assets with order books as a proxy
	// In a real implementation, we'd need a dedicated method on store
	assetsWithFunding := s.store.GetAllAssetsWithFundingRates()
	assetsWithOrderBooks := s.store.GetAllAssetsWithOrderBooks()

	// Combine and deduplicate
	assetMap := make(map[string]portfolio.Asset)
	for _, asset := range assetsWithFunding {
		assetMap[asset.Symbol()] = asset
	}
	for _, asset := range assetsWithOrderBooks {
		assetMap[asset.Symbol()] = asset
	}

	result := make([]portfolio.Asset, 0, len(assetMap))
	for _, asset := range assetMap {
		result = append(result, asset)
	}

	return result
}

// BestBidAsk returns the best bid and ask prices for an asset across all exchanges.
type BestBidAsk struct {
	BestBid     decimal.Decimal
	BestAsk     decimal.Decimal
	BidExchange connector.ExchangeName
	AskExchange connector.ExchangeName
	Spread      decimal.Decimal
	SpreadBps   decimal.Decimal
}

// GetBestBidAsk finds the best bid and ask prices across all exchanges.
func (s *MarketService) GetBestBidAsk(asset portfolio.Asset) (*BestBidAsk, error) {
	obMap := s.store.GetOrderBooks(asset)
	if len(obMap) == 0 {
		return nil, fmt.Errorf("no order book data available for %s", asset.Symbol())
	}

	var bestBid, bestAsk decimal.Decimal
	var bidExchange, askExchange connector.ExchangeName

	for exchange, instruments := range obMap {
		for _, ob := range instruments {
			// Get top bid and ask
			if len(ob.Bids) > 0 {
				topBid := ob.Bids[0].Price
				if bestBid.IsZero() || topBid.GreaterThan(bestBid) {
					bestBid = topBid
					bidExchange = exchange
				}
			}

			if len(ob.Asks) > 0 {
				topAsk := ob.Asks[0].Price
				if bestAsk.IsZero() || topAsk.LessThan(bestAsk) {
					bestAsk = topAsk
					askExchange = exchange
				}
			}
		}
	}

	if bestBid.IsZero() || bestAsk.IsZero() {
		return nil, fmt.Errorf("incomplete order book data for %s", asset.Symbol())
	}

	spread := bestAsk.Sub(bestBid)
	spreadBps := spread.Div(bestBid).Mul(decimal.NewFromInt(10000))

	return &BestBidAsk{
		BestBid:     bestBid,
		BestAsk:     bestAsk,
		BidExchange: bidExchange,
		AskExchange: askExchange,
		Spread:      spread,
		SpreadBps:   spreadBps,
	}, nil
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
