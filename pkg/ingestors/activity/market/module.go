package market

import (
	"context"

	"github.com/backtesting-org/kronos-sdk/pkg/ingestors/activity/market/batch"
	"github.com/backtesting-org/kronos-sdk/pkg/ingestors/activity/market/realtime"
	"github.com/backtesting-org/kronos-sdk/pkg/types/logging"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		realtime.NewIngestor,
		batch.NewBatchIngestor,
		NewCoordinator,
	),
	fx.Invoke(func(
		lc fx.Lifecycle,
		coordinator *Coordinator,
		logger logging.ApplicationLogger,
	) {
		lc.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				logger.Info("Starting asset data ingestion coordinator")

				// Use background context for long-lived streams
				backgroundCtx := context.Background()
				if err := coordinator.StartDataCollection(backgroundCtx); err != nil {
					logger.Error("Failed to start data collection: %v", err)
					return err
				}

				status := coordinator.GetStatus()
				logger.Info("Asset data ingestion started successfully: %+v", status)
				return nil
			},
			OnStop: func(ctx context.Context) error {
				logger.Info("Stopping asset data ingestion coordinator")

				if err := coordinator.StopDataCollection(); err != nil {
					logger.Error("Error stopping data collection: %v", err)
					return err
				}

				logger.Info("Asset data ingestion stopped successfully")
				return nil
			},
		})
	}),
)
