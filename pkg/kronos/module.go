package kronos

import (
	"go.uber.org/fx"

	"github.com/backtesting-org/kronos-sdk/pkg/kronos/analytics"
	"github.com/backtesting-org/kronos-sdk/pkg/kronos/indicators"
	"github.com/backtesting-org/kronos-sdk/pkg/kronos/market"
	"github.com/backtesting-org/kronos-sdk/pkg/kronos/signal"
	"github.com/backtesting-org/kronos-sdk/pkg/kronos/trade"
)

// Module provides the Kronos SDK with all its services wired up via fx DI.
// This should be included in the fx.Options for your application.
var Module = fx.Module("kronos",
	// Provide all the internal services
	fx.Provide(
		indicators.NewIndicatorService,
		market.NewMarketService,
		analytics.NewAnalyticsService,
		signal.NewService,
		trade.NewTradeService,
	),

	// Provide the main Kronos context
	fx.Provide(NewKronos),

	// Provide the executor (only used by orchestrator)
	fx.Provide(NewKronosExecutor),
)
