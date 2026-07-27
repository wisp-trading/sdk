package types

import (
	baseTypes "github.com/wisp-trading/sdk/pkg/markets/base/types"
)

// SpotAction is a pair order action for the spot domain (embeds shared PairAction).
// Prefer constructing via the signal builder — MarketType is stamped there.
type SpotAction = baseTypes.PairAction
