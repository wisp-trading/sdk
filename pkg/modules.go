package packages

import (
	"github.com/wisp-trading/sdk/pkg/activity"
	"github.com/wisp-trading/sdk/pkg/adapters"
	"github.com/wisp-trading/sdk/pkg/analytics"
	"github.com/wisp-trading/sdk/pkg/config"
	"github.com/wisp-trading/sdk/pkg/executor"
	"github.com/wisp-trading/sdk/pkg/lifecycle"
	"github.com/wisp-trading/sdk/pkg/markets/options"
	"github.com/wisp-trading/sdk/pkg/markets/perp"
	"github.com/wisp-trading/sdk/pkg/markets/prediction"
	"github.com/wisp-trading/sdk/pkg/markets/price_feeds"
	"github.com/wisp-trading/sdk/pkg/markets/spot"
	"github.com/wisp-trading/sdk/pkg/monitoring"
	"github.com/wisp-trading/sdk/pkg/registry"
	"github.com/wisp-trading/sdk/pkg/runtime"
	"github.com/wisp-trading/sdk/pkg/signal"
	"go.uber.org/fx"
)

// Module is the full SDK fx graph.
// Strategy packaging is standalone only (no plugin/.so loader in the graph).
// pkg/plugin remains available for optional hook plugins if wired explicitly.
var Module = fx.Options(
	activity.Module,
	adapters.Module,
	analytics.Module,
	config.Module,
	monitoring.Module,
	lifecycle.Module,
	registry.Module,
	runtime.Module,
	signal.Module,
	executor.Module,
	prediction.Module,
	perp.Module,
	spot.Module,
	options.Module,
	price_feeds.Module,
)
