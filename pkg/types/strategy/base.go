package strategy

import (
	"context"
	"sync"
	"time"
)

const signalChannelBufferSize = 64
const statusLogCap = 100

// BaseStrategyConfig holds configuration for creating a base strategy
type BaseStrategyConfig struct {
	Name StrategyName
}

// BaseStrategy provides common lifecycle, signal channel, and status log management
// for all strategies. Concrete strategies embed BaseStrategy and call
// StartWithRunner(ctx, s.run) from their own Start method.
type BaseStrategy struct {
	name StrategyName

	signalCh chan Signal
	ctx      context.Context
	cancel   context.CancelFunc
	wg       sync.WaitGroup

	// marks are pending checkpoints bundled into the next emitted signal's metadata.
	marks   []Mark
	marksMu sync.Mutex

	// statusLog is a fixed-capacity ring buffer owned by the strategy.
	// Written by EmitStatus; read by LatestStatus and StatusLog.
	statusLog  []StrategyStatus
	statusHead int // next write slot
	statusSize int
	statusMu   sync.RWMutex
}

// Mark is a named timestamp checkpoint recorded inside a strategy's run loop.
type Mark struct {
	Label string
	At    int64 // UnixNano
}

// NewBaseStrategy creates a new BaseStrategy suitable for embedding.
func NewBaseStrategy(config BaseStrategyConfig) *BaseStrategy {
	return &BaseStrategy{
		name:      config.Name,
		statusLog: make([]StrategyStatus, statusLogCap),
	}
}

// StartWithRunner initialises the signal channel and context, then launches the
// provided run function in a managed goroutine.
func (b *BaseStrategy) StartWithRunner(ctx context.Context, run func(ctx context.Context)) error {
	b.ctx, b.cancel = context.WithCancel(ctx)
	b.signalCh = make(chan Signal, signalChannelBufferSize)

	b.wg.Add(1)
	go func() {
		defer b.wg.Done()
		defer close(b.signalCh)
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

// Signals returns the read-only signal channel for observing emitted signals.
func (b *BaseStrategy) Signals() <-chan Signal {
	return b.signalCh
}

// EmitStatus records a status snapshot into the strategy's internal ring buffer.
// Non-blocking and safe to call from the run loop at any frequency.
// The At field is set automatically if zero.
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

// Emit publishes a signal to the channel. Non-blocking: drops if buffer is full.
func (b *BaseStrategy) Emit(signal Signal) {
	b.marksMu.Lock()
	b.marks = nil
	b.marksMu.Unlock()

	if b.signalCh == nil {
		return
	}
	select {
	case b.signalCh <- signal:
	default:
	}
}

// Mark records a named checkpoint within the current evaluation cycle.
func (b *BaseStrategy) Mark(label string) {
	b.marksMu.Lock()
	defer b.marksMu.Unlock()
	b.marks = append(b.marks, Mark{Label: label})
}

// GetName returns the strategy name.
func (b *BaseStrategy) GetName() StrategyName { return b.name }
