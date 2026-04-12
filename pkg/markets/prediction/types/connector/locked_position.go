package connector

import "math/big"

// LockedPosition represents a set of CTF ERC-1155 conditional tokens currently
// held on-chain by the signing EOA, grouped by condition ID.
//
// Each LockedPosition is self-contained: pass Market directly to MergePositions
// to recover the locked collateral back to USDC.
type LockedPosition struct {
	// Market holds the condition ID (MarketID) and the outcome token IDs
	// (Outcomes[i].OutcomeID) — everything MergePositions needs.
	Market Market

	// OutcomeBalances is the raw 6-decimal on-chain balance per outcome token.
	// Key is the OutcomeID of the corresponding Market.Outcomes entry.
	OutcomeBalances map[OutcomeID]*big.Int

	// MergeableAmount is min(OutcomeBalances) in raw 6-decimal USDC units — the
	// maximum amount that can be passed to MergePositions right now.
	MergeableAmount *big.Int
}

// IsRecoverable returns true when the position has a non-zero mergeable amount,
// i.e. both YES and NO tokens are held and a merge call would succeed.
func (p LockedPosition) IsRecoverable() bool {
	return p.MergeableAmount != nil && p.MergeableAmount.Sign() > 0
}
