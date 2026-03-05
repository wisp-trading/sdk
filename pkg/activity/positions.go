package activity

import (
	"context"

	"github.com/wisp-trading/sdk/pkg/types/connector"
	storeActivity "github.com/wisp-trading/sdk/pkg/types/data/stores/activity"
	"github.com/wisp-trading/sdk/pkg/types/strategy"
	wispActivity "github.com/wisp-trading/sdk/pkg/types/wisp/activity"
)

type positions struct {
	store storeActivity.Positions
}

func NewPositions(store storeActivity.Positions) wispActivity.Positions {
	return &positions{store: store}
}

func (p *positions) GetStrategyExecution(name strategy.StrategyName) *strategy.StrategyExecution {
	return p.store.GetStrategyExecution(name)
}

func (p *positions) GetTradesForStrategy(name strategy.StrategyName) []connector.Trade {
	return p.store.GetTradesForStrategy(name)
}

func (p *positions) GetAllStrategyExecutions() map[strategy.StrategyName]*strategy.StrategyExecution {
	return p.store.GetAllStrategyExecutions()
}

func (p *positions) GetStrategyForOrder(_ context.Context, orderID string) (strategy.StrategyName, bool) {
	return p.store.GetStrategyForOrder(orderID)
}

func (p *positions) GetTotalOrderCount(_ context.Context) int64 {
	return p.store.GetTotalOrderCount()
}

var _ wispActivity.Positions = (*positions)(nil)
