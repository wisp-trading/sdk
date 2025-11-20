---
sidebar_position: 3
---

# Writing Strategies

Learn how to build sophisticated trading strategies with Kronos.

## Strategy Structure

Every strategy implements the `strategy.Strategy` interface:

```go
type Strategy interface {
    GetSignals() ([]*Signal, error)
    GetName() StrategyName
    GetDescription() string
    GetRiskLevel() RiskLevel
    GetStrategyType() StrategyType
}
```

The most important method is `GetSignals()` - this is where your trading logic lives.

## GetSignals()

Kronos calls `GetSignals()` on each interval (configured in `config.yaml`). Your job is to:
1. Analyze market conditions
2. Return trade signals if conditions are met
3. Return `nil, nil` if no action needed

```go
func (s *Strategy) GetSignals() ([]*strategy.Signal, error) {
    // 1. Get your assets
    btc := s.k.Asset("BTC")
    
    // 2. Analyze market
    rsi := s.k.Indicators().RSI(btc, 14)
    price := s.k.Market().Price(btc)
    
    // 3. Decide
    if rsi.LessThan(decimal.NewFromInt(30)) {
        signal := s.k.Signal(s.GetName()).
            Buy(btc, connector.Binance, decimal.NewFromFloat(0.1)).
            Build()
        return []*strategy.Signal{signal}, nil
    }
    
    // 4. No action
    return nil, nil
}
```

## Multiple Assets

Trade multiple assets in one strategy:

```go
func (s *Strategy) GetSignals() ([]*strategy.Signal, error) {
    btc := s.k.Asset("BTC")
    eth := s.k.Asset("ETH")
    sol := s.k.Asset("SOL")
    
    var signals []*strategy.Signal
    
    // Check each asset
    for _, asset := range []types.Asset{btc, eth, sol} {
        rsi := s.k.Indicators().RSI(asset, 14)
        
        if rsi.LessThan(decimal.NewFromInt(30)) {
            signal := s.k.Signal(s.GetName()).
                Buy(asset, connector.Binance, decimal.NewFromFloat(0.1)).
                Build()
            signals = append(signals, signal)
        }
    }
    
    return signals, nil
}
```

Kronos automatically manages data for all assets in parallel.

## Multiple Timeframes

Use different timeframes for trend and signals:

```go
func (s *Strategy) GetSignals() ([]*strategy.Signal, error) {
    btc := s.k.Asset("BTC")
    
    // Long-term trend (4-hour)
    sma200 := s.k.Indicators().SMA(btc, 200, indicators.IndicatorOptions{
        Interval: "4h",
    })
    
    // Short-term signal (1-hour)
    rsi := s.k.Indicators().RSI(btc, 14, indicators.IndicatorOptions{
        Interval: "1h",
    })
    
    price := s.k.Market().Price(btc)
    
    // Only buy if in uptrend AND oversold
    if price.GreaterThan(sma200) && rsi.LessThan(decimal.NewFromInt(30)) {
        signal := s.k.Signal(s.GetName()).
            Buy(btc, connector.Binance, decimal.NewFromFloat(0.1)).
            Build()
        s.k.Log().Opportunity(string(s.GetName()), "BTC", "Oversold in uptrend")
        return []*strategy.Signal{signal}, nil
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
    rsi := s.k.Indicators().RSI(btc, 14)
    stoch := s.k.Indicators().Stochastic(btc, 14, 3)
    bb := s.k.Indicators().BollingerBands(btc, 20, 2.0)
    price := s.k.Market().Price(btc)
    
    // Require all three to confirm oversold
    oversold := rsi.LessThan(decimal.NewFromInt(30)) &&
                stoch.K.LessThan(decimal.NewFromInt(20)) &&
                price.LessThan(bb.Lower)
    
    if oversold {
        signal := s.k.Signal(s.GetName()).
            Buy(btc, connector.Binance, decimal.NewFromFloat(0.15)).  // Larger size with confirmation
            Build()
        s.k.Log().Opportunity(string(s.GetName()), "BTC", "Triple confirmation: RSI, Stoch, BB")
        return []*strategy.Signal{signal}, nil
    }
    
    return nil, nil
}
```

## Risk Management

### Stop Loss & Take Profit

Use ATR for dynamic stops:

```go
func (s *Strategy) GetSignals() ([]*strategy.Signal, error) {
    btc := s.k.Asset("BTC")
    
    rsi := s.k.Indicators().RSI(btc, 14)
    price := s.k.Market().Price(btc)
    atr := s.k.Indicators().ATR(btc, 14)
    
    if rsi.LessThan(decimal.NewFromInt(30)) {
        // Calculate stops based on ATR
        // Stop loss at 2× ATR below entry
        stopLoss := price.Sub(atr.Mul(decimal.NewFromInt(2)))
        
        // Take profit at 3× ATR above entry
        takeProfit := price.Add(atr.Mul(decimal.NewFromInt(3)))
        
        // Log the stop levels for reference
        s.k.Log().Info(string(s.GetName()), \"BTC\", 
            \"Entry: %s, Stop: %s, Target: %s (R:R 1:1.5)\", 
            price, stopLoss, takeProfit)
        
        signal := s.k.Signal(s.GetName()).
            Buy(btc, connector.Binance, decimal.NewFromFloat(0.1)).
            Build()
        return []*strategy.Signal{signal}, nil
    }
    
    return nil, nil
}
```

### Position Sizing

Size positions based on volatility:

```go
func (s *Strategy) GetSignals() ([]*strategy.Signal, error) {
    btc := s.k.Asset("BTC")
    
    atr := s.k.Indicators().ATR(btc, 14)
    price := s.k.Market().Price(btc)
    
    // Risk 2% of account on this trade
    accountBalance := decimal.NewFromInt(10000)
    riskAmount := accountBalance.Mul(decimal.NewFromFloat(0.02))  // $200 risk
    
    // Stop loss at 2× ATR
    stopDistance := atr.Mul(decimal.NewFromInt(2))
    
    // Position size = risk / stop distance
    quantity := riskAmount.Div(stopDistance)
    
    // Rest of logic...
}
```

### Volatility Filter

Skip trading in extreme volatility:

```go
func (s *Strategy) GetSignals() ([]*strategy.Signal, error) {
    btc := s.k.Asset("BTC")
    
    atr := s.k.Indicators().ATR(btc, 14)
    price := s.k.Market().Price(btc)
    
    // ATR as percentage of price
    atrPercent := atr.Div(price).Mul(decimal.NewFromInt(100))
    
    // Skip if volatility too high
    if atrPercent.GreaterThan(decimal.NewFromInt(5)) {
        s.k.Log().MarketCondition("Volatility too high: %s%%", atrPercent)
        return nil, nil
    }
    
    // Normal trading logic...
}
```

## Cross-Exchange Strategies

### Arbitrage

Find price differences across exchanges:

```go
func (s *Strategy) GetSignals() ([]*strategy.Signal, error) {
    btc := s.k.Asset("BTC")
    
    // Get prices from all exchanges
    prices := s.k.Market().Prices(btc)
    
    binancePrice := prices[connector.Binance]
    bybitPrice := prices[connector.Bybit]
    
    // Calculate spread
    spread := binancePrice.Sub(bybitPrice).Div(bybitPrice).Mul(decimal.NewFromInt(100))
    
    // If spread > 0.5%, arbitrage opportunity
    if spread.GreaterThan(decimal.NewFromFloat(0.5)) {
        s.k.Log().Opportunity(string(s.GetName()), \"BTC\", \"Spread: %.2f%%\", spread)
        
        // Note: You need to specify quantity for each trade
        qty := decimal.NewFromFloat(0.1)
        return []*strategy.Signal{
            s.k.Signal(s.GetName()).
                Buy(btc, connector.Bybit, qty).
                Sell(btc, connector.Binance, qty).
                Build(),
        }, nil
    }
    
    return nil, nil
}
```

### Best Price Execution

Always trade on the exchange with best price:

```go
func (s *Strategy) GetSignals() ([]*strategy.Signal, error) {
    btc := s.k.Asset("BTC")
    rsi := s.k.Indicators().RSI(btc, 14)
    
    if rsi.LessThan(decimal.NewFromInt(30)) {
        // Find exchange with lowest price
        prices := s.k.Market().Prices(btc)
        
        var bestExchange connector.ExchangeName
        var bestPrice decimal.Decimal
        
        for exchange, price := range prices {
            if bestPrice.IsZero() || price.LessThan(bestPrice) {
                bestPrice = price
                bestExchange = exchange
            }
        }
        
        s.k.Log().Info(string(s.GetName()), \"BTC\", \"Best price on %s: %s\", bestExchange, bestPrice)
        
        signal := s.k.Signal(s.GetName()).
            Buy(btc, bestExchange, decimal.NewFromFloat(0.1)).
            Build()
        return []*strategy.Signal{signal}, nil
    }
    
    return nil, nil
}
```

## Advanced Patterns

### Trend Following

Follow the trend with moving average crossovers:

```go
func (s *Strategy) GetSignals() ([]*strategy.Signal, error) {
    btc := s.k.Asset("BTC")
    
    sma50 := s.k.Indicators().SMA(btc, 50)
    sma200 := s.k.Indicators().SMA(btc, 200)
    price := s.k.Market().Price(btc)
    
    // Golden cross: 50 SMA crosses above 200 SMA
    if sma50.GreaterThan(sma200) && price.GreaterThan(sma50) {
        signal := s.k.Signal(s.GetName()).
            Buy(btc, connector.Binance, decimal.NewFromFloat(0.2)).
            Build()
        s.k.Log().Opportunity(string(s.GetName()), "BTC", "Golden cross + price above 50 SMA")
        return []*strategy.Signal{signal}, nil
    }
    
    // Death cross: 50 SMA crosses below 200 SMA
    if sma50.LessThan(sma200) {
        signal := s.k.Signal(s.GetName()).
            Sell(btc, connector.Binance, decimal.NewFromFloat(0.2)).
            Build()
        s.k.Log().Opportunity(string(s.GetName()), "BTC", "Death cross")
        return []*strategy.Signal{signal}, nil
    }
    
    return nil, nil
}
```

### Mean Reversion

Trade bounces from Bollinger Bands:

```go
func (s *Strategy) GetSignals() ([]*strategy.Signal, error) {
    btc := s.k.Asset("BTC")
    
    bb := s.k.Indicators().BollingerBands(btc, 20, 2.0)
    price := s.k.Market().Price(btc)
    rsi := s.k.Indicators().RSI(btc, 14)
    
    // Buy when price touches lower band AND RSI confirms
    if price.LessThan(bb.Lower) && rsi.LessThan(decimal.NewFromInt(30)) {
        signal := s.k.Signal(s.GetName()).
            Buy(btc, connector.Binance, decimal.NewFromFloat(0.1)).
            Build()
        s.k.Log().Opportunity(string(s.GetName()), "BTC", 
            "Mean reversion from lower BB, target middle: %s", bb.Middle)
        return []*strategy.Signal{signal}, nil
    }
    
    return nil, nil
}
```

### Breakout Trading

Trade breakouts with volume confirmation:

```go
func (s *Strategy) GetSignals() ([]*strategy.Signal, error) {
    btc := s.k.Asset("BTC")
    
    bb := s.k.Indicators().BollingerBands(btc, 20, 2.0)
    price := s.k.Market().Price(btc)
    
    // Band width (volatility)
    bandWidth := bb.Upper.Sub(bb.Lower).Div(bb.Middle).Mul(decimal.NewFromInt(100))
    
    // Squeeze: bands narrow (breakout coming)
    if bandWidth.LessThan(decimal.NewFromInt(10)) {
        s.k.Log().MarketCondition("Bollinger squeeze detected")
        
        // Buy on upward breakout
        if price.GreaterThan(bb.Upper) {
            signal := s.k.Signal(s.GetName()).
                Buy(btc, connector.Binance, decimal.NewFromFloat(0.15)).
                Build()
            s.k.Log().Opportunity(string(s.GetName()), "BTC", "Breakout from squeeze")
            return []*strategy.Signal{signal}, nil
        }
    }
    
    return nil, nil
}
```

## State Management

Keep track of positions and state:

```go
type TrendStrategy struct {
    k            *sdk.Kronos
    inPosition   bool
    entryPrice   decimal.Decimal
    trailingStop decimal.Decimal
}

func (s *TrendStrategy) GetSignals() ([]*strategy.Signal, error) {
    btc := s.k.Asset("BTC")
    price := s.k.Market().Price(btc)
    atr := s.k.Indicators().ATR(btc, 14)
    
    // Entry logic
    if !s.inPosition {
        rsi := s.k.Indicators().RSI(btc, 14)
        if rsi.LessThan(decimal.NewFromInt(30)) {
            s.inPosition = true
            s.entryPrice = price
            s.trailingStop = price.Sub(atr.Mul(decimal.NewFromInt(2)))
            
            signal := s.k.Signal(s.GetName()).
                Buy(btc, connector.Binance, decimal.NewFromFloat(0.1)).
                Build()
            return []*strategy.Signal{signal}, nil
        }
    }
    
    // Exit logic (trailing stop)
    if s.inPosition {
        // Update trailing stop as price rises
        newStop := price.Sub(atr.Mul(decimal.NewFromInt(2)))
        if newStop.GreaterThan(s.trailingStop) {
            s.trailingStop = newStop
        }
        
        // Exit if stop hit
        if price.LessThan(s.trailingStop) {
            s.inPosition = false
            
            signal := s.k.Signal(s.GetName()).
                Sell(btc, connector.Binance, decimal.NewFromFloat(0.1)).
                Build()
            s.k.Log().Info(string(s.GetName()), "BTC", "Trailing stop hit")
            return []*strategy.Signal{signal}, nil
        }
    }
    
    return nil, nil
}
```

## Best Practices

### ✅ Do

- **Test thoroughly** - Backtest extensively before live trading
- **Use proper decimals** - Always use `decimal.Decimal` for money
- **Handle errors** - Check for data availability
- **Log decisions** - Use `s.k.Log()` to track behavior
- **Start small** - Begin with small position sizes
- **Use stop losses** - Protect against adverse moves
- **Combine indicators** - Use multiple confirmations

### ❌ Don't

- **Don't use floats** - Never use `float64` for financial calculations
- **Don't overtrade** - Avoid excessive signals
- **Don't overfit** - Don't optimize for past data
- **Don't ignore risk** - Always manage position sizes
- **Don't skip backtesting** - Always test before going live

## Next Steps

- **[Examples](examples)** - See complete strategy implementations
- **[Configuration](configuration)** - Configure exchanges and parameters
- **[Indicators Reference](../api/indicators/rsi)** - Full indicator documentation
