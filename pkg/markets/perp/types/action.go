package types

import (
	baseTypes "github.com/wisp-trading/sdk/pkg/markets/base/types"
)

// PerpAction is a pair order action for the perp domain (shared PairAction + Leverage field).
// Prefer constructing via the signal builder — MarketType and Leverage are stamped there.
type PerpAction = baseTypes.PairAction
