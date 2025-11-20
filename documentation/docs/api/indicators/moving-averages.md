---
sidebar_position: 3
---

# Moving Averages (SMA & EMA)

## SMA (Simple Moving Average)

### Usage

```go
// Basic usage
sma20 := s.k.Indicators().SMA(btc, 20)    // 20-period SMA
sma50 := s.k.Indicators().SMA(btc, 50)    // 50-period SMA
sma200 := s.k.Indicators().SMA(btc, 200)  // 200-period SMA

// With options
sma := s.k.Indicators().SMA(btc, 50, indicators.IndicatorOptions{
    Exchange: connector.Binance,
    Interval: "4h",
})
```

### In a Strategy

```go
func (s *Strategy) GetSignals() ([]*strategy.Signal, error) {
    btc := s.k.Asset("BTC")
    
    price := s.k.Market().Price(btc)
    sma50 := s.k.Indicators().SMA(btc, 50)
    sma200 := s.k.Indicators().SMA(btc, 200)
    
    // Golden cross: SMA50 > SMA200 and price > SMA50
    if sma50.GreaterThan(sma200) && price.GreaterThan(sma50) {
        return []*strategy.Signal{
            s.Signal().
                Buy(btc).
                Quantity(decimal.NewFromFloat(0.1)).
                Reason("Price above SMA50 in golden cross").
                Build(),
        }, nil
    }
    
    return nil, nil
}
```

## EMA (Exponential Moving Average)

### Usage

```go
// Basic usage
ema12 := s.k.Indicators().EMA(btc, 12)  // Fast EMA
ema26 := s.k.Indicators().EMA(btc, 26)  // Slow EMA
ema200 := s.k.Indicators().EMA(btc, 200) // Long-term trend

// With options
ema := s.k.Indicators().EMA(eth, 50, indicators.IndicatorOptions{
    Interval: "1h",
})
```

### In a Strategy

```go
func (s *Strategy) GetSignals() ([]*strategy.Signal, error) {
    eth := s.k.Asset("ETH")
    
    price := s.k.Market().Price(eth)
    ema20 := s.k.Indicators().EMA(eth, 20)
    ema50 := s.k.Indicators().EMA(eth, 50)
    
    // EMA crossover
    if ema20.GreaterThan(ema50) && price.GreaterThan(ema20) {
        return []*strategy.Signal{
            s.Signal().
                Buy(eth).
                Quantity(decimal.NewFromFloat(1.0)).
                Reason("Bullish EMA crossover").
                Build(),
        }, nil
    }
    
    return nil, nil
}
```

## Parameters

### SMA
```go
SMA(asset, period, ...options) decimal.Decimal
```

### EMA
```go
EMA(asset, period, ...options) decimal.Decimal
```

| Parameter | Type | Description |
|-----------|------|-------------|
| `asset` | `types.Asset` | Asset to calculate for |
| `period` | `int` | Number of periods |
| `options` | `...IndicatorOptions` | Optional exchange/interval |

## Common Patterns

### Trend Identification

```go
price := s.k.Market().Price(btc)
sma200 := s.k.Indicators().SMA(btc, 200)

// Uptrend: price above 200 SMA
if price.GreaterThan(sma200) {
    // Bullish trend
}

// Downtrend: price below 200 SMA
if price.LessThan(sma200) {
    // Bearish trend
}
```

### Support/Resistance

```go
price := s.k.Market().Price(btc)
ema50 := s.k.Indicators().EMA(btc, 50)

// EMA acting as support in uptrend
if price.GreaterThan(ema50) {
    distance := price.Sub(ema50)
    
    // Price pulled back to EMA - potential entry
    if distance.LessThan(ema50.Mul(decimal.NewFromFloat(0.02))) {  // Within 2%
        return s.Signal().Buy(btc).Reason("Price at EMA50 support").Build()
    }
}
```

### Golden Cross / Death Cross

```go
sma50 := s.k.Indicators().SMA(btc, 50)
sma200 := s.k.Indicators().SMA(btc, 200)

// Golden cross: SMA50 crosses above SMA200 (bullish)
if sma50.GreaterThan(sma200) {
    return s.Signal().Buy(btc).Reason("Golden cross").Build()
}

// Death cross: SMA50 crosses below SMA200 (bearish)
if sma50.LessThan(sma200) {
    return s.Signal().Sell(btc).Reason("Death cross").Build()
}
```

### Multi-MA Trend System

```go
ema20 := s.k.Indicators().EMA(btc, 20)
ema50 := s.k.Indicators().EMA(btc, 50)
ema200 := s.k.Indicators().EMA(btc, 200)

// Strong uptrend: all EMAs aligned
if ema20.GreaterThan(ema50) && ema50.GreaterThan(ema200) {
    price := s.k.Market().Price(btc)
    
    // Buy pullbacks to fast EMA
    if price.LessThan(ema20.Mul(decimal.NewFromFloat(1.02))) {  // Within 2%
        return s.Signal().Buy(btc).Reason("Pullback in strong uptrend").Build()
    }
}
```

## SMA vs EMA

### Key Differences

| Feature | SMA | EMA |
|---------|-----|-----|
| **Weighting** | Equal weight to all periods | More weight to recent prices |
| **Responsiveness** | Slower to react | Faster to react |
| **Smoothness** | Smoother | More sensitive |
| **Best for** | Long-term trends | Short-term signals |

### When to Use SMA

```go
// Long-term trend identification
sma200 := s.k.Indicators().SMA(btc, 200, indicators.IndicatorOptions{
    Interval: "1d",  // Daily
})

// SMA is better for:
// - Long-term trend filters
// - Major support/resistance levels
// - Reducing noise in volatile markets
```

### When to Use EMA

```go
// Responsive trend following
ema20 := s.k.Indicators().EMA(btc, 20)

// EMA is better for:
// - Short to medium-term trading
// - Faster entry/exit signals
// - Dynamic support/resistance
// - MACD calculations (uses EMA internally)
```

## Common Periods

### SMA Periods

| Period | Use Case |
|--------|----------|
| **20** | Short-term trend |
| **50** | Medium-term trend |
| **100** | Intermediate trend |
| **200** | Long-term trend (most popular) |

### EMA Periods

| Period | Use Case |
|--------|----------|
| **12/26** | MACD components |
| **20** | Short-term trend |
| **50** | Medium-term trend |
| **200** | Long-term trend |

```go
// Common SMA periods
sma20 := s.k.Indicators().SMA(btc, 20)    // Short-term
sma50 := s.k.Indicators().SMA(btc, 50)    // Medium-term
sma200 := s.k.Indicators().SMA(btc, 200)  // Long-term

// Common EMA periods
ema12 := s.k.Indicators().EMA(btc, 12)    // Fast
ema26 := s.k.Indicators().EMA(btc, 26)    // Slow
ema50 := s.k.Indicators().EMA(btc, 50)    // Medium
ema200 := s.k.Indicators().EMA(btc, 200)  // Long-term
```

## What They Measure

### SMA (Simple Moving Average)

The arithmetic mean of prices over a period:

```
SMA = (P1 + P2 + ... + Pn) / n
```

All prices have equal weight.

### EMA (Exponential Moving Average)

Weighted average that gives more importance to recent prices:

```
EMA = (Price × Multiplier) + (Previous EMA × (1 - Multiplier))
Multiplier = 2 / (Period + 1)
```

Recent prices have exponentially more weight.

## Best Practices

### ✅ Do

```go
// Use multiple timeframes
sma200_4h := s.k.Indicators().SMA(btc, 200, indicators.IndicatorOptions{
    Interval: "4h",
})
price := s.k.Market().Price(btc)

// Only trade with the trend
if price.GreaterThan(sma200_4h) {
    // Look for buy opportunities only
}
```

```go
// Combine SMA and EMA
sma200 := s.k.Indicators().SMA(btc, 200)  // Trend filter
ema20 := s.k.Indicators().EMA(btc, 20)    // Entry signal

price := s.k.Market().Price(btc)

// Buy: above SMA200 + pullback to EMA20
if price.GreaterThan(sma200) && price.LessThan(ema20.Mul(decimal.NewFromFloat(1.01))) {
    return s.Signal().Buy(btc).Build()
}
```

### ❌ Don't

```go
// Don't use MA alone without confirmation
sma50 := s.k.Indicators().SMA(btc, 50)
price := s.k.Market().Price(btc)

if price.GreaterThan(sma50) {
    // ❌ Too simplistic, needs confirmation
    return s.Signal().Buy(btc).Build()
}
```

## Complete Example

```go
package main

import (
    sdk "github.com/backtesting-org/kronos-sdk/pkg/kronos"
    "github.com/backtesting-org/kronos-sdk/pkg/types/strategy"
    "github.com/shopspring/decimal"
)

type MAStrategy struct {
    k *sdk.Kronos
}

func NewMA(k *sdk.Kronos) *MAStrategy {
    return &MAStrategy{k: k}
}

func (s *MAStrategy) GetSignals() ([]*strategy.Signal, error) {
    btc := s.k.Asset("BTC")
    
    // Get moving averages
    price := s.k.Market().Price(btc)
    ema20 := s.k.Indicators().EMA(btc, 20)
    ema50 := s.k.Indicators().EMA(btc, 50)
    sma200 := s.k.Indicators().SMA(btc, 200)
    
    // Trend filter: above SMA200
    inUptrend := price.GreaterThan(sma200)
    
    // EMA alignment: EMA20 > EMA50
    emasAligned := ema20.GreaterThan(ema50)
    
    // Buy: uptrend + EMA crossover + price pullback
    if inUptrend && emasAligned {
        distanceToEMA20 := price.Sub(ema20).Div(ema20).Mul(decimal.NewFromInt(100))
        
        // Price within 2% of EMA20 (pullback)
        if distanceToEMA20.LessThan(decimal.NewFromInt(2)) {
            return []*strategy.Signal{
                s.Signal().
                    Buy(btc).
                    Quantity(decimal.NewFromFloat(0.1)).
                    Reason("EMA pullback in uptrend").
                    Build(),
            }, nil
        }
    }
    
    // Sell: death cross forming
    if ema20.LessThan(ema50) {
        return []*strategy.Signal{
            s.Signal().
                Sell(btc).
                Quantity(decimal.NewFromFloat(0.1)).
                Reason("Bearish EMA crossover").
                Build(),
        }, nil
    }
    
    return nil, nil
}

// Interface methods
func (s *MAStrategy) GetName() strategy.StrategyName { return "MA" }
func (s *MAStrategy) GetDescription() string { return "Moving average crossover strategy" }
func (s *MAStrategy) GetRiskLevel() strategy.RiskLevel { return strategy.RiskLevelLow }
func (s *MAStrategy) GetStrategyType() strategy.StrategyType { return strategy.StrategyTypeTechnical }
```

## See Also

- [MACD](macd) - Uses EMAs internally
- [Bollinger Bands](bollinger-bands) - Uses SMA as middle band
- [RSI](rsi) - Momentum indicator

## References

- **Go Package**: [pkg.go.dev - SMA](https://pkg.go.dev/github.com/backtesting-org/kronos-sdk/pkg/analytics/indicators#SMA) | [EMA](https://pkg.go.dev/github.com/backtesting-org/kronos-sdk/pkg/analytics/indicators#EMA)
- **Source Code**: [sma.go](https://github.com/backtesting-org/kronos-sdk/blob/main/pkg/analytics/indicators/sma.go) | [ema.go](https://github.com/backtesting-org/kronos-sdk/blob/main/pkg/analytics/indicators/ema.go)
- **Theory**: [Investopedia - SMA](https://www.investopedia.com/terms/s/sma.asp) | [EMA](https://www.investopedia.com/terms/e/ema.asp)
