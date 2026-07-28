package activity

import (
	"context"

	baseActivity "github.com/wisp-trading/sdk/pkg/markets/base/activity"
	onchainTypes "github.com/wisp-trading/sdk/pkg/markets/onchain/types"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

type onchainPNL struct {
	store onchainTypes.MarketStore
}

func NewOnchainPNL(store onchainTypes.MarketStore) onchainTypes.OnchainPNL {
	return &onchainPNL{store: store}
}

func (s *onchainPNL) Positions(_ context.Context) []onchainTypes.PositionPNL {
	trackers := baseActivity.BuildTrackers(s.store.GetAllTrades())
	results := make([]onchainTypes.PositionPNL, 0, len(trackers))
	for _, t := range trackers {
		results = append(results, onchainTypes.PositionPNL{
			Pair:       t.Pair,
			Exchange:   t.Exchange,
			Realized:   t.Realized,
			Unrealized: numerical.Zero(),
			Fees:       t.Fees,
		})
	}
	return results
}

func (s *onchainPNL) Realized(_ context.Context) numerical.Decimal {
	trades := s.store.GetAllTrades()
	vals := make([]numerical.Decimal, 0)
	for _, t := range baseActivity.BuildTrackers(trades) {
		vals = append(vals, t.Realized)
	}
	return baseActivity.SumDecimals(vals).Sub(baseActivity.SumTradeFees(trades))
}

func (s *onchainPNL) Unrealized(_ context.Context) numerical.Decimal {
	return numerical.Zero()
}

func (s *onchainPNL) Fees(_ context.Context) numerical.Decimal {
	return baseActivity.SumTradeFees(s.store.GetAllTrades())
}
