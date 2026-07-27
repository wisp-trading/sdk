package activity

import (
	"context"

	baseActivity "github.com/wisp-trading/sdk/pkg/markets/base/activity"
	spotTypes "github.com/wisp-trading/sdk/pkg/markets/spot/types"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

type spotPNL struct {
	store spotTypes.MarketStore
}

func NewSpotPNL(store spotTypes.MarketStore) spotTypes.SpotPNL {
	return &spotPNL{store: store}
}

func (s *spotPNL) Positions(_ context.Context) []spotTypes.PositionPNL {
	trackers := baseActivity.BuildTrackers(s.store.GetAllTrades())
	results := make([]spotTypes.PositionPNL, 0, len(trackers))
	for _, t := range trackers {
		results = append(results, spotTypes.PositionPNL{
			Pair:       t.Pair,
			Exchange:   t.Exchange,
			Realized:   t.Realized,
			Unrealized: s.unrealized(t),
			Fees:       t.Fees,
		})
	}
	return results
}

func (s *spotPNL) Realized(_ context.Context) numerical.Decimal {
	trades := s.store.GetAllTrades()
	vals := make([]numerical.Decimal, 0)
	for _, t := range baseActivity.BuildTrackers(trades) {
		vals = append(vals, t.Realized)
	}
	return baseActivity.SumDecimals(vals).Sub(baseActivity.SumTradeFees(trades))
}

func (s *spotPNL) Unrealized(_ context.Context) numerical.Decimal {
	vals := make([]numerical.Decimal, 0)
	for _, t := range baseActivity.BuildTrackers(s.store.GetAllTrades()) {
		vals = append(vals, s.unrealized(t))
	}
	return baseActivity.SumDecimals(vals)
}

func (s *spotPNL) Fees(_ context.Context) numerical.Decimal {
	return baseActivity.SumTradeFees(s.store.GetAllTrades())
}

func (s *spotPNL) unrealized(t *baseActivity.PositionTracker) numerical.Decimal {
	price := s.store.GetPairPrice(t.Pair, t.Exchange)
	if price == nil {
		return numerical.Zero()
	}
	return baseActivity.UnrealizedFromMark(t, price.Price)
}

var _ spotTypes.SpotPNL = (*spotPNL)(nil)
