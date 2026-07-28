package executor

import (
	"fmt"

	baseExecutor "github.com/wisp-trading/sdk/pkg/markets/base/executor"
	onchainTypes "github.com/wisp-trading/sdk/pkg/markets/onchain/types"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/execution"
	"github.com/wisp-trading/sdk/pkg/types/logging"
	"github.com/wisp-trading/sdk/pkg/types/registry"
	"github.com/wisp-trading/sdk/pkg/types/strategy"
	"github.com/wisp-trading/sdk/pkg/types/temporal"
)

type executor struct {
	baseExecutor.Base
	store onchainTypes.MarketStore
}

func NewExecutor(
	connectors registry.ConnectorRegistry,
	store onchainTypes.MarketStore,
	logger logging.ApplicationLogger,
	timeProvider temporal.TimeProvider,
) onchainTypes.SignalExecutor {
	return &executor{
		Base: baseExecutor.Base{
			Connectors:   connectors,
			Logger:       logger,
			TimeProvider: timeProvider,
		},
		store: store,
	}
}

func (e *executor) ExecuteOnchainSignal(
	signal onchainTypes.OnchainSignal,
	ctx *execution.ExecutionContext,
	result *execution.ExecutionResult,
) error {
	_ = ctx
	for i, action := range signal.GetActions() {
		if err := action.Validate(); err != nil {
			return fmt.Errorf("onchain action %d invalid: %w", i, err)
		}

		orderID, err := e.executeAction(&action)
		if err != nil {
			return fmt.Errorf("onchain action %d failed: %w", i, err)
		}

		if orderID != "" {
			result.OrderIDs = append(result.OrderIDs, orderID)
		}
	}
	return nil
}

func (e *executor) HandleTrade(trade connector.Trade) error {
	return e.RecordTradeFill(e.store, trade, "onchain")
}

func (e *executor) executeAction(action *onchainTypes.OnchainAction) (string, error) {
	switch action.ActionType {
	case strategy.ActionHold:
		e.Logger.Info("Holding onchain bag for %s", action.Pair.Symbol())
		return "", nil
	case strategy.ActionClose:
		e.Logger.Info("Close onchain action noted for %s", action.Pair.Symbol())
		return "", nil
	}

	side := baseExecutor.SideFromAction(action.ActionType)

	e.Logger.Info("Executing onchain %s swap: qty=%s pair=%s on %s",
		action.ActionType, action.Quantity.String(), action.Pair.Symbol(), action.Exchange,
	)

	// Zero price = market (AMM exact-in). Non-zero limit rejected by UniV3 pilot connector.
	return e.PlaceOrderAndRecord(
		e.store,
		action.Exchange,
		action.Pair,
		side,
		action.Quantity,
		action.Price,
		"onchain",
	)
}

var _ onchainTypes.SignalExecutor = (*executor)(nil)
