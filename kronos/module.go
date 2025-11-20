package kronos

import (
	"github.com/backtesting-org/kronos-sdk/kronos/trade"
	packages "github.com/backtesting-org/kronos-sdk/pkg"
	"go.uber.org/fx"
)

// Module provides the Kronos SDK with all its services wired up via fx DI.
var Module = fx.Module("kronos",
	// Include all pkg modules
	packages.Module,

	// Provide the trade service
	fx.Provide(trade.NewTradeService),

	// Provide the main Kronos context
	fx.Provide(NewKronos),

	// Provide the executor (only used by orchestrator)
	fx.Provide(NewKronosExecutor),
)
