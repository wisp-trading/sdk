package activity

import (
	"context"

	storeTypes "github.com/wisp-trading/sdk/pkg/markets/base/types/stores/activity"
	spotTypes "github.com/wisp-trading/sdk/pkg/markets/spot/types"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	wispActivity "github.com/wisp-trading/sdk/pkg/types/wisp/activity"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

// spotPNL calculates PNL from the spot trade store and current prices in the spot market store.
type spotPNL struct {
	trades storeTypes.Trades
	store  spotTypes.MarketStore
}

func NewSpotPNL(trades storeTypes.Trades, store spotTypes.MarketStore) wispActivity.SpotPNL {
	return &spotPNL{trades: trades, store: store}
}

func (s *spotPNL) Positions(_ context.Context) []wispActivity.PositionPNL {
	trackers := buildTrackers(s.trades.GetAllTrades())
	results := make([]wispActivity.PositionPNL, 0, len(trackers))
	for _, t := range trackers {
		unrealized := s.unrealizedForTracker(t)
		results = append(results, wispActivity.PositionPNL{
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
	trades := s.trades.GetAllTrades()
	realized := numerical.Zero()
	for _, t := range buildTrackers(trades) {
		realized = realized.Add(t.realized)
	}
	return realized.Sub(sumTradeFees(trades))
}

func (s *spotPNL) Unrealized(_ context.Context) numerical.Decimal {
	total := numerical.Zero()
	for _, t := range buildTrackers(s.trades.GetAllTrades()) {
		total = total.Add(s.unrealizedForTracker(t))
	}
	return total
}

func (s *spotPNL) Fees(_ context.Context) numerical.Decimal {
	return sumTradeFees(s.trades.GetAllTrades())
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

var _ wispActivity.SpotPNL = (*spotPNL)(nil)

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
