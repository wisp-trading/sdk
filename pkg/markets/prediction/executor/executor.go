package executor

import (
	"fmt"

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
	connectors   registry.ConnectorRegistry
	store        predTypes.MarketStore
	logger       logging.ApplicationLogger
	timeProvider temporal.TimeProvider
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
		connectors:   connectors,
		store:        store,
		logger:       logger,
		timeProvider: timeProvider,
	}
}

// ExecutePredictionSignal executes all actions in a prediction signal.
// Satisfies predTypes.SignalExecutor.
func (e *executor) ExecutePredictionSignal(
	signal predTypes.PredictionSignal,
	ctx *execution.ExecutionContext,
	result *execution.ExecutionResult,
) error {
	actions := signal.GetActions()

	for i, action := range actions {
		if err := action.Validate(); err != nil {
			return fmt.Errorf("prediction action %d invalid: %w", i, err)
		}

		orderID, err := e.executeAction(signal.GetStrategy(), action)
		if err != nil {
			return fmt.Errorf("prediction action %d failed: %w", i, err)
		}

		result.OrderIDs = append(result.OrderIDs, orderID)
	}

	return nil
}

func (e *executor) executeAction(strategyName strategy.StrategyName, action *predTypes.PredictionAction) (string, error) {
	exchange, exists := e.connectors.Connector(action.Exchange)
	if !exists {
		return "", fmt.Errorf("exchange %s not available", action.Exchange)
	}

	predExec, ok := exchange.(predconnector.OrderExecutor)
	if !ok {
		return "", fmt.Errorf("exchange %s does not support prediction order execution", action.Exchange)
	}

	side := getSide(action.ActionType)

	resp, err := predExec.PlaceLimitOrder(predconnector.LimitOrder{
		Outcome:    action.Outcome,
		Price:      action.MaxPrice,
		Amount:     action.Shares,
		Side:       side,
		Expiration: action.Expiration,
	})
	if err != nil {
		return "", fmt.Errorf("failed to place prediction order on %s: %w", action.Exchange, err)
	}

	e.store.AddOrder(predTypes.PredictionOrder{
		ID:         resp.OrderID,
		Exchange:   action.Exchange,
		MarketSlug: action.Market.Slug,
		OutcomeID:  action.Outcome.OutcomeID,
		Side:       side,
		Shares:     action.Shares,
		Price:      action.MaxPrice,
		Status:     connector.OrderStatusPending,
		CreatedAt:  e.timeProvider.Now(),
		UpdatedAt:  e.timeProvider.Now(),
	})

	e.logger.Info(
		"Prediction order %s placed (strategy: %s, market: %s, outcome: %s, side: %s, shares: %s @ %s)",
		resp.OrderID,
		strategyName,
		action.Market.Slug,
		action.Outcome.OutcomeID,
		side,
		action.Shares.StringFixed(4),
		action.MaxPrice.StringFixed(4),
	)

	return resp.OrderID, nil
}

func getSide(actionType strategy.ActionType) connector.OrderSide {
	switch actionType {
	case strategy.ActionSell, strategy.ActionSellShort:
		return connector.OrderSideSell
	default:
		return connector.OrderSideBuy
	}
}
