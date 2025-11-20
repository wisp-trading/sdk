---
sidebar_position: 2
---

# Moving Average Crossover

Golden cross / death cross trend following strategy.

## Strategy Overview

- **Type**: Trend Following
- **Indicators**: SMA(50), SMA(200)
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

type MACrossover struct {
	k *sdk.Kronos
}

func NewMACrossover(k *sdk.Kronos) *MACrossover {
	return &MACrossover{k: k}
}

func (s *MACrossover) GetSignals() ([]*strategy.Signal, error) {
	btc := s.k.Asset("BTC")

	sma50, _ := s.k.Indicators().SMA(btc, 50)
	sma200, _ := s.k.Indicators().SMA(btc, 200)
	price, _ := s.k.Market().Price(btc)

	// Golden cross: 50 crosses above 200
	if sma50.GreaterThan(sma200) && price.GreaterThan(sma50) {
		s.k.Log().Opportunity("MA-Crossover", "BTC", "Golden cross detected")
		signal := s.k.Signal(s.GetName()).
			Buy(btc, connector.Binance, decimal.NewFromFloat(0.2)).
			Build()
		return []*strategy.Signal{signal}, nil
	}

	// Death cross: 50 crosses below 200
	if sma50.LessThan(sma200) {
		s.k.Log().Opportunity("MA-Crossover", "BTC", "Death cross detected")
		signal := s.k.Signal(s.GetName()).
			Sell(btc, connector.Binance, decimal.NewFromFloat(0.2)).
			Build()
		return []*strategy.Signal{signal}, nil
	}

	return nil, nil
}

func (s *MACrossover) GetName() strategy.StrategyName         { return "MA-Crossover" }
func (s *MACrossover) GetDescription() string                 { return "Golden/Death cross strategy" }
func (s *MACrossover) GetRiskLevel() strategy.RiskLevel       { return strategy.RiskLevelLow }
func (s *MACrossover) GetStrategyType() strategy.StrategyType { return strategy.StrategyType("Trend")
```

## How It Works

1. **Calculate SMAs**: Get 50-period and 200-period SMAs
2. **Golden Cross**: When SMA(50) > SMA(200) and price is above SMA(50), buy
3. **Death Cross**: When SMA(50) < SMA(200), sell
4. **Confirm**: Only buy when price is also above the fast MA

## Key Concepts

- **Golden Cross**: Bullish signal, fast MA crosses above slow MA
- **Death Cross**: Bearish signal, fast MA crosses below slow MA
- **Price filter**: Ensures momentum in direction of signal
- **Larger position**: Uses 0.2 BTC (more confident with trend confirmation)

## Backtesting

Run with:

```bash
kronos backtest
```

Expected characteristics:
- Low trade frequency (few signals per year)
- Catches major trends
- Lags at trend changes (by design)
- Best in trending markets

## Improvements

Consider adding:
- Volume confirmation
- Volatility filter (avoid low-volume periods)
- Partial exits (scale out of winners)
- Additional timeframe for confirmation

## Related Strategies

- [Bollinger Mean Reversion](bollinger-bands) - Opposite approach
- [MACD Momentum](../intermediate/macd-momentum) - Uses EMAs instead
