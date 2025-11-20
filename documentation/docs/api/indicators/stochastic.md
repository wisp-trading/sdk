---
sidebar_position: 1
---

# Stochastic Oscillator

## Usage

```go
// Basic usage - uses configured exchange and interval
stoch := s.k.Indicators().Stochastic(btc, 14, 3)

k := stoch.K  // Fast stochastic (0-100)
d := stoch.D  // Slow stochastic (0-100)
```

### With Options

```go
// Specify exchange
stoch := s.k.Indicators().Stochastic(btc, 14, 3, indicators.IndicatorOptions{
    Exchange: connector.Binance,
})

// Specify interval
stoch := s.k.Indicators().Stochastic(btc, 14, 3, indicators.IndicatorOptions{
    Interval: "4h",
})

// Both
stoch := s.k.Indicators().Stochastic(btc, 14, 3, indicators.IndicatorOptions{
    Exchange: connector.Bybit,
    Interval: "1h",
})
```

## In a Strategy

```go
func (s *Strategy) GetSignals() ([]*strategy.Signal, error) {
    btc := s.k.Asset("BTC")
    
    // Get stochastic (14, 3 is standard)
    stoch := s.k.Indicators().Stochastic(btc, 14, 3)
    
    // Oversold: both lines below 20
    if stoch.K.LessThan(decimal.NewFromInt(20)) && 
       stoch.D.LessThan(decimal.NewFromInt(20)) {
        return []*strategy.Signal{
            s.Signal().
                Buy(btc).
                Quantity(decimal.NewFromFloat(0.1)).
                Reason(fmt.Sprintf("Stochastic oversold - K: %s, D: %s", 
                    stoch.K.StringFixed(2), stoch.D.StringFixed(2))).
                Build(),
        }, nil
    }
    
    // Overbought: both lines above 80
    if stoch.K.GreaterThan(decimal.NewFromInt(80)) && 
       stoch.D.GreaterThan(decimal.NewFromInt(80)) {
        return []*strategy.Signal{
            s.Signal().
                Sell(btc).
                Quantity(decimal.NewFromFloat(0.1)).
                Reason(fmt.Sprintf("Stochastic overbought - K: %s, D: %s", 
                    stoch.K.StringFixed(2), stoch.D.StringFixed(2))).
                Build(),
        }, nil
    }
    
    return nil, nil
}
```

## Parameters

```go
Stochastic(asset, kPeriod, dPeriod, ...options)
```

| Parameter | Type | Description |
|-----------|------|-------------|
| `asset` | `portfolio.Asset` | Asset to calculate for (e.g., `btc`) |
| `kPeriod` | `int` | Lookback period for %K (typically 14) |
| `dPeriod` | `int` | Smoothing period for %D (typically 3) |
| `options` | `...IndicatorOptions` | Optional exchange/interval |

## Return Value

```go
type StochasticResult struct {
    K decimal.Decimal  // %K line (fast), values 0-100
    D decimal.Decimal  // %D line (slow), values 0-100
}
```

## Common Patterns

### Basic Oversold/Overbought

```go
stoch := s.k.Indicators().Stochastic(btc, 14, 3)

// Buy when oversold
if stoch.K.LessThan(decimal.NewFromInt(20)) {
    return s.Signal().Buy(btc).Build()
}

// Sell when overbought
if stoch.K.GreaterThan(decimal.NewFromInt(80)) {
    return s.Signal().Sell(btc).Build()
}
```

### Multi-Timeframe

```go
// Higher timeframe for trend
stoch4h := s.k.Indicators().Stochastic(btc, 14, 3, indicators.IndicatorOptions{
    Interval: "4h",
})

// Lower timeframe for entry
stoch1h := s.k.Indicators().Stochastic(btc, 14, 3, indicators.IndicatorOptions{
    Interval: "1h",
})

// Only buy if 4h shows uptrend and 1h is oversold
if stoch4h.K.GreaterThan(decimal.NewFromInt(50)) &&
   stoch1h.K.LessThan(decimal.NewFromInt(20)) {
    return s.Signal().Buy(btc).Reason("4h uptrend + 1h oversold").Build()
}
```

### Combined with Other Indicators

```go
stoch := s.k.Indicators().Stochastic(btc, 14, 3)
rsi := s.k.Indicators().RSI(btc, 14)
price := s.k.Market().Price(btc)
sma200 := s.k.Indicators().SMA(btc, 200)

// Strong buy: uptrend + both indicators oversold
if price.GreaterThan(sma200) &&
   stoch.K.LessThan(decimal.NewFromInt(20)) &&
   rsi.LessThan(decimal.NewFromInt(30)) {
    return s.Signal().
        Buy(btc).
        Quantity(decimal.NewFromFloat(0.15)).  // Larger position
        Reason("Uptrend with double oversold confirmation").
        Build()
}
```

## Common Parameter Sets

| Use Case | kPeriod | dPeriod | Description |
|----------|---------|---------|-------------|
| **Standard** | 14 | 3 | Default, balanced sensitivity |
| **Fast** | 5 | 3 | More signals, more noise |
| **Slow** | 14 | 5 | Fewer signals, smoother |
| **Very Slow** | 21 | 5 | Long-term, least noise |

```go
// Standard
stoch := s.k.Indicators().Stochastic(btc, 14, 3)

// Fast (more responsive)
stoch := s.k.Indicators().Stochastic(btc, 5, 3)

// Slow (smoother)
stoch := s.k.Indicators().Stochastic(btc, 14, 5)
```

## Interpretation Guide

### Levels

| Condition | Value | Action |
|-----------|-------|--------|
| **Oversold** | < 20 | Consider buying |
| **Neutral** | 20-80 | Normal range |
| **Overbought** | > 80 | Consider selling |

### Crossovers

**Bullish Signal:**
- %K crosses above %D (especially in oversold zone < 30)

**Bearish Signal:**
- %K crosses below %D (especially in overbought zone > 70)

### Divergence

**Bullish Divergence** (reversal signal):
- Price makes lower lows
- Stochastic makes higher lows

**Bearish Divergence** (reversal signal):
- Price makes higher highs
- Stochastic makes lower highs

## What It Measures

The Stochastic Oscillator compares the current closing price to the price range over a period of time.

### Formula

**%K (Fast Stochastic):**
```
%K = (Current Close - Lowest Low) / (Highest High - Lowest Low) × 100
```

**%D (Slow Stochastic):**
```
%D = SMA(%K, dPeriod)
```

### The Theory

The indicator operates on the principle that:
- In an **uptrend**, prices tend to close near their high
- In a **downtrend**, prices tend to close near their low

When %K and %D are:
- **Below 20** - Price is closing near the period's lows (oversold)
- **Above 80** - Price is closing near the period's highs (overbought)

## Best Practices

### ✅ Do

```go
// Use with trend filter
price := s.k.Market().Price(btc)
sma := s.k.Indicators().SMA(btc, 200)
stoch := s.k.Indicators().Stochastic(btc, 14, 3)

// Only buy oversold signals in uptrend
if price.GreaterThan(sma) && stoch.K.LessThan(decimal.NewFromInt(20)) {
    return s.Signal().Buy(btc).Build()
}
```

```go
// Wait for confirmation from both lines
stoch := s.k.Indicators().Stochastic(btc, 14, 3)

if stoch.K.LessThan(decimal.NewFromInt(20)) && 
   stoch.D.LessThan(decimal.NewFromInt(20)) {
    // Both confirm oversold
    return s.Signal().Buy(btc).Build()
}
```

```go
// Combine with other indicators
stoch := s.k.Indicators().Stochastic(btc, 14, 3)
rsi := s.k.Indicators().RSI(btc, 14)

// Strong signal when both agree
if stoch.K.LessThan(decimal.NewFromInt(20)) && 
   rsi.LessThan(decimal.NewFromInt(30)) {
    return s.Signal().Buy(btc).Build()
}
```

### ❌ Don't

```go
// Don't trade solely on overbought/oversold
if stoch.K.GreaterThan(decimal.NewFromInt(80)) {
    // ❌ Could stay overbought in strong uptrend
    return s.Signal().Sell(btc).Build()
}
```

```go
// Don't ignore the trend
// ❌ Buying oversold in downtrend often fails
stoch := s.k.Indicators().Stochastic(btc, 14, 3)
if stoch.K.LessThan(decimal.NewFromInt(20)) {
    return s.Signal().Buy(btc).Build()  // ❌
}

// ✅ Check trend first
price := s.k.Market().Price(btc)
sma := s.k.Indicators().SMA(btc, 200)
if price.GreaterThan(sma) && stoch.K.LessThan(decimal.NewFromInt(20)) {
    return s.Signal().Buy(btc).Build()  // ✅
}
```

## Common Pitfalls

### 1. False Signals in Trends

**Problem:** In strong trends, stochastic can remain in extreme zones for extended periods.

**Solution:**
```go
// Add trend filter
price := s.k.Market().Price(btc)
ema200 := s.k.Indicators().EMA(btc, 200)
stoch := s.k.Indicators().Stochastic(btc, 14, 3)

// Only take signals aligned with trend
if price.GreaterThan(ema200) {
    // Uptrend - only take buy signals
    if stoch.K.LessThan(decimal.NewFromInt(20)) {
        return s.Signal().Buy(btc).Build()
    }
}
```

### 2. Whipsaws in Choppy Markets

**Problem:** Rapid crossovers generate many false signals.

**Solution:**
```go
// Require both lines in extreme zone
stoch := s.k.Indicators().Stochastic(btc, 14, 3)

if stoch.K.LessThan(decimal.NewFromInt(20)) && 
   stoch.D.LessThan(decimal.NewFromInt(20)) {
    // Both confirm - stronger signal
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

type StochasticStrategy struct {
    k *sdk.Kronos
}

func NewStochastic(k *sdk.Kronos) *StochasticStrategy {
    return &StochasticStrategy{k: k}
}

func (s *StochasticStrategy) GetSignals() ([]*strategy.Signal, error) {
    btc := s.k.Asset("BTC")
    
    // Multi-timeframe stochastic
    stoch4h := s.k.Indicators().Stochastic(btc, 14, 3, indicators.IndicatorOptions{
        Interval: "4h",
    })
    stoch1h := s.k.Indicators().Stochastic(btc, 14, 3, indicators.IndicatorOptions{
        Interval: "1h",
    })
    
    // Trend filter
    price := s.k.Market().Price(btc)
    ema200 := s.k.Indicators().EMA(btc, 200, indicators.IndicatorOptions{
        Interval: "4h",
    })
    
    // Check trend
    inUptrend := price.GreaterThan(ema200)
    
    // Buy conditions:
    // 1. Uptrend (4h)
    // 2. 4h stochastic not overbought
    // 3. 1h stochastic oversold
    if inUptrend &&
       stoch4h.K.LessThan(decimal.NewFromInt(70)) &&
       stoch1h.K.LessThan(decimal.NewFromInt(20)) &&
       stoch1h.D.LessThan(decimal.NewFromInt(20)) {
        
        return []*strategy.Signal{
            s.Signal().
                Buy(btc).
                Quantity(decimal.NewFromFloat(0.1)).
                Reason("4h uptrend, 1h stochastic oversold").
                Build(),
        }, nil
    }
    
    // Sell conditions:
    // 1. 1h stochastic overbought
    // 2. Both K and D confirm
    if stoch1h.K.GreaterThan(decimal.NewFromInt(80)) &&
       stoch1h.D.GreaterThan(decimal.NewFromInt(80)) {
        
        return []*strategy.Signal{
            s.Signal().
                Sell(btc).
                Quantity(decimal.NewFromFloat(0.1)).
                Reason("1h stochastic overbought").
                Build(),
        }, nil
    }
    
    return nil, nil
}

// Interface methods
func (s *StochasticStrategy) GetName() strategy.StrategyName { 
    return "Stochastic" 
}

func (s *StochasticStrategy) GetDescription() string { 
    return "Multi-timeframe stochastic strategy" 
}

func (s *StochasticStrategy) GetRiskLevel() strategy.RiskLevel { 
    return strategy.RiskLevelMedium 
}

func (s *StochasticStrategy) GetStrategyType() strategy.StrategyType { 
    return strategy.StrategyTypeTechnical 
}
```

## See Also

- [RSI](rsi) - Another momentum oscillator
- [MACD](macd) - Trend-following indicator
- [Bollinger Bands](bollinger) - Volatility-based indicator

## References

- **Go Package**: [pkg.go.dev](https://pkg.go.dev/github.com/backtesting-org/kronos-sdk/pkg/analytics/indicators#Stochastic)
- **Source Code**: [stochastic.go](https://github.com/backtesting-org/kronos-sdk/blob/main/pkg/analytics/indicators/stochastic.go)
- **Theory**: [Investopedia - Stochastic Oscillator](https://www.investopedia.com/terms/s/stochasticoscillator.asp)
