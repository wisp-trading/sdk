package packages

import (
	"github.com/wisp-trading/wisp/pkg/activity"
	"github.com/wisp-trading/wisp/pkg/adapters"
	"github.com/wisp-trading/wisp/pkg/analytics"
	"github.com/wisp-trading/wisp/pkg/config"
	"github.com/wisp-trading/wisp/pkg/data/ingestors"
	"github.com/wisp-trading/wisp/pkg/data/stores"
	"github.com/wisp-trading/wisp/pkg/executor"
	"github.com/wisp-trading/wisp/pkg/inference/features"
	"github.com/wisp-trading/wisp/pkg/lifecycle"
	"github.com/wisp-trading/wisp/pkg/monitoring"
	"github.com/wisp-trading/wisp/pkg/plugin"
	"github.com/wisp-trading/wisp/pkg/registry"
	"github.com/wisp-trading/wisp/pkg/runtime"
	"github.com/wisp-trading/wisp/pkg/signal"
	"go.uber.org/fx"
)

var Module = fx.Options(
	activity.Module,
	adapters.Module,
	analytics.Module,
	config.Module,
	monitoring.Module,
	features.Module,
	ingestors.Module,
	lifecycle.Module,
	plugin.Module,
	registry.Module,
	runtime.Module,
	signal.Module,
	stores.Module,
	executor.Module,
)
