package strategy

import "context"

type contextKey struct{}

// WithStrategyName adds strategy name to context
func WithStrategyName(ctx context.Context, name StrategyName) context.Context {
	return context.WithValue(ctx, contextKey{}, name)
}

// FromContext extracts strategy name from context
func FromContext(ctx context.Context) (StrategyName, bool) {
	name, ok := ctx.Value(contextKey{}).(StrategyName)
	return name, ok
}
