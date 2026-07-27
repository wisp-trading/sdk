package executor

import (
	"fmt"

	baseExecutor "github.com/wisp-trading/sdk/pkg/markets/base/executor"
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
	baseExecutor.Base
	store perpTypes.MarketStore
}

// NewExecutor creates a new perp market executor.
func NewExecutor(
	connectors registry.ConnectorRegistry,
	store perpTypes.MarketStore,
	logger logging.ApplicationLogger,
	timeProvider temporal.TimeProvider,
) perpTypes.SignalExecutor {
	return &executor{
		Base: baseExecutor.Base{
			Connectors:   connectors,
			Logger:       logger,
			TimeProvider: timeProvider,
		},
		store: store,
	}
}

// ExecutePerpSignal executes all actions in a perp signal.
func (e *executor) ExecutePerpSignal(
	signal perpTypes.PerpSignal,
	ctx *execution.ExecutionContext,
	result *execution.ExecutionResult,
) error {
	for i, action := range signal.GetActions() {
		if err := action.Validate(); err != nil {
			return fmt.Errorf("perp action %d invalid: %w", i, err)
		}

		orderID, err := e.executeAction(&action)
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
	return e.RecordTradeFill(e.store, trade, "perp")
}

func (e *executor) executeAction(action *perpTypes.PerpAction) (string, error) {
	switch action.ActionType {
	case strategy.ActionHold:
		e.Logger.Info("Holding perp position for %s", action.Pair.Symbol())
		return "", nil
	case strategy.ActionClose:
		e.Logger.Info("Close perp action noted for %s", action.Pair.Symbol())
		return "", nil
	}

	side := baseExecutor.SideFromAction(action.ActionType)

	e.Logger.Info("Executing perp %s order: %s %s @ %s (leverage: %s) on %s",
		action.ActionType, action.Quantity.StringFixed(4), action.Pair.Symbol(),
		action.Price.StringFixed(2), action.Leverage.StringFixed(1), action.Exchange,
	)

	if !action.Leverage.IsZero() {
		if conn, ok := e.Connectors.Connector(action.Exchange); ok {
			if perpConnector, isPerp := conn.(perpConn.Connector); isPerp {
				if err := e.setLeverage(perpConnector, action); err != nil {
					e.Logger.Warn("Could not set leverage for %s on %s: %v",
						action.Pair.Symbol(), action.Exchange, err)
				}
			}
		}
	}

	return e.PlaceLimitAndRecord(
		e.store,
		action.Exchange,
		action.Pair,
		side,
		action.Quantity,
		action.Price,
		"perp",
	)
}

func (e *executor) setLeverage(conn perpConn.Connector, action *perpTypes.PerpAction) error {
	symbol := conn.GetPerpSymbol(action.Pair)
	if symbol == "" {
		return fmt.Errorf("could not resolve perp symbol for %s", action.Pair.Symbol())
	}
	e.Logger.Debug("Leverage %s requested for %s (%s) on %s",
		action.Leverage.StringFixed(1), action.Pair.Symbol(), symbol, action.Exchange)
	return nil
}

var _ perpTypes.SignalExecutor = (*executor)(nil)
