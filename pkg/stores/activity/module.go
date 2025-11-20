package activity

import (
	"github.com/backtesting-org/kronos-sdk/pkg/stores/activity/position"
	"github.com/backtesting-org/kronos-sdk/pkg/stores/activity/trade"
	"go.uber.org/fx"
)

var Module = fx.Options(
	position.Module,
	trade.Module,
)
