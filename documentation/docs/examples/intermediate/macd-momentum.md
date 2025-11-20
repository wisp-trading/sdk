---
sidebar_position: 2
---

# MACD Momentum with Trend Filter

Trade MACD crossovers only with the prevailing trend.

## Strategy Overview

- **Type**: Momentum
- **Indicators**: MACD (12, 26, 9), SMA(200)
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

type MACDMomentum struct {
	k *sdk.Kronos
}

func NewMACDMomentum(k *sdk.Kronos) *MACDMomentum {
	return &MACDMomentum{k: k}
}

func (s *MACDMomentum) GetSignals() ([]*strategy.Signal, error) {
	btc := s.k.Asset("BTC")

	macd, _ := s.k.Indicators().MACD(btc, 12, 26, 9)
	sma200, _ := s.k.Indicators().SMA(btc, 200)
	price, _ := s.k.Market().Price(btc)

	// Only trade with the trend
	inUptrend := price.GreaterThan(sma200)
	inDowntrend := price.LessThan(sma200)

	// Bullish crossover in uptrend
	if macd.MACD.GreaterThan(macd.Signal) &&
		macd.Histogram.GreaterThan(decimal.Zero) &&
		inUptrend {
		s.k.Log().Opportunity("MACD-Momentum", "BTC", "MACD bullish crossover in uptrend")
		signal := s.k.Signal(s.GetName()).
			Buy(btc, connector.Binance, decimal.NewFromFloat(0.15)).
			Build()
		return []*strategy.Signal{signal}, nil
	}

	// Bearish crossover in downtrend
	if macd.MACD.LessThan(macd.Signal) &&
		macd.Histogram.LessThan(decimal.Zero) &&
		inDowntrend {
		s.k.Log().Opportunity("MACD-Momentum", "BTC", "MACD bearish crossover in downtrend")
		signal := s.k.Signal(s.GetName()).
			Sell(btc, connector.Binance, decimal.NewFromFloat(0.15)).
			Build()
		return []*strategy.Signal{signal}, nil
	}

	return nil, nil
}

func (s *MACDMomentum) GetName() strategy.StrategyName { return "MACD-Momentum" }
func (s *MACDMomentum) GetDescription() string { return "MACD with trend filter" }
func (s *MACDMomentum) GetRiskLevel() strategy.RiskLevel { return strategy.RiskLevelMedium }
func (s *MACDMomentum) GetStrategyType() strategy.StrategyType { return strategy.StrategyTypeMomentum }
```

## How It Works

1. **Determine Trend**: Use 200 SMA to identify trend direction
2. **Check MACD**: Look for bullish or bearish crossovers
3. **Filter**: Only trade crossovers aligned with trend
4. **Execute**: Buy in uptrends, sell in downtrends

## Key Concepts

- **MACD Crossover**: When MACD line crosses signal line
- **Histogram**: Shows strength of crossover
- **Trend Filter**: Prevents counter-trend trades
- **With-Trend Only**: Increases win rate

## Backtesting

Run with:

```bash
kronos backtest
```

Expected characteristics:
- Moderate trade frequency
- Better win rate than raw MACD
- Avoids whipsaws in choppy markets
- May miss early trend entries

## Improvements

Consider adding:
- Multiple timeframe confirmation
- Volume analysis
- Stop loss based on recent swing
- Partial exits on histogram divergence

## Related Strategies

- [Multi-Indicator](multi-indicator) - Adds more confirmation
- [MA Crossover](../basic/ma-crossover) - Pure trend following
