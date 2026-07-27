package analytics

import (
	"context"
	"testing"

	market "github.com/wisp-trading/sdk/pkg/markets/base/types/stores/market"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

type stubStore struct {
	prices map[portfolio.Pair]market.PriceMap
	books  map[string]*connector.OrderBook
}

func (s *stubStore) GetPairPrices(pair portfolio.Pair) market.PriceMap {
	if s.prices == nil {
		return market.PriceMap{}
	}
	return s.prices[pair]
}

func (s *stubStore) GetOrderBook(pair portfolio.Pair, exchange connector.ExchangeName) *connector.OrderBook {
	if s.books == nil {
		return nil
	}
	return s.books[string(exchange)+":"+pair.Symbol()]
}

func (s *stubStore) GetKlines(portfolio.Pair, connector.ExchangeName, string, int) []connector.Kline {
	return nil
}

func TestPairServicePrice(t *testing.T) {
	pair := portfolio.NewPair(portfolio.NewAsset("BTC"), portfolio.NewAsset("USD"))
	ex := connector.ExchangeName("hyperliquid")
	store := &stubStore{
		prices: map[portfolio.Pair]market.PriceMap{
			pair: {ex: connector.Price{Price: numerical.NewFromFloat(100)}},
		},
	}
	svc := NewService(store)

	got, err := svc.Price(context.Background(), pair, ex)
	if err != nil {
		t.Fatal(err)
	}
	if !got.Equal(numerical.NewFromFloat(100)) {
		t.Fatalf("got %v", got)
	}

	all := svc.Prices(context.Background(), pair)
	if len(all) != 1 || !all[ex].Equal(numerical.NewFromFloat(100)) {
		t.Fatalf("prices %#v", all)
	}

	_, err = svc.Price(context.Background(), pair, "missing")
	if err == nil {
		t.Fatal("want error for missing exchange")
	}
}

func TestPairServiceOrderBook(t *testing.T) {
	pair := portfolio.NewPair(portfolio.NewAsset("BTC"), portfolio.NewAsset("USD"))
	ex := connector.ExchangeName("hyperliquid")
	ob := &connector.OrderBook{}
	store := &stubStore{
		books: map[string]*connector.OrderBook{
			string(ex) + ":" + pair.Symbol(): ob,
		},
	}
	svc := NewService(store)
	got, err := svc.OrderBook(context.Background(), pair, ex)
	if err != nil || got != ob {
		t.Fatalf("got %v err %v", got, err)
	}
	_, err = svc.OrderBook(context.Background(), pair, "missing")
	if err == nil {
		t.Fatal("want error")
	}
}
