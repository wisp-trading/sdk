package types

import (
	"context"

	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

// PredictionPositionPNL holds PNL for a single prediction market order.
// Unrealized is current implied value based on the orderbook mid-price.
// Realized is non-zero once the market resolves and winnings are redeemed.
type PredictionPositionPNL struct {
	Order      PredictionOrder
	Realized   numerical.Decimal
	Unrealized numerical.Decimal
}

// PredictionPNL provides PNL for the prediction market domain.
type PredictionPNL interface {
	// Positions returns PNL per prediction market order.
	Positions(ctx context.Context) []PredictionPositionPNL

	// Realized returns total realized PNL from resolved and redeemed markets.
	Realized(ctx context.Context) numerical.Decimal

	// Unrealized returns total implied PNL from open prediction positions.
	Unrealized(ctx context.Context) numerical.Decimal

	// Fees returns total fees paid on prediction market orders.
	Fees(ctx context.Context) numerical.Decimal
}
