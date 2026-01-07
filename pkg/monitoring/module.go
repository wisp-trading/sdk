package monitoring

import (
	"github.com/backtesting-org/kronos-sdk/pkg/monitoring/health"
	"github.com/backtesting-org/kronos-sdk/pkg/monitoring/profiling"
	"go.uber.org/fx"
)

var Module = fx.Module("monitoring",
	health.Module,
	profiling.Module,

	fx.Provide(NewServer),
)
