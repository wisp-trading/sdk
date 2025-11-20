package kronos

import (
	"github.com/backtesting-org/kronos-sdk/kronos/market"
	"github.com/backtesting-org/kronos-sdk/kronos/signal"
	"github.com/backtesting-org/kronos-sdk/kronos/trade"
	"github.com/backtesting-org/kronos-sdk/pkg/analytics"
	"github.com/backtesting-org/kronos-sdk/pkg/analytics/indicators"
	"go.uber.org/fx"
)

// Module provides the Kronos SDK with all its services wired up via fx DI.
var Module = fx.Module("kronos",
	// Provide all the internal services
	fx.Provide(
		indicators.NewIndicators,
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
