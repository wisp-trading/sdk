package packages

import (
	"github.com/wisp-trading/sdk/pkg/activity"
	"github.com/wisp-trading/sdk/pkg/adapters"
	"github.com/wisp-trading/sdk/pkg/analytics"
	"github.com/wisp-trading/sdk/pkg/config"
	"github.com/wisp-trading/sdk/pkg/data"
	"github.com/wisp-trading/sdk/pkg/executor"
	"github.com/wisp-trading/sdk/pkg/lifecycle"
	"github.com/wisp-trading/sdk/pkg/monitoring"
	"github.com/wisp-trading/sdk/pkg/plugin"
	"github.com/wisp-trading/sdk/pkg/registry"
	"github.com/wisp-trading/sdk/pkg/runtime"
	"github.com/wisp-trading/sdk/pkg/signal"
	"go.uber.org/fx"
)

var Module = fx.Options(
	activity.Module,
	adapters.Module,
	analytics.Module,
	config.Module,
	monitoring.Module,
	lifecycle.Module,
	plugin.Module,
	registry.Module,
	runtime.Module,
	signal.Module,
	executor.Module,
	data.Module,
)
