---
sidebar_position: 4
---

# MACD (Moving Average Convergence Divergence)

## Usage

```go
// Basic usage (12, 26, 9 is standard)
macd := s.k.Indicators().MACD(btc, 12, 26, 9)

s.k.Log().Debug("MACD", btc.Symbol(), "MACD: %s, Signal: %s, Histogram: %s",
    macd.MACD, macd.Signal, macd.Histogram)

// With options
macd := s.k.Indicators().MACD(btc, 12, 26, 9, indicators.IndicatorOptions{
    Interval: "4h",
})
```

## In a Strategy

```go
func (s *Strategy) GetSignals() ([]*strategy.Signal, error) {
    btc := s.k.Asset("BTC")
    
    macd := s.k.Indicators().MACD(btc, 12, 26, 9)
    
    // Bullish crossover: MACD crosses above Signal
    if macd.MACD.GreaterThan(macd.Signal) && macd.Histogram.GreaterThan(decimal.Zero) {
        return []*strategy.Signal{
            s.Signal().
                Buy(btc).
                Quantity(decimal.NewFromFloat(0.1)).
                Reason("MACD bullish crossover").
                Build(),
        }, nil
    }
    
    // Bearish crossover: MACD crosses below Signal
    if macd.MACD.LessThan(macd.Signal) && macd.Histogram.LessThan(decimal.Zero) {
        return []*strategy.Signal{
            s.Signal().
                Sell(btc).
                Quantity(decimal.NewFromFloat(0.1)).
                Reason("MACD bearish crossover").
                Build(),
        }, nil
    }
    
    return nil, nil
}
```

## Parameters

```go
MACD(asset, fastPeriod, slowPeriod, signalPeriod, ...options) *MACDResult
```

| Parameter | Type | Description |
|-----------|------|-------------|
| `asset` | `types.Asset` | Asset to calculate for |
| `fastPeriod` | `int` | Fast EMA period (typically 12) |
| `slowPeriod` | `int` | Slow EMA period (typically 26) |
| `signalPeriod` | `int` | Signal line period (typically 9) |
| `options` | `...IndicatorOptions` | Optional exchange/interval |

## Return Value

```go
type MACDResult struct {
    MACD      decimal.Decimal  // MACD line (fast - slow)
    Signal    decimal.Decimal  // Signal line (EMA of MACD)
    Histogram decimal.Decimal  // MACD - Signal
}
```

## Common Patterns

### Crossover Signals

```go
macd := s.k.Indicators().MACD(btc, 12, 26, 9)

// Bullish: MACD crosses above Signal
if macd.MACD.GreaterThan(macd.Signal) {
    return s.Signal().Buy(btc).Build()
}

// Bearish: MACD crosses below Signal
if macd.MACD.LessThan(macd.Signal) {
    return s.Signal().Sell(btc).Build()
}
```

### Histogram Momentum

```go
macd := s.k.Indicators().MACD(btc, 12, 26, 9)

// Histogram growing = momentum strengthening
if macd.Histogram.GreaterThan(decimal.Zero) {
    // Bullish momentum
}

// Histogram shrinking = momentum weakening
if macd.Histogram.LessThan(decimal.Zero) {
    // Bearish momentum
}
```

### Zero Line Cross

```go
macd := s.k.Indicators().MACD(btc, 12, 26, 9)

// MACD above zero = uptrend
if macd.MACD.GreaterThan(decimal.Zero) {
    // Bullish trend
}

// MACD below zero = downtrend
if macd.MACD.LessThan(decimal.Zero) {
    // Bearish trend
}
```

## What It Measures

MACD shows the relationship between two moving averages:

### Formulas

```
MACD Line = EMA(12) - EMA(26)
Signal Line = EMA(9) of MACD Line
Histogram = MACD Line - Signal Line
```

### Interpretation

- **MACD > Signal**: Bullish momentum
- **MACD < Signal**: Bearish momentum
- **Histogram growing**: Momentum strengthening
- **Histogram shrinking**: Momentum weakening

## See Also

- [RSI](rsi) - Momentum oscillator
- [Moving Averages](moving-averages) - MACD uses EMAs
- [Stochastic](stochastic) - Another momentum indicator

## References

- **Go Package**: [pkg.go.dev](https://pkg.go.dev/github.com/wisp-trading/sdk/pkg/analytics/indicators#MACD)
- **Source Code**: [macd.go](https://github.com/wisp-trading/sdk/blob/main/pkg/analytics/indicators/macd.go)
- **Theory**: [Investopedia - MACD](https://www.investopedia.com/terms/m/macd.asp)
