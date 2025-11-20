package ingestors

import (
	"github.com/backtesting-org/kronos-sdk/pkg/ingestors/activity/market"
	"github.com/backtesting-org/kronos-sdk/pkg/ingestors/discovery"
	"github.com/backtesting-org/kronos-sdk/pkg/ingestors/position"
	"go.uber.org/fx"
)

var Module = fx.Options(
	discovery.Module,
	position.Module,
	market.Module,
)
