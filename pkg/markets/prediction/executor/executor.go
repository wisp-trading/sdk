package executor

import (
	"fmt"

	baseExecutor "github.com/wisp-trading/sdk/pkg/markets/base/executor"
	predTypes "github.com/wisp-trading/sdk/pkg/markets/prediction/types"
	predconnector "github.com/wisp-trading/sdk/pkg/markets/prediction/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/execution"
	"github.com/wisp-trading/sdk/pkg/types/logging"
	"github.com/wisp-trading/sdk/pkg/types/registry"
	"github.com/wisp-trading/sdk/pkg/types/strategy"
	"github.com/wisp-trading/sdk/pkg/types/temporal"
)

// executor handles execution of prediction market signals.
// It owns position tracking via the prediction store — no shared activity.Positions.
type executor struct {
	baseExecutor.Base
	store predTypes.MarketStore
}

// NewExecutor creates a new prediction market executor.
func NewExecutor(
	connectors registry.ConnectorRegistry,
	store predTypes.MarketStore,
	logger logging.ApplicationLogger,
	timeProvider temporal.TimeProvider,
) predTypes.SignalExecutor {
	logger.Info("Initializing prediction executor")
	return &executor{
		Base: baseExecutor.Base{
			Connectors:   connectors,
			Logger:       logger,
			TimeProvider: timeProvider,
		},
		store: store,
	}
}

// ExecutePredictionSignal executes all actions in a prediction signal.
// When multiple actions target the same exchange, they are batched into a single
// PlaceLimitOrders call so the CLOB receives them in one HTTP request.
// Satisfies predTypes.SignalExecutor.
func (e *executor) ExecutePredictionSignal(
	signal predTypes.PredictionSignal,
	ctx *execution.ExecutionContext,
	result *execution.ExecutionResult,
) error {
	actions := signal.GetActions()

	// Validate all actions first.
	for i, action := range actions {
		if err := action.Validate(); err != nil {
			return fmt.Errorf("prediction action %d invalid: %w", i, err)
		}
	}

	// Group actions by exchange for batch submission.
	type actionGroup struct {
		exchange connector.ExchangeName
		actions  []predTypes.PredictionAction
	}
	var groups []actionGroup
	groupIdx := make(map[connector.ExchangeName]int)

	for _, action := range actions {
		idx, exists := groupIdx[action.Exchange]
		if !exists {
			idx = len(groups)
			groupIdx[action.Exchange] = idx
			groups = append(groups, actionGroup{exchange: action.Exchange})
		}
		groups[idx].actions = append(groups[idx].actions, action)
	}

	// Execute each exchange group as a batch.
	for _, group := range groups {
		orderIDs, err := e.executeBatch(signal.GetStrategy(), group.exchange, group.actions)
		if err != nil {
			return err
		}
		result.OrderIDs = append(result.OrderIDs, orderIDs...)
	}

	return nil
}

// executeBatch submits all actions for a single exchange in one PlaceLimitOrders call.
func (e *executor) executeBatch(strategyName strategy.StrategyName, exchangeName connector.ExchangeName, actions []predTypes.PredictionAction) ([]string, error) {
	exchange, exists := e.Connectors.Connector(exchangeName)
	if !exists {
		return nil, fmt.Errorf("exchange %s not available", exchangeName)
	}

	predExec, ok := exchange.(predconnector.OrderExecutor)
	if !ok {
		return nil, fmt.Errorf("exchange %s does not support prediction order execution", exchangeName)
	}

	// Build the order slice.
	orders := make([]predconnector.LimitOrder, len(actions))
	for i, action := range actions {
		orders[i] = predconnector.LimitOrder{
			Outcome:     action.Outcome,
			Price:       action.MaxPrice,
			Amount:      action.Shares,
			Side:        getSide(action.ActionType),
			Expiration:  action.Expiration,
			TimeInForce: action.TimeInForce,
		}
	}

	// Single-order path for backwards compatibility with connectors that only
	// implement PlaceLimitOrder properly.
	if len(orders) == 1 {
		resp, err := predExec.PlaceLimitOrder(orders[0])
		if err != nil {
			return nil, fmt.Errorf("prediction action 0 failed: failed to place prediction order on %s: %w", exchangeName, err)
		}
		e.recordOrder(strategyName, &actions[0], resp.OrderID, orders[0].Side)
		return []string{resp.OrderID}, nil
	}

	// Batch path.
	responses, err := predExec.PlaceLimitOrders(orders)
	if err != nil {
		return nil, fmt.Errorf("prediction batch failed on %s: %w", exchangeName, err)
	}

	orderIDs := make([]string, 0, len(responses))
	for i, resp := range responses {
		if i < len(actions) {
			e.recordOrder(strategyName, &actions[i], resp.OrderID, orders[i].Side)
		}
		orderIDs = append(orderIDs, resp.OrderID)
	}

	return orderIDs, nil
}

func (e *executor) recordOrder(strategyName strategy.StrategyName, action *predTypes.PredictionAction, orderID string, side connector.OrderSide) {
	e.store.AddOrder(predTypes.PredictionOrder{
		ID:        orderID,
		Exchange:  action.Exchange,
		MarketID:  action.Market.MarketID,
		OutcomeID: action.Outcome.OutcomeID,
		Side:      side,
		Shares:    action.Shares,
		Price:     action.MaxPrice,
		Status:    connector.OrderStatusPending,
		CreatedAt: e.TimeProvider.Now(),
		UpdatedAt: e.TimeProvider.Now(),
	})

	e.Logger.Info(
		"Prediction order %s placed (strategy: %s, market: %s, outcome: %s, side: %s, shares: %s @ %s)",
		orderID,
		strategyName,
		action.Market.Slug,
		action.Outcome.OutcomeID,
		side,
		action.Shares.StringFixed(4),
		action.MaxPrice.StringFixed(4),
	)
}

func getSide(actionType strategy.ActionType) connector.OrderSide {
	switch actionType {
	case strategy.ActionSell, strategy.ActionSellShort:
		return connector.OrderSideSell
	default:
		return connector.OrderSideBuy
	}
}

var _ predTypes.SignalExecutor = (*executor)(nil)
