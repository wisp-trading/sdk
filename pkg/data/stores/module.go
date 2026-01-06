package stores

import (
	"github.com/backtesting-org/kronos-sdk/pkg/data/stores/activity"
	"github.com/backtesting-org/kronos-sdk/pkg/data/stores/market"
	"go.uber.org/fx"
)

var Module = fx.Options(
	activity.Module,
	market.Module,
)
