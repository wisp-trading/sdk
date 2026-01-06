package market

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/backtesting-org/kronos-sdk/pkg/types/data/ingestors"
	"github.com/backtesting-org/kronos-sdk/pkg/types/logging"
	"github.com/backtesting-org/kronos-sdk/pkg/types/temporal"
)

type coordinator struct {
	realtimeIngestor ingestors.RealtimeIngestor
	batchIngestor    ingestors.BatchIngestor
	logger           logging.ApplicationLogger
	timeProvider     temporal.TimeProvider

	// State management
	isRunning bool
	mu        sync.RWMutex
}

func NewCoordinator(
	realtimeIngestor ingestors.RealtimeIngestor,
	batchIngestor ingestors.BatchIngestor,
	timeProvider temporal.TimeProvider,
	logger logging.ApplicationLogger,
) ingestors.MarketDataCoordinator {
	return &coordinator{
		realtimeIngestor: realtimeIngestor,
		batchIngestor:    batchIngestor,
		timeProvider:     timeProvider,
		logger:           logger,
		isRunning:        false,
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

	// Backfill current data before starting streams
	dic.logger.Info("Collecting initial market data snapshot...")
	dic.batchIngestor.CollectNow()
	dic.timeProvider.Sleep(2 * time.Second) // Give it time to collect

	// Start real-time ingestion
	if err := dic.realtimeIngestor.Start(ctx); err != nil {
		dic.logger.Error("Failed to start real-time ingestion: %v", err)
		// Continue with batch fallback
	}

	// Start batch ingestion as backup
	if err := dic.batchIngestor.Start(30 * time.Second); err != nil {
		dic.logger.Error("Failed to start batch ingestion: %v", err)
		return err
	}

	dic.isRunning = true
	dic.logger.Info("Started hybrid data ingestion (WebSocket + REST backup)")
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

	if err := dic.realtimeIngestor.Stop(); err != nil {
		errs = append(errs, fmt.Errorf("realtime stop error: %w", err))
	}

	if err := dic.batchIngestor.Stop(); err != nil {
		errs = append(errs, fmt.Errorf("batch stop error: %w", err))
	}

	dic.isRunning = false
	dic.logger.Info("Stopped all data ingestion")

	if len(errs) > 0 {
		return fmt.Errorf("errors stopping ingestion: %v", errs)
	}

	return nil
}

func (dic *coordinator) GetStatus() map[string]interface{} {
	dic.mu.RLock()
	defer dic.mu.RUnlock()

	status := map[string]interface{}{
		"realtime":            dic.realtimeIngestor.IsActive(),
		"batch":               dic.batchIngestor.IsActive(),
		"coordinator_running": dic.isRunning,
	}

	// Add connection details for realtime
	if dic.realtimeIngestor.IsActive() {
		connections := dic.realtimeIngestor.GetActiveConnections()
		connectionStatus := make(map[string]bool)
		for name := range connections {
			connectionStatus[string(name)] = true
		}
		status["connections"] = connectionStatus
	}

	return status
}

func (dic *coordinator) ForceCollectNow() {
	dic.logger.Info("Forcing immediate data collection")
	dic.batchIngestor.CollectNow()
}

func (dic *coordinator) RestartRealtime(ctx context.Context) error {
	dic.logger.Info("Restarting realtime ingestion")

	if err := dic.realtimeIngestor.Stop(); err != nil {
		dic.logger.Error("Error stopping realtime ingestor: %v", err)
	}

	dic.timeProvider.Sleep(2 * time.Second) //  Brief pause for cleanup

	if err := dic.realtimeIngestor.Start(ctx); err != nil {
		dic.logger.Error("Error restarting realtime ingestor: %v", err)
		return err
	}

	dic.logger.Info("Successfully restarted realtime ingestion")
	return nil
}
