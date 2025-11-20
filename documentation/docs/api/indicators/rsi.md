---
sidebar_position: 2
---

# RSI (Relative Strength Index)

## Usage

```go
// Basic usage
rsi := s.k.Indicators().RSI(btc, 14)  // 14-period RSI

// With options
rsi := s.k.Indicators().RSI(btc, 14, indicators.IndicatorOptions{
    Exchange: connector.Binance,
    Interval: "4h",
})
```

## In a Strategy

```go
func (s *Strategy) GetSignals() ([]*strategy.Signal, error) {
    btc := s.k.Asset("BTC")
    
    // Get 14-period RSI (standard)
    rsi := s.k.Indicators().RSI(btc, 14)
    
    // Oversold: RSI < 30
    if rsi.LessThan(decimal.NewFromInt(30)) {
        return []*strategy.Signal{
            s.Signal().
                Buy(btc).
                Quantity(decimal.NewFromFloat(0.1)).
                Reason(fmt.Sprintf("RSI oversold at %s", rsi.StringFixed(2))).
                Build(),
        }, nil
    }
    
    // Overbought: RSI > 70
    if rsi.GreaterThan(decimal.NewFromInt(70)) {
        return []*strategy.Signal{
            s.Signal().
                Sell(btc).
                Quantity(decimal.NewFromFloat(0.1)).
                Reason(fmt.Sprintf("RSI overbought at %s", rsi.StringFixed(2))).
                Build(),
        }, nil
    }
    
    return nil, nil
}
```

## Parameters

```go
RSI(asset, period, ...options) decimal.Decimal
```

| Parameter | Type | Description |
|-----------|------|-------------|
| `asset` | `types.Asset` | Asset to calculate for |
| `period` | `int` | Number of periods (typically 14) |
| `options` | `...IndicatorOptions` | Optional exchange/interval |

## Return Value

Returns a `decimal.Decimal` value between 0 and 100.

## Common Patterns

### Basic Oversold/Overbought

```go
rsi := s.k.Indicators().RSI(btc, 14)

// Buy when oversold
if rsi.LessThan(decimal.NewFromInt(30)) {
    return s.Signal().Buy(btc).Build()
}

// Sell when overbought
if rsi.GreaterThan(decimal.NewFromInt(70)) {
    return s.Signal().Sell(btc).Build()
}
```

### With Trend Filter

```go
price := s.k.Market().Price(btc)
sma200 := s.k.Indicators().SMA(btc, 200)
rsi := s.k.Indicators().RSI(btc, 14)

// Only buy oversold signals in uptrend
if price.GreaterThan(sma200) && rsi.LessThan(decimal.NewFromInt(30)) {
    return s.Signal().Buy(btc).Reason("Oversold in uptrend").Build()
}
```

### Multi-Timeframe

```go
// Higher timeframe for trend
rsi4h := s.k.Indicators().RSI(btc, 14, indicators.IndicatorOptions{
    Interval: "4h",
})

// Lower timeframe for entry
rsi1h := s.k.Indicators().RSI(btc, 14, indicators.IndicatorOptions{
    Interval: "1h",
})

// 4h not overbought, 1h oversold
if rsi4h.LessThan(decimal.NewFromInt(70)) && 
   rsi1h.LessThan(decimal.NewFromInt(30)) {
    return s.Signal().Buy(btc).Build()
}
```

### Divergence Detection

```go
// Bullish divergence: price makes lower lows, RSI makes higher lows
prevPrice := s.k.Market().PricePrevious(btc, 1)
prevRSI := s.k.Indicators().RSIPrevious(btc, 14, 1)

price := s.k.Market().Price(btc)
rsi := s.k.Indicators().RSI(btc, 14)

if price.LessThan(prevPrice) && rsi.GreaterThan(prevRSI) {
    // Bullish divergence detected
    return s.Signal().Buy(btc).Reason("RSI bullish divergence").Build()
}
```

## Interpretation Guide

### Levels

| RSI Value | Condition | Action |
|-----------|-----------|--------|
| **> 70** | Overbought | Consider selling |
| **50-70** | Bullish | Hold/accumulate |
| **30-50** | Bearish | Cautious |
| **< 30** | Oversold | Consider buying |

### Advanced Levels

Some traders use tighter levels:

```go
// Extreme oversold/overbought
if rsi.LessThan(decimal.NewFromInt(20)) {
    // Extremely oversold - strong buy
}

if rsi.GreaterThan(decimal.NewFromInt(80)) {
    // Extremely overbought - strong sell
}
```

## What It Measures

RSI measures the magnitude of recent price changes to evaluate overbought or oversold conditions.

### Formula

```
RS = Average Gain / Average Loss
RSI = 100 - (100 / (1 + RS))
```

Where:
- **Average Gain** = Average of gains over the period
- **Average Loss** = Average of losses over the period

### The Theory

- **RSI > 70**: Momentum is strong to the upside, potentially overbought
- **RSI < 30**: Momentum is strong to the downside, potentially oversold
- **RSI = 50**: Neutral, no clear momentum

## Best Practices

### ✅ Do

```go
// Use with trend confirmation
price := s.k.Market().Price(btc)
ema200 := s.k.Indicators().EMA(btc, 200)
rsi := s.k.Indicators().RSI(btc, 14)

// Only buy oversold in uptrend
if price.GreaterThan(ema200) && rsi.LessThan(decimal.NewFromInt(30)) {
    return s.Signal().Buy(btc).Build()
}
```

```go
// Combine with other indicators
rsi := s.k.Indicators().RSI(btc, 14)
macd := s.k.Indicators().MACD(btc, 12, 26, 9)

// Both confirm bullish
if rsi.GreaterThan(decimal.NewFromInt(50)) && 
   macd.MACD.GreaterThan(macd.Signal) {
    return s.Signal().Buy(btc).Build()
}
```

### ❌ Don't

```go
// Don't trade solely on RSI
if rsi.LessThan(decimal.NewFromInt(30)) {
    // ❌ Can stay oversold in downtrend
    return s.Signal().Buy(btc).Build()
}
```

```go
// Don't use same levels for all assets
// ❌ Volatile assets may need different thresholds
rsi := s.k.Indicators().RSI(volatileAsset, 14)
if rsi.LessThan(decimal.NewFromInt(30)) {
    return s.Signal().Buy(volatileAsset).Build()
}

// ✅ Adjust for asset characteristics
if rsi.LessThan(decimal.NewFromInt(20)) {  // Tighter threshold
    return s.Signal().Buy(volatileAsset).Build()
}
```

## Common Periods

| Period | Use Case |
|--------|----------|
| **14** | Standard, balanced |
| **9** | More sensitive, faster signals |
| **21** | Less sensitive, smoother |
| **7** | Day trading, very responsive |

```go
// Standard
rsi14 := s.k.Indicators().RSI(btc, 14)

// More sensitive
rsi9 := s.k.Indicators().RSI(btc, 9)

// Smoother
rsi21 := s.k.Indicators().RSI(btc, 21)
```

## Common Pitfalls

### 1. False Signals in Trends

**Problem:** RSI can remain overbought/oversold during strong trends.

**Solution:**
```go
// Add trend filter
ema := s.k.Indicators().EMA(btc, 200)
price := s.k.Market().Price(btc)
rsi := s.k.Indicators().RSI(btc, 14)

inUptrend := price.GreaterThan(ema)

if inUptrend {
    // Only take buy signals in uptrend
    if rsi.LessThan(decimal.NewFromInt(30)) {
        return s.Signal().Buy(btc).Build()
    }
}
```

### 2. Ignoring Market Context

**Problem:** Using same levels in all market conditions.

**Solution:**
```go
// Adjust levels based on volatility
vol := s.k.Analytics().Volatility(btc, 24)
rsi := s.k.Indicators().RSI(btc, 14)

oversoldLevel := decimal.NewFromInt(30)
if vol.GreaterThan(decimal.NewFromInt(50)) {
    // High volatility - use tighter levels
    oversoldLevel = decimal.NewFromInt(20)
}

if rsi.LessThan(oversoldLevel) {
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

type RSIStrategy struct {
    k *sdk.Kronos
}

func NewRSI(k *sdk.Kronos) *RSIStrategy {
    return &RSIStrategy{k: k}
}

func (s *RSIStrategy) GetSignals() ([]*strategy.Signal, error) {
    btc := s.k.Asset("BTC")
    
    // Multi-timeframe RSI
    rsi4h := s.k.Indicators().RSI(btc, 14, indicators.IndicatorOptions{
        Interval: "4h",
    })
    rsi1h := s.k.Indicators().RSI(btc, 14, indicators.IndicatorOptions{
        Interval: "1h",
    })
    
    // Trend filter
    price := s.k.Market().Price(btc)
    ema200 := s.k.Indicators().EMA(btc, 200, indicators.IndicatorOptions{
        Interval: "4h",
    })
    
    inUptrend := price.GreaterThan(ema200)
    
    // Buy: uptrend + 4h not overbought + 1h oversold
    if inUptrend &&
       rsi4h.LessThan(decimal.NewFromInt(70)) &&
       rsi1h.LessThan(decimal.NewFromInt(30)) {
        
        return []*strategy.Signal{
            s.Signal().
                Buy(btc).
                Quantity(decimal.NewFromFloat(0.1)).
                Reason("Uptrend with 1h RSI oversold").
                Build(),
        }, nil
    }
    
    // Sell: 1h RSI overbought
    if rsi1h.GreaterThan(decimal.NewFromInt(70)) {
        return []*strategy.Signal{
            s.Signal().
                Sell(btc).
                Quantity(decimal.NewFromFloat(0.1)).
                Reason("1h RSI overbought").
                Build(),
        }, nil
    }
    
    return nil, nil
}

// Interface methods
func (s *RSIStrategy) GetName() strategy.StrategyName { return "RSI" }
func (s *RSIStrategy) GetDescription() string { return "Multi-timeframe RSI strategy" }
func (s *RSIStrategy) GetRiskLevel() strategy.RiskLevel { return strategy.RiskLevelMedium }
func (s *RSIStrategy) GetStrategyType() strategy.StrategyType { return strategy.StrategyTypeTechnical }
```

## See Also

- [Stochastic](stochastic) - Another momentum oscillator
- [MACD](macd) - Trend-following momentum
- [Bollinger Bands](bollinger-bands) - Volatility-based signals

## References

- **Go Package**: [pkg.go.dev](https://pkg.go.dev/github.com/backtesting-org/kronos-sdk/pkg/analytics/indicators#RSI)
- **Source Code**: [rsi.go](https://github.com/backtesting-org/kronos-sdk/blob/main/pkg/analytics/indicators/rsi.go)
- **Theory**: [Investopedia - RSI](https://www.investopedia.com/terms/r/rsi.asp)
