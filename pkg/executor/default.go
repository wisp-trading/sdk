package executor

import (
	"fmt"

	"github.com/wisp-trading/sdk/pkg/markets/prediction/types"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/data/stores/activity"
	"github.com/wisp-trading/sdk/pkg/types/execution"
	"github.com/wisp-trading/sdk/pkg/types/logging"
	"github.com/wisp-trading/sdk/pkg/types/registry"
	"github.com/wisp-trading/sdk/pkg/types/strategy"
	"github.com/wisp-trading/sdk/pkg/types/temporal"
)

// executor is the standard executor implementation with hook support
type executor struct {
	connectors   registry.ConnectorRegistry
	positions    activity.Positions
	logger       logging.ApplicationLogger
	timeProvider temporal.TimeProvider

	// Hook registry for runtime hook management
	hookRegistry registry.Hooks
}

// NewExecutor creates a new default executor
func NewExecutor(
	connectors registry.ConnectorRegistry,
	positions activity.Positions,
	logger logging.ApplicationLogger,
	timeProvider temporal.TimeProvider,
	hookRegistry registry.Hooks,
) execution.Executor {
	logger.Info("📌 Initializing executor with hook registry")

	return &executor{
		connectors:   connectors,
		positions:    positions,
		logger:       logger,
		timeProvider: timeProvider,
		hookRegistry: hookRegistry,
	}
}

// ExecuteSignal processes a signal and executes the associated actions
func (e *executor) ExecuteSignal(signal strategy.Signal) error {
	ctx := &execution.ExecutionContext{
		Signal:    signal,
		Timestamp: e.timeProvider.Now(),
		Metadata:  make(map[string]interface{}),
	}

	e.logger.Info("🎯 Executing signal %s with strategy %s", signal.GetID(), signal.GetStrategy())

	// Get hooks from registry at execution time
	hooks := e.hookRegistry.GetHooks()

	// Run BeforeExecute hooks
	for _, hook := range hooks {
		if err := hook.BeforeExecute(ctx); err != nil {
			e.logger.Warn("🚫 Hook blocked execution: %v", err)
			e.handleError(ctx, err, hooks)
			return err
		}
	}

	// Execute core logic, dispatching on concrete signal type
	result := &execution.ExecutionResult{
		OrderIDs: make([]string, 0),
		Success:  true,
	}

	var execErr error
	switch s := signal.(type) {
	case strategy.SpotSignal:
		execErr = e.executeSpotSignal(ctx, s, result)
	case strategy.PerpSignal:
		execErr = e.executePerpSignal(ctx, s, result)
	case types.PredictionSignal:
		execErr = e.executePredictionSignal(ctx, s, result)
	default:
		execErr = fmt.Errorf("unsupported signal type: %T", signal)
	}

	if execErr != nil {
		result.Error = execErr
		result.Success = false
		e.handleError(ctx, execErr, hooks)
		return execErr
	}

	// Run AfterExecute hooks
	for _, hook := range hooks {
		if err := hook.AfterExecute(ctx, result); err != nil {
			e.logger.Error("Hook AfterExecute failed: %v", err)
		}
	}

	e.logger.Info("✅ Successfully executed all actions for signal %s", signal.GetID())
	return nil
}

// executeSpotSignal executes all actions in a spot signal
func (e *executor) executeSpotSignal(ctx *execution.ExecutionContext, signal strategy.SpotSignal, result *execution.ExecutionResult) error {
	for i, action := range signal.GetActions() {
		orderID, err := e.executeSpotAction(signal.GetStrategy(), action)
		if err != nil {
			e.logger.Error("Failed to execute spot action %d for signal %s: %v", i, signal.GetID(), err)
			return err
		}
		if orderID != "" {
			result.OrderIDs = append(result.OrderIDs, orderID)
		}
	}
	return nil
}

// executePerpSignal executes all actions in a perp signal
func (e *executor) executePerpSignal(ctx *execution.ExecutionContext, signal strategy.PerpSignal, result *execution.ExecutionResult) error {
	for i, action := range signal.GetActions() {
		orderID, err := e.executePerpAction(signal.GetStrategy(), action)
		if err != nil {
			e.logger.Error("Failed to execute perp action %d for signal %s: %v", i, signal.GetID(), err)
			return err
		}
		if orderID != "" {
			result.OrderIDs = append(result.OrderIDs, orderID)
		}
	}
	return nil
}

// executePredictionSignal executes all actions in a prediction signal
func (e *executor) executePredictionSignal(ctx *execution.ExecutionContext, signal types.PredictionSignal, result *execution.ExecutionResult) error {
	for i, action := range signal.GetActions() {
		e.logger.Info("🔮 Prediction action %d: %s on market %s", i, action.ActionType, action.Market.MarketID.String())
	}
	return nil
}

// executeSpotAction executes a single spot action
func (e *executor) executeSpotAction(strategyName strategy.StrategyName, action *strategy.SpotAction) (string, error) {
	switch action.ActionType {
	case strategy.ActionHold:
		e.logger.Info("📊 Holding position as instructed for %s", action.Pair.Symbol())
		return "", nil
	case strategy.ActionClose:
		e.logger.Info("🔚 Close action noted for %s", action.Pair.Symbol())
		return "", nil
	}

	exchange, exists := e.connectors.Connector(action.Exchange)
	if !exists {
		return "", fmt.Errorf("exchange %s not available", action.Exchange)
	}

	e.logger.Info(
		"📈 Executing %s order: %s %s at price %s on %s",
		action.ActionType,
		action.Quantity.StringFixed(4),
		action.Pair.Symbol(),
		action.Price.StringFixed(2),
		action.Exchange,
	)

	exec, ok := exchange.(connector.OrderExecutor)
	if !ok {
		return "", fmt.Errorf("exchange %s does not support order execution", action.Exchange)
	}

	orderResponse, err := e.placeSpotOrder(exec, action)
	if err != nil {
		return "", fmt.Errorf("failed to place order: %w", err)
	}

	order := connector.Order{
		Pair:      action.Pair,
		ID:        orderResponse.OrderID,
		Side:      e.getOrderSide(action.ActionType),
		Quantity:  action.Quantity,
		Price:     action.Price,
		Status:    connector.OrderStatusPending,
		Type:      connector.OrderTypeLimit,
		CreatedAt: e.timeProvider.Now(),
		UpdatedAt: e.timeProvider.Now(),
	}

	e.positions.AddOrderToStrategy(strategyName, order)
	e.logger.Info("✅ Order recorded for strategy %s: %s", strategyName, orderResponse.OrderID)
	return orderResponse.OrderID, nil
}

// executePerpAction executes a single perp action
func (e *executor) executePerpAction(strategyName strategy.StrategyName, action *strategy.PerpAction) (string, error) {
	switch action.ActionType {
	case strategy.ActionHold:
		e.logger.Info("📊 Holding perp position for %s", action.Pair.Symbol())
		return "", nil
	case strategy.ActionClose:
		e.logger.Info("🔚 Close perp action noted for %s", action.Pair.Symbol())
		return "", nil
	}

	exchange, exists := e.connectors.Connector(action.Exchange)
	if !exists {
		return "", fmt.Errorf("exchange %s not available", action.Exchange)
	}

	e.logger.Info(
		"📈 Executing perp %s order: %s %s at price %s (leverage: %s) on %s",
		action.ActionType,
		action.Quantity.StringFixed(4),
		action.Pair.Symbol(),
		action.Price.StringFixed(2),
		action.Leverage.StringFixed(1),
		action.Exchange,
	)

	exec, ok := exchange.(connector.OrderExecutor)
	if !ok {
		return "", fmt.Errorf("exchange %s does not support order execution", action.Exchange)
	}

	side := e.getOrderSide(action.ActionType)
	orderResponse, err := exec.PlaceLimitOrder(action.Pair, side, action.Quantity, action.Price)
	if err != nil {
		return "", fmt.Errorf("failed to place perp order: %w", err)
	}

	order := connector.Order{
		Pair:      action.Pair,
		ID:        orderResponse.OrderID,
		Side:      side,
		Quantity:  action.Quantity,
		Price:     action.Price,
		Status:    connector.OrderStatusPending,
		Type:      connector.OrderTypeLimit,
		CreatedAt: e.timeProvider.Now(),
		UpdatedAt: e.timeProvider.Now(),
	}

	e.positions.AddOrderToStrategy(strategyName, order)
	e.logger.Info("✅ Perp order recorded for strategy %s: %s", strategyName, orderResponse.OrderID)
	return orderResponse.OrderID, nil
}

// placeSpotOrder places a spot order on the exchange
func (e *executor) placeSpotOrder(exchange connector.OrderExecutor, action *strategy.SpotAction) (*connector.OrderResponse, error) {
	switch action.ActionType {
	case strategy.ActionBuy, strategy.ActionCover:
		return exchange.PlaceLimitOrder(action.Pair, connector.OrderSideBuy, action.Quantity, action.Price)
	case strategy.ActionSell, strategy.ActionSellShort:
		return exchange.PlaceLimitOrder(action.Pair, connector.OrderSideSell, action.Quantity, action.Price)
	default:
		return nil, fmt.Errorf("unsupported action type: %s", action.ActionType)
	}
}

// getOrderSide converts an action to an order side
func (e *executor) getOrderSide(actionType strategy.ActionType) connector.OrderSide {
	switch actionType {
	case strategy.ActionBuy, strategy.ActionCover:
		return connector.OrderSideBuy
	case strategy.ActionSell, strategy.ActionSellShort:
		return connector.OrderSideSell
	default:
		return connector.OrderSideBuy // Default fallback
	}
}

// HandleTradeExecution is called when a trade executes to record it for the strategy
func (e *executor) HandleTradeExecution(trade connector.Trade) error {
	// Use the trade's order ID if available, otherwise fall back to trade ID
	// The connector should populate the OrderID field to link trades to orders
	orderID := trade.OrderID
	if orderID == "" {
		orderID = trade.ID
		e.logger.Debug("Trade %s has no OrderID field, using trade ID as fallback", trade.ID)
	}

	// Find which strategy owns this order
	strategyName, exists := e.positions.GetStrategyForOrder(orderID)
	if !exists {
		e.logger.Debug("Trade %s (order %s) could not be matched to any strategy order", trade.ID, orderID)
		return nil
	}

	// Record the trade for the strategy
	e.positions.AddTradeToStrategy(strategyName, trade)

	// Update the corresponding order to mark it as filled
	err := e.positions.UpdateOrderStatus(strategyName, orderID, connector.OrderStatusFilled)
	if err != nil {
		e.logger.Debug("Could not update order %s: %v", orderID, err)
	}

	e.logger.Info("✅ Trade executed and recorded for strategy %s: %s", strategyName, trade.ID)
	return nil
}

// handleError calls OnError hooks
func (e *executor) handleError(ctx *execution.ExecutionContext, err error, hooks []execution.ExecutionHook) {
	for _, hook := range hooks {
		if hookErr := hook.OnError(ctx, err); hookErr != nil {
			e.logger.Error("Hook OnError failed: %v", hookErr)
		}
	}
}
