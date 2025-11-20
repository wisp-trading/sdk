---
sidebar_position: 2
---

# Quick Reference

Essential concepts for writing Kronos strategies.

## The Kronos Context

Every strategy has access to the Kronos SDK through `k`:

```go
type MyStrategy struct {
    k *sdk.Kronos  // Your gateway to everything
}
```

The `k` instance provides:
- **`k.Asset(symbol)`** - Get asset references
- **`k.Indicators`** - Technical indicators
- **`k.Market`** - Market data (prices, order books, funding)
- **`k.Analytics`** - Market analytics
- **`k.Log()`** - Structured logging

## Assets

Get references to assets you want to trade:

```go
btc := s.k.Asset("BTC")
eth := s.k.Asset("ETH")
sol := s.k.Asset("SOL")
```

Assets are just references. Kronos knows which exchange to use based on your config.

## Indicators

All indicators follow the same pattern: pass the asset and parameters.

```go
// RSI - Relative Strength Index
rsi := s.k.Indicators().RSI(btc, 14)

// SMA - Simple Moving Average  
sma := s.k.Indicators().SMA(btc, 20)

// EMA - Exponential Moving Average
ema := s.k.Indicators().EMA(btc, 50)

// MACD - Moving Average Convergence Divergence
macd := s.k.Indicators().MACD(btc, 12, 26, 9)
s.k.Log().Debug("MACD", btc.Symbol(), "MACD: %s, Signal: %s, Histogram: %s", 
    macd.MACD, macd.Signal, macd.Histogram)

// Bollinger Bands
bb := s.k.Indicators().BollingerBands(btc, 20, 2.0)
s.k.Log().Debug("BB", btc.Symbol(), "Upper: %s, Middle: %s, Lower: %s",
    bb.Upper, bb.Middle, bb.Lower)

// Stochastic Oscillator
stoch := s.k.Indicators().Stochastic(btc, 14, 3)
s.k.Log().Debug("Stochastic", btc.Symbol(), "K: %s, D: %s", stoch.K, stoch.D)

// ATR - Average True Range
atr := s.k.Indicators().ATR(btc, 14)
```

Kronos automatically:
1. Fetches the required price data
2. Calculates the indicator
3. Returns the latest value

### Indicator Options

Specify exchange or interval when needed:

```go
import "github.com/backtesting-org/kronos-sdk/pkg/kronos/indicators"
import "github.com/backtesting-org/kronos-sdk/pkg/types/connector"

// Use specific exchange
rsi := s.k.Indicators().RSI(btc, 14, indicators.IndicatorOptions{
    Exchange: connector.Binance,
})

// Use different timeframe
sma := s.k.Indicators().SMA(btc, 200, indicators.IndicatorOptions{
    Interval: "4h",
})

// Both
ema := s.k.Indicators().EMA(btc, 50, indicators.IndicatorOptions{
    Exchange: connector.Bybit,
    Interval: "1h",
})
```

## Market Data

Access real-time market data:

```go
// Current price
price := s.k.Market().Price(btc)

// Price from specific exchange
price := s.k.Market().Price(btc, market.MarketOptions{
    Exchange: connector.Binance,
})

// Prices from all exchanges
prices := s.k.Market().Prices(btc)
for exchange, price := range prices {
    s.k.Log().Info("%s: %s", exchange, price)
}

// Order book
book := s.k.Market().OrderBook(btc)
topBid := book.Bids[0]  // Best bid
topAsk := book.Asks[0]  // Best ask

// Funding rate (perpetuals)
funding := s.k.Market().FundingRate(btc)

// Historical klines
klines := s.k.Market().Klines(btc, "1h", 100)  // Last 100 1h candles
```

## Signals

Create trading signals with the fluent API:

```go
import "github.com/backtesting-org/kronos-sdk/pkg/types/connector"

// Buy signal - specify asset, exchange, and quantity
signal := s.k.Signal(s.GetName()).
    Buy(btc, connector.Binance, decimal.NewFromFloat(0.1)).
    Build()

// Sell signal
signal := s.k.Signal(s.GetName()).
    Sell(btc, connector.Binance, decimal.NewFromFloat(0.1)).
    Build()

// Short sell signal
signal := s.k.Signal(s.GetName()).
    SellShort(btc, connector.Binance, decimal.NewFromFloat(0.1)).
    Build()

// Limit orders with specific price
signal := s.k.Signal(s.GetName()).
    BuyLimit(btc, connector.Binance, decimal.NewFromFloat(0.1), decimal.NewFromInt(45000)).
    Build()

// Multiple trades in one signal (e.g., cash and carry)
signal := s.k.Signal(s.GetName()).
    Buy(btc, connector.Binance, decimal.NewFromFloat(0.1)).
    SellShort(btc, connector.Bybit, decimal.NewFromFloat(0.1)).
    Build()
```

Return signals from `GetSignals()`:

```go
func (s *Strategy) GetSignals() ([]*strategy.Signal, error) {
    btc := s.k.Asset("BTC")
    rsi := s.k.Indicators().RSI(btc, 14)
    
    if rsi.LessThan(decimal.NewFromInt(30)) {
        signal := s.k.Signal(s.GetName()).
            Buy(btc, connector.Binance, decimal.NewFromFloat(0.1)).
            Build()
        return []*strategy.Signal{signal}, nil
    }
    
    return nil, nil
}
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
s.k.Log().Info("Price: %s", price.String())              // "50000.5"
s.k.Log().Info("Price fixed: %s", price.StringFixed(2)) // "50000.50"
```

:::danger Never use float64 for money
Always use `decimal.Decimal` to avoid floating-point precision errors.
:::

## Logging

Use structured logging to track your strategy:

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

Always handle errors properly:

```go
func (s *Strategy) GetSignals() ([]*strategy.Signal, error) {
    btc := s.k.Asset("BTC")
    
    // Check if indicators have valid data
    rsi := s.k.Indicators().RSI(btc, 14)
    if rsi.IsZero() {
        // No data yet, skip this cycle
        return nil, nil
    }
    
    // Your logic here...
    
    return signals, nil
}
```

## Complete Example

A simple RSI strategy:

```go
package main

import (
    sdk "github.com/backtesting-org/kronos-sdk/pkg/kronos"
    "github.com/backtesting-org/kronos-sdk/pkg/types/strategy"
    "github.com/shopspring/decimal"
)

type RSIStrategy struct {
    k *sdk.Kronos
}

func NewRSI(k *sdk.Kronos) *RSIStrategy {
    return &RSIStrategy{k: k}
}

func (s *RSIStrategy) GetSignals() ([]*strategy.Signal, error) {
    btc := s.k.Asset("BTC")
    
    // Get RSI - Kronos handles everything
    rsi := s.k.Indicators().RSI(btc, 14)
    
    // Buy when oversold
    if rsi.LessThan(decimal.NewFromInt(30)) {
        return []*strategy.Signal{
            s.Signal().
                Buy(btc).
                Quantity(decimal.NewFromFloat(0.1)).
                Reason("RSI oversold at " + rsi.String()).
                Build(),
        }, nil
    }
    
    // Sell when overbought
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

// Required interface methods
func (s *RSIStrategy) GetName() strategy.StrategyName { 
    return "RSI" 
}

func (s *RSIStrategy) GetDescription() string { 
    return "Simple RSI momentum strategy" 
}

func (s *RSIStrategy) GetRiskLevel() strategy.RiskLevel { 
    return strategy.RiskLevelMedium 
}

func (s *RSIStrategy) GetStrategyType() strategy.StrategyType { 
    return strategy.StrategyTypeTechnical 
}
```

## Next Steps

- **[Writing Strategies](writing-strategies)** - Deep dive into strategy patterns
- **[Examples](examples)** - Common strategy implementations
- **[Indicators Reference](../api/indicators/rsi)** - Full indicator documentation
