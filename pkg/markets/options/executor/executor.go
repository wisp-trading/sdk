package executor

import (
	"fmt"

	baseExecutor "github.com/wisp-trading/sdk/pkg/markets/base/executor"
	optionsTypes "github.com/wisp-trading/sdk/pkg/markets/options/types"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	optionsconnector "github.com/wisp-trading/sdk/pkg/types/connector/options"
	"github.com/wisp-trading/sdk/pkg/types/execution"
	"github.com/wisp-trading/sdk/pkg/types/logging"
	registryTypes "github.com/wisp-trading/sdk/pkg/types/registry"
	"github.com/wisp-trading/sdk/pkg/types/strategy"
	"github.com/wisp-trading/sdk/pkg/types/temporal"
)

type executor struct {
	baseExecutor.Base
	store optionsTypes.OptionsStore
}

// NewExecutor creates a new options executor.
func NewExecutor(
	connectorRegistry registryTypes.ConnectorRegistry,
	store optionsTypes.OptionsStore,
	logger logging.ApplicationLogger,
	timeProvider temporal.TimeProvider,
) optionsTypes.SignalExecutor {
	return &executor{
		Base: baseExecutor.Base{
			Connectors:   connectorRegistry,
			Logger:       logger,
			TimeProvider: timeProvider,
		},
		store: store,
	}
}

// ExecuteOptionsSignal executes all actions in an options signal.
// Satisfies optionsTypes.SignalExecutor.
func (e *executor) ExecuteOptionsSignal(
	signal optionsTypes.OptionsSignal,
	ctx *execution.ExecutionContext,
	result *execution.ExecutionResult,
) error {
	for i, action := range signal.GetActions() {
		if err := action.Validate(); err != nil {
			return fmt.Errorf("options action %d invalid: %w", i, err)
		}

		orderID, err := e.executeAction(&action)
		if err != nil {
			return fmt.Errorf("options action %d failed: %w", i, err)
		}

		if orderID != "" {
			result.OrderIDs = append(result.OrderIDs, orderID)
		}
	}
	return nil
}

func (e *executor) executeAction(action *optionsTypes.OptionsAction) (string, error) {
	switch action.ActionType {
	case strategy.ActionHold:
		e.Logger.Info("Holding options position for %s", action.Contract.Pair.Symbol())
		return "", nil
	case strategy.ActionClose:
		e.Logger.Info("Close options action noted for %s", action.Contract.Pair.Symbol())
		return "", nil
	}

	conn, exists := e.Connectors.Connector(action.Exchange)
	if !exists {
		return "", fmt.Errorf("exchange %s not available", action.Exchange)
	}

	optionsConn, ok := conn.(optionsconnector.Connector)
	if !ok {
		return "", fmt.Errorf("exchange %s does not support options", action.Exchange)
	}

	exec, ok := optionsConn.(connector.OrderExecutor)
	if !ok {
		return "", fmt.Errorf("exchange %s does not support order execution", action.Exchange)
	}

	side := connector.OrderSideBuy
	if action.ActionType == strategy.ActionSell || action.ActionType == strategy.ActionSellShort {
		side = connector.OrderSideSell
	}

	e.Logger.Info("Executing options %s order: %s %s @ %s on %s",
		action.ActionType, action.Quantity.StringFixed(4),
		action.Contract.Pair.Symbol(), action.Price.StringFixed(2), action.Exchange,
	)

	var (
		resp    *connector.OrderResponse
		execErr error
	)

	if action.Price.IsZero() {
		resp, execErr = exec.PlaceMarketOrder(action.Contract.Pair, side, action.Quantity)
	} else {
		resp, execErr = exec.PlaceLimitOrder(action.Contract.Pair, side, action.Quantity, action.Price)
	}

	if execErr != nil {
		return "", fmt.Errorf("failed to place options order on %s: %w", action.Exchange, execErr)
	}

	e.store.AddOrder(connector.Order{
		Pair:      action.Contract.Pair,
		ID:        resp.OrderID,
		Side:      side,
		Quantity:  action.Quantity,
		Price:     action.Price,
		Status:    connector.OrderStatusPending,
		Type:      connector.OrderTypeLimit,
		CreatedAt: e.TimeProvider.Now(),
		UpdatedAt: e.TimeProvider.Now(),
	})

	e.Logger.Info("Options order placed: %s (contract: %s, side: %s)",
		resp.OrderID, action.Contract.Pair.Symbol(), side)
	return resp.OrderID, nil
}

var _ optionsTypes.SignalExecutor = (*executor)(nil)
