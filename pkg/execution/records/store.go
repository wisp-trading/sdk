package records

import (
	"sync"

	"github.com/google/uuid"
	"github.com/wisp-trading/sdk/pkg/types/execution"
	"github.com/wisp-trading/sdk/pkg/types/strategy"
)

type store struct {
	mu      sync.RWMutex
	records []execution.ExecutionRecord
	byID    map[uuid.UUID]int // signalID -> index in records slice
}

// NewStore creates a new in-memory ExecutionRecords store.
func NewStore() execution.ExecutionRecords {
	return &store{
		records: make([]execution.ExecutionRecord, 0),
		byID:    make(map[uuid.UUID]int),
	}
}

func (s *store) Add(record execution.ExecutionRecord) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.byID[record.SignalID] = len(s.records)
	s.records = append(s.records, record)
}

func (s *store) GetAll() []execution.ExecutionRecord {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]execution.ExecutionRecord, len(s.records))
	copy(out, s.records)
	return out
}

func (s *store) GetBySignalID(id uuid.UUID) *execution.ExecutionRecord {
	s.mu.RLock()
	defer s.mu.RUnlock()
	idx, ok := s.byID[id]
	if !ok {
		return nil
	}
	r := s.records[idx]
	return &r
}

func (s *store) GetByStrategy(name strategy.StrategyName) []execution.ExecutionRecord {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var out []execution.ExecutionRecord
	for _, r := range s.records {
		if r.Strategy == name {
			out = append(out, r)
		}
	}
	return out
}

var _ execution.ExecutionRecords = (*store)(nil)
