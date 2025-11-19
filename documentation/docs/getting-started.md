---
sidebar_position: 2
---

# Getting Started

This guide shows you how to build and deploy a trading strategy with Kronos. We'll use the Kronos CLI to set up your project, test it with backtesting, and deploy it live.

## Installation

Install Kronos via Homebrew:

```bash
brew install backtesting-org/tap/kronos
```

Verify the installation:

```bash
kronos --version
```

## Prerequisites

- **Go 1.21+** - [Install Go](https://go.dev/doc/install)
- Basic Go knowledge
- Understanding of technical indicators (RSI, SMA, etc.)

## Create Your First Strategy

Initialize a new Kronos project:

```bash
kronos init my-strategy
cd my-strategy
```

This creates a project structure with:
- `strategy.go` - Your strategy implementation
- `go.mod` - Go module dependencies
- `config.yaml` - Kronos configuration
- `README.md` - Project documentation

### The Generated Strategy

Create `strategy.go`:

```go
package main

import (
    sdk "github.com/backtesting-org/kronos-sdk/pkg/kronos"
    "github.com/backtesting-org/kronos-sdk/pkg/types/strategy"
    "github.com/shopspring/decimal"
)

type MomentumStrategy struct {
    k *sdk.Kronos
}

func NewMomentum(k *sdk.Kronos) *MomentumStrategy {
    return &MomentumStrategy{k: k}
}

func (s *MomentumStrategy) GetSignals() ([]*strategy.Signal, error) {
    // Get asset reference
    btc := s.k.Asset("BTC")
    
    // Get RSI - Kronos fetches data and calculates automatically
    rsi := s.k.Indicators.RSI(btc, 14)
    
    // Buy when RSI shows oversold
    if rsi.LessThan(decimal.NewFromInt(30)) {
        return []*strategy.Signal{
            s.Signal().
                Buy(btc).
                Quantity(decimal.NewFromFloat(0.1)).
                Reason("RSI oversold at " + rsi.String()).
                Build(),
        }, nil
    }
    
    // Sell when RSI shows overbought
    if rsi.GreaterThan(decimal.NewFromInt(70)) {
        return []*strategy.Signal{
            s.Signal().
                Sell(btc).
                Quantity(decimal.NewFromFloat(0.1)).
                Reason("RSI overbought at " + rsi.String()).
                Build(),
        }, nil
    }
    
    return nil, nil
}

// Interface methods
func (s *MomentumStrategy) GetName() strategy.StrategyName { 
    return "Momentum" 
}

func (s *MomentumStrategy) GetDescription() string { 
    return "RSI momentum strategy" 
}

func (s *MomentumStrategy) GetRiskLevel() strategy.RiskLevel { 
    return strategy.RiskLevelMedium 
}

func (s *MomentumStrategy) GetStrategyType() strategy.StrategyType { 
    return strategy.StrategyTypeTechnical 
}
```

That's it. No exchange setup, no data management, no indicator implementation.

## Understanding the Code

### The Kronos Context

```go
type MomentumStrategy struct {
    k *sdk.Kronos  // Your gateway to everything
}
```

The `*sdk.Kronos` instance gives you access to:
- `k.Asset(symbol)` - Asset references
- `k.Indicators` - All technical indicators
- `k.Market` - Market data (prices, order books, funding)
- `k.Analytics` - Market analytics (volatility, trend, volume)
- `k.Log()` - Logging

### Assets

```go
btc := s.k.Asset("BTC")
eth := s.k.Asset("ETH")
```

Assets are just references. Kronos knows which exchange to use based on your configuration.

### Indicators

```go
// Just ask for what you need
rsi := s.k.Indicators.RSI(btc, 14)
```

Kronos automatically:
1. Fetches price data from the exchange
2. Calculates the RSI
3. Returns the current value

You don't manage any of this.

### Signals

```go
s.Signal().
    Buy(btc).
    Quantity(decimal.NewFromFloat(0.1)).
    Reason("RSI oversold").
    Build()
```

Fluent API for building trade signals. Kronos validates and executes them.

## Using Indicators

All indicators follow the same pattern: pass the asset and parameters.

```go
// RSI - Relative Strength Index
rsi := s.k.Indicators.RSI(btc, 14)

// SMA - Simple Moving Average  
sma := s.k.Indicators.SMA(btc, 20)

// EMA - Exponential Moving Average
ema := s.k.Indicators.EMA(btc, 50)

// MACD - Moving Average Convergence Divergence
macd := s.k.Indicators.MACD(btc, 12, 26, 9)

// Bollinger Bands
bb := s.k.Indicators.BollingerBands(btc, 20, 2.0)

// Stochastic Oscillator
stoch := s.k.Indicators.Stochastic(btc, 14, 3)

// ATR - Average True Range
atr := s.k.Indicators.ATR(btc, 14)
```

### Specifying Exchange or Interval

By default, Kronos uses your configured exchange and interval. Override when needed:

```go
// Use Binance specifically
rsi := s.k.Indicators.RSI(btc, 14, indicators.IndicatorOptions{
    Exchange: connector.Binance,
})

// Use 4-hour interval
sma := s.k.Indicators.SMA(btc, 200, indicators.IndicatorOptions{
    Interval: "4h",
})

// Both
ema := s.k.Indicators.EMA(btc, 50, indicators.IndicatorOptions{
    Exchange: connector.Bybit,
    Interval: "1h",
})
```

## Market Data

Access market data the same way:

```go
// Current price
price := s.k.Market.Price(btc)

// Price from specific exchange
price := s.k.Market.Price(btc, market.MarketOptions{
    Exchange: connector.Binance,
})

// Prices from all exchanges
prices := s.k.Market.Prices(btc)
for exchange, price := range prices {
    fmt.Printf("%s: %s\n", exchange, price)
}

// Order book
book := s.k.Market.OrderBook(btc)
topBid := book.Bids[0]  // Best bid
topAsk := book.Asks[0]  // Best ask

// Funding rate (perpetuals)
funding := s.k.Market.FundingRate(btc)
fmt.Printf("Rate: %s, Next: %s\n", 
    funding.CurrentRate, funding.NextFundingTime)

// Historical klines
klines := s.k.Market.Klines(btc, "1h", 100)  // Last 100 1h candles
```

## Working with Decimals

Kronos uses `decimal.Decimal` for all financial calculations:

```go
import "github.com/shopspring/decimal"

// Create decimals
price := decimal.NewFromFloat(50000.50)
quantity := decimal.NewFromInt(2)
pct := decimal.NewFromString("0.025")  // 2.5%

// Math operations
total := price.Mul(quantity)
fee := total.Mul(pct)
net := total.Sub(fee)

// Comparisons
if price.GreaterThan(decimal.NewFromInt(50000)) {
    // Price above 50k
}

if rsi.LessThan(decimal.NewFromInt(30)) {
    // RSI oversold
}

// String conversion
fmt.Println(price.String())           // "50000.5"
fmt.Println(price.StringFixed(2))     // "50000.50"
```

**Never use `float64` for money!** Always use `decimal.Decimal`.

## Multiple Assets

Trade multiple assets in one strategy:

```go
func (s *Strategy) GetSignals() ([]*strategy.Signal, error) {
    btc := s.k.Asset("BTC")
    eth := s.k.Asset("ETH")
    sol := s.k.Asset("SOL")
    
    var signals []*strategy.Signal
    
    // Check BTC
    btcRSI := s.k.Indicators.RSI(btc, 14)
    if btcRSI.LessThan(decimal.NewFromInt(30)) {
        signals = append(signals, 
            s.Signal().Buy(btc).Quantity(decimal.NewFromFloat(0.1)).Build())
    }
    
    // Check ETH
    ethRSI := s.k.Indicators.RSI(eth, 14)
    if ethRSI.GreaterThan(decimal.NewFromInt(70)) {
        signals = append(signals, 
            s.Signal().Sell(eth).Quantity(decimal.NewFromFloat(1.0)).Build())
    }
    
    // Check SOL
    solRSI := s.k.Indicators.RSI(sol, 14)
    if solRSI.LessThan(decimal.NewFromInt(30)) {
        signals = append(signals, 
            s.Signal().Buy(sol).Quantity(decimal.NewFromFloat(5.0)).Build())
    }
    
    return signals, nil
}
```

Kronos automatically manages data for all assets in parallel.

## Multiple Timeframes

Use different timeframes for different purposes:

```go
func (s *Strategy) GetSignals() ([]*strategy.Signal, error) {
    btc := s.k.Asset("BTC")
    
    // Long-term trend (4-hour)
    sma200 := s.k.Indicators.SMA(btc, 200, indicators.IndicatorOptions{
        Interval: "4h",
    })
    
    // Short-term signal (1-hour)
    rsi := s.k.Indicators.RSI(btc, 14, indicators.IndicatorOptions{
        Interval: "1h",
    })
    
    price := s.k.Market.Price(btc)
    
    // Only buy if:
    // - Price is above 4h SMA200 (uptrend)
    // - 1h RSI is oversold
    if price.GreaterThan(sma200) && rsi.LessThan(decimal.NewFromInt(30)) {
        return []*strategy.Signal{
            s.Signal().
                Buy(btc).
                Quantity(decimal.NewFromFloat(0.1)).
                Reason("Oversold in uptrend").
                Build(),
        }, nil
    }
    
    return nil, nil
}
```

## Combining Indicators

Use multiple indicators for confirmation:

```go
func (s *Strategy) GetSignals() ([]*strategy.Signal, error) {
    btc := s.k.Asset("BTC")
    
    // Get multiple indicators
    rsi := s.k.Indicators.RSI(btc, 14)
    stoch := s.k.Indicators.Stochastic(btc, 14, 3)
    bb := s.k.Indicators.BollingerBands(btc, 20, 2.0)
    price := s.k.Market.Price(btc)
    
    // Strong buy signal: all three confirm oversold
    if rsi.LessThan(decimal.NewFromInt(30)) &&         // RSI oversold
       stoch.K.LessThan(decimal.NewFromInt(20)) &&     // Stochastic oversold
       price.LessThan(bb.Lower) {                      // Price below lower band
        
        return []*strategy.Signal{
            s.Signal().
                Buy(btc).
                Quantity(decimal.NewFromFloat(0.15)).  // Larger position
                Reason("Triple confirmation: RSI, Stoch, BB").
                Build(),
        }, nil
    }
    
    return nil, nil
}
```

## Cross-Exchange Strategies

Find arbitrage opportunities across exchanges:

```go
func (s *Strategy) GetSignals() ([]*strategy.Signal, error) {
    btc := s.k.Asset("BTC")
    
    // Get prices from all exchanges
    prices := s.k.Market.Prices(btc)
    
    binancePrice := prices[connector.Binance]
    bybitPrice := prices[connector.Bybit]
    
    // Calculate spread
    spread := binancePrice.Sub(bybitPrice).Div(bybitPrice).Mul(decimal.NewFromInt(100))
    
    // If spread > 0.5%, arbitrage opportunity
    if spread.GreaterThan(decimal.NewFromFloat(0.5)) {
        s.k.Log().Opportunity("Arbitrage", btc.Symbol(),
            "Buy Bybit @ %s, Sell Binance @ %s, Spread: %s%%",
            bybitPrice, binancePrice, spread)
        
        return []*strategy.Signal{
            s.Signal().Buy(btc).Exchange(connector.Bybit).Build(),
            s.Signal().Sell(btc).Exchange(connector.Binance).Build(),
        }, nil
    }
    
    // Or use the built-in helper
    arbOpps := s.k.Market.FindArbitrage(btc, decimal.NewFromInt(50)) // 50 bps min
    for _, opp := range arbOpps {
        s.k.Log().Opportunity("Arbitrage", btc.Symbol(),
            "%s → %s: %s bps", opp.BuyExchange, opp.SellExchange, opp.SpreadBps)
    }
    
    return nil, nil
}
```

## Analytics

Beyond indicators, use market analytics:

```go
// Volatility
vol := s.k.Analytics.Volatility(btc, 24)
if vol.GreaterThan(decimal.NewFromInt(50)) {
    s.k.Log().MarketCondition("High volatility: %s%%", vol)
}

// Trend analysis
trend := s.k.Analytics.Trend(btc, 50)
if trend.Direction == "up" && trend.Strength.GreaterThan(decimal.NewFromInt(70)) {
    s.k.Log().MarketCondition("Strong uptrend")
}

// Volume analysis
volAnalysis := s.k.Analytics.VolumeAnalysis(btc, 24)
if volAnalysis.IsVolumeSpike {
    s.k.Log().MarketCondition("Volume spike detected")
}

// Price change
change := s.k.Analytics.GetPriceChange(btc, 24)
s.k.Log().MarketCondition("24h change: %s%%", change.ChangePercent)
```

## Logging

Kronos provides structured logging:

```go
// Info level
s.k.Log().Info("Strategy initialized")

// Debug level
s.k.Log().Debug("Strategy", btc.Symbol(), "RSI: %s", rsi)

// Market conditions
s.k.Log().MarketCondition("Price: %s, RSI: %s", price, rsi)

// Opportunities
s.k.Log().Opportunity("Strategy", btc.Symbol(), 
    "Buy signal - RSI oversold at %s", rsi)

// Failures
s.k.Log().Failed("Strategy", btc.Symbol(), 
    "Failed to calculate RSI: %v", err)
```

## Error Handling

Handle errors properly:

```go
func (s *Strategy) GetSignals() ([]*strategy.Signal, error) {
    btc := s.k.Asset("BTC")
    
    // Indicators can fail if insufficient data
    rsi, err := s.k.Indicators.RSI(btc, 14)
    if err != nil {
        return nil, fmt.Errorf("RSI calculation failed: %w", err)
    }
    
    // Market data can fail on API errors
    price, err := s.k.Market.Price(btc)
    if err != nil {
        return nil, fmt.Errorf("price fetch failed: %w", err)
    }
    
    // Your logic here...
    
    return signals, nil
}
```

## Complete Example

Here's a full strategy using multiple features:

```go
package main

import (
    sdk "github.com/backtesting-org/kronos-sdk/pkg/kronos"
    "github.com/backtesting-org/kronos-sdk/pkg/types/strategy"
    "github.com/shopspring/decimal"
)

type AdvancedStrategy struct {
    k *sdk.Kronos
}

func NewAdvanced(k *sdk.Kronos) *AdvancedStrategy {
    return &AdvancedStrategy{k: k}
}

func (s *AdvancedStrategy) GetSignals() ([]*strategy.Signal, error) {
    btc := s.k.Asset("BTC")
    
    // Multi-timeframe analysis
    sma200_4h := s.k.Indicators.SMA(btc, 200, indicators.IndicatorOptions{Interval: "4h"})
    rsi_1h := s.k.Indicators.RSI(btc, 14, indicators.IndicatorOptions{Interval: "1h"})
    stoch_1h := s.k.Indicators.Stochastic(btc, 14, 3, indicators.IndicatorOptions{Interval: "1h"})
    
    price := s.k.Market.Price(btc)
    vol := s.k.Analytics.Volatility(btc, 24)
    
    // Check trend
    inUptrend := price.GreaterThan(sma200_4h)
    
    // Check oversold
    rsiOversold := rsi_1h.LessThan(decimal.NewFromInt(30))
    stochOversold := stoch_1h.K.LessThan(decimal.NewFromInt(20))
    
    // Check volatility isn't too high
    volOk := vol.LessThan(decimal.NewFromInt(60))
    
    // Buy if all conditions met
    if inUptrend && rsiOversold && stochOversold && volOk {
        s.k.Log().Opportunity("Advanced", btc.Symbol(),
            "Buy signal - Uptrend + Oversold + Normal vol")
        
        return []*strategy.Signal{
            s.Signal().
                Buy(btc).
                Quantity(decimal.NewFromFloat(0.1)).
                Reason("Multi-factor confirmation").
                Build(),
        }, nil
    }
    
    return nil, nil
}

// Interface methods...
func (s *AdvancedStrategy) GetName() strategy.StrategyName { return "Advanced" }
func (s *AdvancedStrategy) GetDescription() string { return "Multi-factor strategy" }
func (s *AdvancedStrategy) GetRiskLevel() strategy.RiskLevel { return strategy.RiskLevelMedium }
func (s *AdvancedStrategy) GetStrategyType() strategy.StrategyType { return strategy.StrategyTypeTechnical }
```

## Next Steps

- **[API Reference](api/indicators/stochastic)** - Explore all indicators and market data functions
- **Examples** - See more strategy patterns
- **Testing** - Learn how to test your strategies

## Key Takeaways

1. **Focus on logic** - Kronos handles data, exchanges, and calculations
2. **Just ask** - Call `s.k.Indicators.RSI(btc, 14)` and Kronos does the rest
3. **Exchange agnostic** - Same code works on any exchange
4. **Type safe** - Compiler catches errors before runtime
5. **Decimal precision** - No floating-point errors

You write strategy logic. Kronos handles everything else.
