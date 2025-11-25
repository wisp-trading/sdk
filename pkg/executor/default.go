package executor

import (
	"fmt"
	"sync"

	"github.com/backtesting-org/kronos-sdk/pkg/types/connector"
	"github.com/backtesting-org/kronos-sdk/pkg/types/data/stores/activity"
	"github.com/backtesting-org/kronos-sdk/pkg/types/execution"
	"github.com/backtesting-org/kronos-sdk/pkg/types/logging"
	"github.com/backtesting-org/kronos-sdk/pkg/types/registry"
	"github.com/backtesting-org/kronos-sdk/pkg/types/strategy"
	"github.com/backtesting-org/kronos-sdk/pkg/types/temporal"
)

// executor is the standard executor implementation with hook support
type executor struct {
	connectors   registry.ConnectorRegistry
	positions    activity.Positions
	logger       logging.ApplicationLogger
	timeProvider temporal.TimeProvider

	// Track which strategy placed each order
	orderToStrategy map[string]strategy.StrategyName
	mu              sync.RWMutex

	// Execution hooks
	hooks   []execution.ExecutionHook
	hooksMu sync.RWMutex
}

// NewExecutor creates a new default executor
func NewExecutor(
	connectors registry.ConnectorRegistry,
	positions activity.Positions,
	logger logging.ApplicationLogger,
	timeProvider temporal.TimeProvider,
) execution.Executor {
	return &executor{
		connectors:      connectors,
		positions:       positions,
		logger:          logger,
		timeProvider:    timeProvider,
		orderToStrategy: make(map[string]strategy.StrategyName),
		hooks:           make([]execution.ExecutionHook, 0),
	}
}

// RegisterHook adds an execution hook
func (e *executor) RegisterHook(hook execution.ExecutionHook) {
	e.hooksMu.Lock()
	defer e.hooksMu.Unlock()
	e.hooks = append(e.hooks, hook)
	e.logger.Info("📌 Registered execution hook: %T", hook)
}

// ExecuteSignal processes a signal and executes the associated actions
func (e *executor) ExecuteSignal(signal *strategy.Signal) error {
	ctx := &ExecutionContext{
		Signal:    signal,
		Timestamp: e.timeProvider.Now(),
		Metadata:  make(map[string]interface{}),
	}

	e.logger.Info("🎯 Executing signal %s with %d actions", signal.ID, len(signal.Actions))

	// Get hooks snapshot
	e.hooksMu.RLock()
	hooks := make([]execution.ExecutionHook, len(e.hooks))
	copy(hooks, e.hooks)
	e.hooksMu.RUnlock()

	// Run BeforeExecute hooks
	for _, hook := range hooks {
		if err := hook.BeforeExecute(ctx); err != nil {
			e.logger.Warn("🚫 Hook blocked execution: %v", err)
			e.handleError(ctx, err, hooks)
			return err
		}
	}

	// Execute core logic
	result := &ExecutionResult{
		OrderIDs: make([]string, 0),
		Success:  true,
	}

	for i, action := range signal.Actions {
		orderID, err := e.executeAction(signal, action)
		if err != nil {
			e.logger.Error("Failed to execute action %d for signal %s: %v", i, signal.ID, err)
			result.Error = err
			result.Success = false
			e.handleError(ctx, err, hooks)
			return err
		}
		if orderID != "" {
			result.OrderIDs = append(result.OrderIDs, orderID)
		}
	}

	// Run AfterExecute hooks
	for _, hook := range hooks {
		if err := hook.AfterExecute(ctx, result); err != nil {
			e.logger.Error("Hook AfterExecute failed: %v", err)
			// Don't fail the execution if post-execution hooks fail
		}
	}

	e.logger.Info("✅ Successfully executed all actions for signal %s", signal.ID)
	return nil
}

// executeAction executes a single trade action
func (e *executor) executeAction(signal *strategy.Signal, action strategy.TradeAction) (string, error) {
	switch action.Action {
	case strategy.ActionBuy, strategy.ActionSell, strategy.ActionSellShort, strategy.ActionCover:
		return e.executeTradeAction(signal, action)
	case strategy.ActionHold:
		e.logger.Info("📊 Holding position as instructed for %s", action.Asset.Symbol())
		return "", nil
	case strategy.ActionClose:
		e.logger.Info("🔚 Close action noted for %s", action.Asset.Symbol())
		return "", nil
	default:
		e.logger.Warn("Unknown action type: %s for signal %s", action.Action, signal.ID)
		return "", nil
	}
}

// executeTradeAction executes a buy/sell trade action
func (e *executor) executeTradeAction(signal *strategy.Signal, action strategy.TradeAction) (string, error) {
	// Get exchange connector
	exchange, exists := e.connectors.GetConnector(action.Exchange)
	if !exists {
		return "", fmt.Errorf("exchange %s not available", action.Exchange)
	}

	e.logger.Info(
		"📈 Executing %s order: %s %s at price %s on %s",
		action.Action,
		action.Quantity.StringFixed(4),
		action.Asset.Symbol(),
		action.Price.StringFixed(2),
		action.Exchange,
	)

	// Place order on exchange
	orderResponse, err := e.placeOrder(exchange, action)
	if err != nil {
		return "", fmt.Errorf("failed to place order: %w", err)
	}

	// Create order record
	order := connector.Order{
		ID:        orderResponse.OrderID,
		Symbol:    action.Asset.Symbol(),
		Side:      e.getOrderSide(action.Action),
		Quantity:  action.Quantity,
		Price:     action.Price,
		Status:    connector.OrderStatusPending,
		Type:      connector.OrderTypeLimit,
		CreatedAt: e.timeProvider.Now(),
		UpdatedAt: e.timeProvider.Now(),
	}

	// Add order to strategy execution
	e.positions.AddOrderToStrategy(signal.Strategy, order)

	// Track which strategy placed this order
	e.mu.Lock()
	e.orderToStrategy[orderResponse.OrderID] = signal.Strategy
	e.mu.Unlock()

	e.logger.Info("✅ Order recorded for strategy %s: %s", signal.Strategy, orderResponse.OrderID)
	return orderResponse.OrderID, nil
}

// placeOrder places an order on the exchange
func (e *executor) placeOrder(exchange connector.Connector, action strategy.TradeAction) (*connector.OrderResponse, error) {
	switch action.Action {
	case strategy.ActionBuy, strategy.ActionCover:
		return exchange.PlaceLimitOrder(action.Asset.Symbol(), connector.OrderSideBuy, action.Quantity, action.Price)
	case strategy.ActionSell, strategy.ActionSellShort:
		return exchange.PlaceLimitOrder(action.Asset.Symbol(), connector.OrderSideSell, action.Quantity, action.Price)
	default:
		return nil, fmt.Errorf("unsupported trade action: %s", action.Action)
	}
}

// getOrderSide converts an action to an order side
func (e *executor) getOrderSide(action strategy.Action) connector.OrderSide {
	switch action {
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
	// Extract order ID from trade ID (trades are typically named after orders)
	orderID := trade.ID
	if len(trade.ID) > 6 && trade.ID[:6] == "trade_" {
		orderID = trade.ID[6:] // Remove "trade_" prefix
	}

	e.mu.RLock()
	strategyName, exists := e.orderToStrategy[orderID]
	e.mu.RUnlock()

	if exists {
		// Record the trade for the strategy
		e.positions.AddTradeToStrategy(strategyName, trade)

		// Update the corresponding order to mark it as filled
		err := e.positions.UpdateOrderStatus(strategyName, orderID, connector.OrderStatusFilled)
		if err != nil {
			// Order might not exist in position store yet, which is okay
			e.logger.Debug("Could not update order %s: %v", orderID, err)
		}

		// Clean up the tracking
		e.mu.Lock()
		delete(e.orderToStrategy, orderID)
		e.mu.Unlock()

		e.logger.Info("✅ Trade executed and recorded for strategy %s: %s", strategyName, trade.ID)
	}

	return nil
}

// handleError calls OnError hooks
func (e *executor) handleError(ctx *ExecutionContext, err error, hooks []execution.ExecutionHook) {
	for _, hook := range hooks {
		if hookErr := hook.OnError(ctx, err); hookErr != nil {
			e.logger.Error("Hook OnError failed: %v", hookErr)
		}
	}
}
