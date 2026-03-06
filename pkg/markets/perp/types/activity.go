package types

import storeActivity "github.com/wisp-trading/sdk/pkg/markets/base/types/stores/activity"

// PerpPositions is the perp-domain-typed positions store.
type PerpPositions interface {
	storeActivity.Positions
}

// PerpTrades is the perp-domain-typed trades store.
type PerpTrades interface {
	storeActivity.Trades
}
