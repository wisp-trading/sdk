package main

import (
	"fmt"

	"github.com/backtesting-org/kronos-sdk/pkg/kronos"
	"github.com/backtesting-org/kronos-sdk/pkg/types/connector"
	"github.com/backtesting-org/kronos-sdk/pkg/types/logging"
	"github.com/backtesting-org/kronos-sdk/pkg/types/portfolio/store"
	"github.com/backtesting-org/kronos-sdk/pkg/types/strategy"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// ExampleStrategy demonstrates how to use the Kronos context in a strategy
type ExampleStrategy struct {
	*strategy.BaseStrategy
	k *kronos.Kronos // Injected Kronos context
}

// NewExampleStrategy creates a new example strategy with Kronos injection
func NewExampleStrategy(k *kronos.Kronos) *ExampleStrategy {
	return &ExampleStrategy{
		k: k,
	}
}

// GetSignals demonstrates the user-friendly Kronos API
func (s *ExampleStrategy) GetSignals() ([]*strategy.Signal, error) {
	// Define the assets we're trading - simple one-line creation
	btc := s.k.Asset("BTC-PERP")
	eth := s.k.Asset("ETH-PERP")

	s.k.Log().Info("Starting signal generation")

	// === INDICATOR EXAMPLES ===

	// Get SMA - simple, one-line call
	sma20, err := s.k.Indicators.SMA(btc, 20)
	if err != nil {
		s.k.Log().Warn("Failed to calculate SMA", "error", err)
	} else {
		s.k.Log().Info("BTC SMA(20)", "value", sma20.String())
	}

	// Get EMA with custom exchange
	ema50, err := s.k.Indicators.EMA(btc, 50, kronos.IndicatorOptions{
		Exchange: connector.Binance,
		Interval: "4h",
	})
	if err != nil {
		s.k.Log().Warn("Failed to calculate EMA", "error", err)
	} else {
		s.k.Log().Info("BTC EMA(50) on Binance", "value", ema50.String())
	}

	// Get RSI
	rsi, err := s.k.Indicators.RSI(btc, 14)
	if err != nil {
		s.k.Log().Warn("Failed to calculate RSI", "error", err)
	} else {
		s.k.Log().Info("BTC RSI(14)", "value", rsi.String())
	}

	// Get MACD
	macd, err := s.k.Indicators.MACD(btc, 12, 26, 9)
	if err != nil {
		s.k.Log().Warn("Failed to calculate MACD", "error", err)
	} else {
		s.k.Log().Info("BTC MACD",
			"macd", macd.MACD.String(),
			"signal", macd.Signal.String(),
			"histogram", macd.Histogram.String(),
		)
	}

	// Get Bollinger Bands
	bb, err := s.k.Indicators.BollingerBands(btc, 20, 2.0)
	if err != nil {
		s.k.Log().Warn("Failed to calculate Bollinger Bands", "error", err)
	} else {
		s.k.Log().Info("BTC Bollinger Bands",
			"upper", bb.Upper.String(),
			"middle", bb.Middle.String(),
			"lower", bb.Lower.String(),
		)
	}

	// === MARKET DATA EXAMPLES ===

	// Get current price - simple
	price, err := s.k.Market.Price(btc)
	if err != nil {
		s.k.Log().Warn("Failed to get price", "error", err)
	} else {
		s.k.Log().Info("BTC Price", "price", price.String())
	}

	// Get prices across all exchanges
	prices := s.k.Market.Prices(btc)
	for exchange, p := range prices {
		s.k.Log().Info("BTC Price by exchange", "exchange", exchange, "price", p.String())
	}

	// Get funding rates
	fundingRates := s.k.Market.FundingRates(btc)
	for exchange, rate := range fundingRates {
		s.k.Log().Info("BTC Funding Rate",
			"exchange", exchange,
			"rate", rate.CurrentRate.String(),
			"next_funding", rate.NextFundingTime.String(),
		)
	}

	// Find arbitrage opportunities
	arbOpps := s.k.Market.FindArbitrage(btc)
	for _, opp := range arbOpps {
		s.k.Log().Info("Arbitrage Opportunity Found",
			"buy_exchange", opp.BuyExchange,
			"sell_exchange", opp.SellExchange,
			"spread_bps", opp.SpreadBps.String(),
			"estimated_profit_bps", opp.EstimatedProfit.String(),
		)
	}

	// Get best bid/ask across exchanges
	bestBidAsk, err := s.k.Market.GetBestBidAsk(btc)
	if err != nil {
		s.k.Log().Warn("Failed to get best bid/ask", "error", err)
	} else {
		s.k.Log().Info("Best Bid/Ask",
			"bid", bestBidAsk.BestBid.String(),
			"ask", bestBidAsk.BestAsk.String(),
			"spread_bps", bestBidAsk.SpreadBps.String(),
		)
	}

	// === ANALYTICS EXAMPLES ===

	// Calculate volatility
	vol, err := s.k.Analytics.Volatility(btc, 24)
	if err != nil {
		s.k.Log().Warn("Failed to calculate volatility", "error", err)
	} else {
		s.k.Log().Info("BTC Volatility (24h)", "volatility", vol.String()+"%")
	}

	// Analyze trend
	trend, err := s.k.Analytics.Trend(btc, 50)
	if err != nil {
		s.k.Log().Warn("Failed to analyze trend", "error", err)
	} else {
		s.k.Log().Info("BTC Trend Analysis",
			"direction", trend.Direction,
			"strength", trend.Strength.String(),
			"slope", trend.Slope.String(),
		)
	}

	// Analyze volume
	volumeAnalysis, err := s.k.Analytics.VolumeAnalysis(btc, 24)
	if err != nil {
		s.k.Log().Warn("Failed to analyze volume", "error", err)
	} else {
		s.k.Log().Info("BTC Volume Analysis",
			"current", volumeAnalysis.CurrentVolume.String(),
			"average", volumeAnalysis.AverageVolume.String(),
			"ratio", volumeAnalysis.VolumeRatio.String(),
			"is_spike", volumeAnalysis.IsVolumeSpike,
			"trend", volumeAnalysis.VolumeTrend,
		)
	}

	// Get price change
	priceChange, err := s.k.Analytics.GetPriceChange(btc, 24)
	if err != nil {
		s.k.Log().Warn("Failed to get price change", "error", err)
	} else {
		s.k.Log().Info("BTC Price Change (24h)",
			"change_percent", priceChange.ChangePercent.String()+"%",
			"high", priceChange.HighPrice.String(),
			"low", priceChange.LowPrice.String(),
		)
	}

	// === SIGNAL GENERATION LOGIC ===

	var signals []*strategy.Signal

	// Example: Simple SMA crossover strategy
	if !sma20.IsZero() && !price.IsZero() {
		if price.GreaterThan(sma20) {
			// Price above SMA - bullish signal
			signal := &strategy.Signal{
				ID:       uuid.New(),
				Strategy: strategy.StrategyName("Example Strategy"),
				Actions: []strategy.TradeAction{
					{
						Action:   strategy.ActionBuy,
						Asset:    btc,
						Exchange: connector.Binance,
						Quantity: decimal.NewFromInt(1),
						Price:    decimal.Zero, // Market order
					},
				},
			}
			signals = append(signals, signal)
		}
	}

	// Example: RSI oversold strategy
	if !rsi.IsZero() {
		oversoldThreshold := decimal.NewFromInt(30)
		if rsi.LessThan(oversoldThreshold) {
			signal := &strategy.Signal{
				ID:       uuid.New(),
				Strategy: strategy.StrategyName("Example Strategy"),
				Actions: []strategy.TradeAction{
					{
						Action:   strategy.ActionBuy,
						Asset:    eth,
						Exchange: connector.Binance,
						Quantity: decimal.NewFromInt(10),
						Price:    decimal.Zero,
					},
				},
			}
			signals = append(signals, signal)
		}
	}

	s.k.Log().Info("Signal generation complete", "signal_count", len(signals))

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

// Example of how the orchestrator would use KronosExecutor
func demonstrateExecutor(store store.Store, logger logging.ApplicationLogger) {
	// Create executor with trade capabilities
	executor := kronos.NewKronosExecutor(store, logger)

	// Can use all read operations from base Kronos
	btc := executor.Asset("BTC-PERP")
	price, _ := executor.Market.Price(btc)
	fmt.Printf("BTC Price: %s\n", price.String())

	// Can also execute trades
	result, err := executor.Trade.Buy(
		btc,
		connector.Binance,
		decimal.NewFromInt(1),
		kronos.TradeOptions{
			OrderType: kronos.OrderTypeMarket,
		},
	)
	if err != nil {
		fmt.Printf("Trade failed: %v\n", err)
	} else {
		fmt.Printf("Trade executed: %s\n", result.OrderID)
	}
}

func main() {
	fmt.Println("Kronos SDK Example - See strategy implementation above")
	fmt.Println("This demonstrates the user-friendly Kronos context API")
}
