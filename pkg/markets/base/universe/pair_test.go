package universe

import (
	"testing"

	baseTypes "github.com/wisp-trading/sdk/pkg/markets/base/types"
	"github.com/wisp-trading/sdk/pkg/markets/base/watchlist"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
)

func TestBuildPairUniverse(t *testing.T) {
	wl := watchlist.NewPairWatchlist()
	ex := connector.ExchangeName("hyperliquid")
	pair := portfolio.NewPair(portfolio.NewAsset("BTC"), portfolio.NewAsset("USD"))
	wl.RequirePair(ex, pair)

	uni := BuildPairUniverse([]connector.ExchangeName{ex}, connector.MarketTypePerp, wl)
	if len(uni.Exchanges) != 1 || uni.Exchanges[0].MarketType != connector.MarketTypePerp {
		t.Fatalf("exchanges %#v", uni.Exchanges)
	}
	if len(uni.Assets[ex]) != 1 {
		t.Fatalf("assets %#v", uni.Assets)
	}

	// empty watchlist still lists exchange
	empty := BuildPairUniverse([]connector.ExchangeName{ex}, connector.MarketTypeSpot, watchlist.NewPairWatchlist())
	if len(empty.Exchanges) != 1 || len(empty.Assets) != 0 {
		t.Fatalf("empty %#v", empty)
	}

	_ = baseTypes.MarketWatchlist(wl)
}
