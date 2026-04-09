package prediction

import (
	"context"
	"fmt"
	"sync"

	"github.com/wisp-trading/sdk/pkg/markets/prediction/types"
	predictionconnector "github.com/wisp-trading/sdk/pkg/markets/prediction/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/logging"
	"github.com/wisp-trading/sdk/pkg/types/registry"
)

// marketLoader handles background pagination and loading of markets into the store
type marketLoader struct {
	store               types.MarketStore
	connectorRegistry   registry.ConnectorRegistry
	logger              logging.ApplicationLogger
	batchSize           int
	mu                  sync.RWMutex
	loadingState        map[connector.ExchangeName]bool
	loadProgress        map[connector.ExchangeName]int
	activeCancelFuncs   map[connector.ExchangeName]context.CancelFunc
}

// NewMarketLoader creates a new market loader with default batch size of 100.
// This is exported for fx dependency injection.
func NewMarketLoader(
	store types.MarketStore,
	connectorRegistry registry.ConnectorRegistry,
	logger logging.ApplicationLogger,
) types.MarketLoader {
	return &marketLoader{
		store:             store,
		connectorRegistry: connectorRegistry,
		logger:            logger,
		batchSize:         100,
		loadingState:      make(map[connector.ExchangeName]bool),
		loadProgress:      make(map[connector.ExchangeName]int),
		activeCancelFuncs: make(map[connector.ExchangeName]context.CancelFunc),
	}
}

func (ml *marketLoader) LoadMarkets(
	exchange connector.ExchangeName,
	filter *predictionconnector.MarketsFilter,
) error {
	ml.mu.Lock()
	defer ml.mu.Unlock()

	// Already loading for this exchange
	if ml.loadingState[exchange] {
		return fmt.Errorf("market loading already in progress for %s", exchange)
	}

	// Get connector
	conn, exists := ml.connectorRegistry.Prediction(exchange)
	if !exists {
		return fmt.Errorf("connector not found for exchange: %s", exchange)
	}

	// Mark as loading
	ml.loadingState[exchange] = true
	ml.loadProgress[exchange] = 0

	// Create cancellable context
	ctx, cancel := context.WithCancel(context.Background())
	ml.activeCancelFuncs[exchange] = cancel

	// Start background pagination loop
	go ml.paginateAndLoad(ctx, exchange, conn, filter)

	ml.logger.Info("Started background market loading for %s", exchange)
	return nil
}

func (ml *marketLoader) paginateAndLoad(
	ctx context.Context,
	exchange connector.ExchangeName,
	conn predictionconnector.Connector,
	filter *predictionconnector.MarketsFilter,
) {
	defer func() {
		ml.mu.Lock()
		ml.loadingState[exchange] = false
		delete(ml.activeCancelFuncs, exchange)
		ml.mu.Unlock()
		ml.logger.Info("Completed background market loading for %s", exchange)
	}()

	offset := 0
	totalFetched := 0

	for {
		select {
		case <-ctx.Done():
			ml.logger.Warn("Market loading cancelled for %s (fetched %d markets)", exchange, totalFetched)
			return
		default:
		}

		// Create a copy of the filter for this batch
		batchFilter := *filter
		batchFilter.Offset = &offset
		batchFilter.Limit = intPtr(ml.batchSize)

		// Fetch batch
		markets, err := conn.Markets(&batchFilter)
		if err != nil {
			ml.logger.Error("Failed to fetch market batch for %s at offset %d: %v", exchange, offset, err)
			return
		}

		if len(markets) == 0 {
			ml.logger.Info("Market loading complete for %s (fetched %d markets total)", exchange, totalFetched)
			return
		}

		// Store batch
		ml.store.UpdateMarkets(exchange, markets)
		totalFetched += len(markets)

		// Fetch orderbooks for all markets in this batch
		if err := ml.fetchAndStoreOrderBooks(ctx, exchange, conn, markets); err != nil {
			ml.logger.Warn("Failed to fetch orderbooks for batch: %v", err)
			// Continue loading markets even if orderbook fetch fails
		}

		// Update progress
		ml.mu.Lock()
		ml.loadProgress[exchange] = totalFetched
		ml.mu.Unlock()

		ml.logger.Debug("Loaded batch for %s: %d markets (total: %d)", exchange, len(markets), totalFetched)

		// Move to next batch
		offset += ml.batchSize

		// If we got fewer markets than batch size, we've reached the end
		if len(markets) < ml.batchSize {
			ml.logger.Info("Market loading complete for %s (fetched %d markets total)", exchange, totalFetched)
			return
		}
	}
}

func (ml *marketLoader) fetchAndStoreOrderBooks(
	ctx context.Context,
	exchange connector.ExchangeName,
	conn predictionconnector.Connector,
	markets []predictionconnector.Market,
) error {
	// Fetch orderbooks for each market using batch API
	for _, market := range markets {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		// Use batch API to fetch all outcomes for this market
		orderbooks, err := conn.FetchOrderBooksForMarket(market)
		if err != nil {
			ml.logger.Warn("Failed to fetch orderbooks for market %s: %v", market.MarketID, err)
			continue
		}

		// Store all orderbooks
		for outcomeID, orderbook := range orderbooks {
			ml.store.UpdateOrderBook(exchange, market.MarketID, predictionconnector.OutcomeID(outcomeID), orderbook.OrderBook)
		}
	}

	return nil
}

func (ml *marketLoader) IsLoading(exchange connector.ExchangeName) bool {
	ml.mu.RLock()
	defer ml.mu.RUnlock()
	return ml.loadingState[exchange]
}

func (ml *marketLoader) GetLoadProgress(exchange connector.ExchangeName) int {
	ml.mu.RLock()
	defer ml.mu.RUnlock()
	return ml.loadProgress[exchange]
}

func intPtr(i int) *int {
	return &i
}
