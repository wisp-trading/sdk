package types

import (
	"context"

	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

// PositionPNL holds PNL for a single onchain bag.
type PositionPNL struct {
	Pair       portfolio.Pair
	Exchange   connector.ExchangeName
	Realized   numerical.Decimal
	Unrealized numerical.Decimal
	Fees       numerical.Decimal
}

// OnchainPNL provides PNL calculations for the onchain domain.
type OnchainPNL interface {
	Positions(ctx context.Context) []PositionPNL
	Realized(ctx context.Context) numerical.Decimal
	Unrealized(ctx context.Context) numerical.Decimal
	Fees(ctx context.Context) numerical.Decimal
}
