package position

import (
	"github.com/backtesting-org/kronos-sdk/pkg/types/connector"
	"github.com/backtesting-org/kronos-sdk/pkg/types/strategy"
)

func (ds *dataStore) LinkTradeToStrategy(strategyName strategy.StrategyName, tradeID string) {
	ds.mutex.Lock()
	defer ds.mutex.Unlock()

	current := ds.getExecutions()
	execution := current[strategyName]

	if execution == nil {
		execution = &strategy.StrategyExecution{
			Orders:   []connector.Order{},
			TradeIDs: []string{},
		}
	}

	// Check if trade ID already linked
	for _, id := range execution.TradeIDs {
		if id == tradeID {
			return // Already linked
		}
	}

	// Add trade ID
	execution.TradeIDs = append(execution.TradeIDs, tradeID)

	// Store updated map
	updated := make(map[strategy.StrategyName]*strategy.StrategyExecution, len(current)+1)
	for k, v := range current {
		updated[k] = v
	}
	updated[strategyName] = execution
	ds.executions.Store(updated)
}

func (ds *dataStore) GetTradeIDsForStrategy(strategyName strategy.StrategyName) []string {
	executions := ds.getExecutions()
	execution := executions[strategyName]

	if execution == nil {
		return []string{}
	}

	return execution.TradeIDs
}
