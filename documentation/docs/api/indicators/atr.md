---
sidebar_position: 6
---

# ATR (Average True Range)

## Usage

```go
// Basic usage
atr. _ := s.k.Indicators().ATR(btc, 14)  // 14-period ATR

// With options
atr, _ := s.k.Indicators().ATR(btc, 14, indicators.IndicatorOptions{
    Interval: "4h",
})
```

## In a Strategy

```go
func (s *Strategy) GetSignals() ([]*strategy.Signal, error) {
    btc := s.k.Asset("BTC")
    
    price, _ := s.k.Market().Price(btc)
    atr, _ := s.k.Indicators().ATR(btc, 14)
    
    // Set stop loss at 2× ATR below entry
    stopLoss := price.Sub(atr.Mul(decimal.NewFromInt(2)))
    
    // Buy signal with dynamic stop loss
    rsi := s.k.Indicators().RSI(btc, 14)
    if rsi.LessThan(decimal.NewFromInt(30)) {
        return []*strategy.Signal{
            s.Signal().
                Buy(btc).
                Quantity(decimal.NewFromFloat(0.1)).
                StopLoss(stopLoss).
                Reason("RSI oversold with ATR stop").
                Build(),
        }, nil
    }
    
    return nil, nil
}
```

## Parameters

```go
ATR(asset, period, ...options) decimal.Decimal
```

| Parameter | Type | Description |
|-----------|------|-------------|
| `asset` | `types.Asset` | Asset to calculate for |
| `period` | `int` | Number of periods (typically 14) |
| `options` | `...IndicatorOptions` | Optional exchange/interval |

## Common Patterns

### Dynamic Stop Loss

```go
price := s.k.Market().Price(btc)
atr := s.k.Indicators().ATR(btc, 14)

// Stop loss at 2× ATR
stopLoss := price.Sub(atr.Mul(decimal.NewFromInt(2)))

// Take profit at 3× ATR
takeProfit := price.Add(atr.Mul(decimal.NewFromInt(3)))
```

### Position Sizing

```go
atr := s.k.Indicators().ATR(btc, 14)
accountBalance := decimal.NewFromInt(10000)
riskPerTrade := accountBalance.Mul(decimal.NewFromFloat(0.02))  // 2% risk

// Position size based on ATR
positionSize := riskPerTrade.Div(atr.Mul(decimal.NewFromInt(2)))
```

### Volatility Filter

```go
atr := s.k.Indicators().ATR(btc, 14)
price := s.k.Market().Price(btc)

// ATR as percentage of price
atrPercent := atr.Div(price).Mul(decimal.NewFromInt(100))

// Only trade if volatility is reasonable
if atrPercent.LessThan(decimal.NewFromInt(5)) {
    // Low to medium volatility - trade
} else {
    // High volatility - skip
}
```

## What It Measures

ATR measures market volatility by calculating the average range between high and low prices:

### Formula

```
True Range = max(
    High - Low,
    |High - Previous Close|,
    |Low - Previous Close|
)

ATR = Average of True Range over period
```

### Interpretation

- **High ATR**: High volatility, larger price swings
- **Low ATR**: Low volatility, price consolidation
- **Rising ATR**: Volatility increasing
- **Falling ATR**: Volatility decreasing

## Best Practices

### ✅ Do

```go
// Use ATR for stop losses
atr := s.k.Indicators().ATR(btc, 14)
stopDistance := atr.Mul(decimal.NewFromInt(2))  // 2× ATR stop

// Adjust position size based on volatility
positionSize := baseSize.Div(atr)  // Smaller positions in high volatility
```

### ❌ Don't

```go
// Don't use fixed stop losses
// ❌ Same stop for all market conditions
stopLoss := price.Sub(decimal.NewFromInt(1000))

// ✅ Dynamic stop based on ATR
atr := s.k.Indicators().ATR(btc, 14)
stopLoss := price.Sub(atr.Mul(decimal.NewFromInt(2)))
```

## Complete Example

```go
func (s *Strategy) GetSignals() ([]*strategy.Signal, error) {
    btc := s.k.Asset("BTC")
    
    price := s.k.Market().Price(btc)
    atr := s.k.Indicators().ATR(btc, 14)
    rsi := s.k.Indicators().RSI(btc, 14)
    
    // Only trade in reasonable volatility
    atrPercent := atr.Div(price).Mul(decimal.NewFromInt(100))
    if atrPercent.GreaterThan(decimal.NewFromInt(6)) {
        return nil, nil  // Too volatile
    }
    
    // Buy on RSI oversold
    if rsi.LessThan(decimal.NewFromInt(30)) {
        stopLoss := price.Sub(atr.Mul(decimal.NewFromInt(2)))
        takeProfit := price.Add(atr.Mul(decimal.NewFromInt(3)))
        
        return []*strategy.Signal{
            s.Signal().
                Buy(btc).
                Quantity(decimal.NewFromFloat(0.1)).
                StopLoss(stopLoss).
                TakeProfit(takeProfit).
                Reason("RSI oversold with ATR risk management").
                Build(),
        }, nil
    }
    
    return nil, nil
}
```

## See Also

- [Bollinger Bands](bollinger-bands) - Another volatility indicator
- [RSI](rsi) - Combine with ATR for entries
- [Stochastic](stochastic) - Momentum with ATR stops

## References

- **Go Package**: [pkg.go.dev](https://pkg.go.dev/github.com/wisp-trading/sdk/pkg/analytics/indicators#ATR)
- **Source Code**: [atr.go](https://github.com/wisp-trading/sdk/blob/main/pkg/analytics/indicators/atr.go)
- **Theory**: [Investopedia - ATR](https://www.investopedia.com/terms/a/atr.asp)
