package batch

import (
	"github.com/wisp-trading/sdk/pkg/markets/base/types"
	perpTypes "github.com/wisp-trading/sdk/pkg/markets/perp/types"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
)

// watchlistAdapter adapts PerpWatchlist to the data.MarketWatchlist interface
// expected by the base batch ingestor.
type watchlistAdapter struct {
	inner perpTypes.PerpWatchlist
}

func newWatchlistAdapter(w perpTypes.PerpWatchlist) types.MarketWatchlist {
	return &watchlistAdapter{inner: w}
}

func (a *watchlistAdapter) RequirePair(exchange connector.ExchangeName, pair portfolio.Pair) {
	a.inner.RequirePair(exchange, pair)
}

func (a *watchlistAdapter) ReleasePair(exchange connector.ExchangeName, pair portfolio.Pair) {
	a.inner.ReleasePair(exchange, pair)
}

func (a *watchlistAdapter) GetRequiredPairs(exchange connector.ExchangeName) []portfolio.Pair {
	return a.inner.GetRequiredPairs(exchange)
}

func (a *watchlistAdapter) Subscribe(exchange connector.ExchangeName) chan types.MarketWatchEvent {
	perpCh := a.inner.Subscribe(exchange)
	out := make(chan types.MarketWatchEvent, 128)

	go func() {
		defer close(out)
		for ev := range perpCh {
			var evType types.MarketWatchEventType
			if ev.Type == perpTypes.PerpPairAdded {
				evType = types.PairAdded
			} else {
				evType = types.PairRemoved
			}
			out <- types.MarketWatchEvent{
				Requirement: types.PairRequirement{
					Exchange: ev.Exchange,
					Pair:     ev.Pair,
				},
				Type: evType,
			}
		}
	}()

	return out
}

func (a *watchlistAdapter) Unsubscribe(exchange connector.ExchangeName) {
	a.inner.Unsubscribe(exchange)
}

var _ types.MarketWatchlist = (*watchlistAdapter)(nil)
