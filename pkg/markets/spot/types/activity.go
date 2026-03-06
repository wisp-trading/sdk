package types

import storeActivity "github.com/wisp-trading/sdk/pkg/markets/base/types/stores/activity"

// SpotPositions is the spot-domain-typed positions store.
// Distinct type prevents fx clashing with perp/prediction stores.
type SpotPositions interface {
	storeActivity.Positions
}

// SpotTrades is the spot-domain-typed trades store.
type SpotTrades interface {
	storeActivity.Trades
}
