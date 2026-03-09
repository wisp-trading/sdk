package monitoring

import (
	"encoding/json"

	prediction "github.com/wisp-trading/sdk/pkg/markets/prediction/types/connector"
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

	// GetOrderbook delegates to spot or perp based on the registered connector type for the exchange.
	GetOrderbook(exchange connector.ExchangeName, pair portfolio.Pair) *connector.OrderBook

	GetPredictionOrderbookView(
		exchange connector.ExchangeName,
		marketID prediction.MarketID,
		outcomeID prediction.OutcomeID,
	) *connector.OrderBook

	// GetKlines delegates to spot or perp based on the registered connector type for the exchange.
	GetKlines(exchange connector.ExchangeName, pair portfolio.Pair, interval string, limit int) []connector.Kline

	// GetStrategyStatus returns the latest status snapshot for each running strategy.
	GetStrategyStatus() []StrategyStatusView

	// GetStrategyStatusLog returns the last N status snapshots for each strategy,
	// ordered oldest-first. N is capped at StatusLogMaxEntries.
	GetStrategyStatusLog() []StrategyStatusView
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
	Exchange connector.ExchangeName `json:"exchange"`
	Pair     portfolio.Pair         `json:"-"`
}

type spotMarketViewWire struct {
	Exchange connector.ExchangeName `json:"exchange"`
	Symbol   string                 `json:"symbol"`
}

func (v SpotMarketView) MarshalJSON() ([]byte, error) {
	return json.Marshal(spotMarketViewWire{
		Exchange: v.Exchange,
		Symbol:   v.Pair.Symbol(),
	})
}

func (v *SpotMarketView) UnmarshalJSON(data []byte) error {
	var w spotMarketViewWire
	if err := json.Unmarshal(data, &w); err != nil {
		return err
	}
	v.Exchange = w.Exchange
	v.Pair = pairFromSymbol(w.Symbol)
	return nil
}

// PerpMarketView represents a perp pair watched on an exchange.
type PerpMarketView struct {
	Exchange connector.ExchangeName `json:"exchange"`
	Pair     portfolio.Pair         `json:"-"`
}

type perpMarketViewWire struct {
	Exchange connector.ExchangeName `json:"exchange"`
	Symbol   string                 `json:"symbol"`
}

func (v PerpMarketView) MarshalJSON() ([]byte, error) {
	return json.Marshal(perpMarketViewWire{
		Exchange: v.Exchange,
		Symbol:   v.Pair.Symbol(),
	})
}

func (v *PerpMarketView) UnmarshalJSON(data []byte) error {
	var w perpMarketViewWire
	if err := json.Unmarshal(data, &w); err != nil {
		return err
	}
	v.Exchange = w.Exchange
	v.Pair = pairFromSymbol(w.Symbol)
	return nil
}

// PredictionMarketView represents a prediction market with its full outcome list.
type PredictionMarketView struct {
	Exchange connector.ExchangeName  `json:"exchange"`
	MarketID prediction.MarketID     `json:"market_id"`
	Slug     string                  `json:"slug"`
	Outcomes []PredictionOutcomeView `json:"outcomes"`
}

// PredictionOutcomeView is a single tradeable outcome within a prediction market.
type PredictionOutcomeView struct {
	OutcomeID prediction.OutcomeID `json:"outcome_id"`
	Name      string               `json:"name"`
}

// ViewQuerier queries views from running strategy instances via Unix socket.
type ViewQuerier interface {
	QueryPnL(instanceID string) (*PnLView, error)
	QueryPositions(instanceID string) (*strategy.StrategyExecution, error)
	QueryRecentTrades(instanceID string, limit int) ([]connector.Trade, error)
	QueryMetrics(instanceID string) (*StrategyMetrics, error)
	QueryProfilingStats(instanceID string) (*ProfilingStats, error)
	QueryRecentExecutions(instanceID string, limit int) ([]ProfilingMetrics, error)

	// QueryMarkets returns the live market tree for building the CLI navigation hierarchy.
	QueryMarkets(instanceID string) (*MarketViews, error)

	// QueryOrderbook retrieves a spot/perp order book — exchange determines market type automatically.
	QueryOrderbook(instanceID string, exchange connector.ExchangeName, pair portfolio.Pair) (*connector.OrderBook, error)

	// QueryKlines retrieves kline/candlestick data — exchange determines spot vs perp automatically.
	QueryKlines(
		instanceID string,
		exchange connector.ExchangeName,
		pair portfolio.Pair,
		interval string,
		limit int,
	) ([]connector.Kline, error)

	// QueryPredictionOrderbook retrieves an order book for a specific prediction market outcome.
	QueryPredictionOrderbook(instanceID string, marketID prediction.MarketID, outcomeID prediction.OutcomeID) (*connector.OrderBook, error)

	// QueryStatus returns the latest status snapshot for each strategy in the instance.
	QueryStatus(instanceID string) ([]StrategyStatusView, error)

	// QueryStatusLog returns the full status history (up to StatusLogMaxEntries) for each strategy.
	QueryStatusLog(instanceID string) ([]StrategyStatusView, error)

	HealthCheck(instanceID string) error
	Shutdown(instanceID string) error
	ListInstances() ([]string, error)
}

// pairFromSymbol reconstructs a portfolio.Pair from a symbol string like "BTC-USD".
// Splits on the last "-" to handle assets that may contain hyphens.
func pairFromSymbol(symbol string) portfolio.Pair {
	for i := len(symbol) - 1; i >= 0; i-- {
		if symbol[i] == '-' {
			base := portfolio.NewAsset(symbol[:i])
			quote := portfolio.NewAsset(symbol[i+1:])
			return portfolio.NewPair(base, quote, "-")
		}
	}
	// No separator — treat the whole symbol as the base with empty quote
	return portfolio.NewPair(portfolio.NewAsset(symbol), portfolio.NewAsset(""), "")
}
