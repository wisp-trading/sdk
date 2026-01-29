package market

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/wisp-trading/sdk/pkg/types/data/ingestors"
	"github.com/wisp-trading/sdk/pkg/types/data/ingestors/batch"
	"github.com/wisp-trading/sdk/pkg/types/data/ingestors/realtime"
	"github.com/wisp-trading/sdk/pkg/types/logging"
	"github.com/wisp-trading/sdk/pkg/types/temporal"
)

type coordinator struct {
	batchFactories    []batch.BatchIngestorFactory
	realtimeFactories []realtime.RealtimeIngestorFactory
	logger            logging.ApplicationLogger
	timeProvider      temporal.TimeProvider

	// State management
	batchIngestors    []batch.BatchIngestor
	realtimeIngestors []realtime.RealtimeIngestor
	isRunning         bool
	mu                sync.RWMutex
}

func NewCoordinator(
	batchFactories []batch.BatchIngestorFactory,
	realtimeFactories []realtime.RealtimeIngestorFactory,
	timeProvider temporal.TimeProvider,
	logger logging.ApplicationLogger,
) ingestors.MarketDataCoordinator {
	return &coordinator{
		batchFactories:    batchFactories,
		realtimeFactories: realtimeFactories,
		timeProvider:      timeProvider,
		logger:            logger,
		isRunning:         false,
	}
}

func (dic *coordinator) IsRunning() bool {
	dic.mu.RLock()
	defer dic.mu.RUnlock()
	return dic.isRunning
}

func (dic *coordinator) StartDataCollection(ctx context.Context) error {
	dic.mu.Lock()
	defer dic.mu.Unlock()

	if dic.isRunning {
		dic.logger.Debug("Data collection already running, skipping start")
		return nil
	}

	// Create ingestors from factories at runtime based on registered connectors
	dic.logger.Info("Creating ingestors from factories...")

	dic.batchIngestors = nil
	for _, factory := range dic.batchFactories {
		ingestors := factory.CreateIngestors()
		dic.batchIngestors = append(dic.batchIngestors, ingestors...)
	}

	dic.realtimeIngestors = nil
	for _, factory := range dic.realtimeFactories {
		ingestors := factory.CreateIngestors()
		dic.realtimeIngestors = append(dic.realtimeIngestors, ingestors...)
	}

	dic.logger.Info("Created %d batch ingestors and %d realtime ingestors",
		len(dic.batchIngestors), len(dic.realtimeIngestors))

	// Backfill current data before starting streams (all market types)
	dic.logger.Info("Collecting initial market data snapshot across all market types...")
	for _, batchIngestor := range dic.batchIngestors {
		marketType := batchIngestor.GetMarketType()
		dic.logger.Info("Collecting %s market data snapshot...", marketType)
		batchIngestor.CollectNow()
	}
	dic.logger.Info("Initial market data snapshot complete")

	// Start all real-time ingestors
	for _, realtimeIngestor := range dic.realtimeIngestors {
		marketType := realtimeIngestor.GetMarketType()
		if err := realtimeIngestor.Start(ctx); err != nil {
			dic.logger.Error("Failed to start %s real-time ingestion: %v", marketType, err)
			// Continue with batch fallback
		} else {
			dic.logger.Info("Started %s real-time ingestion", marketType)
		}
	}

	// Start all batch ingestors as backup
	for _, batchIngestor := range dic.batchIngestors {
		marketType := batchIngestor.GetMarketType()
		if err := batchIngestor.Start(30 * time.Second); err != nil {
			dic.logger.Error("Failed to start %s batch ingestion: %v", marketType, err)
			return err
		}
		dic.logger.Info("Started %s batch ingestion", marketType)
	}

	dic.isRunning = true
	dic.logger.Info("Started hybrid data ingestion across all market types (WebSocket + REST backup)")
	return nil
}

func (dic *coordinator) StopDataCollection() error {
	dic.mu.Lock()
	defer dic.mu.Unlock()

	if !dic.isRunning {
		dic.logger.Debug("Data collection not running, skipping stop")
		return nil
	}

	var errs []error

	// Stop all realtime ingestors
	for _, realtimeIngestor := range dic.realtimeIngestors {
		marketType := realtimeIngestor.GetMarketType()
		if err := realtimeIngestor.Stop(); err != nil {
			errs = append(errs, fmt.Errorf("%s realtime stop error: %w", marketType, err))
		}
	}

	// Stop all batch ingestors
	for _, batchIngestor := range dic.batchIngestors {
		marketType := batchIngestor.GetMarketType()
		if err := batchIngestor.Stop(); err != nil {
			errs = append(errs, fmt.Errorf("%s batch stop error: %w", marketType, err))
		}
	}

	dic.isRunning = false
	dic.logger.Info("Stopped all data ingestion across all market types")

	if len(errs) > 0 {
		return fmt.Errorf("errors stopping ingestion: %v", errs)
	}

	return nil
}

func (dic *coordinator) GetStatus() map[string]interface{} {
	dic.mu.RLock()
	defer dic.mu.RUnlock()

	status := map[string]interface{}{
		"coordinator_running": dic.isRunning,
		"market_types":        make(map[string]interface{}),
	}

	marketTypes := status["market_types"].(map[string]interface{})

	// Status for each realtime ingestor
	for _, realtimeIngestor := range dic.realtimeIngestors {
		marketType := string(realtimeIngestor.GetMarketType())
		if _, exists := marketTypes[marketType]; !exists {
			marketTypes[marketType] = make(map[string]interface{})
		}
		mt := marketTypes[marketType].(map[string]interface{})
		mt["realtime"] = realtimeIngestor.IsActive()

		// Add connection details if active
		if realtimeIngestor.IsActive() {
			connections := realtimeIngestor.GetActiveConnections()
			connectionStatus := make(map[string]bool)
			for name := range connections {
				connectionStatus[string(name)] = true
			}
			mt["connections"] = connectionStatus
		}
	}

	// Status for each batch ingestor
	for _, batchIngestor := range dic.batchIngestors {
		marketType := string(batchIngestor.GetMarketType())
		if _, exists := marketTypes[marketType]; !exists {
			marketTypes[marketType] = make(map[string]interface{})
		}
		mt := marketTypes[marketType].(map[string]interface{})
		mt["batch"] = batchIngestor.IsActive()
	}

	return status
}

func (dic *coordinator) ForceCollectNow() {
	dic.logger.Info("Forcing immediate data collection across all market types")
	for _, batchIngestor := range dic.batchIngestors {
		marketType := batchIngestor.GetMarketType()
		dic.logger.Info("Forcing %s market data collection", marketType)
		batchIngestor.CollectNow()
	}
}

func (dic *coordinator) RestartRealtime(ctx context.Context) error {
	dic.logger.Info("Restarting realtime ingestion across all market types")

	var errs []error

	// Stop all realtime ingestors
	for _, realtimeIngestor := range dic.realtimeIngestors {
		marketType := realtimeIngestor.GetMarketType()
		if err := realtimeIngestor.Stop(); err != nil {
			dic.logger.Error("Error stopping %s realtime ingestor: %v", marketType, err)
		}
	}

	dic.timeProvider.Sleep(2 * time.Second) // Brief pause for cleanup

	// Start all realtime ingestors
	for _, realtimeIngestor := range dic.realtimeIngestors {
		marketType := realtimeIngestor.GetMarketType()
		if err := realtimeIngestor.Start(ctx); err != nil {
			dic.logger.Error("Error restarting %s realtime ingestor: %v", marketType, err)
			errs = append(errs, fmt.Errorf("%s restart error: %w", marketType, err))
		} else {
			dic.logger.Info("Successfully restarted %s realtime ingestion", marketType)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors restarting realtime ingestion: %v", errs)
	}

	dic.logger.Info("Successfully restarted realtime ingestion across all market types")
	return nil
}
