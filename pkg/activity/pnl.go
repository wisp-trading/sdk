package activity

import (
	"context"

	perpTypes "github.com/wisp-trading/sdk/pkg/markets/perp/types"
	predTypes "github.com/wisp-trading/sdk/pkg/markets/prediction/types"
	wispActivity "github.com/wisp-trading/sdk/pkg/types/wisp/activity"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

// pnl aggregates PNL across all market domains.
// Each domain service is responsible for its own calculation.
type pnl struct {
	spot    wispActivity.SpotPNL
	perp    perpTypes.PerpPNL
	predict predTypes.PredictionPNL
}

func NewPNL(
	spot wispActivity.SpotPNL,
	perp perpTypes.PerpPNL,
	predict predTypes.PredictionPNL,
) wispActivity.PNL {
	return &pnl{spot: spot, perp: perp, predict: predict}
}

func (p *pnl) Spot() wispActivity.SpotPNL { return p.spot }

func (p *pnl) TotalRealized(ctx context.Context) numerical.Decimal {
	return p.spot.Realized(ctx).
		Add(p.perp.Realized(ctx)).
		Add(p.predict.Realized(ctx))
}

func (p *pnl) TotalUnrealized(ctx context.Context) numerical.Decimal {
	return p.spot.Unrealized(ctx).
		Add(p.perp.Unrealized(ctx)).
		Add(p.predict.Unrealized(ctx))
}

func (p *pnl) TotalFees(ctx context.Context) numerical.Decimal {
	return p.spot.Fees(ctx).
		Add(p.perp.Fees(ctx)).
		Add(p.predict.Fees(ctx))
}

var _ wispActivity.PNL = (*pnl)(nil)
