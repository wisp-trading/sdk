package strategy

import (
	"context"
	"sync"
	"time"
)

const statusLogCap = 100

// BaseStrategyConfig holds configuration for creating a base strategy
type BaseStrategyConfig struct {
	Name StrategyName
}

// BaseStrategy provides common lifecycle and status log management for strategies.
// Concrete strategies embed BaseStrategy and call StartWithRunner(ctx, s.run)
// from their own Start method.
//
// Order placement is NOT here — use market-scoped Emit only:
//
//	sig, err := k.Perp().Signal(name).BuyLimit(...).Build()
//	k.Perp().Emit(sig) // Spot / Predict / Options likewise
type BaseStrategy struct {
	name StrategyName

	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup

	// statusLog is a fixed-capacity ring buffer for operator/monitor views.
	// Written by EmitStatus; read by LatestStatus and StatusLog.
	statusLog  []StrategyStatus
	statusHead int // next write slot
	statusSize int
	statusMu   sync.RWMutex
}

// NewBaseStrategy creates a new BaseStrategy suitable for embedding.
func NewBaseStrategy(config BaseStrategyConfig) *BaseStrategy {
	return &BaseStrategy{
		name:      config.Name,
		statusLog: make([]StrategyStatus, statusLogCap),
	}
}

// StartWithRunner initialises the context, then launches the provided run
// function in a managed goroutine.
func (b *BaseStrategy) StartWithRunner(ctx context.Context, run func(ctx context.Context)) error {
	b.ctx, b.cancel = context.WithCancel(ctx)

	b.wg.Add(1)
	go func() {
		defer b.wg.Done()
		run(b.ctx)
	}()

	return nil
}

// Stop signals the strategy to shut down and waits for the run goroutine to exit.
func (b *BaseStrategy) Stop(_ context.Context) error {
	if b.cancel != nil {
		b.cancel()
	}
	b.wg.Wait()
	return nil
}

// EmitStatus records a status snapshot into the strategy's internal ring buffer.
// Non-blocking. Used by monitoring (LatestStatus / StatusLog) — not order routing.
func (b *BaseStrategy) EmitStatus(s StrategyStatus) {
	if s.At.IsZero() {
		s.At = time.Now()
	}
	b.statusMu.Lock()
	b.statusLog[b.statusHead] = s
	b.statusHead = (b.statusHead + 1) % statusLogCap
	if b.statusSize < statusLogCap {
		b.statusSize++
	}
	b.statusMu.Unlock()
}

// LatestStatus returns the most recently emitted status snapshot, or the zero
// value if none has been emitted yet.
func (b *BaseStrategy) LatestStatus() StrategyStatus {
	b.statusMu.RLock()
	defer b.statusMu.RUnlock()
	if b.statusSize == 0 {
		return StrategyStatus{}
	}
	idx := (b.statusHead - 1 + statusLogCap) % statusLogCap
	return b.statusLog[idx]
}

// StatusLog returns up to the last 100 status snapshots, oldest-first.
func (b *BaseStrategy) StatusLog() []StrategyStatus {
	b.statusMu.RLock()
	defer b.statusMu.RUnlock()
	if b.statusSize == 0 {
		return nil
	}
	out := make([]StrategyStatus, b.statusSize)
	start := (b.statusHead - b.statusSize + statusLogCap) % statusLogCap
	for i := 0; i < b.statusSize; i++ {
		out[i] = b.statusLog[(start+i)%statusLogCap]
	}
	return out
}

// GetName returns the strategy name.
func (b *BaseStrategy) GetName() StrategyName { return b.name }
