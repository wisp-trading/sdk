package config

import (
	"github.com/wisp-trading/sdk/pkg/config/settings"
	"github.com/wisp-trading/sdk/pkg/config/settings/connectors"
	"github.com/wisp-trading/sdk/pkg/config/startup"
	"github.com/wisp-trading/sdk/pkg/config/strategy"
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
