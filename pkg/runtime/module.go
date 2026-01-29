package runtime

import (
	"github.com/wisp-trading/sdk/pkg/runtime/time"
	"go.uber.org/fx"
)

var Module = fx.Options(
	time.Module,
	fx.Provide(NewRuntime),
)
