package main

import (
	"fmt"

	"github.com/backtesting-org/kronos-sdk/pkg/types/connector"
	"github.com/backtesting-org/kronos-sdk/pkg/types/kronos"
	"github.com/backtesting-org/kronos-sdk/pkg/types/kronos/analytics"
	"github.com/backtesting-org/kronos-sdk/pkg/types/strategy"
	"github.com/shopspring/decimal"
)

// ExampleStrategy demonstrates using the Kronos SDK
type ExampleStrategy struct {
	strategy.BaseStrategy
	k kronos.Kronos
}

// NewExampleStrategy creates a new example strategy
func NewExampleStrategy(k kronos.Kronos) strategy.Strategy {
	return &ExampleStrategy{k: k}
}

// GetSignals demonstrates all Kronos API features
func (s *ExampleStrategy) GetSignals() ([]*strategy.Signal, error) {
	s.k.Log().Info("🔍 Starting signal generation for Example Strategy")

	// Create assets using helper
	btc := s.k.Asset("BTC")
	eth := s.k.Asset("ETH")

	// === INDICATOR EXAMPLES ===

	// Get SMA - simple, one-line call
	sma20, err := s.k.Indicators().SMA(btc, 20)
	if err != nil {
		s.k.Log().Failed("ExampleStrategy", btc.Symbol(), "Failed to calculate SMA: %v", err)
	} else {
		s.k.Log().Debug("ExampleStrategy", btc.Symbol(), "BTC SMA(20): %s", sma20.String())
	}

	// Get EMA with custom exchange
	ema50, err := s.k.Indicators().EMA(btc, 50, analytics.IndicatorOptions{
		Exchange: connector.Binance,
		Interval: "4h",
	})
	if err != nil {
		s.k.Log().Failed("ExampleStrategy", btc.Symbol(), "Failed to calculate EMA: %v", err)
	} else {
		s.k.Log().Debug("ExampleStrategy", btc.Symbol(), "BTC EMA(50) on Binance: %s", ema50.String())
	}

	// Get RSI
	rsi, err := s.k.Indicators().RSI(btc, 14)
	if err != nil {
		s.k.Log().Failed("ExampleStrategy", btc.Symbol(), "Failed to calculate RSI: %v", err)
	} else {
		s.k.Log().Debug("ExampleStrategy", btc.Symbol(), "BTC RSI(14): %s", rsi.String())
	}

	// Get MACD
	macd, err := s.k.Indicators().MACD(btc, 12, 26, 9)
	if err != nil {
		s.k.Log().Failed("ExampleStrategy", btc.Symbol(), "Failed to calculate MACD: %v", err)
	} else {
		s.k.Log().Debug("ExampleStrategy", btc.Symbol(), "BTC MACD - MACD: %s, Signal: %s, Histogram: %s",
			macd.MACD.String(), macd.Signal.String(), macd.Histogram.String())
	}

	// Get Bollinger Bands
	bb, err := s.k.Indicators().BollingerBands(btc, 20, 2.0)
	if err != nil {
		s.k.Log().Failed("ExampleStrategy", btc.Symbol(), "Failed to calculate Bollinger Bands: %v", err)
	} else {
		s.k.Log().Debug("ExampleStrategy", btc.Symbol(), "BTC Bollinger Bands - Upper: %s, Middle: %s, Lower: %s",
			bb.Upper.String(), bb.Middle.String(), bb.Lower.String())
	}

	// === MARKET DATA EXAMPLES ===

	// Get current price - simple
	price, err := s.k.Market().Price(btc)
	if err != nil {
		s.k.Log().Failed("ExampleStrategy", btc.Symbol(), "Failed to get price: %v", err)
	} else {
		s.k.Log().MarketCondition("BTC Price: %s", price.String())
	}

	// Get prices across all exchanges
	prices := s.k.Market().Prices(btc)
	for exchange, p := range prices {
		s.k.Log().Debug("ExampleStrategy", btc.Symbol(), "Price on %s: %s", exchange, p.String())
	}

	// Get funding rates
	fundingRates := s.k.Market().FundingRates(btc)
	for exchange, rate := range fundingRates {
		s.k.Log().Debug("ExampleStrategy", btc.Symbol(), "Funding rate on %s: %s (Next: %s)",
			exchange, rate.CurrentRate.String(), rate.NextFundingTime.String())
	}

	// Find arbitrage opportunities (minimum 10 bps spread)
	arbOpps := s.k.Market().FindArbitrage(btc, decimal.NewFromInt(10))
	for _, opp := range arbOpps {
		s.k.Log().Opportunity("ExampleStrategy", btc.Symbol(),
			"Arbitrage: Buy %s @ %s, Sell %s @ %s, Spread: %s bps",
			opp.BuyExchange, opp.BuyPrice.String(),
			opp.SellExchange, opp.SellPrice.String(),
			opp.SpreadBps.String())
	}

	// === ANALYTICS EXAMPLES ===

	// Calculate volatility
	vol, err := s.k.Analytics().Volatility(btc, 24)
	if err != nil {
		s.k.Log().Failed("ExampleStrategy", btc.Symbol(), "Failed to calculate volatility: %v", err)
	} else {
		s.k.Log().MarketCondition("BTC Volatility (24h): %s%%", vol.String())
	}

	// Analyze trend
	trend, err := s.k.Analytics().Trend(btc, 50)
	if err != nil {
		s.k.Log().Failed("ExampleStrategy", btc.Symbol(), "Failed to analyze trend: %v", err)
	} else {
		s.k.Log().MarketCondition("BTC Trend: %s (Strength: %s%%, Slope: %s)",
			trend.Direction, trend.Strength.String(), trend.Slope.String())
	}

	// Analyze volume
	volumeAnalysis, err := s.k.Analytics().VolumeAnalysis(btc, 24)
	if err != nil {
		s.k.Log().Failed("ExampleStrategy", btc.Symbol(), "Failed to analyze volume: %v", err)
	} else {
		spikeStr := ""
		if volumeAnalysis.IsVolumeSpike {
			spikeStr = " [SPIKE]"
		}
		s.k.Log().MarketCondition("BTC Volume: Current %s, Avg %s, Ratio: %sx%s",
			volumeAnalysis.CurrentVolume.String(), volumeAnalysis.AverageVolume.String(),
			volumeAnalysis.VolumeRatio.String(), spikeStr)
	}

	// Get price change
	priceChange, err := s.k.Analytics().GetPriceChange(btc, 24)
	if err != nil {
		s.k.Log().Failed("ExampleStrategy", btc.Symbol(), "Failed to get price change: %v", err)
	} else {
		s.k.Log().MarketCondition("BTC 24h Change: %s%% (High: %s, Low: %s)",
			priceChange.ChangePercent.String(), priceChange.HighPrice.String(), priceChange.LowPrice.String())
	}

	// === SIGNAL GENERATION LOGIC ===

	var signals []*strategy.Signal

	// Example: Simple SMA crossover strategy
	if !sma20.IsZero() && !price.IsZero() {
		if price.GreaterThan(sma20) {
			// Price above SMA - bullish signal
			s.k.Log().Opportunity("ExampleStrategy", btc.Symbol(),
				"Price %s > SMA(20) %s - Bullish signal", price.String(), sma20.String())

			signal := s.k.Signal(s.GetName()).
				Buy(btc, connector.Binance, decimal.NewFromInt(1)).
				Build()

			signals = append(signals, signal)
		}
	}

	// Example: RSI oversold strategy
	if !rsi.IsZero() {
		oversoldThreshold := decimal.NewFromInt(30)
		if rsi.LessThan(oversoldThreshold) {
			s.k.Log().Opportunity("ExampleStrategy", eth.Symbol(),
				"RSI %s < 30 - Oversold signal", rsi.String())

			signal := s.k.Signal(s.GetName()).
				Buy(eth, connector.Binance, decimal.NewFromInt(5)).
				Build()
			signals = append(signals, signal)
		}
	}

	s.k.Log().Info("✅ Signal generation complete - Generated %d signals", len(signals))

	return signals, nil
}

// Interface compliance
func (s *ExampleStrategy) GetName() strategy.StrategyName {
	return strategy.StrategyName("Example Strategy")
}

func (s *ExampleStrategy) GetDescription() string {
	return "Example strategy demonstrating Kronos API usage"
}

func (s *ExampleStrategy) GetRiskLevel() strategy.RiskLevel {
	return strategy.RiskLevelMedium
}

func (s *ExampleStrategy) GetStrategyType() strategy.StrategyType {
	return strategy.StrategyTypeTechnical
}

func main() {
	fmt.Println("Kronos SDK Example - See strategy implementation above")
	fmt.Println("This demonstrates the user-friendly Kronos context API")
}
