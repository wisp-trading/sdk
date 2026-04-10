package executor

import (
	"go.uber.org/fx"

	"github.com/wisp-trading/sdk/pkg/execution/records"
	"github.com/wisp-trading/sdk/pkg/types/execution"
	"github.com/wisp-trading/sdk/pkg/types/logging"
	profileTypes "github.com/wisp-trading/sdk/pkg/types/monitoring/profiling"
)

// routerParams allows optional profiling deps to be injected without failing
// the fx graph when they are absent.
type routerParams struct {
	fx.In

	Executor        execution.Executor
	Logger          logging.ApplicationLogger
	ProfilingStore  profileTypes.ProfilingStore  `optional:"true"`
	AnomalyDetector profileTypes.AnomalyDetector `optional:"true"`
}

// newSignalRouter provides a SignalRouter to the fx graph.
// When a ProfilingStore is wired, signals are wrapped with the profilingRouter
// decorator so metrics are captured automatically — invisible to strategies.
func newSignalRouter(p routerParams) execution.SignalRouter {
	base := NewExecutorRouter(p.Executor, p.Logger)
	if p.ProfilingStore != nil {
		return NewProfilingRouter(base, p.ProfilingStore, p.AnomalyDetector)
	}
	return base
}

// Module provides the executor functionality to the application
var Module = fx.Module("executor",
	fx.Provide(
		NewExecutor,
		newSignalRouter,
		records.NewStore,
	),
)
