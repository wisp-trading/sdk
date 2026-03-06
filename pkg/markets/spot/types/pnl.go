package types

import (
	"context"

	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

// PositionPNL holds PNL for a single open spot position.
type PositionPNL struct {
	Pair       portfolio.Pair
	Exchange   connector.ExchangeName
	Realized   numerical.Decimal
	Unrealized numerical.Decimal
	Fees       numerical.Decimal
}

// SpotPNL provides PNL calculations for the spot domain.
// Accessed via wisp.Spot().PNL() — never exposed at the top-level wisp.Activity().
type SpotPNL interface {
	// Positions returns PNL broken down per open spot position.
	Positions(ctx context.Context) []PositionPNL

	// Realized returns total realized PNL across all spot trades, net of fees.
	Realized(ctx context.Context) numerical.Decimal

	// Unrealized returns total unrealized PNL across all open spot positions.
	Unrealized(ctx context.Context) numerical.Decimal

	// Fees returns total fees paid on spot trades.
	Fees(ctx context.Context) numerical.Decimal
}
