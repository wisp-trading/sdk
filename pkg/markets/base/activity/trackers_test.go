package activity

import (
	"testing"

	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

func TestBuildTrackersBuyThenSellRealized(t *testing.T) {
	pair := portfolio.NewPair(portfolio.NewAsset("BTC"), portfolio.NewAsset("USD"))
	ex := connector.ExchangeName("test")
	trades := []connector.Trade{
		{Exchange: ex, Pair: pair, Side: connector.OrderSideBuy, Quantity: numerical.NewFromFloat(1), Price: numerical.NewFromFloat(100), Fee: numerical.NewFromFloat(0.1)},
		{Exchange: ex, Pair: pair, Side: connector.OrderSideSell, Quantity: numerical.NewFromFloat(1), Price: numerical.NewFromFloat(110), Fee: numerical.NewFromFloat(0.1)},
	}
	tr := BuildTrackers(trades)
	if len(tr) != 1 {
		t.Fatalf("trackers=%d", len(tr))
	}
	for _, p := range tr {
		if !p.Realized.Equal(numerical.NewFromFloat(10)) {
			t.Fatalf("realized=%v want 10", p.Realized)
		}
		if !p.Size.IsZero() {
			t.Fatalf("size should be flat, got %v", p.Size)
		}
		if !p.Fees.Equal(numerical.NewFromFloat(0.2)) {
			t.Fatalf("fees=%v", p.Fees)
		}
	}
}
