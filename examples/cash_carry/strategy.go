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
		rsi, _ := s.k.Indicators.RSI(a.asset, 14)
		sma200, _ := s.k.Indicators.SMA(a.asset, 200)
		price, _ := s.k.Market.Price(a.asset)

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
