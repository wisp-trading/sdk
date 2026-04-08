package executor

import (
	"fmt"

	optionsTypes "github.com/wisp-trading/sdk/pkg/markets/options/types"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	optionsconnector "github.com/wisp-trading/sdk/pkg/types/connector/options"
	"github.com/wisp-trading/sdk/pkg/types/logging"
	registryTypes "github.com/wisp-trading/sdk/pkg/types/registry"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

type executor struct {
	connectorRegistry registryTypes.ConnectorRegistry
	logger            logging.ApplicationLogger
}

// NewExecutor creates a new options executor
func NewExecutor(
	connectorRegistry registryTypes.ConnectorRegistry,
	logger logging.ApplicationLogger,
) optionsTypes.OptionsExecutor {
	return &executor{
		connectorRegistry: connectorRegistry,
		logger:            logger,
	}
}

// PlaceOrder places an order for an options contract
func (e *executor) PlaceOrder(order optionsTypes.OptionOrder) (*connector.OrderResponse, error) {
	if err := e.validateOrder(order); err != nil {
		return nil, fmt.Errorf("invalid option order: %w", err)
	}

	if order.Exchange == "" {
		return nil, fmt.Errorf("order must specify an exchange")
	}

	conn, err := e.getConnectorByExchange(order.Exchange)
	if err != nil {
		return nil, err
	}

	orderExecutor, ok := conn.(connector.OrderExecutor)
	if !ok {
		return nil, fmt.Errorf("connector does not support order execution")
	}

	quantity := numerical.NewFromFloat(order.Quantity)
	var resp *connector.OrderResponse
	var execErr error

	if order.Price > 0 {
		price := numerical.NewFromFloat(order.Price)
		resp, execErr = orderExecutor.PlaceLimitOrder(order.Contract.Pair, order.Side, quantity, price)
	} else {
		resp, execErr = orderExecutor.PlaceMarketOrder(order.Contract.Pair, order.Side, quantity)
	}

	if execErr != nil {
		e.logger.Errorf("failed to place option order: %v", execErr)
		return nil, execErr
	}

	e.logger.Infof("Option order placed: %s (contract: %s/%s, side: %s, quantity: %f)",
		resp.OrderID, order.Contract.Pair.Symbol(), order.Contract.OptionType, order.Side, order.Quantity)

	return resp, nil
}

// CancelOrder cancels an order by ID
func (e *executor) CancelOrder(orderID string, exchange connector.ExchangeName) (*connector.CancelResponse, error) {
	if orderID == "" {
		return nil, fmt.Errorf("order ID cannot be empty")
	}

	conn, err := e.getConnectorByExchange(exchange)
	if err != nil {
		return nil, err
	}

	orderExecutor, ok := conn.(connector.OrderExecutor)
	if !ok {
		return nil, fmt.Errorf("connector does not support order cancellation")
	}

	resp, err := orderExecutor.CancelOrder(orderID)
	if err != nil {
		e.logger.Errorf("failed to cancel order %s: %v", orderID, err)
		return nil, err
	}

	e.logger.Infof("Option order cancelled: %s", orderID)
	return resp, nil
}

func (e *executor) getConnectorByExchange(exchange connector.ExchangeName) (optionsconnector.Connector, error) {
	conn, exists := e.connectorRegistry.Connector(exchange)
	if !exists {
		return nil, fmt.Errorf("no connector found for exchange: %s", exchange)
	}

	optionsConn, ok := conn.(optionsconnector.Connector)
	if !ok {
		return nil, fmt.Errorf("connector %s does not support options", exchange)
	}

	return optionsConn, nil
}

func (e *executor) validateOrder(order optionsTypes.OptionOrder) error {
	if order.Contract.Strike <= 0 {
		return fmt.Errorf("contract strike must be positive")
	}
	if order.Contract.OptionType != "CALL" && order.Contract.OptionType != "PUT" {
		return fmt.Errorf("contract option type must be CALL or PUT")
	}
	if order.Quantity <= 0 {
		return fmt.Errorf("order quantity must be positive")
	}
	if order.Price < 0 {
		return fmt.Errorf("order price cannot be negative")
	}
	return nil
}

var _ optionsTypes.OptionsExecutor = (*executor)(nil)
