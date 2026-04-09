package realtime

import (
	"context"
	"fmt"
	"sync"

	priceFeedTypes "github.com/wisp-trading/sdk/pkg/markets/price_feeds/types"
	"github.com/wisp-trading/sdk/pkg/types/logging"
)

type ingestor struct {
	connector Connector
	store     priceFeedTypes.PriceFeedsStore
	logger    logging.ApplicationLogger
	isActive  bool
	ctx       context.Context
	cancel    context.CancelFunc
	wg        sync.WaitGroup
	mu        sync.RWMutex

	// subscriptions tracks active feed channels
	subscriptions map[string]<-chan priceFeedTypes.PriceFeedUpdate
}

// New creates a new Pyth price feed ingestor.
func New(connector Connector, store priceFeedTypes.PriceFeedsStore, logger logging.ApplicationLogger) Ingestor {
	return &ingestor{
		connector:     connector,
		store:         store,
		logger:        logger,
		subscriptions: make(map[string]<-chan priceFeedTypes.PriceFeedUpdate),
	}
}

// Start begins listening to the connector and persisting updates to the store.
func (i *ingestor) Start(ctx context.Context) error {
	i.mu.Lock()
	defer i.mu.Unlock()

	if i.isActive {
		return nil
	}

	if !i.connector.IsConnected() {
		return fmt.Errorf("connector not connected")
	}

	i.isActive = true
	i.ctx, i.cancel = context.WithCancel(ctx)

	// Start monitoring connector errors
	i.wg.Add(1)
	go i.monitorConnectorErrors()

	// Start processing price updates
	i.wg.Add(1)
	go i.processPriceUpdates()

	i.logger.Infof("Pyth price feed ingestor started")
	return nil
}

// Stop ceases consuming and closes all subscriber channels.
func (i *ingestor) Stop() error {
	i.mu.Lock()
	defer i.mu.Unlock()

	if !i.isActive {
		return nil
	}

	i.isActive = false

	if i.cancel != nil {
		i.cancel()
	}

	// Wait for all goroutines to finish
	i.wg.Wait()

	i.logger.Infof("Pyth price feed ingestor stopped")
	return nil
}

// monitorConnectorErrors logs any errors from the connector.
func (i *ingestor) monitorConnectorErrors() {
	defer i.wg.Done()

	errCh := i.connector.ErrorChannel()
	for {
		select {
		case <-i.ctx.Done():
			return
		case err := <-errCh:
			if err != nil {
				i.logger.Errorf("Pyth connector error: %v", err)
			}
		}
	}
}

// processPriceUpdates listens for price updates and persists them.
func (i *ingestor) processPriceUpdates() {
	defer i.wg.Done()

	// For now, subscribe to a generic "pyth" feed channel.
	// In production, you'd subscribe to specific feeds based on watchlist.
	updateCh := i.connector.Subscribe("pyth")

	for {
		select {
		case <-i.ctx.Done():
			return
		case update := <-updateCh:
			snap := priceFeedTypes.PriceSnapshot{
				FeedID:    update.FeedID,
				Price:     update.Price,
				Timestamp: update.Timestamp,
			}
			if err := i.store.RecordPrice(snap); err != nil {
				i.logger.Errorf("failed to record price for feed %s: %v", update.FeedID, err)
			}
		}
	}
}
