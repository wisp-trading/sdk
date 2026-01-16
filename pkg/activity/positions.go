package activity

import (
	"github.com/backtesting-org/kronos-sdk/pkg/types/connector"
	storeActivity "github.com/backtesting-org/kronos-sdk/pkg/types/data/stores/activity"
	kronosActivity "github.com/backtesting-org/kronos-sdk/pkg/types/kronos/activity"
	"github.com/backtesting-org/kronos-sdk/pkg/types/strategy"
)

type positions struct {
	store storeActivity.Positions
}

func NewPositions(store storeActivity.Positions) kronosActivity.Positions {
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

var _ kronosActivity.Positions = (*positions)(nil)
