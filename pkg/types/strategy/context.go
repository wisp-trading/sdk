package strategy

import "context"

// StrategyContext is a context that guarantees strategy name is present
type StrategyContext struct {
	context.Context
	strategyName StrategyName
}

// NewStrategyContext creates a context with guaranteed strategy name
func NewStrategyContext(parent context.Context, name StrategyName) StrategyContext {
	return StrategyContext{
		Context:      parent,
		strategyName: name,
	}
}

// StrategyName returns the strategy name (guaranteed to exist)
func (sc *StrategyContext) StrategyName() StrategyName {
	return sc.strategyName
}
