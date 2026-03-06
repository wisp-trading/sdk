package activity

import (
	"context"

	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

// PNL aggregates PNL totals across all market domains.
type PNL interface {

	// TotalRealized returns realized PNL summed across all domains.
	TotalRealized(ctx context.Context) numerical.Decimal

	// TotalUnrealized returns unrealized PNL summed across all domains.
	TotalUnrealized(ctx context.Context) numerical.Decimal

	// TotalFees returns fees summed across all domains.
	TotalFees(ctx context.Context) numerical.Decimal
}
