package position

import (
	"sync"
	"sync/atomic"

	"github.com/backtesting-org/kronos-sdk/pkg/types/temporal"
)

type dataStore struct {
	timeProvider temporal.TimeProvider
	executions   atomic.Value // portfolioTypes.StrategyExecutionMap
	lastUpdated  atomic.Value // portfolioTypes.LastUpdatedMap
	mutex        sync.RWMutex
}
