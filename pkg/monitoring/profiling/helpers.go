package profiling

import (
	"context"

	"github.com/backtesting-org/kronos-sdk/pkg/types/profiling"
)

// contextKey is a private type for storing profiling context in context.Context
type contextKey struct{}

var key = contextKey{}

// WithContext embeds a profiling.Context into a Go context.Context
// This allows profiling data to flow through the execution chain
func WithContext(ctx context.Context, profilingCtx profiling.Context) context.Context {
	return context.WithValue(ctx, key, profilingCtx)
}

// FromContext extracts the profiling.Context from a Go context.Context
// Returns nil if no profiling context is present (profiling disabled)
// This is safe to call even when context wrapping is deep (telemetry, timeout, etc.)
func FromContext(ctx context.Context) profiling.Context {
	if v := ctx.Value(key); v != nil {
		if profilingCtx, ok := v.(profiling.Context); ok {
			return profilingCtx
		}
	}
	return nil
}
