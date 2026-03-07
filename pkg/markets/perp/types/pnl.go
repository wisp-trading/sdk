package types

import (
	"context"

	perpConn "github.com/wisp-trading/sdk/pkg/types/connector/perp"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

// PerpPositionPNL holds PNL for a single perp position.
// Sourced directly from the connector — the exchange owns realized/unrealized.
type PerpPositionPNL struct {
	Position   perpConn.Position
	Realized   numerical.Decimal
	Unrealized numerical.Decimal
}

// PerpPNL provides PNL for the perp domain.
// Sourced from live connector positions — the exchange owns the numbers.
type PerpPNL interface {
	// Positions returns PNL per open perp position, as reported by the exchange.
	Positions(ctx context.Context) []PerpPositionPNL

	// Realized returns total realized PNL across all perp positions.
	Realized(ctx context.Context) numerical.Decimal

	// Unrealized returns total unrealized PNL across all open perp positions.
	Unrealized(ctx context.Context) numerical.Decimal

	// Fees returns total fees paid on perp trades.
	Fees(ctx context.Context) numerical.Decimal
}
