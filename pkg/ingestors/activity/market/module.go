package market

import (
	"github.com/backtesting-org/kronos-sdk/pkg/ingestors/activity/market/batch"
	"github.com/backtesting-org/kronos-sdk/pkg/ingestors/activity/market/realtime"
	"go.uber.org/fx"
)

// Module provides market data ingestion components.
var Module = fx.Options(
	fx.Provide(
		realtime.NewIngestor,
		batch.NewBatchIngestor,
		NewCoordinator,
	),
)
