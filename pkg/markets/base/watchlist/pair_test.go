package watchlist

import (
	"testing"

	baseTypes "github.com/wisp-trading/sdk/pkg/markets/base/types"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
)

func TestPairWatchlistRequireAndRelease(t *testing.T) {
	wl := NewPairWatchlist()
	ex := connector.ExchangeName("hyperliquid")
	pair := portfolio.NewPair(portfolio.NewAsset("BTC"), portfolio.NewAsset("USD"))

	ch := wl.Subscribe(ex)
	wl.RequirePair(ex, pair)

	select {
	case ev := <-ch:
		if ev.Type != baseTypes.PairAdded {
			t.Fatalf("want PairAdded, got %v", ev.Type)
		}
	default:
		t.Fatal("expected PairAdded event")
	}

	got := wl.GetRequiredPairs(ex)
	if len(got) != 1 || got[0].Symbol() != pair.Symbol() {
		t.Fatalf("pairs: %#v", got)
	}

	// idempotent require
	wl.RequirePair(ex, pair)
	if len(wl.GetRequiredPairs(ex)) != 1 {
		t.Fatal("duplicate require should no-op")
	}

	wl.ReleasePair(ex, pair)
	select {
	case ev := <-ch:
		if ev.Type != baseTypes.PairRemoved {
			t.Fatalf("want PairRemoved, got %v", ev.Type)
		}
	default:
		t.Fatal("expected PairRemoved event")
	}
	if len(wl.GetRequiredPairs(ex)) != 0 {
		t.Fatal("expected empty after release")
	}
}
