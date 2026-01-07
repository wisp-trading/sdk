package lifecycle

import (
	"context"

	"github.com/backtesting-org/kronos-sdk/pkg/types/strategy"
)

// State represents the current state of the SDK
type State int

const (
	StateCreated State = iota
	StateStarting
	StateReady
	StateStopping
	StateStopped
)

func (s State) String() string {
	switch s {
	case StateCreated:
		return "Created"
	case StateStarting:
		return "Starting"
	case StateReady:
		return "Ready"
	case StateStopping:
		return "Stopping"
	case StateStopped:
		return "Stopped"
	default:
		return "Unknown"
	}
}

// Controller controls the lifecycle of SDK internal components.
// This is used by the application layer (orchestrator), never exposed to strategy authors.
// It operates at the infrastructure level, below the Kronos/KronosExecutor APIs.
type Controller interface {
	// Start starts the SDK and all its components
	Start(ctx context.Context, name strategy.StrategyName) error

	// Stop gracefully shuts down the SDK
	Stop(ctx context.Context) error

	// WaitUntilReady blocks until the SDK is ready or context is cancelled
	WaitUntilReady(ctx context.Context) error

	// State returns the current lifecycle state
	State() State

	// IsReady returns true if the SDK is ready
	IsReady() bool
}
