package market

import (
	"github.com/backtesting-org/kronos-sdk/pkg/data/stores/market/perp"
	"github.com/backtesting-org/kronos-sdk/pkg/data/stores/market/spot"
	"go.uber.org/fx"
)

var Module = fx.Options(
	perp.Module,
	spot.Module,
)
