package executor

import (
	"fmt"

	perpTypes "github.com/wisp-trading/sdk/pkg/markets/perp/types"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	perpConn "github.com/wisp-trading/sdk/pkg/types/connector/perp"
	"github.com/wisp-trading/sdk/pkg/types/execution"
	"github.com/wisp-trading/sdk/pkg/types/logging"
	"github.com/wisp-trading/sdk/pkg/types/registry"
	"github.com/wisp-trading/sdk/pkg/types/strategy"
	"github.com/wisp-trading/sdk/pkg/types/temporal"
)

type executor struct {
	connectors   registry.ConnectorRegistry
	positions    perpTypes.PerpPositions
	trades       perpTypes.PerpTrades
	logger       logging.ApplicationLogger
	timeProvider temporal.TimeProvider
}

// NewExecutor creates a new perp market executor.
func NewExecutor(
	connectors registry.ConnectorRegistry,
	positions perpTypes.PerpPositions,
	trades perpTypes.PerpTrades,
	logger logging.ApplicationLogger,
	timeProvider temporal.TimeProvider,
) perpTypes.SignalExecutor {
	logger.Info("Initializing perp executor")
	return &executor{
		connectors:   connectors,
		positions:    positions,
		trades:       trades,
		logger:       logger,
		timeProvider: timeProvider,
	}
}

// ExecutePerpSignal executes all actions in a perp signal.
// Satisfies perpTypes.SignalExecutor.
func (e *executor) ExecutePerpSignal(
	signal strategy.PerpSignal,
	ctx *execution.ExecutionContext,
	result *execution.ExecutionResult,
) error {
	for i, action := range signal.GetActions() {
		if err := action.Validate(); err != nil {
			return fmt.Errorf("perp action %d invalid: %w", i, err)
		}

		orderID, err := e.executeAction(action)
		if err != nil {
			return fmt.Errorf("perp action %d failed: %w", i, err)
		}

		if orderID != "" {
			result.OrderIDs = append(result.OrderIDs, orderID)
		}
	}

	return nil
}

// HandleTrade records an inbound perp trade fill and marks the order filled.
func (e *executor) HandleTrade(trade connector.Trade) error {
	e.trades.AddTrade(trade)

	orderID := trade.OrderID
	if orderID == "" {
		orderID = trade.ID
	}

	if err := e.positions.UpdateOrderStatus(orderID, connector.OrderStatusFilled); err != nil {
		e.logger.Debug("Could not mark perp order %s filled: %v", orderID, err)
	}

	e.logger.Info("Perp trade recorded: %s (order: %s, pair: %s)", trade.ID, orderID, trade.Pair.Symbol())
	return nil
}

func (e *executor) executeAction(action *strategy.PerpAction) (string, error) {
	switch action.ActionType {
	case strategy.ActionHold:
		e.logger.Info("Holding perp position for %s", action.Pair.Symbol())
		return "", nil
	case strategy.ActionClose:
		e.logger.Info("Close perp action noted for %s", action.Pair.Symbol())
		return "", nil
	}

	conn, exists := e.connectors.Connector(action.Exchange)
	if !exists {
		return "", fmt.Errorf("exchange %s not available", action.Exchange)
	}

	perpConnector, isPerpConn := conn.(perpConn.Connector)

	side := getSide(action.ActionType)

	e.logger.Info("Executing perp %s order: %s %s @ %s (leverage: %s) on %s",
		action.ActionType, action.Quantity.StringFixed(4), action.Pair.Symbol(),
		action.Price.StringFixed(2), action.Leverage.StringFixed(1), action.Exchange,
	)

	if isPerpConn && !action.Leverage.IsZero() {
		if err := e.setLeverage(perpConnector, action); err != nil {
			e.logger.Warn("Could not set leverage for %s on %s: %v", action.Pair.Symbol(), action.Exchange, err)
		}
	}

	exec, ok := conn.(connector.OrderExecutor)
	if !ok {
		return "", fmt.Errorf("exchange %s does not support order execution", action.Exchange)
	}

	resp, err := exec.PlaceLimitOrder(action.Pair, side, action.Quantity, action.Price)
	if err != nil {
		return "", fmt.Errorf("failed to place perp order on %s: %w", action.Exchange, err)
	}

	e.positions.AddOrder(connector.Order{
		Pair:      action.Pair,
		ID:        resp.OrderID,
		Side:      side,
		Quantity:  action.Quantity,
		Price:     action.Price,
		Status:    connector.OrderStatusPending,
		Type:      connector.OrderTypeLimit,
		CreatedAt: e.timeProvider.Now(),
		UpdatedAt: e.timeProvider.Now(),
	})

	e.logger.Info("Perp order placed: %s (pair: %s, side: %s)", resp.OrderID, action.Pair.Symbol(), side)

	return resp.OrderID, nil
}

func (e *executor) setLeverage(conn perpConn.Connector, action *strategy.PerpAction) error {
	symbol := conn.GetPerpSymbol(action.Pair)
	if symbol == "" {
		return fmt.Errorf("could not resolve perp symbol for %s", action.Pair.Symbol())
	}
	e.logger.Debug("Leverage %s requested for %s (%s) on %s",
		action.Leverage.StringFixed(1), action.Pair.Symbol(), symbol, action.Exchange)
	return nil
}

func getSide(actionType strategy.ActionType) connector.OrderSide {
	switch actionType {
	case strategy.ActionSell, strategy.ActionSellShort:
		return connector.OrderSideSell
	default:
		return connector.OrderSideBuy
	}
}

var _ perpTypes.SignalExecutor = (*executor)(nil)
