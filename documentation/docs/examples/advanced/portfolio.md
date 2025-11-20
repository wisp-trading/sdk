---
sidebar_position: 1
---

# Multi-Asset Portfolio Strategy

Trade multiple assets with individual analysis and allocation.

## Strategy Overview

- **Type**: Portfolio
- **Indicators**: RSI (14), SMA(200)
- **Risk Level**: Medium
- **Assets**: Multiple (BTC, ETH, SOL)

## Complete Code

```go
package main

import (
	sdk "github.com/backtesting-org/kronos-sdk/pkg/kronos"
	"github.com/backtesting-org/kronos-sdk/pkg/types/connector"
	"github.com/backtesting-org/kronos-sdk/pkg/types/portfolio"
	"github.com/backtesting-org/kronos-sdk/pkg/types/strategy"
	"github.com/shopspring/decimal"
)

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
		asset portfolio.Asset
		size  float64
	}{
		{btc, 0.1},
		{eth, 1.0},
		{sol, 10.0},
	}

	for _, a := range assets {
		rsi, _ := s.k.Indicators().RSI(a.asset, 14)
		sma200, _ := s.k.Indicators().SMA(a.asset, 200)
		price, _ := s.k.Market().Price(a.asset)

		// Buy if oversold and in uptrend
		if rsi.LessThan(decimal.NewFromInt(30)) && price.GreaterThan(sma200) {
			s.k.Log().Opportunity("Portfolio", a.asset.Symbol(), "Oversold in uptrend")
			signal := s.k.Signal(s.GetName()).
				Buy(a.asset, connector.Binance, decimal.NewFromFloat(a.size)).
				Build()
			signals = append(signals, signal)
		}
	}

	return signals, nil
}

func (s *Portfolio) GetName() strategy.StrategyName         { return "Portfolio" }
func (s *Portfolio) GetDescription() string                 { return "Multi-asset portfolio strategy" }
func (s *Portfolio) GetRiskLevel() strategy.RiskLevel       { return strategy.RiskLevelMedium }
func (s *Portfolio) GetStrategyType() strategy.StrategyType { return strategy.StrategyTypeTechnical }
```

## How It Works

1. **Define Universe**: Set up multiple assets to trade
2. **Individual Analysis**: Check each asset independently
3. **Entry Criteria**: Buy when oversold AND in uptrend
4. **Return Multiple Signals**: Can trade multiple assets simultaneously

## Key Concepts

- **Parallel Analysis**: Kronos handles data for all assets automatically
- **Different Sizes**: Position sizes vary by asset (0.1 BTC, 1.0 ETH, 10.0 SOL)
- **Diversification**: Spreads risk across multiple assets
- **Same Logic**: Each asset uses identical entry criteria

## Backtesting

Run with:

```bash
kronos backtest
```

Expected characteristics:
- More trading opportunities (multiple assets)
- Better diversification
- Reduced overall portfolio risk
- May require more capital

## Improvements

Consider adding:
- Dynamic allocation (adjust sizes based on volatility)
- Correlation filtering (avoid highly correlated positions)
- Total exposure limits
- Individual stops per asset
- Rebalancing logic

## Related Strategies

- [ATR Risk Management](../intermediate/atr-risk) - Add per-asset risk control
- [Arbitrage](arbitrage) - Cross-asset opportunity
