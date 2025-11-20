---
sidebar_position: 3
---

# Bollinger Bands Mean Reversion

Buy at lower band, sell at upper band - fade extremes.

## Strategy Overview

- **Type**: Mean Reversion
- **Indicators**: Bollinger Bands (20, 2.0), RSI (14)
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

type BollingerMeanReversion struct {
	k *sdk.Kronos
}

func NewBollingerMR(k *sdk.Kronos) *BollingerMeanReversion {
	return &BollingerMeanReversion{k: k}
}

func (s *BollingerMeanReversion) GetSignals() ([]*strategy.Signal, error) {
	btc := s.k.Asset("BTC")

	bb, _ := s.k.Indicators().BollingerBands(btc, 20, 2.0)
	price, _ := s.k.Market().Price(btc)
	rsi, _ := s.k.Indicators().RSI(btc, 14)

	// Buy at lower band with RSI confirmation
	if price.LessThan(bb.Lower) && rsi.LessThan(decimal.NewFromInt(35)) {
		s.k.Log().Opportunity("BB-MeanReversion", "BTC",
			"Mean reversion from lower band, target middle: %s", bb.Middle)
		signal := s.k.Signal(s.GetName()).
			Buy(btc, connector.Binance, decimal.NewFromFloat(0.1)).
			Build()
		return []*strategy.Signal{signal}, nil
	}

	// Sell at upper band
	if price.GreaterThan(bb.Upper) && rsi.GreaterThan(decimal.NewFromInt(65)) {
		s.k.Log().Opportunity("BB-MeanReversion", "BTC",
			"Mean reversion from upper band, target middle: %s", bb.Middle)
		signal := s.k.Signal(s.GetName()).
			Sell(btc, connector.Binance, decimal.NewFromFloat(0.1)).
			Build()
		return []*strategy.Signal{signal}, nil
	}

	return nil, nil
}

func (s *BollingerMeanReversion) GetName() strategy.StrategyName { return "BB-MeanReversion" }
func (s *BollingerMeanReversion) GetDescription() string { return "Bollinger Bands mean reversion" }
func (s *BollingerMeanReversion) GetRiskLevel() strategy.RiskLevel { return strategy.RiskLevelMedium }
func (s *BollingerMeanReversion) GetStrategyType() strategy.StrategyType { return strategy.StrategyTypeMeanReversion }
```

## How It Works

1. **Calculate Bands**: Get Bollinger Bands (20-period SMA ± 2 std dev)
2. **Lower Band Touch**: When price < lower band and RSI < 35, buy
3. **Upper Band Touch**: When price > upper band and RSI > 65, sell
4. **Target Middle**: Take profit at the middle band (mean reversion)

## Key Concepts

- **Bollinger Bands**: Measure volatility and identify extremes
- **Mean Reversion**: Assumes price returns to average
- **RSI Confirmation**: Filters false signals
- **Automatic Take Profit**: Exits at middle band

## Backtesting

Run with:

```bash
kronos backtest
```

Expected characteristics:
- Moderate trade frequency
- Works best in ranging markets
- Struggles in strong trends
- Quick wins, defined exits

## Improvements

Consider adding:
- Trend filter (avoid counter-trend trades)
- Band squeeze detection (potential breakouts)
- Dynamic take profit (not always middle)
- Stop loss at opposite band

## Related Strategies

- [RSI](rsi) - Similar momentum confirmation
- [MA Crossover](ma-crossover) - Opposite (trend following)
