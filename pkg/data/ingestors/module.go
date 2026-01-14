package ingestors

import (
	"github.com/backtesting-org/kronos-sdk/pkg/data/ingestors/activity"
	"github.com/backtesting-org/kronos-sdk/pkg/data/ingestors/position"
	"go.uber.org/fx"
)

var Module = fx.Options(
	position.Module,
	activity.Module,
)
