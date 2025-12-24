package batch

import (
	"fmt"
	"sync"
	"time"

	"github.com/backtesting-org/kronos-sdk/pkg/types/connector"
	"github.com/backtesting-org/kronos-sdk/pkg/types/data/ingestors"
	"github.com/backtesting-org/kronos-sdk/pkg/types/data/stores/market"
	"github.com/backtesting-org/kronos-sdk/pkg/types/health"
	"github.com/backtesting-org/kronos-sdk/pkg/types/logging"
	"github.com/backtesting-org/kronos-sdk/pkg/types/portfolio"
	"github.com/backtesting-org/kronos-sdk/pkg/types/registry"
	"github.com/backtesting-org/kronos-sdk/pkg/types/temporal"
)

type ingestor struct {
	store            market.MarketData
	exchangeRegistry registry.ConnectorRegistry
	assetRegistry    registry.AssetRegistry
	logger           logging.ApplicationLogger
	timeProvider     temporal.TimeProvider
	healthStore      health.CoordinatorHealthStore
	notifier         ingestors.DataUpdateNotifier

	// Scheduling
	ticker   temporal.Ticker
	stopChan chan struct{}
	isActive bool
	mutex    sync.RWMutex
}

func NewBatchIngestor(
	store market.MarketData,
	exchangeRegistry registry.ConnectorRegistry,
	assetRegistry registry.AssetRegistry,
	logger logging.ApplicationLogger,
	timeProvider temporal.TimeProvider,
	healthStore health.CoordinatorHealthStore,
	notifier ingestors.DataUpdateNotifier,
) ingestors.BatchIngestor {
	return &ingestor{
		store:            store,
		exchangeRegistry: exchangeRegistry,
		assetRegistry:    assetRegistry,
		logger:           logger,
		timeProvider:     timeProvider,
		healthStore:      healthStore,
		notifier:         notifier,
		stopChan:         make(chan struct{}),
	}
}

func (bi *ingestor) Start(interval time.Duration) error {
	bi.mutex.Lock()
	defer bi.mutex.Unlock()

	if bi.isActive {
		return fmt.Errorf("batch ingestor already active")
	}

	bi.ticker = bi.timeProvider.NewTicker(interval)
	bi.isActive = true

	go bi.collectLoop()

	bi.logger.Info("Started REST batch ingestion with %v interval", interval)
	return nil
}

func (bi *ingestor) collectLoop() {
	// Run initial collection immediately
	bi.collectMarketData()

	for {
		select {
		case <-bi.ticker.C():
			bi.collectMarketData()
		case <-bi.stopChan:
			return
		}
	}
}

func (bi *ingestor) collectMarketData() {
	// Get required assets from strategy configs
	requiredAssets := bi.assetRegistry.GetRequiredAssets()

	if len(requiredAssets) == 0 {
		bi.logger.Debug("No assets required by enabled strategies")
		return
	}

	var wg sync.WaitGroup

	for _, conn := range bi.exchangeRegistry.GetReadyConnectors() {
		wg.Add(1)

		go func(conn connector.Connector) {
			defer wg.Done()

			exchangeName := conn.GetConnectorInfo().Name

			for _, asset := range requiredAssets {
				// Collect orderbooks
				supportedTypes := bi.getSupportedInstrumentTypes(conn)

				for _, instrumentType := range supportedTypes {
					// REST API call to fetch order book
					orderBook, err := conn.FetchOrderBook(asset, instrumentType, 20)
					if err != nil {
						bi.logger.Debug("Failed to fetch %s orderbook for %s on %s: %v",
							instrumentType, asset.Symbol(), string(exchangeName), err)
						bi.healthStore.RecordDataError(exchangeName, health.DataTypeOrderbooks, err)
						continue
					}

					bi.store.UpdateOrderBook(asset, exchangeName, instrumentType, *orderBook)
					bi.healthStore.RecordDataReceived(exchangeName, health.DataTypeOrderbooks, health.SourceBatch, 0)

					if len(orderBook.Bids) == 0 || len(orderBook.Asks) == 0 {
						bi.logger.Debug("Empty %s orderbook for %s on %s - no bids or asks",
							instrumentType, asset.Symbol(), string(exchangeName))
						continue
					}

					bi.logger.Debug("REST updated %s orderbook for %s on %s - bid: %s, ask: %s",
						instrumentType, asset.Symbol(), string(exchangeName),
						orderBook.Bids[0].Price.StringFixed(2),
						orderBook.Asks[0].Price.StringFixed(2))
				}

				// Collect klines for multiple intervals
				bi.collectKlines(conn, exchangeName, asset)
			}
		}(conn)
	}

	wg.Wait()

	// Notify that data was updated
	bi.notifyDataUpdate()
}

func (bi *ingestor) collectKlines(conn connector.Connector, exchangeName connector.ExchangeName, asset portfolio.Asset) {
	intervals := []string{"1m", "5m", "15m", "1h", "4h", "1d"}

	for _, interval := range intervals {
		klines, err := conn.FetchKlines(asset.Symbol(), interval, 100)
		if err != nil {
			bi.logger.Debug("Failed to fetch %s klines for %s on %s: %v",
				interval, asset.Symbol(), string(exchangeName), err)
			continue
		}

		if len(klines) == 0 {
			bi.logger.Debug("No %s klines for %s on %s", interval, asset.Symbol(), string(exchangeName))
			continue
		}

		// Store all klines
		for _, kline := range klines {
			bi.store.UpdateKline(asset, exchangeName, kline)
		}

		bi.logger.Debug("REST updated %d %s klines for %s on %s",
			len(klines), interval, asset.Symbol(), string(exchangeName))
	}
}

// notifyDataUpdate signals that data was updated
func (bi *ingestor) notifyDataUpdate() {
	bi.notifier.Notify()
}

// Get supported instrument types from exchange capabilities
func (bi *ingestor) getSupportedInstrumentTypes(conn connector.Connector) []connector.Instrument {
	var types []connector.Instrument

	if conn.SupportsPerpetuals() {
		types = append(types, connector.TypePerpetual)
	}

	if conn.SupportsSpot() {
		types = append(types, connector.TypeSpot)
	}

	return types
}

func (bi *ingestor) CollectNow() {
	if bi.IsActive() {
		bi.logger.Info("Triggering immediate data collection")
		go bi.collectMarketData()
	} else {
		bi.logger.Warn("Cannot collect now - batch ingestor is not active")
	}
}

func (bi *ingestor) Stop() error {
	bi.mutex.Lock()
	defer bi.mutex.Unlock()

	if !bi.isActive {
		return nil
	}

	if bi.ticker != nil {
		bi.ticker.Stop()
	}

	close(bi.stopChan)
	bi.isActive = false

	bi.logger.Info("Stopped REST batch ingestion")
	return nil
}

func (bi *ingestor) IsActive() bool {
	bi.mutex.RLock()
	defer bi.mutex.RUnlock()
	return bi.isActive
}
