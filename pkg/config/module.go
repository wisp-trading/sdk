package config

import (
	"github.com/wisp-trading/wisp/pkg/config/settings"
	"github.com/wisp-trading/wisp/pkg/config/settings/connectors"
	"github.com/wisp-trading/wisp/pkg/config/startup"
	"github.com/wisp-trading/wisp/pkg/config/strategy"
	"go.uber.org/fx"
)

var Module = fx.Module("config",
	fx.Options(
		settings.Module,
		connectors.Module,
		startup.Module,
	),
	fx.Provide(
		strategy.NewStrategyConfigService,
	),
)
