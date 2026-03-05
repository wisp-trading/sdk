package batch

import (
	perpTypes "github.com/wisp-trading/sdk/pkg/markets/perp/types"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/data"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
)

// watchlistAdapter adapts PerpWatchlist to the data.MarketWatchlist interface
// expected by the base batch ingestor.
type watchlistAdapter struct {
	inner perpTypes.PerpWatchlist
}

func newWatchlistAdapter(w perpTypes.PerpWatchlist) data.MarketWatchlist {
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

func (a *watchlistAdapter) Subscribe(exchange connector.ExchangeName) chan data.MarketWatchEvent {
	perpCh := a.inner.Subscribe(exchange)
	out := make(chan data.MarketWatchEvent, 128)

	go func() {
		defer close(out)
		for ev := range perpCh {
			var evType data.MarketWatchEventType
			if ev.Type == perpTypes.PerpPairAdded {
				evType = data.PairAdded
			} else {
				evType = data.PairRemoved
			}
			out <- data.MarketWatchEvent{
				Requirement: data.PairRequirement{
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

var _ data.MarketWatchlist = (*watchlistAdapter)(nil)
