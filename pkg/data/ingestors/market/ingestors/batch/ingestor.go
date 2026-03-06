package batch

import (
	"fmt"
	"sync"
	"time"

	"github.com/wisp-trading/sdk/pkg/markets/base/types"
	"github.com/wisp-trading/sdk/pkg/markets/base/types/ingestors/batch"
	marketTypes "github.com/wisp-trading/sdk/pkg/markets/base/types/stores/market"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/logging"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	"github.com/wisp-trading/sdk/pkg/types/temporal"
)

// batchIngestor is a generic base implementation for REST batch data collection
type batchIngestor struct {
	conn            connector.Connector
	exchangeName    connector.ExchangeName
	marketType      connector.MarketType
	marketWatchlist types.MarketWatchlist
	store           marketTypes.MarketStore
	logger          logging.ApplicationLogger
	timeProvider    temporal.TimeProvider

	// Kline configuration
	klineLimits map[string]int

	// Scheduling
	ticker   temporal.Ticker
	stopChan chan struct{}
	isActive bool
	mu       sync.RWMutex

	// Extension point for market-specific logic
	extensions []batch.CollectionExtension
}

func NewBatchIngestor(
	conn connector.Connector,
	exchangeName connector.ExchangeName,
	marketType connector.MarketType,
	marketWatchlist types.MarketWatchlist,
	store marketTypes.MarketStore,
	timeProvider temporal.TimeProvider,
	logger logging.ApplicationLogger,
	extensions ...batch.CollectionExtension,
) batch.BatchIngestor {
	return &batchIngestor{
		conn:            conn,
		exchangeName:    exchangeName,
		marketType:      marketType,
		marketWatchlist: marketWatchlist,
		store:           store,
		timeProvider:    timeProvider,
		logger:          logger,
		stopChan:        make(chan struct{}),
		extensions:      extensions,
	}
}

func (bi *batchIngestor) Start(interval time.Duration) error {
	bi.mu.Lock()
	defer bi.mu.Unlock()

	if bi.isActive {
		return fmt.Errorf("batch ingestor for %s already active", bi.exchangeName)
	}

	bi.ticker = bi.timeProvider.NewTicker(interval)
	bi.isActive = true

	go bi.collectLoop()

	bi.logger.Info("Started %s batch ingestion for %s with %v interval", bi.marketType, bi.exchangeName, interval)
	return nil
}

func (bi *batchIngestor) CollectNow() {
	bi.logger.Debug("Starting %s market data collection for %s", bi.marketType, bi.exchangeName)

	requirements := bi.marketWatchlist.GetRequiredPairs(bi.exchangeName)
	if len(requirements) == 0 {
		bi.logger.Debug("No assets required for collection")
		return
	}

	requiredPairs := make([]portfolio.Pair, len(requirements))
	for i, req := range requirements {
		requiredPairs[i] = req
	}

	// Run any market-specific collection extensions
	for _, ext := range bi.extensions {
		ext.Collect(bi.conn, bi.exchangeName, requiredPairs)

	}

	bi.logger.Debug("Completed %s market data collection for %s", bi.marketType, bi.exchangeName)
}

func (bi *batchIngestor) collectLoop() {
	// Run initial collection immediately
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

func (bi *batchIngestor) Stop() error {
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

	bi.logger.Info("Stopped %s batch ingestion for %s", bi.marketType, bi.exchangeName)
	return nil
}

func (bi *batchIngestor) IsActive() bool {
	bi.mu.RLock()
	defer bi.mu.RUnlock()
	return bi.isActive
}

func (bi *batchIngestor) GetMarketType() connector.MarketType {
	return bi.marketType
}

var _ batch.BatchIngestor = (*batchIngestor)(nil)
