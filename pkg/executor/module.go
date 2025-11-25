// Package executor provides execution functionality for trading signals.
// It includes a default executor implementation with support for custom hooks,
// allowing users to extend execution behavior through plugins.
package executor

import (
	"go.uber.org/fx"
)

// Module provides the executor functionality to the application
var Module = fx.Module("executor",
	fx.Provide(
		NewExecutor,
	),
)
