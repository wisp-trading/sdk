package batch

//type batchIngestor struct {
//	connector     perpConn.Connector     // Single perp connector
//	exchangeName  connector.ExchangeName // Name of this exchange
//	assetRegistry registry.AssetRegistry // Assets to collect data for
//	store         perpStore.MarketStore  // Where to store perp data
//	logger        logging.ApplicationLogger
//	timeProvider  temporal.TimeProvider
//
//	// Kline configuration
//	klineLimits map[string]int // Number of candles to fetch per interval
//
//	// Scheduling
//	ticker   temporal.Ticker
//	stopChan chan struct{}
//	isActive bool
//	mu       sync.RWMutex
//}
//
//// NewBatchIngestor creates a batch ingestor for a single perp connector
//func NewBatchIngestor(
//	conn perpConn.Connector,
//	exchangeName connector.ExchangeName,
//	assetRegistry registry.AssetRegistry,
//	store perpStore.MarketStore,
//	timeProvider temporal.TimeProvider,
//	logger logging.ApplicationLogger,
//) ingestors.BatchIngestor {
//	return &batchIngestor{
//		connector:     conn,
//		exchangeName:  exchangeName,
//		assetRegistry: assetRegistry,
//		store:         store,
//		timeProvider:  timeProvider,
//		logger:        logger,
//		stopChan:      make(chan struct{}),
//
//		klineLimits: map[string]int{
//			"1m":  500,
//			"5m":  300,
//			"15m": 200,
//			"1h":  168,
//			"4h":  180,
//			"1d":  90,
//		},
//	}
//}
//
//func (bi *batchIngestor) Start(interval time.Duration) error {
//	bi.mu.Lock()
//	defer bi.mu.Unlock()
//
//	if bi.isActive {
//		return fmt.Errorf("batch ingestor for %s already active", bi.exchangeName)
//	}
//
//	bi.ticker = bi.timeProvider.NewTicker(interval)
//	bi.isActive = true
//
//	go bi.collectLoop()
//
//	bi.logger.Info("Started perp batch ingestion for %s with %v interval", bi.exchangeName, interval)
//	return nil
//}
//
//func (bi *batchIngestor) collectLoop() {
//	// Run initial collection immediately
//	bi.CollectNow()
//
//	for {
//		select {
//		case <-bi.ticker.C():
//			bi.CollectNow()
//		case <-bi.stopChan:
//			return
//		}
//	}
//}
//
//func (bi *batchIngestor) CollectNow() {
//	bi.logger.Debug("Starting perp market data collection for %s", bi.exchangeName)
//
//	// Collect shared market data
//	bi.collectOrderBooks()
//	bi.collectPrices()
//	bi.collectKlines()
//
//	// Collect perp-specific data
//	bi.collectFundingRates()
//
//	bi.logger.Debug("Completed perp market data collection for %s", bi.exchangeName)
//}
//
//// collectOrderBooks fetches order books for all registered assets
//func (bi *batchIngestor) collectOrderBooks() {
//	assets := bi.assetRegistry.GetRequiredAssets()
//	if len(assets) == 0 {
//		return
//	}
//
//	var wg sync.WaitGroup
//
//	for _, asset := range assets {
//		wg.Add(1)
//		go func(a portfolio.Asset) {
//			defer wg.Done()
//
//			orderBook, err := bi.connector.FetchOrderBook(a, 20)
//			if err != nil {
//				bi.logger.Debug("Failed to fetch order book for %s on %s: %v", a.Symbol(), bi.exchangeName, err)
//				return
//			}
//
//			if len(orderBook.Bids) == 0 || len(orderBook.Asks) == 0 {
//				bi.logger.Debug("Empty order book for %s on %s", a.Symbol(), bi.exchangeName)
//				return
//			}
//
//			bi.store.UpdateOrderBook(a, bi.exchangeName, *orderBook)
//			bi.store.UpdateLastUpdated(commonStore.UpdateKey{
//				DataType: commonStore.DataKeyOrderBooks,
//				Asset:    a,
//				Exchange: bi.exchangeName,
//			})
//
//			bi.logger.Debug("Updated order book for %s on %s - bid: %s, ask: %s",
//				a.Symbol(), bi.exchangeName,
//				orderBook.Bids[0].Price.StringFixed(2),
//				orderBook.Asks[0].Price.StringFixed(2))
//		}(asset)
//	}
//
//	wg.Wait()
//}
//
//// collectPrices fetches current prices for all registered assets
//func (bi *batchIngestor) collectPrices() {
//	assets := bi.assetRegistry.GetRequiredAssets()
//	if len(assets) == 0 {
//		return
//	}
//
//	var wg sync.WaitGroup
//
//	for _, asset := range assets {
//		wg.Add(1)
//		go func(a portfolio.Asset) {
//			defer wg.Done()
//
//			price, err := bi.connector.FetchPrice(a.Symbol())
//			if err != nil {
//				bi.logger.Debug("Failed to fetch price for %s on %s: %v", a.Symbol(), bi.exchangeName, err)
//				return
//			}
//
//			bi.store.UpdateAssetPrice(a, bi.exchangeName, *price)
//			bi.store.UpdateLastUpdated(commonStore.UpdateKey{
//				DataType: commonStore.DataKeyAssetPrice,
//				Asset:    a,
//				Exchange: bi.exchangeName,
//			})
//
//			bi.logger.Debug("Updated price for %s on %s = %s",
//				a.Symbol(), bi.exchangeName, price.Price.String())
//		}(asset)
//	}
//
//	wg.Wait()
//}
//
//// collectKlines fetches historical klines for all registered assets
//func (bi *batchIngestor) collectKlines() {
//	assets := bi.assetRegistry.GetRequiredAssets()
//	if len(assets) == 0 {
//		return
//	}
//
//	intervals := []string{"1m", "5m", "15m", "1h", "4h", "1d"}
//	var wg sync.WaitGroup
//
//	for _, asset := range assets {
//		for _, interval := range intervals {
//			wg.Add(1)
//			go func(a portfolio.Asset, iv string) {
//				defer wg.Done()
//
//				limit := bi.klineLimits[iv]
//				if limit == 0 {
//					limit = 100
//				}
//
//				klines, err := bi.connector.FetchKlines(a.Symbol(), iv, limit)
//				if err != nil {
//					bi.logger.Debug("Failed to fetch %s klines for %s on %s: %v", iv, a.Symbol(), bi.exchangeName, err)
//					return
//				}
//
//				if len(klines) == 0 {
//					return
//				}
//
//				// Store all klines
//				for _, kline := range klines {
//					bi.store.UpdateKline(a, bi.exchangeName, kline)
//				}
//
//				bi.logger.Debug("Updated %d %s klines for %s on %s", len(klines), iv, a.Symbol(), bi.exchangeName)
//			}(asset, interval)
//		}
//	}
//
//	wg.Wait()
//}
//
//// collectFundingRates fetches current funding rates (perp-specific)
//func (bi *batchIngestor) collectFundingRates() {
//	rates, err := bi.connector.FetchCurrentFundingRates()
//	if err != nil {
//		bi.logger.Error("Failed to fetch funding rates from %s: %v", bi.exchangeName, err)
//		return
//	}
//
//	// Update all funding rates from this connector
//	bi.store.UpdateFundingRates(bi.exchangeName, rates)
//
//	for asset, rate := range rates {
//		bi.store.UpdateLastUpdated(commonStore.UpdateKey{
//			DataType: perpStore.DataKeyFundingRates,
//			Asset:    asset,
//			Exchange: bi.exchangeName,
//		})
//
//		bi.logger.Debug("Updated funding rate for %s on %s = %s",
//			asset.Symbol(), bi.exchangeName, rate.CurrentRate.String())
//	}
//}
//
//func (bi *batchIngestor) Stop() error {
//	bi.mu.Lock()
//	defer bi.mu.Unlock()
//
//	if !bi.isActive {
//		return nil
//	}
//
//	if bi.ticker != nil {
//		bi.ticker.Stop()
//	}
//
//	close(bi.stopChan)
//	bi.isActive = false
//
//	bi.logger.Info("Stopped perp batch ingestion for %s", bi.exchangeName)
//	return nil
//}
//
//func (bi *batchIngestor) IsActive() bool {
//	bi.mu.RLock()
//	defer bi.mu.RUnlock()
//	return bi.isActive
//}
//
//func (bi *batchIngestor) GetMarketType() connector.MarketType {
//	return connector.MarketTypePerp
//}
//
//var _ ingestors.BatchIngestor = (*batchIngestor)(nil)
