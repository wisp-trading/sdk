---
sidebar_position: 3
---

# ATR-Based Risk Management

Dynamic stops and position sizing based on volatility.

## Strategy Overview

- **Type**: Technical with Risk Management
- **Indicators**: RSI (14), ATR (14)
- **Risk Level**: Low
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

type ATRRiskManaged struct {
	k *sdk.Kronos
}

func NewATRRiskManaged(k *sdk.Kronos) *ATRRiskManaged {
	return &ATRRiskManaged{k: k}
}

func (s *ATRRiskManaged) GetSignals() ([]*strategy.Signal, error) {
	btc := s.k.Asset("BTC")

	rsi, _ := s.k.Indicators().RSI(btc, 14)
	price, _ := s.k.Market().Price(btc)
	atr, _ := s.k.Indicators().ATR(btc, 14)

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
		riskAmount := accountBalance.Mul(decimal.NewFromFloat(0.02)) // 2% risk
		stopDistance := atr.Mul(decimal.NewFromInt(2))
		quantity := riskAmount.Div(stopDistance)

		s.k.Log().Opportunity("ATR-Risk", "BTC",
			"ATR-managed entry - Stop: %s, Target: %s, Size: %s",
			stopLoss, takeProfit, quantity)

		signal := s.k.Signal(s.GetName()).
			Buy(btc, connector.Binance, quantity).
			Build()
		return []*strategy.Signal{signal}, nil
	}

	return nil, nil
}

func (s *ATRRiskManaged) GetName() strategy.StrategyName   { return "ATR-Risk" }
func (s *ATRRiskManaged) GetDescription() string           { return "ATR-based risk management" }
func (s *ATRRiskManaged) GetRiskLevel() strategy.RiskLevel { return strategy.RiskLevelLow }
func (s *ATRRiskManaged) GetStrategyType() strategy.StrategyType {
	return strategy.StrategyTypeTechnical
}
```

## How It Works

1. **Check Volatility**: Calculate ATR as % of price
2. **Volatility Filter**: Skip if market too volatile (>5%)
3. **Dynamic Stops**: Stop loss at 2× ATR, take profit at 3× ATR
4. **Risk-Based Sizing**: Position size to risk exactly 2% of account

## Key Concepts

- **ATR (Average True Range)**: Measures volatility
- **Dynamic Stops**: Adapt to market conditions
- **Position Sizing**: Larger positions in low volatility, smaller in high
- **Fixed Risk**: Always risk 2% of account per trade
- **Risk:Reward**: 1:1.5 ratio (2× ATR stop, 3× ATR target)

## Backtesting

Run with:

```bash
kronos backtest
```

Expected characteristics:
- Consistent risk per trade
- Better capital preservation
- Avoids trading in extreme volatility
- Position sizes vary with market conditions

## Improvements

Consider adding:
- Trailing stops (move stop as profit grows)
- Partial exits (scale out winners)
- ATR-based entry timing (enter on low ATR)
- Multiple asset support

## Related Strategies

- [RSI](../basic/rsi) - Same entry, no risk management
- [Portfolio](../advanced/portfolio) - Multi-asset version
