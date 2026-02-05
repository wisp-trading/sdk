package ingestors

import (
	"fmt"
	"sync"
	"time"

	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/data/ingestors/batch"
	marketTypes "github.com/wisp-trading/sdk/pkg/types/data/stores/market"
	"github.com/wisp-trading/sdk/pkg/types/logging"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	"github.com/wisp-trading/sdk/pkg/types/registry"
	"github.com/wisp-trading/sdk/pkg/types/temporal"
)

// BatchIngestor is a generic base implementation for REST batch data collection
type BatchIngestor struct {
	conn          connector.Connector
	marketData    connector.MarketDataReader
	orderExecutor connector.OrderExecutor
	exchangeName  connector.ExchangeName
	marketType    connector.MarketType
	assetRegistry registry.PairRegistry
	store         marketTypes.MarketStore
	logger        logging.ApplicationLogger
	timeProvider  temporal.TimeProvider

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
	assetRegistry registry.PairRegistry,
	store marketTypes.MarketStore,
	timeProvider temporal.TimeProvider,
	logger logging.ApplicationLogger,
	extensions ...batch.CollectionExtension,
) batch.BatchIngestor {
	// Type assert to get market data capabilities
	marketData, _ := conn.(connector.MarketDataReader)
	orderExecutor, _ := conn.(connector.OrderExecutor)

	return &BatchIngestor{
		conn:          conn,
		marketData:    marketData,
		orderExecutor: orderExecutor,
		exchangeName:  exchangeName,
		marketType:    marketType,
		assetRegistry: assetRegistry,
		store:         store,
		timeProvider:  timeProvider,
		logger:        logger,
		stopChan:      make(chan struct{}),
		extensions:    extensions,
		klineLimits: map[string]int{
			"1m":  500,
			"5m":  300,
			"15m": 200,
			"1h":  168,
			"4h":  180,
			"1d":  90,
		},
	}
}

func (bi *BatchIngestor) Start(interval time.Duration) error {
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

func (bi *BatchIngestor) collectLoop() {
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

func (bi *BatchIngestor) CollectNow() {
	bi.logger.Debug("Starting %s market data collection for %s", bi.marketType, bi.exchangeName)

	assets := bi.assetRegistry.GetRequiredPairs()
	if len(assets) == 0 {
		bi.logger.Debug("No assets required for collection")
		return
	}

	// Collect general market data
	bi.collectOrderBooks(assets)
	bi.collectPrices(assets)
	bi.collectKlines(assets)

	// Run any market-specific collection extensions
	for _, ext := range bi.extensions {
		ext.Collect(bi.conn, bi.exchangeName, assets)
	}

	bi.logger.Debug("Completed %s market data collection for %s", bi.marketType, bi.exchangeName)
}

func (bi *BatchIngestor) collectOrderBooks(assets []portfolio.Pair) {
	if bi.marketData == nil {
		return
	}

	var wg sync.WaitGroup

	for _, asset := range assets {
		wg.Add(1)
		go func(a portfolio.Pair) {
			defer wg.Done()

			orderBook, err := bi.marketData.FetchOrderBook(a, 20)
			if err != nil {
				bi.logger.Debug("Failed to fetch order book for %s on %s: %v", a.Symbol(), bi.exchangeName, err)
				return
			}

			if len(orderBook.Bids) == 0 || len(orderBook.Asks) == 0 {
				bi.logger.Debug("Empty order book for %s on %s", a.Symbol(), bi.exchangeName)
				return
			}

			bi.store.UpdateOrderBook(a, bi.exchangeName, *orderBook)
			bi.store.UpdateLastUpdated(marketTypes.UpdateKey{
				DataType: marketTypes.DataKeyOrderBooks,
				Asset:    a,
				Exchange: bi.exchangeName,
			})

			bi.logger.Debug("Updated order book for %s on %s - bid: %s, ask: %s",
				a.Symbol(), bi.exchangeName,
				orderBook.Bids[0].Price.StringFixed(2),
				orderBook.Asks[0].Price.StringFixed(2))
		}(asset)
	}

	wg.Wait()
}

func (bi *BatchIngestor) collectPrices(pairs []portfolio.Pair) {
	if bi.marketData == nil {
		return
	}

	var wg sync.WaitGroup

	for _, pair := range pairs {
		wg.Add(1)
		go func(p portfolio.Pair) {
			defer wg.Done()

			price, err := bi.marketData.FetchPrice(p)
			if err != nil {
				bi.logger.Debug("Failed to fetch price for %s on %s: %v", p.Symbol(), bi.exchangeName, err)
				return
			}

			bi.store.UpdateAssetPrice(p, bi.exchangeName, *price)
			bi.store.UpdateLastUpdated(marketTypes.UpdateKey{
				DataType: marketTypes.DataKeyAssetPrice,
				Asset:    p,
				Exchange: bi.exchangeName,
			})

			bi.logger.Debug(
				"Updated price for %s on %s = %s",
				p.Symbol(),
				bi.exchangeName,
				price.Price.String(),
			)
		}(pair)
	}

	wg.Wait()
}

func (bi *BatchIngestor) collectKlines(pairs []portfolio.Pair) {
	if bi.marketData == nil {
		return
	}

	intervals := []string{"1m", "5m", "15m", "1h", "4h", "1d"}
	var wg sync.WaitGroup

	for _, pair := range pairs {
		for _, interval := range intervals {
			wg.Add(1)
			go func(p portfolio.Pair, iv string) {
				defer wg.Done()

				limit := bi.klineLimits[iv]
				if limit == 0 {
					limit = 100
				}

				klines, err := bi.marketData.FetchKlines(p, iv, limit)
				if err != nil {
					bi.logger.Debug("Failed to fetch %s klines for %s on %s: %v", iv, p.Symbol(), bi.exchangeName, err)
					return
				}

				if len(klines) == 0 {
					return
				}

				// Store all klines
				for _, kline := range klines {
					bi.store.UpdateKline(p, bi.exchangeName, kline)
				}

				bi.logger.Debug("Updated %d %s klines for %s on %s", len(klines), iv, p.Symbol(), bi.exchangeName)
			}(pair, interval)
		}
	}

	wg.Wait()
}

func (bi *BatchIngestor) Stop() error {
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

func (bi *BatchIngestor) IsActive() bool {
	bi.mu.RLock()
	defer bi.mu.RUnlock()
	return bi.isActive
}

func (bi *BatchIngestor) GetMarketType() connector.MarketType {
	return bi.marketType
}

var _ batch.BatchIngestor = (*BatchIngestor)(nil)
