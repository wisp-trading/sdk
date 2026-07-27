package activity

import (
	"context"

	baseActivity "github.com/wisp-trading/sdk/pkg/markets/base/activity"
	perpTypes "github.com/wisp-trading/sdk/pkg/markets/perp/types"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

type perpPNL struct {
	store perpTypes.MarketStore
}

func NewPerpPNL(store perpTypes.MarketStore) perpTypes.PerpPNL {
	return &perpPNL{store: store}
}

func (p *perpPNL) Positions(_ context.Context) []perpTypes.PerpPositionPNL {
	positions := p.store.GetPositions()
	results := make([]perpTypes.PerpPositionPNL, 0, len(positions))
	for _, pos := range positions {
		results = append(results, perpTypes.PerpPositionPNL{
			Position:   pos,
			Realized:   pos.RealizedPnL,
			Unrealized: pos.UnrealizedPnL,
		})
	}
	return results
}

func (p *perpPNL) Realized(_ context.Context) numerical.Decimal {
	vals := make([]numerical.Decimal, 0)
	for _, pos := range p.store.GetPositions() {
		vals = append(vals, pos.RealizedPnL)
	}
	return baseActivity.SumDecimals(vals)
}

func (p *perpPNL) Unrealized(_ context.Context) numerical.Decimal {
	vals := make([]numerical.Decimal, 0)
	for _, pos := range p.store.GetPositions() {
		vals = append(vals, pos.UnrealizedPnL)
	}
	return baseActivity.SumDecimals(vals)
}

func (p *perpPNL) Fees(_ context.Context) numerical.Decimal {
	return baseActivity.SumTradeFees(p.store.GetAllTrades())
}

var _ perpTypes.PerpPNL = (*perpPNL)(nil)
