package market

import (
	"github.com/backtesting-org/kronos-sdk/pkg/data/ingestors/market/perp"
	"github.com/backtesting-org/kronos-sdk/pkg/data/ingestors/market/spot"
	"github.com/backtesting-org/kronos-sdk/pkg/types/data/ingestors"
	"go.uber.org/fx"
)

var Module = fx.Module("market_ingestor",
	fx.Options(
		spot.Module,
		perp.Module,

		fx.Provide(
			NewCoordinator,
			newDataUpdateNotifier,
		),
	),
)

func newDataUpdateNotifier() ingestors.DataUpdateNotifier {
	return NewDataUpdateNotifier(100)
}
