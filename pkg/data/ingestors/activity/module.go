package activity

import (
	"github.com/backtesting-org/kronos-sdk/pkg/data/ingestors/activity/market"
	"github.com/backtesting-org/kronos-sdk/pkg/types/data/ingestors"
	"go.uber.org/fx"
)

// Module provides market data ingestion components.
var Module = fx.Options(
	market.Module,
	fx.Provide(
		NewCoordinator,
		newDataUpdateNotifier,
	),
)

func newDataUpdateNotifier() ingestors.DataUpdateNotifier {
	return NewDataUpdateNotifier(100)
}
