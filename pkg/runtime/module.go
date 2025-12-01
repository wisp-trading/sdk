package runtime

import (
	"github.com/backtesting-org/kronos-sdk/pkg/runtime/time"
	"go.uber.org/fx"
)

var Module = fx.Options(
	time.Module,
	fx.Provide(NewRuntime),
)
