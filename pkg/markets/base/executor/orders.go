package executor

import (
	"fmt"

	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	"github.com/wisp-trading/sdk/pkg/types/strategy"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

// OrderBook is the store surface for recording placed orders and fills.
type OrderBook interface {
	AddTrade(trade connector.Trade)
	AddOrder(order connector.Order)
	UpdateOrderStatus(orderID string, status connector.OrderStatus) error
}

// SideFromAction maps strategy action types to order side.
func SideFromAction(actionType strategy.ActionType) connector.OrderSide {
	switch actionType {
	case strategy.ActionSell, strategy.ActionSellShort:
		return connector.OrderSideSell
	default:
		return connector.OrderSideBuy
	}
}

// RecordTradeFill stores a fill and marks the related order filled when possible.
func (b *Base) RecordTradeFill(store OrderBook, trade connector.Trade, domain string) error {
	store.AddTrade(trade)

	orderID := trade.OrderID
	if orderID == "" {
		orderID = trade.ID
	}

	if err := store.UpdateOrderStatus(orderID, connector.OrderStatusFilled); err != nil {
		b.Logger.Debug("Could not mark %s order %s filled: %v", domain, orderID, err)
	}

	b.Logger.Info("%s trade recorded: %s (order: %s, pair: %s)",
		domain, trade.ID, orderID, trade.Pair.Symbol())
	return nil
}

// PlaceLimitAndRecord places a limit order and records it on the store.
// domain is used only for log/error messages ("spot", "perp", …).
func (b *Base) PlaceLimitAndRecord(
	store OrderBook,
	exchange connector.ExchangeName,
	pair portfolio.Pair,
	side connector.OrderSide,
	quantity, price numerical.Decimal,
	domain string,
) (string, error) {
	exec, err := b.GetOrderExecutor(exchange)
	if err != nil {
		return "", err
	}

	resp, err := exec.PlaceLimitOrder(pair, side, quantity, price)
	if err != nil {
		return "", fmt.Errorf("failed to place %s order on %s: %w", domain, exchange, err)
	}

	now := b.TimeProvider.Now()
	store.AddOrder(connector.Order{
		Pair:      pair,
		ID:        resp.OrderID,
		Side:      side,
		Quantity:  quantity,
		Price:     price,
		Status:    connector.OrderStatusPending,
		Type:      connector.OrderTypeLimit,
		CreatedAt: now,
		UpdatedAt: now,
	})

	b.Logger.Info("%s order placed: %s (pair: %s, side: %s)",
		domain, resp.OrderID, pair.Symbol(), side)
	return resp.OrderID, nil
}
