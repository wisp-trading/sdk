package batch

import (
	"fmt"
	"sync"
	"time"

	batchTypes "github.com/wisp-trading/sdk/pkg/markets/base/types/ingestors/batch"
	"github.com/wisp-trading/sdk/pkg/markets/prediction/types"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/logging"
	"github.com/wisp-trading/sdk/pkg/types/temporal"
)

// predictionBatchIngestor polls REST endpoints for prediction market data on a ticker.
type predictionBatchIngestor struct {
	conn         interface{}
	exchangeName connector.ExchangeName
	logger       logging.ApplicationLogger
	timeProvider temporal.TimeProvider
	extensions   []types.PredictionCollectionExtension

	ticker   temporal.Ticker
	stopChan chan struct{}
	isActive bool
	mu       sync.RWMutex
}

func NewPredictionBatchIngestor(
	conn interface{},
	exchangeName connector.ExchangeName,
	logger logging.ApplicationLogger,
	timeProvider temporal.TimeProvider,
	extensions ...types.PredictionCollectionExtension,
) batchTypes.BatchIngestor {
	return &predictionBatchIngestor{
		conn:         conn,
		exchangeName: exchangeName,
		logger:       logger,
		timeProvider: timeProvider,
		extensions:   extensions,
		stopChan:     make(chan struct{}),
	}
}

func (bi *predictionBatchIngestor) Start(interval time.Duration) error {
	bi.mu.Lock()
	defer bi.mu.Unlock()

	if bi.isActive {
		return fmt.Errorf("prediction batch ingestor for %s already active", bi.exchangeName)
	}

	bi.ticker = bi.timeProvider.NewTicker(interval)
	bi.isActive = true

	go bi.collectLoop()

	bi.logger.Info("Started prediction batch ingestion for %s with %v interval", bi.exchangeName, interval)
	return nil
}

func (bi *predictionBatchIngestor) Stop() error {
	bi.mu.Lock()
	defer bi.mu.Unlock()

	if !bi.isActive {
		return nil
	}

	if bi.ticker != nil {
		bi.ticker.Stop()
	}

	close(bi.stopChan)
	bi.isActive = false

	bi.logger.Info("Stopped prediction batch ingestion for %s", bi.exchangeName)
	return nil
}

func (bi *predictionBatchIngestor) IsActive() bool {
	bi.mu.RLock()
	defer bi.mu.RUnlock()
	return bi.isActive
}

func (bi *predictionBatchIngestor) GetMarketType() connector.MarketType {
	return connector.MarketTypePrediction
}

func (bi *predictionBatchIngestor) CollectNow() {
	bi.logger.Debug("Starting prediction batch collection for %s", bi.exchangeName)

	for _, ext := range bi.extensions {
		ext.Collect(bi.conn, bi.exchangeName)
	}

	bi.logger.Debug("Completed prediction batch collection for %s", bi.exchangeName)
}

func (bi *predictionBatchIngestor) collectLoop() {
	// Run immediately on start, then on each tick
	bi.CollectNow()

	for {
		select {
		case <-bi.ticker.C():
			bi.CollectNow()
		case <-bi.stopChan:
			return
		}
	}
}

// Compile-time check
var _ batchTypes.BatchIngestor = (*predictionBatchIngestor)(nil)
