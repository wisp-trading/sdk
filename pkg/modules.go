package pkg

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
	"go.uber.org/fx"
)

// Module is the full SDK fx graph.
// Strategy packaging is standalone only: main + StartStandalone + Wait.
//
// Market shell (clone this to add a market):
//
//	pkg/markets/<domain>/
//	  module.go      — fx Module; Provide facade.New*, lifecycle, stores, ingestors
//	  facade/        — strategy-facing context (wisp.Spot / Perp / …)
//	  types/         — interfaces
//	  signal/        — builders
//	  executor/      — place orders for domain signals
//	  store/         — domain state (+ base/store/extensions)
//	  ingestor/      — batch + realtime
//	  activity/ views/ … as needed
//
// Then add <domain>.Module here and register connectors in the connectors repo.
var Module = fx.Options(
	activity.Module,
	adapters.Module,
	analytics.Module,
	config.Module,
	monitoring.Module,
	lifecycle.Module,
	registry.Module,
	runtime.Module,
	executor.Module,
	prediction.Module,
	perp.Module,
	spot.Module,
	options.Module,
	price_feeds.Module,
)
