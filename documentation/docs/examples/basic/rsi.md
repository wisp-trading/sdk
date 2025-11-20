---
sidebar_position: 1
---

# Simple RSI Strategy

Classic momentum strategy using RSI oversold/overbought levels.

## Strategy Overview

- **Type**: Momentum
- **Indicators**: RSI (14 periods)
- **Risk Level**: Medium
- **Assets**: Single asset (BTC)

## Complete Code

```go
package main

import (
	sdk "github.com/backtesting-org/kronos-sdk/pkg/kronos"
	"github.com/backtesting-org/kronos-sdk/pkg/types/connector"
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
	rsi, _ := s.k.Indicators().RSI(btc, 14)

	// Buy oversold
	if rsi.LessThan(decimal.NewFromInt(30)) {
		s.k.Log().Opportunity("RSI", "BTC", "RSI oversold: %s", rsi.String())
		signal := s.k.Signal(s.GetName()).
			Buy(btc, connector.Binance, decimal.NewFromFloat(0.1)).
			Build()
		return []*strategy.Signal{signal}, nil
	}

	// Sell overbought
	if rsi.GreaterThan(decimal.NewFromInt(70)) {
		s.k.Log().Opportunity("RSI", "BTC", "RSI overbought: %s", rsi.String())
		signal := s.k.Signal(s.GetName()).
			Sell(btc, connector.Binance, decimal.NewFromFloat(0.1)).
			Build()
		return []*strategy.Signal{signal}, nil
	}

	return nil, nil
}

func (s *RSIStrategy) GetName() strategy.StrategyName { return "RSI" }
func (s *RSIStrategy) GetDescription() string { return "Simple RSI momentum" }
func (s *RSIStrategy) GetRiskLevel() strategy.RiskLevel { return strategy.RiskLevelMedium }
func (s *RSIStrategy) GetStrategyType() strategy.StrategyType { return strategy.StrategyTypeTechnical }
```

## How It Works

1. **Check RSI**: Get the 14-period RSI for BTC
2. **Oversold**: When RSI < 30, buy
3. **Overbought**: When RSI > 70, sell
4. **Wait**: Otherwise, do nothing

## Key Concepts

- **RSI < 30**: Asset is oversold, potential reversal up
- **RSI > 70**: Asset is overbought, potential reversal down
- **Fixed quantity**: Always trades 0.1 BTC

## Backtesting

Run with:

```bash
kronos backtest
```

Expected characteristics:
- Moderate trade frequency
- Works best in ranging markets
- May whipsaw in strong trends

## Improvements

Consider adding:
- Trend filter (only buy in uptrend)
- Stop loss protection
- Multiple timeframe confirmation
- Dynamic position sizing

## Related Strategies

- [Multi-Indicator Confirmation](../intermediate/multi-indicator) - Adds more signals
- [ATR Risk Management](../intermediate/atr-risk) - Adds dynamic stops
