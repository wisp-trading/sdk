package trade

import (
	"fmt"

	"github.com/backtesting-org/kronos-sdk/pkg/types/connector"
	"github.com/backtesting-org/kronos-sdk/pkg/types/kronos/numerical"
	"github.com/backtesting-org/kronos-sdk/pkg/types/logging"
	"github.com/backtesting-org/kronos-sdk/pkg/types/portfolio"
)

// TradeService provides trade execution methods.
// This is only available in KronosExecutor, not in the base Kronos.
type TradeService struct {
	tradingLogger logging.TradingLogger
	// TODO: Add reference to actual trade executor when implemented
}

// NewTradeService creates a new TradeService
func NewTradeService(tradingLogger logging.TradingLogger) *TradeService {
	return &TradeService{
		tradingLogger: tradingLogger,
	}
}

// OrderType represents the type of order
type OrderType string

const (
	OrderTypeMarket OrderType = "market"
	OrderTypeLimit  OrderType = "limit"
)

// TradeOptions configures trade execution
type TradeOptions struct {
	OrderType   OrderType         // Market or Limit
	Price       numerical.Decimal // Required for limit orders
	TimeInForce string            // GTC, IOC, FOK, etc.
	ReduceOnly  bool              // Only reduce position
}

// TradeResult holds the result of a trade execution
type TradeResult struct {
	OrderID  string
	Asset    portfolio.Asset
	Exchange connector.ExchangeName
	Side     connector.OrderSide
	Quantity numerical.Decimal
	Price    numerical.Decimal // Filled price
	Status   string
	Message  string
}

// Buy executes a buy order for an asset.
// Default is market order unless price is specified in options.
func (s *TradeService) Buy(asset portfolio.Asset, exchange connector.ExchangeName, quantity numerical.Decimal, opts ...TradeOptions) (*TradeResult, error) {
	return s.executeTrade(asset, exchange, connector.OrderSideBuy, quantity, opts...)
}

// Sell executes a sell order for an asset.
// Default is market order unless price is specified in options.
func (s *TradeService) Sell(asset portfolio.Asset, exchange connector.ExchangeName, quantity numerical.Decimal, opts ...TradeOptions) (*TradeResult, error) {
	return s.executeTrade(asset, exchange, connector.OrderSideSell, quantity, opts...)
}

// Short opens a short position (sell to open).
func (s *TradeService) Short(asset portfolio.Asset, exchange connector.ExchangeName, quantity numerical.Decimal, opts ...TradeOptions) (*TradeResult, error) {
	// In perpetual futures, this is typically a sell order
	return s.executeTrade(asset, exchange, connector.OrderSideSell, quantity, opts...)
}

// CloseShort closes a short position (buy to close).
func (s *TradeService) CloseShort(asset portfolio.Asset, exchange connector.ExchangeName, quantity numerical.Decimal, opts ...TradeOptions) (*TradeResult, error) {
	// In perpetual futures, this is typically a buy order to close
	return s.executeTrade(asset, exchange, connector.OrderSideBuy, quantity, opts...)
}

// executeTrade is the internal method that handles trade execution
func (s *TradeService) executeTrade(
	asset portfolio.Asset,
	exchange connector.ExchangeName,
	side connector.OrderSide,
	quantity numerical.Decimal,
	opts ...TradeOptions,
) (*TradeResult, error) {
	// Parse options
	options := s.parseOptions(opts...)

	// Validate inputs
	if quantity.LessThanOrEqual(numerical.Zero()) {
		return nil, fmt.Errorf("quantity must be greater than zero")
	}

	if options.OrderType == OrderTypeLimit && options.Price.IsZero() {
		return nil, fmt.Errorf("price must be specified for limit orders")
	}

	// TODO: Call the actual trade executor here
	// For now, we return a placeholder result

	// Log the trade attempt using trading logger
	orderTypeStr := string(options.OrderType)
	sideStr := string(side)

	// Placeholder order ID
	orderID := fmt.Sprintf("ORDER_%s_%s", asset.Symbol(), exchange)

	result := &TradeResult{
		OrderID:  orderID,
		Asset:    asset,
		Exchange: exchange,
		Side:     side,
		Quantity: quantity,
		Price:    options.Price, // This would be the filled price from the executor
		Status:   "pending",
		Message:  fmt.Sprintf("%s %s order placed", orderTypeStr, sideStr),
	}

	return result, nil
}

// parseOptions extracts options with defaults
func (s *TradeService) parseOptions(opts ...TradeOptions) TradeOptions {
	if len(opts) > 0 {
		options := opts[0]
		if options.OrderType == "" {
			options.OrderType = OrderTypeMarket
		}
		if options.TimeInForce == "" {
			options.TimeInForce = "GTC"
		}
		return options
	}
	return TradeOptions{
		OrderType:   OrderTypeMarket,
		TimeInForce: "GTC",
	}
}
