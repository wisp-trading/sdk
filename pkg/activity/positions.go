package activity

import (
	"context"

	perpTypes "github.com/wisp-trading/sdk/pkg/markets/perp/types"
	spotTypes "github.com/wisp-trading/sdk/pkg/markets/spot/types"
	wispActivity "github.com/wisp-trading/sdk/pkg/types/wisp/activity"
)

// positions aggregates order counts across spot and perp domains.
type positions struct {
	spot spotTypes.SpotPositions
	perp perpTypes.PerpPositions
}

func NewPositions(
	spot spotTypes.SpotPositions,
	perp perpTypes.PerpPositions,
) wispActivity.Positions {
	return &positions{spot: spot, perp: perp}
}

func (p *positions) GetOrderCount(_ context.Context) int64 {
	return p.spot.GetTotalOrderCount() + p.perp.GetTotalOrderCount()
}

var _ wispActivity.Positions = (*positions)(nil)
