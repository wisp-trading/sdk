package activity

import (
	"context"

	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

// PositionPNL holds PNL for a single open position.
type PositionPNL struct {
	Pair       portfolio.Pair
	Exchange   connector.ExchangeName
	Realized   numerical.Decimal
	Unrealized numerical.Decimal
	Fees       numerical.Decimal
}

// SpotPNL provides PNL for the spot domain.
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

// PNL aggregates PNL across all market domains.
// Each domain owns its own PNL type — this is purely additive.
type PNL interface {
	// Spot returns PNL for the spot domain.
	Spot() SpotPNL

	// TotalRealized returns realized PNL summed across all domains.
	TotalRealized(ctx context.Context) numerical.Decimal

	// TotalUnrealized returns unrealized PNL summed across all domains.
	TotalUnrealized(ctx context.Context) numerical.Decimal

	// TotalFees returns fees summed across all domains.
	TotalFees(ctx context.Context) numerical.Decimal
}
