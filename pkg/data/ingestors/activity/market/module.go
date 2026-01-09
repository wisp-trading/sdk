package market

import (
	"github.com/backtesting-org/kronos-sdk/pkg/data/ingestors/activity/market/batch"
	"github.com/backtesting-org/kronos-sdk/pkg/data/ingestors/activity/market/realtime"
	"github.com/backtesting-org/kronos-sdk/pkg/types/data/ingestors"
	"go.uber.org/fx"
)

// Module provides market data ingestion components.
var Module = fx.Options(
	fx.Provide(
		realtime.NewIngestor,
		batch.NewBatchIngestor,
		NewCoordinator,
		newDataUpdateNotifier,
	),
)

func newDataUpdateNotifier() ingestors.DataUpdateNotifier {
	return NewDataUpdateNotifier(100)
}
