package activity

import "context"

// Positions provides read-only access to order and trade data for this strategy instance.
type Positions interface {
	GetOrderCount(ctx context.Context) int64
}
