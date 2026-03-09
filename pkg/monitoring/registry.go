package monitoring

import (
	"context"

	perpTypes "github.com/wisp-trading/sdk/pkg/markets/perp/types"
	predTypes "github.com/wisp-trading/sdk/pkg/markets/prediction/types"
	predictionconnector "github.com/wisp-trading/sdk/pkg/markets/prediction/types/connector"
	spotTypes "github.com/wisp-trading/sdk/pkg/markets/spot/types"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/monitoring"
	"github.com/wisp-trading/sdk/pkg/types/monitoring/health"
	"github.com/wisp-trading/sdk/pkg/types/monitoring/profiling"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	"github.com/wisp-trading/sdk/pkg/types/registry"
	"github.com/wisp-trading/sdk/pkg/types/strategy"
	"github.com/wisp-trading/sdk/pkg/types/wisp"
)

type viewRegistry struct {
	wisp              wisp.Wisp
	health            health.HealthStore
	strategyRegistry  registry.StrategyRegistry
	profilingStore    profiling.ProfilingStore
	predictionViews   predTypes.PredictionViews
	perpViews         perpTypes.PerpViews
	spotViews         spotTypes.SpotViews
	connectorRegistry registry.ConnectorRegistry
}

func NewViewRegistry(
	health health.HealthStore,
	wisp wisp.Wisp,
	strategyRegistry registry.StrategyRegistry,
	profilingStore profiling.ProfilingStore,
	predictionViews predTypes.PredictionViews,
	perpViews perpTypes.PerpViews,
	spotViews spotTypes.SpotViews,
	connectorRegistry registry.ConnectorRegistry,
) monitoring.ViewRegistry {
	return &viewRegistry{
		health:            health,
		wisp:              wisp,
		strategyRegistry:  strategyRegistry,
		profilingStore:    profilingStore,
		predictionViews:   predictionViews,
		perpViews:         perpViews,
		spotViews:         spotViews,
		connectorRegistry: connectorRegistry,
	}
}

// getStrategyName returns the single registered strategy name
func (r *viewRegistry) getStrategyName() strategy.StrategyName {
	strategies := r.strategyRegistry.GetAllStrategies()
	if len(strategies) == 0 {
		return ""
	}
	return strategies[0].GetName()
}

func (r *viewRegistry) GetPnLView() *monitoring.PnLView {
	name := r.getStrategyName()
	ctx := context.Background()
	if name == "" {
		return nil
	}

	pnl := r.wisp.Activity().PNL()
	realizedPnL := pnl.TotalRealized(ctx)
	unrealizedPnL := pnl.TotalUnrealized(ctx)
	totalPnL := realizedPnL.Add(unrealizedPnL)
	totalFees := pnl.TotalFees(ctx)

	return &monitoring.PnLView{
		StrategyName:  string(name),
		RealizedPnL:   realizedPnL,
		UnrealizedPnL: unrealizedPnL,
		TotalPnL:      totalPnL,
		TotalFees:     totalFees,
	}
}

// todo need to rework cli flow
func (r *viewRegistry) GetPositionsView() *strategy.StrategyExecution {
	return nil
}

// GetOrderbook delegates to the correct domain views based on exchange type.
func (r *viewRegistry) GetOrderbook(exchange connector.ExchangeName, pair portfolio.Pair) *connector.OrderBook {
	marketType, ok := r.connectorRegistry.ConnectorType(exchange)
	if !ok {
		return nil
	}

	switch marketType {
	case connector.MarketTypeSpot:
		return r.spotViews.GetOrderbook(exchange, pair)
	case connector.MarketTypePerp:
		return r.perpViews.GetOrderbook(exchange, pair)
	default:
		return nil
	}
}

// GetKlines delegates to the correct domain views based on the registered connector type for the exchange.
func (r *viewRegistry) GetKlines(exchange connector.ExchangeName, pair portfolio.Pair, interval string, limit int) []connector.Kline {
	marketType, ok := r.connectorRegistry.ConnectorType(exchange)
	if !ok {
		return nil
	}

	switch marketType {
	case connector.MarketTypeSpot:
		return r.spotViews.GetKlines(exchange, pair, interval, limit)
	case connector.MarketTypePerp:
		return r.perpViews.GetKlines(exchange, pair, interval, limit)
	default:
		return nil
	}
}

// / todo need to rework cli flow
func (r *viewRegistry) GetRecentTrades(limit int) []connector.Trade {
	return nil
}

func (r *viewRegistry) GetMetrics() *monitoring.StrategyMetrics {
	name := r.getStrategyName()
	return &monitoring.StrategyMetrics{
		StrategyName: string(name),
		Status:       "running",
	}
}

func (r *viewRegistry) GetHealth() *health.SystemHealthReport {
	return r.health.GetSystemHealth()
}

// GetMarketViews returns the live market tree across all market types.
// Spot comes from wisp.Universe(); perp and prediction are delegated to their
// respective views packages which own those domains.
func (r *viewRegistry) GetMarketViews() *monitoring.MarketViews {
	views := &monitoring.MarketViews{}

	views.Spot = r.spotViews.GetMarketViews()

	// Perp is owned by the perp views package
	views.Perp = r.perpViews.GetMarketViews()

	// Prediction is owned by the prediction views package
	views.Prediction = r.predictionViews.GetMarketViews()

	return views
}

// GetPredictionOrderbookView delegates to the prediction views package.
func (r *viewRegistry) GetPredictionOrderbookView(
	exchange connector.ExchangeName,
	marketID predictionconnector.MarketID,
	outcomeID predictionconnector.OutcomeID,
) *connector.OrderBook {
	return r.predictionViews.GetOrderBook(
		exchange,
		marketID,
		outcomeID,
	)
}

func (r *viewRegistry) GetProfilingStats() *monitoring.ProfilingStats {
	if r.profilingStore == nil {
		return nil
	}

	name := r.getStrategyName()
	if name == "" {
		return nil
	}

	stats := r.profilingStore.GetStats(string(name))
	return &stats
}

func (r *viewRegistry) GetRecentExecutions(limit int) []monitoring.ProfilingMetrics {
	if r.profilingStore == nil {
		return nil
	}

	name := r.getStrategyName()
	if name == "" {
		return nil
	}

	return r.profilingStore.GetRecentMetrics(string(name), limit)
}

// GetStrategyStatus returns the latest status snapshot for each registered strategy.
// Delegates directly to each strategy's own log — no caching in the registry.
func (r *viewRegistry) GetStrategyStatus() []monitoring.StrategyStatusView {
	strategies := r.strategyRegistry.GetAllStrategies()
	out := make([]monitoring.StrategyStatusView, 0, len(strategies))
	for _, strat := range strategies {
		s := strat.LatestStatus()
		if s.At.IsZero() {
			continue // no status emitted yet
		}
		out = append(out, monitoring.StrategyStatusView{
			StrategyName: string(strat.GetName()),
			Phase:        string(s.Phase),
			Summary:      s.Summary,
			Fields:       s.Fields,
			At:           s.At,
		})
	}
	return out
}

// GetStrategyStatusLog returns the full status history for each registered strategy,
// oldest-first. Delegates directly to each strategy's own ring buffer.
func (r *viewRegistry) GetStrategyStatusLog() []monitoring.StrategyStatusView {
	strategies := r.strategyRegistry.GetAllStrategies()
	var out []monitoring.StrategyStatusView
	for _, strat := range strategies {
		for _, s := range strat.StatusLog() {
			out = append(out, monitoring.StrategyStatusView{
				StrategyName: string(strat.GetName()),
				Phase:        string(s.Phase),
				Summary:      s.Summary,
				Fields:       s.Fields,
				At:           s.At,
			})
		}
	}
	return out
}
