package activity

import (
	"context"

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

func (p *positions) GetStrategyExecution(ctx context.Context) *strategy.StrategyExecution {
	name, ok := strategy.FromContext(ctx)
	if !ok {
		return nil
	}
	return p.store.GetStrategyExecution(name)
}

func (p *positions) GetTradesForStrategy(ctx context.Context) []connector.Trade {
	name, ok := strategy.FromContext(ctx)
	if !ok {
		return nil
	}
	return p.store.GetTradesForStrategy(name)
}

func (p *positions) GetAllStrategyExecutions(ctx context.Context) map[strategy.StrategyName]*strategy.StrategyExecution {
	return p.store.GetAllStrategyExecutions()
}

func (p *positions) GetStrategyForOrder(ctx context.Context, orderID string) (strategy.StrategyName, bool) {
	return p.store.GetStrategyForOrder(orderID)
}

func (p *positions) GetTotalOrderCount(ctx context.Context) int64 {
	return p.store.GetTotalOrderCount()
}

var _ kronosActivity.Positions = (*positions)(nil)
