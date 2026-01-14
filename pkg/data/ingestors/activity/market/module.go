package market

import (
	"github.com/backtesting-org/kronos-sdk/pkg/data/ingestors/activity/market/perp"
	"github.com/backtesting-org/kronos-sdk/pkg/data/ingestors/activity/market/spot"
	"go.uber.org/fx"
)

var Module = fx.Module("market_ingestor",
	fx.Options(
		spot.Module,
		perp.Module,
	),
)
