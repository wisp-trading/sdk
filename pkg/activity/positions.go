package activity

import (
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

func (p *positions) GetStrategyExecution(ctx strategy.StrategyContext) *strategy.StrategyExecution {
	return p.store.GetStrategyExecution(ctx.StrategyName())
}

func (p *positions) GetTradesForStrategy(ctx strategy.StrategyContext) []connector.Trade {
	return p.store.GetTradesForStrategy(ctx.StrategyName())
}

func (p *positions) GetAllStrategyExecutions(ctx strategy.StrategyContext) map[strategy.StrategyName]*strategy.StrategyExecution {
	return p.store.GetAllStrategyExecutions()
}

func (p *positions) GetStrategyForOrder(ctx strategy.StrategyContext, orderID string) (strategy.StrategyName, bool) {
	return p.store.GetStrategyForOrder(orderID)
}

func (p *positions) GetTotalOrderCount(ctx strategy.StrategyContext) int64 {
	return p.store.GetTotalOrderCount()
}

var _ wispActivity.Positions = (*positions)(nil)
