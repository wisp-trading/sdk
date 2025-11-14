package position

import (
	"fmt"

	"github.com/backtesting-org/kronos-sdk/pkg/types/connector"
	"github.com/backtesting-org/kronos-sdk/pkg/types/strategy"
)

func (ds *dataStore) ReconcilePosition(strategyName strategy.StrategyName, exchangePos connector.Position) error {
	ds.mutex.RLock()
	defer ds.mutex.RUnlock()

	// Get strategy execution
	executions := ds.getExecutions()
	execution := executions[strategyName]

	if execution == nil {
		return fmt.Errorf("strategy execution not found: %s", strategyName)
	}

	// TODO: Implement position reconciliation logic
	// This would compute position from trades and compare with exchangePos
	// For now, just validate that the strategy exists

	return nil
}
