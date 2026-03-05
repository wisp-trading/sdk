package monitoring

import (
	"context"

	"github.com/wisp-trading/sdk/pkg/markets/prediction/types"
	predictionconnector "github.com/wisp-trading/sdk/pkg/markets/prediction/types/connector"
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
	wisp             wisp.Wisp
	health           health.HealthStore
	strategyRegistry registry.StrategyRegistry
	profilingStore   profiling.ProfilingStore
	predictionViews  types.PredictionViews
}

func NewViewRegistry(
	health health.HealthStore,
	wisp wisp.Wisp,
	strategyRegistry registry.StrategyRegistry,
	profilingStore profiling.ProfilingStore,
	predictionViews types.PredictionViews,
) monitoring.ViewRegistry {
	return &viewRegistry{
		health:           health,
		wisp:             wisp,
		strategyRegistry: strategyRegistry,
		profilingStore:   profilingStore,
		predictionViews:  predictionViews,
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
	realizedPnL := pnl.GetRealizedPNL(ctx, name)
	unrealizedPnL, _ := pnl.GetUnrealizedPNL(ctx, name)
	totalPnL, _ := pnl.GetTotalPNL(ctx)
	totalFees := pnl.GetFeesByStrategy(ctx, name)

	return &monitoring.PnLView{
		StrategyName:  string(name),
		RealizedPnL:   realizedPnL,
		UnrealizedPnL: unrealizedPnL,
		TotalPnL:      totalPnL,
		TotalFees:     totalFees,
	}
}

func (r *viewRegistry) GetPositionsView() *strategy.StrategyExecution {
	name := r.getStrategyName()

	if name == "" {
		return nil
	}

	return r.wisp.Activity().Positions().GetStrategyExecution(name)
}

func (r *viewRegistry) GetOrderbookView(pair portfolio.Pair) *connector.OrderBook {
	ctx := context.Background()

	ob, err := r.wisp.Market().OrderBook(ctx, pair)
	if err != nil {
		return nil
	}
	return ob
}

func (r *viewRegistry) GetRecentTrades(limit int) []connector.Trade {
	name := r.getStrategyName()
	if name == "" {
		return nil
	}

	trades := r.wisp.Activity().Positions().GetTradesForStrategy(name)
	if len(trades) <= limit {
		return trades
	}
	return trades[len(trades)-limit:]
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
// Spot/perp come from wisp.Universe(); prediction is delegated to predictionViews
// which owns that domain and queries the watchlist directly.
func (r *viewRegistry) GetMarketViews() *monitoring.MarketViews {
	universe := r.wisp.Universe()
	views := &monitoring.MarketViews{}

	for _, ex := range universe.Exchanges {
		pairs := universe.Assets[ex.Name]
		for _, pair := range pairs {
			switch ex.MarketType {
			case connector.MarketTypePerp:
				views.Perp = append(views.Perp, monitoring.PerpMarketView{
					Exchange: string(ex.Name),
					Pair:     pair.Symbol(),
				})
			default: // MarketTypeSpot
				views.Spot = append(views.Spot, monitoring.SpotMarketView{
					Exchange: string(ex.Name),
					Pair:     pair.Symbol(),
				})
			}
		}
	}

	// Prediction is entirely owned by the prediction views package
	views.Prediction = r.predictionViews.GetMarketViews()

	return views
}

// GetPredictionOrderbookView delegates to the prediction views package.
func (r *viewRegistry) GetPredictionOrderbookView(exchange, marketID, outcomeID string) *connector.OrderBook {
	return r.predictionViews.GetOrderBook(
		connector.ExchangeName(exchange),
		predictionconnector.MarketID(marketID),
		predictionconnector.OutcomeID(outcomeID),
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
