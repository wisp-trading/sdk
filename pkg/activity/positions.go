package activity

import (
	"context"

	storeActivity "github.com/wisp-trading/sdk/pkg/markets/base/types/stores/activity"
	wispActivity "github.com/wisp-trading/sdk/pkg/types/wisp/activity"
)

type positions struct {
	store storeActivity.Positions
}

func NewPositions(store storeActivity.Positions) wispActivity.Positions {
	return &positions{store: store}
}

func (p *positions) GetOrderCount(_ context.Context) int64 {
	return p.store.GetTotalOrderCount()
}

var _ wispActivity.Positions = (*positions)(nil)
