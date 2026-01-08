package config

import (
	"github.com/backtesting-org/kronos-sdk/pkg/config/settings"
	"github.com/backtesting-org/kronos-sdk/pkg/config/settings/connectors"
	"github.com/backtesting-org/kronos-sdk/pkg/config/startup"
	"github.com/backtesting-org/kronos-sdk/pkg/config/strategy"
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
