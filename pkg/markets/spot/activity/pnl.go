package activity

import (
	"context"

	spotTypes "github.com/wisp-trading/sdk/pkg/markets/spot/types"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

type spotPNL struct {
	store spotTypes.MarketStore
}

func NewSpotPNL(store spotTypes.MarketStore) spotTypes.SpotPNL {
	return &spotPNL{store: store}
}

func (s *spotPNL) Positions(_ context.Context) []spotTypes.PositionPNL {
	trackers := buildTrackers(s.store.GetAllTrades())
	results := make([]spotTypes.PositionPNL, 0, len(trackers))
	for _, t := range trackers {
		unrealized := s.unrealizedForTracker(t)
		results = append(results, spotTypes.PositionPNL{
			Pair:       t.pair,
			Exchange:   t.exchange,
			Realized:   t.realized,
			Unrealized: unrealized,
			Fees:       t.fees,
		})
	}
	return results
}

func (s *spotPNL) Realized(_ context.Context) numerical.Decimal {
	trades := s.store.GetAllTrades()
	realized := numerical.Zero()
	for _, t := range buildTrackers(trades) {
		realized = realized.Add(t.realized)
	}
	return realized.Sub(sumTradeFees(trades))
}

func (s *spotPNL) Unrealized(_ context.Context) numerical.Decimal {
	total := numerical.Zero()
	for _, t := range buildTrackers(s.store.GetAllTrades()) {
		total = total.Add(s.unrealizedForTracker(t))
	}
	return total
}

func (s *spotPNL) Fees(_ context.Context) numerical.Decimal {
	return sumTradeFees(s.store.GetAllTrades())
}

func (s *spotPNL) unrealizedForTracker(t *positionTracker) numerical.Decimal {
	if t.size.IsZero() {
		return numerical.Zero()
	}
	price := s.store.GetPairPrice(t.pair, t.exchange)
	if price == nil {
		return numerical.Zero()
	}
	if t.size.IsPositive() {
		return price.Price.Sub(t.avgEntry).Mul(t.size)
	}
	return t.avgEntry.Sub(price.Price).Mul(t.size.Abs())
}

var _ spotTypes.SpotPNL = (*spotPNL)(nil)

// ============================================================
// Shared position tracking helpers — used by spot (perp owns its own via connector)
// ============================================================

type positionTracker struct {
	pair     portfolio.Pair
	exchange connector.ExchangeName
	size     numerical.Decimal
	avgEntry numerical.Decimal
	realized numerical.Decimal
	fees     numerical.Decimal
}

func buildTrackers(trades []connector.Trade) map[string]*positionTracker {
	open := make(map[string]*positionTracker)
	for _, trade := range trades {
		key := string(trade.Exchange) + ":" + trade.Pair.Symbol()
		t, exists := open[key]
		if !exists {
			t = &positionTracker{
				pair:     trade.Pair,
				exchange: trade.Exchange,
				size:     numerical.Zero(),
				avgEntry: numerical.Zero(),
			}
			open[key] = t
		}
		t.fees = t.fees.Add(trade.Fee)
		t.realized = t.realized.Add(applyTrade(t, trade))
	}
	return open
}

func applyTrade(t *positionTracker, trade connector.Trade) numerical.Decimal {
	qty := trade.Quantity
	price := trade.Price

	signedQty := qty
	if trade.Side == connector.OrderSideSell {
		signedQty = qty.Neg()
	}

	if t.size.IsZero() {
		t.size = signedQty
		t.avgEntry = price
		return numerical.Zero()
	}

	sameDir := (t.size.IsPositive() && signedQty.IsPositive()) ||
		(t.size.IsNegative() && signedQty.IsNegative())

	if sameDir {
		totalValue := t.avgEntry.Mul(t.size.Abs()).Add(price.Mul(qty))
		t.size = t.size.Add(signedQty)
		if !t.size.IsZero() {
			t.avgEntry = totalValue.Div(t.size.Abs())
		}
		return numerical.Zero()
	}

	closeQty := qty
	if closeQty.GreaterThan(t.size.Abs()) {
		closeQty = t.size.Abs()
	}

	wasPositive := t.size.IsPositive()
	var realized numerical.Decimal
	if wasPositive {
		realized = price.Sub(t.avgEntry).Mul(closeQty)
	} else {
		realized = t.avgEntry.Sub(price).Mul(closeQty)
	}

	newSize := t.size.Add(signedQty)
	if !newSize.IsZero() &&
		((wasPositive && newSize.IsNegative()) || (!wasPositive && newSize.IsPositive())) {
		t.avgEntry = price
	}
	t.size = newSize
	return realized
}

func sumTradeFees(trades []connector.Trade) numerical.Decimal {
	total := numerical.Zero()
	for _, t := range trades {
		total = total.Add(t.Fee)
	}
	return total
}
