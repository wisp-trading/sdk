---
sidebar_position: 1
---

# Multi-Indicator Confirmation

Require multiple indicators to agree before trading.

## Strategy Overview

- **Type**: Technical
- **Indicators**: RSI, Stochastic, Bollinger Bands, MACD
- **Risk Level**: Low
- **Assets**: Single asset (BTC)

## Complete Code

```go
package main

import (
    "fmt"
    sdk "github.com/backtesting-org/kronos-sdk/pkg/kronos"
    "github.com/backtesting-org/kronos-sdk/pkg/types/strategy"
    "github.com/shopspring/decimal"
)

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

## How It Works

1. **Poll Indicators**: Get RSI, Stochastic, Bollinger Bands, MACD
2. **Check Each**: Count how many show bullish signals
3. **Require Consensus**: Need 3 out of 4 to agree
4. **Trade with Confidence**: Larger position size due to multiple confirmations

## Key Concepts

- **Consensus Approach**: Multiple indicators reduce false signals
- **Lower Risk**: Less likely to be wrong when multiple signals align
- **Larger Position**: 0.2 BTC (more confident)
- **Flexible Threshold**: Easy to adjust (2/4, 3/4, 4/4)

## Backtesting

Run with:

```bash
kronos backtest
```

Expected characteristics:
- Low trade frequency (high bar for entry)
- Higher win rate
- Catches only the best setups
- May miss some opportunities

## Improvements

Consider adding:
- Bearish signal detection (for sells)
- Weighted voting (some indicators more important)
- Time-based confirmation (signals must persist)
- Dynamic threshold based on market regime

## Related Strategies

- [MACD Momentum](macd-momentum) - Uses one indicator
- [ATR Risk Management](atr-risk) - Adds position sizing
