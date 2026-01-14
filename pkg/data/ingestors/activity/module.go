package activity

import (
	"github.com/backtesting-org/kronos-sdk/pkg/types/data/ingestors"
	"go.uber.org/fx"
)

// Module provides market data ingestion components.
var Module = fx.Options(
	fx.Provide(
		NewCoordinator,
		newDataUpdateNotifier,
	),
)

func newDataUpdateNotifier() ingestors.DataUpdateNotifier {
	return NewDataUpdateNotifier(100)
}
