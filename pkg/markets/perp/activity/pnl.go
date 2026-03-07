package activity

import (
	"context"

	perpTypes "github.com/wisp-trading/sdk/pkg/markets/perp/types"
	"github.com/wisp-trading/sdk/pkg/types/connector"
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
	total := numerical.Zero()
	for _, pos := range p.store.GetPositions() {
		total = total.Add(pos.RealizedPnL)
	}
	return total
}

func (p *perpPNL) Unrealized(_ context.Context) numerical.Decimal {
	total := numerical.Zero()
	for _, pos := range p.store.GetPositions() {
		total = total.Add(pos.UnrealizedPnL)
	}
	return total
}

func (p *perpPNL) Fees(_ context.Context) numerical.Decimal {
	return sumTradeFees(p.store.GetAllTrades())
}

func sumTradeFees(trades []connector.Trade) numerical.Decimal {
	total := numerical.Zero()
	for _, t := range trades {
		total = total.Add(t.Fee)
	}
	return total
}

var _ perpTypes.PerpPNL = (*perpPNL)(nil)
