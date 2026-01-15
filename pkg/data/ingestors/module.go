package ingestors

import (
	"github.com/backtesting-org/kronos-sdk/pkg/data/ingestors/market"
	"github.com/backtesting-org/kronos-sdk/pkg/data/ingestors/position"
	"go.uber.org/fx"
)

var Module = fx.Options(
	position.Module,
	market.Module,
)
