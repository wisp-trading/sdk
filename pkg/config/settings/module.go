package settings

import "go.uber.org/fx"

var Module = fx.Module("settings",
	fx.Provide(
		// Provide default ConfigOptions - callers can override with fx.Replace
		func() ConfigOptions {
			return ConfigOptions{}
		},
		NewConfiguration,
	),
)
