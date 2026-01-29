package monitoring

import (
	"github.com/wisp-trading/wisp/pkg/monitoring/health"
	"github.com/wisp-trading/wisp/pkg/monitoring/profiling"
	"go.uber.org/fx"
)

var Module = fx.Module("monitoring",
	health.Module,
	profiling.Module,

	fx.Provide(
		NewServer,
		NewViewRegistry,
	),
)
