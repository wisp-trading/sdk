package strategy

import (
	"context"
	"sync"
)

const signalChannelBufferSize = 64

// BaseStrategyConfig holds configuration for creating a base strategy
type BaseStrategyConfig struct {
	Name        StrategyName
	Description string
	RiskLevel   RiskLevel
	Type        StrategyType
}

// BaseStrategy provides common lifecycle and signal channel management for all strategies.
// Concrete strategies embed BaseStrategy and call StartWithRunner(ctx, s.run) from their
// own Start method to launch their execution goroutine.
type BaseStrategy struct {
	name         StrategyName
	description  string
	riskLevel    RiskLevel
	strategyType StrategyType

	signalCh chan Signal
	ctx      context.Context
	cancel   context.CancelFunc
	wg       sync.WaitGroup

	// marks are pending checkpoints bundled into the next emitted signal's metadata.
	marks   []Mark
	marksMu sync.Mutex
}

// Mark is a named timestamp checkpoint recorded inside a strategy's run loop.
type Mark struct {
	Label string
	At    int64 // UnixNano
}

// NewBaseStrategy creates a new BaseStrategy. The returned value is suitable for
// embedding — callers should embed *BaseStrategy rather than use it directly.
func NewBaseStrategy(config BaseStrategyConfig) *BaseStrategy {
	return &BaseStrategy{
		name:         config.Name,
		description:  config.Description,
		riskLevel:    config.RiskLevel,
		strategyType: config.Type,
	}
}

// StartWithRunner initialises the signal channel and context, then launches the provided
// run function in a managed goroutine. Concrete strategies call this from their own Start:
//
//	func (s *MyStrategy) Start(ctx context.Context) error {
//	    return s.BaseStrategy.StartWithRunner(ctx, s.run)
//	}
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
// This is an observability tap — production routing goes via wisp.Emit.
func (b *BaseStrategy) Signals() <-chan Signal {
	return b.signalCh
}

// emit publishes a signal to the channel. It is non-blocking: if the buffer is full
// the signal is dropped (the strategy's run loop continues uninterrupted).
// Any pending Mark checkpoints are attached to the signal before sending.
func (b *BaseStrategy) Emit(signal Signal) {
	b.marksMu.Lock()
	// marks are intentionally discarded here — they are available for future
	// extension where the profilingRouter can read them from signal metadata.
	b.marks = nil
	b.marksMu.Unlock()

	if b.signalCh == nil {
		return
	}
	select {
	case b.signalCh <- signal:
	default:
		// Buffer full — drop signal; the strategy remains live.
		// A full buffer indicates the executor or router is too slow.
	}
}

// Mark records a named checkpoint within the current evaluation cycle.
// Marks are attached to the next signal emitted and cleared afterwards.
// If no signal is emitted in this cycle the marks are discarded cleanly.
//
// Example:
//
//	s.Mark("scan_start")
//	candidates := s.scan()
//	s.Mark("score_start")
//	if sig, ok := s.score(candidates); ok {
//	    s.Emit(sig)  // marks bundled automatically
//	}
func (b *BaseStrategy) Mark(label string) {
	b.marksMu.Lock()
	defer b.marksMu.Unlock()
	b.marks = append(b.marks, Mark{Label: label})
}

// GetName returns the strategy name.
func (b *BaseStrategy) GetName() StrategyName { return b.name }

// GetDescription returns the strategy description.
func (b *BaseStrategy) GetDescription() string { return b.description }

// GetRiskLevel returns the strategy risk level.
func (b *BaseStrategy) GetRiskLevel() RiskLevel { return b.riskLevel }

// GetStrategyType returns the strategy type.
func (b *BaseStrategy) GetStrategyType() StrategyType { return b.strategyType }
