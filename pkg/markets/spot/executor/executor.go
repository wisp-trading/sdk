package executor

import (
	"fmt"

	baseExecutor "github.com/wisp-trading/sdk/pkg/markets/base/executor"
	spotTypes "github.com/wisp-trading/sdk/pkg/markets/spot/types"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/execution"
	"github.com/wisp-trading/sdk/pkg/types/logging"
	"github.com/wisp-trading/sdk/pkg/types/registry"
	"github.com/wisp-trading/sdk/pkg/types/strategy"
	"github.com/wisp-trading/sdk/pkg/types/temporal"
)

type executor struct {
	baseExecutor.Base
	store spotTypes.MarketStore
}

func NewExecutor(
	connectors registry.ConnectorRegistry,
	store spotTypes.MarketStore,
	logger logging.ApplicationLogger,
	timeProvider temporal.TimeProvider,
) spotTypes.SignalExecutor {
	return &executor{
		Base: baseExecutor.Base{
			Connectors:   connectors,
			Logger:       logger,
			TimeProvider: timeProvider,
		},
		store: store,
	}
}

func (e *executor) ExecuteSpotSignal(
	signal spotTypes.SpotSignal,
	ctx *execution.ExecutionContext,
	result *execution.ExecutionResult,
) error {
	for i, action := range signal.GetActions() {
		if err := action.Validate(); err != nil {
			return fmt.Errorf("spot action %d invalid: %w", i, err)
		}

		orderID, err := e.executeAction(&action)
		if err != nil {
			return fmt.Errorf("spot action %d failed: %w", i, err)
		}

		if orderID != "" {
			result.OrderIDs = append(result.OrderIDs, orderID)
		}
	}
	return nil
}

// HandleTrade records an inbound spot trade fill and marks the order filled.
func (e *executor) HandleTrade(trade connector.Trade) error {
	e.store.AddTrade(trade)

	orderID := trade.OrderID
	if orderID == "" {
		orderID = trade.ID
	}

	if err := e.store.UpdateOrderStatus(orderID, connector.OrderStatusFilled); err != nil {
		e.Logger.Debug("Could not mark spot order %s filled: %v", orderID, err)
	}

	e.Logger.Info("Spot trade recorded: %s (order: %s, pair: %s)", trade.ID, orderID, trade.Pair.Symbol())
	return nil
}

func (e *executor) executeAction(action *spotTypes.SpotAction) (string, error) {
	switch action.ActionType {
	case strategy.ActionHold:
		e.Logger.Info("Holding spot position for %s", action.Pair.Symbol())
		return "", nil
	case strategy.ActionClose:
		e.Logger.Info("Close spot action noted for %s", action.Pair.Symbol())
		return "", nil
	}

	exec, err := e.GetOrderExecutor(action.Exchange)
	if err != nil {
		return "", err
	}

	side := connector.OrderSideBuy
	if action.ActionType == strategy.ActionSell || action.ActionType == strategy.ActionSellShort {
		side = connector.OrderSideSell
	}

	e.Logger.Info("Executing spot %s order: %s %s @ %s on %s",
		action.ActionType, action.Quantity.StringFixed(4), action.Pair.Symbol(),
		action.Price.StringFixed(2), action.Exchange,
	)

	resp, err := exec.PlaceLimitOrder(action.Pair, side, action.Quantity, action.Price)
	if err != nil {
		return "", fmt.Errorf("failed to place spot order on %s: %w", action.Exchange, err)
	}

	e.store.AddOrder(connector.Order{
		Pair:      action.Pair,
		ID:        resp.OrderID,
		Side:      side,
		Quantity:  action.Quantity,
		Price:     action.Price,
		Status:    connector.OrderStatusPending,
		Type:      connector.OrderTypeLimit,
		CreatedAt: e.TimeProvider.Now(),
		UpdatedAt: e.TimeProvider.Now(),
	})

	e.Logger.Info("Spot order placed: %s (pair: %s, side: %s)", resp.OrderID, action.Pair.Symbol(), side)
	return resp.OrderID, nil
}

var _ spotTypes.SignalExecutor = (*executor)(nil)
