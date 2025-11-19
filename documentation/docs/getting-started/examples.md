---
sidebar_position: 4
---

# Examples

Complete strategy implementations for common trading patterns.

## Simple RSI Strategy

Classic momentum strategy using RSI:

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
    rsi := s.k.Indicators.RSI(btc, 14)
    
    // Buy oversold
    if rsi.LessThan(decimal.NewFromInt(30)) {
        return []*strategy.Signal{
            s.Signal().
                Buy(btc).
                Quantity(decimal.NewFromFloat(0.1)).
                Reason("RSI oversold: " + rsi.String()).
                Build(),
        }, nil
    }
    
    // Sell overbought
    if rsi.GreaterThan(decimal.NewFromInt(70)) {
        return []*strategy.Signal{
            s.Signal().
                Sell(btc).
                Quantity(decimal.NewFromFloat(0.1)).
                Reason("RSI overbought: " + rsi.String()).
                Build(),
        }, nil
    }
    
    return nil, nil
}

func (s *RSIStrategy) GetName() strategy.StrategyName { return "RSI" }
func (s *RSIStrategy) GetDescription() string { return "Simple RSI momentum" }
func (s *RSIStrategy) GetRiskLevel() strategy.RiskLevel { return strategy.RiskLevelMedium }
func (s *RSIStrategy) GetStrategyType() strategy.StrategyType { return strategy.StrategyTypeTechnical }
```

## Moving Average Crossover

Golden cross / death cross strategy:

```go
type MACrossover struct {
    k *sdk.Kronos
}

func NewMACrossover(k *sdk.Kronos) *MACrossover {
    return &MACrossover{k: k}
}

func (s *MACrossover) GetSignals() ([]*strategy.Signal, error) {
    btc := s.k.Asset("BTC")
    
    sma50 := s.k.Indicators.SMA(btc, 50)
    sma200 := s.k.Indicators.SMA(btc, 200)
    price := s.k.Market.Price(btc)
    
    // Golden cross: 50 crosses above 200
    if sma50.GreaterThan(sma200) && price.GreaterThan(sma50) {
        return []*strategy.Signal{
            s.Signal().
                Buy(btc).
                Quantity(decimal.NewFromFloat(0.2)).
                Reason("Golden cross").
                Build(),
        }, nil
    }
    
    // Death cross: 50 crosses below 200
    if sma50.LessThan(sma200) {
        return []*strategy.Signal{
            s.Signal().
                Sell(btc).
                Quantity(decimal.NewFromFloat(0.2)).
                Reason("Death cross").
                Build(),
        }, nil
    }
    
    return nil, nil
}

func (s *MACrossover) GetName() strategy.StrategyName { return "MA-Crossover" }
func (s *MACrossover) GetDescription() string { return "Golden/Death cross strategy" }
func (s *MACrossover) GetRiskLevel() strategy.RiskLevel { return strategy.RiskLevelLow }
func (s *MACrossover) GetStrategyType() strategy.StrategyType { return strategy.StrategyTypeTrend }
```

## Bollinger Bands Mean Reversion

Buy at lower band, sell at upper band:

```go
type BollingerMeanReversion struct {
    k *sdk.Kronos
}

func NewBollingerMR(k *sdk.Kronos) *BollingerMeanReversion {
    return &BollingerMeanReversion{k: k}
}

func (s *BollingerMeanReversion) GetSignals() ([]*strategy.Signal, error) {
    btc := s.k.Asset("BTC")
    
    bb := s.k.Indicators.BollingerBands(btc, 20, 2.0)
    price := s.k.Market.Price(btc)
    rsi := s.k.Indicators.RSI(btc, 14)
    
    // Buy at lower band with RSI confirmation
    if price.LessThan(bb.Lower) && rsi.LessThan(decimal.NewFromInt(35)) {
        return []*strategy.Signal{
            s.Signal().
                Buy(btc).
                Quantity(decimal.NewFromFloat(0.1)).
                TakeProfit(bb.Middle).
                Reason("Mean reversion from lower band").
                Build(),
        }, nil
    }
    
    // Sell at upper band
    if price.GreaterThan(bb.Upper) && rsi.GreaterThan(decimal.NewFromInt(65)) {
        return []*strategy.Signal{
            s.Signal().
                Sell(btc).
                Quantity(decimal.NewFromFloat(0.1)).
                TakeProfit(bb.Middle).
                Reason("Mean reversion from upper band").
                Build(),
        }, nil
    }
    
    return nil, nil
}

func (s *BollingerMeanReversion) GetName() strategy.StrategyName { return "BB-MeanReversion" }
func (s *BollingerMeanReversion) GetDescription() string { return "Bollinger Bands mean reversion" }
func (s *BollingerMeanReversion) GetRiskLevel() strategy.RiskLevel { return strategy.RiskLevelMedium }
func (s *BollingerMeanReversion) GetStrategyType() strategy.StrategyType { return strategy.StrategyTypeMeanReversion }
```

## MACD Momentum

Trade MACD crossovers with trend filter:

```go
type MACDMomentum struct {
    k *sdk.Kronos
}

func NewMACDMomentum(k *sdk.Kronos) *MACDMomentum {
    return &MACDMomentum{k: k}
}

func (s *MACDMomentum) GetSignals() ([]*strategy.Signal, error) {
    btc := s.k.Asset("BTC")
    
    macd := s.k.Indicators.MACD(btc, 12, 26, 9)
    sma200 := s.k.Indicators.SMA(btc, 200)
    price := s.k.Market.Price(btc)
    
    // Only trade with the trend
    inUptrend := price.GreaterThan(sma200)
    inDowntrend := price.LessThan(sma200)
    
    // Bullish crossover in uptrend
    if macd.MACD.GreaterThan(macd.Signal) && 
       macd.Histogram.GreaterThan(decimal.Zero) && 
       inUptrend {
        return []*strategy.Signal{
            s.Signal().
                Buy(btc).
                Quantity(decimal.NewFromFloat(0.15)).
                Reason("MACD bullish crossover in uptrend").
                Build(),
        }, nil
    }
    
    // Bearish crossover in downtrend
    if macd.MACD.LessThan(macd.Signal) && 
       macd.Histogram.LessThan(decimal.Zero) && 
       inDowntrend {
        return []*strategy.Signal{
            s.Signal().
                Sell(btc).
                Quantity(decimal.NewFromFloat(0.15)).
                Reason("MACD bearish crossover in downtrend").
                Build(),
        }, nil
    }
    
    return nil, nil
}

func (s *MACDMomentum) GetName() strategy.StrategyName { return "MACD-Momentum" }
func (s *MACDMomentum) GetDescription() string { return "MACD with trend filter" }
func (s *MACDMomentum) GetRiskLevel() strategy.RiskLevel { return strategy.RiskLevelMedium }
func (s *MACDMomentum) GetStrategyType() strategy.StrategyType { return strategy.StrategyTypeMomentum }
```

## Multi-Indicator Confirmation

Require multiple indicators to agree:

```go
type MultiConfirmation struct {
    k *sdk.Kronos
}

func NewMultiConfirmation(k *sdk.Kronos) *MultiConfirmation {
    return &MultiConfirmation{k: k}
}

func (s *MultiConfirmation) GetSignals() ([]*strategy.Signal, error) {
    btc := s.k.Asset("BTC")
    
    // Get all indicators
    rsi := s.k.Indicators.RSI(btc, 14)
    stoch := s.k.Indicators.Stochastic(btc, 14, 3)
    bb := s.k.Indicators.BollingerBands(btc, 20, 2.0)
    macd := s.k.Indicators.MACD(btc, 12, 26, 9)
    price := s.k.Market.Price(btc)
    
    // Count bullish signals
    bullishSignals := 0
    
    if rsi.LessThan(decimal.NewFromInt(30)) {
        bullishSignals++
    }
    if stoch.K.LessThan(decimal.NewFromInt(20)) {
        bullishSignals++
    }
    if price.LessThan(bb.Lower) {
        bullishSignals++
    }
    if macd.MACD.GreaterThan(macd.Signal) {
        bullishSignals++
    }
    
    // Require at least 3 of 4 indicators to agree
    if bullishSignals >= 3 {
        return []*strategy.Signal{
            s.Signal().
                Buy(btc).
                Quantity(decimal.NewFromFloat(0.2)).
                Reason(fmt.Sprintf("%d indicators confirm buy", bullishSignals)).
                Build(),
        }, nil
    }
    
    return nil, nil
}

func (s *MultiConfirmation) GetName() strategy.StrategyName { return "Multi-Confirmation" }
func (s *MultiConfirmation) GetDescription() string { return "Multi-indicator confirmation" }
func (s *MultiConfirmation) GetRiskLevel() strategy.RiskLevel { return strategy.RiskLevelLow }
func (s *MultiConfirmation) GetStrategyType() strategy.StrategyType { return strategy.StrategyTypeTechnical }
```

## ATR-Based Risk Management

Dynamic stops and position sizing:

```go
type ATRRiskManaged struct {
    k *sdk.Kronos
}

func NewATRRiskManaged(k *sdk.Kronos) *ATRRiskManaged {
    return &ATRRiskManaged{k: k}
}

func (s *ATRRiskManaged) GetSignals() ([]*strategy.Signal, error) {
    btc := s.k.Asset("BTC")
    
    rsi := s.k.Indicators.RSI(btc, 14)
    price := s.k.Market.Price(btc)
    atr := s.k.Indicators.ATR(btc, 14)
    
    // Check volatility
    atrPercent := atr.Div(price).Mul(decimal.NewFromInt(100))
    if atrPercent.GreaterThan(decimal.NewFromInt(5)) {
        s.k.Log().MarketCondition("Volatility too high: %s%%", atrPercent)
        return nil, nil
    }
    
    if rsi.LessThan(decimal.NewFromInt(30)) {
        // Dynamic stops based on ATR
        stopLoss := price.Sub(atr.Mul(decimal.NewFromInt(2)))
        takeProfit := price.Add(atr.Mul(decimal.NewFromInt(3)))
        
        // Position size based on risk
        accountBalance := decimal.NewFromInt(10000)
        riskAmount := accountBalance.Mul(decimal.NewFromFloat(0.02))  // 2% risk
        stopDistance := atr.Mul(decimal.NewFromInt(2))
        quantity := riskAmount.Div(stopDistance)
        
        return []*strategy.Signal{
            s.Signal().
                Buy(btc).
                Quantity(quantity).
                StopLoss(stopLoss).
                TakeProfit(takeProfit).
                Reason("ATR-managed entry").
                Build(),
        }, nil
    }
    
    return nil, nil
}

func (s *ATRRiskManaged) GetName() strategy.StrategyName { return "ATR-Risk" }
func (s *ATRRiskManaged) GetDescription() string { return "ATR-based risk management" }
func (s *ATRRiskManaged) GetRiskLevel() strategy.RiskLevel { return strategy.RiskLevelLow }
func (s *ATRRiskManaged) GetStrategyType() strategy.StrategyType { return strategy.StrategyTypeTechnical }
```

## Multi-Asset Portfolio

Trade multiple assets with allocation:

```go
type Portfolio struct {
    k *sdk.Kronos
}

func NewPortfolio(k *sdk.Kronos) *Portfolio {
    return &Portfolio{k: k}
}

func (s *Portfolio) GetSignals() ([]*strategy.Signal, error) {
    btc := s.k.Asset("BTC")
    eth := s.k.Asset("ETH")
    sol := s.k.Asset("SOL")
    
    var signals []*strategy.Signal
    
    // Check each asset
    assets := []struct {
        asset types.Asset
        size  float64
    }{
        {btc, 0.1},
        {eth, 1.0},
        {sol, 10.0},
    }
    
    for _, a := range assets {
        rsi := s.k.Indicators.RSI(a.asset, 14)
        sma200 := s.k.Indicators.SMA(a.asset, 200)
        price := s.k.Market.Price(a.asset)
        
        // Buy if oversold and in uptrend
        if rsi.LessThan(decimal.NewFromInt(30)) && price.GreaterThan(sma200) {
            signals = append(signals,
                s.Signal().
                    Buy(a.asset).
                    Quantity(decimal.NewFromFloat(a.size)).
                    Reason("Oversold in uptrend").
                    Build())
        }
    }
    
    return signals, nil
}

func (s *Portfolio) GetName() strategy.StrategyName { return "Portfolio" }
func (s *Portfolio) GetDescription() string { return "Multi-asset portfolio strategy" }
func (s *Portfolio) GetRiskLevel() strategy.RiskLevel { return strategy.RiskLevelMedium }
func (s *Portfolio) GetStrategyType() strategy.StrategyType { return strategy.StrategyTypeTechnical }
```

## Arbitrage

Find price differences across exchanges:

```go
type Arbitrage struct {
    k *sdk.Kronos
}

func NewArbitrage(k *sdk.Kronos) *Arbitrage {
    return &Arbitrage{k: k}
}

func (s *Arbitrage) GetSignals() ([]*strategy.Signal, error) {
    btc := s.k.Asset("BTC")
    
    // Get prices from all exchanges
    prices := s.k.Market.Prices(btc)
    
    // Find min and max
    var minPrice, maxPrice decimal.Decimal
    var minExchange, maxExchange connector.ExchangeType
    
    for exchange, price := range prices {
        if minPrice.IsZero() || price.LessThan(minPrice) {
            minPrice = price
            minExchange = exchange
        }
        if maxPrice.IsZero() || price.GreaterThan(maxPrice) {
            maxPrice = price
            maxExchange = exchange
        }
    }
    
    // Calculate spread
    spread := maxPrice.Sub(minPrice).Div(minPrice).Mul(decimal.NewFromInt(100))
    
    // If spread > 0.5%, arbitrage
    if spread.GreaterThan(decimal.NewFromFloat(0.5)) {
        s.k.Log().Opportunity("Arbitrage", btc.Symbol(),
            "Buy %s @ %s, Sell %s @ %s, Spread: %s%%",
            minExchange, minPrice, maxExchange, maxPrice, spread)
        
        return []*strategy.Signal{
            s.Signal().Buy(btc).Exchange(minExchange).Quantity(decimal.NewFromFloat(0.1)).Build(),
            s.Signal().Sell(btc).Exchange(maxExchange).Quantity(decimal.NewFromFloat(0.1)).Build(),
        }, nil
    }
    
    return nil, nil
}

func (s *Arbitrage) GetName() strategy.StrategyName { return "Arbitrage" }
func (s *Arbitrage) GetDescription() string { return "Cross-exchange arbitrage" }
func (s *Arbitrage) GetRiskLevel() strategy.RiskLevel { return strategy.RiskLevelLow }
func (s *Arbitrage) GetStrategyType() strategy.StrategyType { return strategy.StrategyTypeArbitrage }
```

## Next Steps

- **[Configuration](configuration)** - Configure your strategy for production
- **[Writing Strategies](writing-strategies)** - Learn advanced patterns
- **[Indicators Reference](../api/indicators/rsi)** - Full indicator docs
