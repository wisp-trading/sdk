package ingestor

import (
	"context"
	"fmt"
	"sync"
	"time"

	batchTypes "github.com/wisp-trading/sdk/pkg/markets/base/types/ingestors/batch"
	realtimeTypes "github.com/wisp-trading/sdk/pkg/markets/base/types/ingestors/realtime"
	lifecycleTypes "github.com/wisp-trading/sdk/pkg/types/lifecycle"
	"github.com/wisp-trading/sdk/pkg/types/logging"
)

const defaultBatchInterval = 30 * time.Second

type domainCoordinator struct {
	name            string
	batchFactory    batchTypes.BatchIngestorFactory
	realtimeFactory realtimeTypes.RealtimeIngestorFactory
	logger          logging.ApplicationLogger

	mu             sync.Mutex
	batchIngestors []batchTypes.BatchIngestor
	rtIngestors    []realtimeTypes.RealtimeIngestor
	rtCancel       context.CancelFunc
}

func NewDomainCoordinator(
	name string,
	batchFactory batchTypes.BatchIngestorFactory,
	realtimeFactory realtimeTypes.RealtimeIngestorFactory,
	logger logging.ApplicationLogger,
) lifecycleTypes.DomainLifecycle {
	return &domainCoordinator{
		name:            name,
		batchFactory:    batchFactory,
		realtimeFactory: realtimeFactory,
		logger:          logger,
	}
}

func (c *domainCoordinator) Name() string { return c.name }

func (c *domainCoordinator) Start(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.logger.Info("Starting %s ingestors...", c.name)

	// Batch ingestors — start with a periodic interval, then fire an immediate collection.
	c.batchIngestors = c.batchFactory.CreateIngestors()
	for _, ingestor := range c.batchIngestors {
		if err := ingestor.Start(defaultBatchInterval); err != nil {
			return fmt.Errorf("%s batch ingestor start failed: %w", c.name, err)
		}
		ingestor.CollectNow()
	}

	// Realtime ingestors — run in background goroutines managed by their own context.
	rtCtx, cancel := context.WithCancel(ctx)
	c.rtCancel = cancel
	c.rtIngestors = c.realtimeFactory.CreateIngestors()
	for _, ingestor := range c.rtIngestors {
		if err := ingestor.Start(rtCtx); err != nil {
			return fmt.Errorf("%s realtime ingestor start failed: %w", c.name, err)
		}
	}

	c.logger.Info("%s ingestors started (%d batch, %d realtime)", c.name, len(c.batchIngestors), len(c.rtIngestors))
	return nil
}

func (c *domainCoordinator) Stop() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.logger.Info("Stopping %s ingestors...", c.name)

	if c.rtCancel != nil {
		c.rtCancel()
	}
	for _, ingestor := range c.rtIngestors {
		if err := ingestor.Stop(); err != nil {
			c.logger.Warn("%s realtime ingestor stop error: %v", c.name, err)
		}
	}
	for _, ingestor := range c.batchIngestors {
		if err := ingestor.Stop(); err != nil {
			c.logger.Warn("%s batch ingestor stop error: %v", c.name, err)
		}
	}
	return nil
}

var _ lifecycleTypes.DomainLifecycle = (*domainCoordinator)(nil)
