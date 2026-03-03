package monitoring

import (
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/monitoring/health"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	"github.com/wisp-trading/sdk/pkg/types/strategy"
)

// ViewRegistry aggregates runtime data from SDK stores and exposes it for monitoring.
type ViewRegistry interface {
	GetPnLView() *PnLView
	GetPositionsView() *strategy.StrategyExecution
	GetRecentTrades(limit int) []connector.Trade
	GetMetrics() *StrategyMetrics
	GetHealth() *health.SystemHealthReport
	GetProfilingStats() *ProfilingStats
	GetRecentExecutions(limit int) []ProfilingMetrics

	// GetMarketViews returns the live market tree across all market types.
	// Each domain's views package populates its own section.
	GetMarketViews() *MarketViews

	// Orderbook queries — one per market type
	GetOrderbookView(pair portfolio.Pair) *connector.OrderBook
	GetPredictionOrderbookView(exchange, marketID, outcomeID string) *connector.OrderBook
}

// MarketViews is the top-level response for /api/markets.
// Structured by market type so the CLI can build its navigation tree without
// any type-switch logic baked into shared code.
type MarketViews struct {
	Spot       []SpotMarketView       `json:"spot"`
	Perp       []PerpMarketView       `json:"perp"`
	Prediction []PredictionMarketView `json:"prediction"`
}

// SpotMarketView represents a spot pair watched on an exchange.
type SpotMarketView struct {
	Exchange string `json:"exchange"`
	Pair     string `json:"pair"`
}

// PerpMarketView represents a perp pair watched on an exchange.
type PerpMarketView struct {
	Exchange string `json:"exchange"`
	Pair     string `json:"pair"`
}

// PredictionMarketView represents a prediction market with its full outcome list.
type PredictionMarketView struct {
	Exchange string                  `json:"exchange"`
	MarketID string                  `json:"market_id"`
	Slug     string                  `json:"slug"`
	Outcomes []PredictionOutcomeView `json:"outcomes"`
}

// PredictionOutcomeView is a single tradeable outcome within a prediction market.
type PredictionOutcomeView struct {
	OutcomeID string `json:"outcome_id"`
	Name      string `json:"name"`
}

// ViewQuerier queries views from running strategy instances via Unix socket.
// Implemented in the CLI.
type ViewQuerier interface {
	QueryPnL(instanceID string) (*PnLView, error)
	QueryPositions(instanceID string) (*strategy.StrategyExecution, error)
	QueryRecentTrades(instanceID string, limit int) ([]connector.Trade, error)
	QueryMetrics(instanceID string) (*StrategyMetrics, error)
	QueryProfilingStats(instanceID string) (*ProfilingStats, error)
	QueryRecentExecutions(instanceID string, limit int) ([]ProfilingMetrics, error)

	// QueryMarkets returns the live market tree for building the CLI navigation hierarchy.
	QueryMarkets(instanceID string) (*MarketViews, error)

	// QueryOrderbook retrieves a spot/perp order book — pair is "BTC-USDT"
	QueryOrderbook(instanceID, pair, exchange string) (*connector.OrderBook, error)

	// QueryPredictionOrderbook retrieves an order book for a specific prediction outcome
	QueryPredictionOrderbook(instanceID, exchange, marketID, outcomeID string) (*connector.OrderBook, error)

	HealthCheck(instanceID string) error
	Shutdown(instanceID string) error
	ListInstances() ([]string, error)
}
