package perp

import (
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

// AssetBalance represents perpetual futures account information with margin and PnL.
type AssetBalance struct {
	connector.AssetBalance
	UsedMargin    numerical.Decimal `json:"used_margin"`    // Margin used by open positions
	UnrealizedPnL numerical.Decimal `json:"unrealized_pnl"` // Unrealized profit/loss from positions
}
