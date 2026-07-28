package executor

import (
	"testing"
	"time"

	predictionconnector "github.com/wisp-trading/sdk/pkg/markets/prediction/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	optionsconnector "github.com/wisp-trading/sdk/pkg/types/connector/options"
	"github.com/wisp-trading/sdk/pkg/types/connector/perp"
	"github.com/wisp-trading/sdk/pkg/types/connector/spot"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	"github.com/wisp-trading/sdk/pkg/types/registry"
	"github.com/wisp-trading/sdk/pkg/types/temporal"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

// Integration-style regression: zero price must hit PlaceMarketOrder.
// Catches the go-live bug where market signals were always routed as limit@0.

type countingExec struct {
	limits, markets int
}

func (c *countingExec) PlaceLimitOrder(pair portfolio.Pair, side connector.OrderSide, quantity, price numerical.Decimal) (*connector.OrderResponse, error) {
	c.limits++
	return &connector.OrderResponse{OrderID: "L1"}, nil
}
func (c *countingExec) PlaceMarketOrder(pair portfolio.Pair, side connector.OrderSide, quantity numerical.Decimal) (*connector.OrderResponse, error) {
	c.markets++
	return &connector.OrderResponse{OrderID: "M1"}, nil
}
func (c *countingExec) CancelOrder(orderID string, pair ...portfolio.Pair) (*connector.CancelResponse, error) {
	return nil, nil
}
func (c *countingExec) GetOpenOrders(pair ...portfolio.Pair) ([]connector.Order, error) {
	return nil, nil
}
func (c *countingExec) GetOrderStatus(orderID string, pair ...portfolio.Pair) (*connector.Order, error) {
	return nil, nil
}

// orderExecConn is both Connector and OrderExecutor (type assert in GetOrderExecutor).
type orderExecConn struct {
	*countingExec
}

func (orderExecConn) GetConnectorInfo() *connector.Info { return &connector.Info{Name: "test"} }
func (orderExecConn) Initialize(config connector.Config) error {
	return nil
}
func (orderExecConn) IsInitialized() bool             { return true }
func (orderExecConn) NewConfig() connector.Config     { return nil }
func (orderExecConn) SupportsTradingOperations() bool { return true }
func (orderExecConn) SupportsRealTimeData() bool      { return false }

type stubReg struct {
	conn connector.Connector
}

func (s *stubReg) RegisterSpot(connector.ExchangeName, spot.Connector)                               {}
func (s *stubReg) RegisterPerp(connector.ExchangeName, perp.Connector)                               {}
func (s *stubReg) RegisterPrediction(connector.ExchangeName, predictionconnector.Connector)          {}
func (s *stubReg) RegisterOptions(connector.ExchangeName, optionsconnector.Connector)                {}
func (s *stubReg) Connector(name connector.ExchangeName) (connector.Connector, bool)                 { return s.conn, true }
func (s *stubReg) ConnectorType(connector.ExchangeName) (connector.MarketType, bool)                 { return connector.MarketTypePerp, true }
func (s *stubReg) Spot(connector.ExchangeName) (spot.Connector, bool)                                { return nil, false }
func (s *stubReg) Perp(connector.ExchangeName) (perp.Connector, bool)                                { return nil, false }
func (s *stubReg) Prediction(connector.ExchangeName) (predictionconnector.Connector, bool)           { return nil, false }
func (s *stubReg) Options(connector.ExchangeName) (optionsconnector.Connector, bool)                 { return nil, false }
func (s *stubReg) Filter(registry.FilterOptions) []connector.Connector                               { return nil }
func (s *stubReg) FilterSpot(registry.FilterOptions) []spot.Connector                                { return nil }
func (s *stubReg) FilterPerp(registry.FilterOptions) []perp.Connector                                { return nil }
func (s *stubReg) FilterPrediction(registry.FilterOptions) []predictionconnector.Connector           { return nil }
func (s *stubReg) FilterOptions(registry.FilterOptions) []optionsconnector.Connector                 { return nil }
func (s *stubReg) SpotWebSocket(connector.ExchangeName) (spot.WebSocketConnector, bool)               { return nil, false }
func (s *stubReg) PerpWebSocket(connector.ExchangeName) (perp.WebSocketConnector, bool)               { return nil, false }
func (s *stubReg) PredictionWebSocket(connector.ExchangeName) (predictionconnector.WebSocketConnector, bool) {
	return nil, false
}
func (s *stubReg) OptionsWebSocket(connector.ExchangeName) (optionsconnector.WebSocketConnector, bool) {
	return nil, false
}
func (s *stubReg) MarkReady(connector.ExchangeName) error { return nil }
func (s *stubReg) IsReady(connector.ExchangeName) bool    { return true }

type memBook struct {
	orders []connector.Order
}

func (m *memBook) AddTrade(connector.Trade)                               {}
func (m *memBook) AddOrder(o connector.Order)                             { m.orders = append(m.orders, o) }
func (m *memBook) UpdateOrderStatus(string, connector.OrderStatus) error { return nil }

type silentLog struct{}

func (silentLog) Info(string, ...interface{})                   {}
func (silentLog) Debug(string, ...interface{})                  {}
func (silentLog) Warn(string, ...interface{})                   {}
func (silentLog) Error(string, ...interface{})                  {}
func (silentLog) Fatal(string, ...interface{})                  {}
func (silentLog) ErrorWithDebug(string, []byte, ...interface{}) {}
func (silentLog) Infof(string, ...interface{})                  {}
func (silentLog) Debugf(string, ...interface{})                 {}
func (silentLog) Warnf(string, ...interface{})                  {}
func (silentLog) Errorf(string, ...interface{})                 {}

type realClock struct{}

func (realClock) Now() time.Time                         { return time.Unix(1_700_000_000, 0).UTC() }
func (realClock) After(d time.Duration) <-chan time.Time  { return time.After(d) }
func (realClock) Since(t time.Time) time.Duration         { return time.Since(t) }
func (realClock) Sleep(d time.Duration)                   { time.Sleep(d) }
func (realClock) NewTimer(d time.Duration) temporal.Timer {
	t := time.NewTimer(d)
	return timerWrap{t}
}
func (realClock) NewTicker(d time.Duration) temporal.Ticker {
	tk := time.NewTicker(d)
	return tickerWrap{tk}
}

type timerWrap struct{ t *time.Timer }

func (w timerWrap) C() <-chan time.Time        { return w.t.C }
func (w timerWrap) Reset(d time.Duration) bool { return w.t.Reset(d) }
func (w timerWrap) Stop() bool                 { return w.t.Stop() }

type tickerWrap struct{ t *time.Ticker }

func (w tickerWrap) C() <-chan time.Time   { return w.t.C }
func (w tickerWrap) Reset(d time.Duration) { w.t.Reset(d) }
func (w tickerWrap) Stop()                 { w.t.Stop() }

func TestPlaceOrderAndRecord_ZeroPriceIsMarket(t *testing.T) {
	exec := &countingExec{}
	conn := orderExecConn{countingExec: exec}
	b := Base{
		Connectors:   &stubReg{conn: conn},
		Logger:       silentLog{},
		TimeProvider: realClock{},
	}
	book := &memBook{}
	pair := portfolio.NewPair(portfolio.NewAsset("BTC"), portfolio.NewAsset("USD"))

	id, err := b.PlaceOrderAndRecord(book, "test", pair, connector.OrderSideBuy, numerical.NewFromFloat(0.001), numerical.Zero(), "perp")
	if err != nil {
		t.Fatal(err)
	}
	if id != "M1" || exec.markets != 1 || exec.limits != 0 {
		t.Fatalf("market path: id=%s markets=%d limits=%d", id, exec.markets, exec.limits)
	}
	if len(book.orders) != 1 || book.orders[0].Type != connector.OrderTypeMarket {
		t.Fatalf("recorded type want MARKET got %+v", book.orders)
	}
}

func TestPlaceOrderAndRecord_NonZeroIsLimit(t *testing.T) {
	exec := &countingExec{}
	conn := orderExecConn{countingExec: exec}
	b := Base{
		Connectors:   &stubReg{conn: conn},
		Logger:       silentLog{},
		TimeProvider: realClock{},
	}
	book := &memBook{}
	pair := portfolio.NewPair(portfolio.NewAsset("BTC"), portfolio.NewAsset("USD"))

	id, err := b.PlaceOrderAndRecord(book, "test", pair, connector.OrderSideSell, numerical.NewFromFloat(0.001), numerical.NewFromFloat(100_000), "perp")
	if err != nil {
		t.Fatal(err)
	}
	if id != "L1" || exec.limits != 1 || exec.markets != 0 {
		t.Fatalf("limit path: id=%s markets=%d limits=%d", id, exec.markets, exec.limits)
	}
	if book.orders[0].Type != connector.OrderTypeLimit {
		t.Fatalf("recorded type want LIMIT got %s", book.orders[0].Type)
	}
}
