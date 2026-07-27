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
	return e.RecordTradeFill(e.store, trade, "spot")
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

	side := baseExecutor.SideFromAction(action.ActionType)

	e.Logger.Info("Executing spot %s order: %s %s @ %s on %s",
		action.ActionType, action.Quantity.StringFixed(4), action.Pair.Symbol(),
		action.Price.StringFixed(2), action.Exchange,
	)

	return e.PlaceLimitAndRecord(
		e.store,
		action.Exchange,
		action.Pair,
		side,
		action.Quantity,
		action.Price,
		"spot",
	)
}

var _ spotTypes.SignalExecutor = (*executor)(nil)
