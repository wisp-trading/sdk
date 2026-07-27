package activity

import (
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

// PositionTracker accumulates size/avg entry/realized/fees from a trade stream.
// Used by spot PNL (perp prefers live exchange positions).
type PositionTracker struct {
	Pair     portfolio.Pair
	Exchange connector.ExchangeName
	Size     numerical.Decimal
	AvgEntry numerical.Decimal
	Realized numerical.Decimal
	Fees     numerical.Decimal
}

// BuildTrackers folds trades into per exchange+pair trackers.
func BuildTrackers(trades []connector.Trade) map[string]*PositionTracker {
	open := make(map[string]*PositionTracker)
	for _, trade := range trades {
		key := string(trade.Exchange) + ":" + trade.Pair.Symbol()
		t, exists := open[key]
		if !exists {
			t = &PositionTracker{
				Pair:     trade.Pair,
				Exchange: trade.Exchange,
				Size:     numerical.Zero(),
				AvgEntry: numerical.Zero(),
			}
			open[key] = t
		}
		t.Fees = t.Fees.Add(trade.Fee)
		t.Realized = t.Realized.Add(ApplyTrade(t, trade))
	}
	return open
}

// ApplyTrade updates tracker size/avg entry from one trade; returns realized for this fill.
func ApplyTrade(t *PositionTracker, trade connector.Trade) numerical.Decimal {
	qty := trade.Quantity
	price := trade.Price

	signedQty := qty
	if trade.Side == connector.OrderSideSell {
		signedQty = qty.Neg()
	}

	if t.Size.IsZero() {
		t.Size = signedQty
		t.AvgEntry = price
		return numerical.Zero()
	}

	sameDir := (t.Size.IsPositive() && signedQty.IsPositive()) ||
		(t.Size.IsNegative() && signedQty.IsNegative())

	if sameDir {
		totalValue := t.AvgEntry.Mul(t.Size.Abs()).Add(price.Mul(qty))
		t.Size = t.Size.Add(signedQty)
		if !t.Size.IsZero() {
			t.AvgEntry = totalValue.Div(t.Size.Abs())
		}
		return numerical.Zero()
	}

	closeQty := qty
	if closeQty.GreaterThan(t.Size.Abs()) {
		closeQty = t.Size.Abs()
	}

	wasPositive := t.Size.IsPositive()
	var realized numerical.Decimal
	if wasPositive {
		realized = price.Sub(t.AvgEntry).Mul(closeQty)
	} else {
		realized = t.AvgEntry.Sub(price).Mul(closeQty)
	}

	newSize := t.Size.Add(signedQty)
	if !newSize.IsZero() &&
		((wasPositive && newSize.IsNegative()) || (!wasPositive && newSize.IsPositive())) {
		t.AvgEntry = price
	}
	t.Size = newSize
	return realized
}

// UnrealizedFromMark computes unrealized PnL for a tracker given a mark price.
func UnrealizedFromMark(t *PositionTracker, mark numerical.Decimal) numerical.Decimal {
	if t.Size.IsZero() {
		return numerical.Zero()
	}
	if t.Size.IsPositive() {
		return mark.Sub(t.AvgEntry).Mul(t.Size)
	}
	return t.AvgEntry.Sub(mark).Mul(t.Size.Abs())
}
