---
sidebar_position: 5
---

# Bollinger Bands

## Usage

```go
// Basic usage (20, 2.0 is standard)
bb := s.k.Indicators().BollingerBands(btc, 20, 2.0)

s.k.Log().Debug("Bollinger Bands", btc.Symbol(), "Upper: %s, Middle: %s, Lower: %s",
    bb.Upper, bb.Middle, bb.Lower)

// With options
bb := s.k.Indicators().BollingerBands(btc, 20, 2.0, indicators.IndicatorOptions{
    Interval: "1h",
})
```

## In a Strategy

```go
func (s *Strategy) GetSignals() ([]*strategy.Signal, error) {
    btc := s.k.Asset("BTC")
    
    price := s.k.Market().Price(btc)
    bb := s.k.Indicators().BollingerBands(btc, 20, 2.0)
    
    // Buy when price touches lower band
    if price.LessThan(bb.Lower) {
        return []*strategy.Signal{
            s.Signal().
                Buy(btc).
                Quantity(decimal.NewFromFloat(0.1)).
                Reason("Price below lower Bollinger Band").
                Build(),
        }, nil
    }
    
    // Sell when price touches upper band
    if price.GreaterThan(bb.Upper) {
        return []*strategy.Signal{
            s.Signal().
                Sell(btc).
                Quantity(decimal.NewFromFloat(0.1)).
                Reason("Price above upper Bollinger Band").
                Build(),
        }, nil
    }
    
    return nil, nil
}
```

## Parameters

```go
BollingerBands(asset, period, stdDev, ...options) *BollingerBandsResult
```

| Parameter | Type | Description |
|-----------|------|-------------|
| `asset` | `types.Asset` | Asset to calculate for |
| `period` | `int` | SMA period (typically 20) |
| `stdDev` | `float64` | Std dev multiplier (typically 2.0) |
| `options` | `...IndicatorOptions` | Optional exchange/interval |

## Return Value

```go
type BollingerBandsResult struct {
    Upper  decimal.Decimal  // Upper band (SMA + stdDev)
    Middle decimal.Decimal  // Middle band (SMA)
    Lower  decimal.Decimal  // Lower band (SMA - stdDev)
}
```

## Common Patterns

### Mean Reversion

```go
price := s.k.Market().Price(btc)
bb := s.k.Indicators().BollingerBands(btc, 20, 2.0)

// Buy at lower band
if price.LessThan(bb.Lower) {
    return s.Signal().Buy(btc).Build()
}

// Sell at upper band
if price.GreaterThan(bb.Upper) {
    return s.Signal().Sell(btc).Build()
}
```

### Breakout Detection

```go
bb := s.k.Indicators().BollingerBands(btc, 20, 2.0)

// Band width (volatility measure)
bandWidth := bb.Upper.Sub(bb.Lower).Div(bb.Middle).Mul(decimal.NewFromInt(100))

// Squeeze: bands narrowing (low volatility)
if bandWidth.LessThan(decimal.NewFromInt(10)) {
    // Potential breakout coming
}
```

## What It Measures

Bollinger Bands measure volatility and identify overbought/oversold conditions:

### Formulas

```
Middle Band = SMA(period)
Upper Band = Middle + (stdDev × Standard Deviation)
Lower Band = Middle - (stdDev × Standard Deviation)
```

### Interpretation

- **Price near upper band**: Overbought
- **Price near lower band**: Oversold
- **Bands narrow**: Low volatility, potential breakout
- **Bands wide**: High volatility

## See Also

- [ATR](atr) - Another volatility indicator
- [RSI](rsi) - Overbought/oversold momentum
- [Moving Averages](moving-averages) - BB uses SMA

## References

- **Go Package**: [pkg.go.dev](https://pkg.go.dev/github.com/wisp-trading/sdk/pkg/analytics/indicators#BollingerBands)
- **Source Code**: [bollinger.go](https://github.com/wisp-trading/sdk/blob/main/pkg/analytics/indicators/bollinger.go)
- **Theory**: [Investopedia - Bollinger Bands](https://www.investopedia.com/terms/b/bollingerbands.asp)
