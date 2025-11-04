package kronos

import (
	"fmt"

	"github.com/backtesting-org/kronos-sdk/pkg/types/connector"
	"github.com/backtesting-org/kronos-sdk/pkg/types/logging"
	"github.com/backtesting-org/kronos-sdk/pkg/types/portfolio"
	"github.com/shopspring/decimal"
)

// TradeService provides trade execution methods.
// This is only available in KronosExecutor, not in the base Kronos.
type TradeService struct {
	logger logging.ApplicationLogger
}

// OrderType represents the type of order
type OrderType string

const (
	OrderTypeMarket OrderType = "market"
	OrderTypeLimit  OrderType = "limit"
)

// TradeOptions configures trade execution
type TradeOptions struct {
	OrderType   OrderType // Market or Limit
	Price       decimal.Decimal
	TimeInForce string // GTC, IOC, FOK, etc.
	ReduceOnly  bool   // Only reduce position
}

// TradeResult holds the result of a trade execution
type TradeResult struct {
	OrderID  string
	Asset    portfolio.Asset
	Exchange connector.ExchangeName
	Side     connector.OrderSide
	Quantity decimal.Decimal
	Price    decimal.Decimal // Filled price
	Status   string
	Message  string
}

// Buy executes a buy order for an asset.
// Default is market order unless price is specified in options.
func (s *TradeService) Buy(asset portfolio.Asset, exchange connector.ExchangeName, quantity decimal.Decimal, opts ...TradeOptions) (*TradeResult, error) {
	return s.executeTrade(asset, exchange, connector.OrderSideBuy, quantity, opts...)
}

// Sell executes a sell order for an asset.
// Default is market order unless price is specified in options.
func (s *TradeService) Sell(asset portfolio.Asset, exchange connector.ExchangeName, quantity decimal.Decimal, opts ...TradeOptions) (*TradeResult, error) {
	return s.executeTrade(asset, exchange, connector.OrderSideSell, quantity, opts...)
}

// BuyShort opens a short position (sell to open).
func (s *TradeService) BuyShort(asset portfolio.Asset, exchange connector.ExchangeName, quantity decimal.Decimal, opts ...TradeOptions) (*TradeResult, error) {
	// In perpetual futures, this is typically a sell order
	return s.executeTrade(asset, exchange, connector.OrderSideSell, quantity, opts...)
}

// SellShort closes a short position (buy to close).
func (s *TradeService) SellShort(asset portfolio.Asset, exchange connector.ExchangeName, quantity decimal.Decimal, opts ...TradeOptions) (*TradeResult, error) {
	// In perpetual futures, this is typically a buy order to close
	return s.executeTrade(asset, exchange, connector.OrderSideBuy, quantity, opts...)
}

// executeTrade is the internal method that handles trade execution
func (s *TradeService) executeTrade(
	asset portfolio.Asset,
	exchange connector.ExchangeName,
	side connector.OrderSide,
	quantity decimal.Decimal,
	opts ...TradeOptions,
) (*TradeResult, error) {
	// Parse options
	options := s.parseOptions(opts...)

	// Validate inputs
	if quantity.LessThanOrEqual(decimal.Zero) {
		return nil, fmt.Errorf("quantity must be greater than zero")
	}

	// Log the trade attempt
	s.logger.Info("Executing trade",
		"asset", asset.Symbol(),
		"exchange", exchange,
		"side", side,
		"quantity", quantity.String(),
		"type", options.OrderType,
	)

	// In a real implementation, this would call the actual trade executor
	// For now, we return a placeholder result
	result := &TradeResult{
		OrderID:  "ORDER_" + asset.Symbol(), // Placeholder
		Asset:    asset,
		Exchange: exchange,
		Side:     side,
		Quantity: quantity,
		Price:    options.Price,
		Status:   "pending",
		Message:  "Trade execution not implemented - this is a placeholder",
	}

	s.logger.Info("Trade executed",
		"orderID", result.OrderID,
		"status", result.Status,
	)

	return result, nil
}

// CancelOrder cancels an open order
func (s *TradeService) CancelOrder(orderID string, exchange connector.ExchangeName) error {
	s.logger.Info("Canceling order", "orderID", orderID, "exchange", exchange)

	// In a real implementation, this would call the connector to cancel the order
	// For now, just log it

	return nil
}

// GetOpenOrders retrieves all open orders
func (s *TradeService) GetOpenOrders(asset portfolio.Asset, exchange connector.ExchangeName) ([]TradeResult, error) {
	s.logger.Debug("Fetching open orders", "asset", asset.Symbol(), "exchange", exchange)

	// In a real implementation, this would fetch from the connector
	// For now, return empty slice

	return []TradeResult{}, nil
}

// parseOptions extracts options with defaults
func (s *TradeService) parseOptions(opts ...TradeOptions) TradeOptions {
	if len(opts) > 0 {
		options := opts[0]
		// If price is set, default to limit order
		if !options.Price.IsZero() && options.OrderType == "" {
			options.OrderType = OrderTypeLimit
		}
		// Otherwise default to market
		if options.OrderType == "" {
			options.OrderType = OrderTypeMarket
		}
		return options
	}
	return TradeOptions{
		OrderType: OrderTypeMarket,
	}
}
