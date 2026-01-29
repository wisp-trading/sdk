package position

import (
	"sync"
	"sync/atomic"

	"github.com/wisp-trading/wisp/pkg/types/temporal"
)

type dataStore struct {
	timeProvider temporal.TimeProvider
	executions   atomic.Value // portfolioTypes.StrategyExecutionMap
	lastUpdated  atomic.Value // portfolioTypes.LastUpdatedMap
	mutex        sync.RWMutex
}
