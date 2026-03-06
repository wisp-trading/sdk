package activity

import (
	"context"
	"time"

	perpTypes "github.com/wisp-trading/sdk/pkg/markets/perp/types"
	spotTypes "github.com/wisp-trading/sdk/pkg/markets/spot/types"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	wispActivity "github.com/wisp-trading/sdk/pkg/types/wisp/activity"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

// trades aggregates trade history across spot and perp domains.
type trades struct {
	spot spotTypes.SpotTrades
	perp perpTypes.PerpTrades
}

func NewTrades(
	spot spotTypes.SpotTrades,
	perp perpTypes.PerpTrades,
) wispActivity.Trades {
	return &trades{spot: spot, perp: perp}
}

func (t *trades) GetAllTrades(_ context.Context) []connector.Trade {
	all := t.spot.GetAllTrades()
	return append(all, t.perp.GetAllTrades()...)
}

func (t *trades) GetTradesByExchange(_ context.Context, exchange connector.ExchangeName) []connector.Trade {
	all := t.spot.GetTradesByExchange(exchange)
	return append(all, t.perp.GetTradesByExchange(exchange)...)
}

func (t *trades) GetTradesByPair(_ context.Context, pair portfolio.Pair) []connector.Trade {
	all := t.spot.GetTradesByPair(pair)
	return append(all, t.perp.GetTradesByPair(pair)...)
}

func (t *trades) GetTradesSince(_ context.Context, since time.Time) []connector.Trade {
	all := t.spot.GetTradesSince(since)
	return append(all, t.perp.GetTradesSince(since)...)
}

func (t *trades) GetTradeByID(_ context.Context, tradeID string) *connector.Trade {
	if tr := t.spot.GetTradeByID(tradeID); tr != nil {
		return tr
	}
	return t.perp.GetTradeByID(tradeID)
}

func (t *trades) GetTradeCount(_ context.Context) int {
	return t.spot.GetTradeCount() + t.perp.GetTradeCount()
}

func (t *trades) GetTotalVolume(_ context.Context, pair portfolio.Pair) numerical.Decimal {
	return t.spot.GetTotalVolume(pair).Add(t.perp.GetTotalVolume(pair))
}

var _ wispActivity.Trades = (*trades)(nil)
