package connector

import (
	"time"

	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

// MarketDataReader provides read-only market data access
type MarketDataReader interface {
	FetchPrice(pair portfolio.Pair) (*Price, error)
	FetchKlines(pair portfolio.Pair, interval string, limit int) ([]Kline, error)
	FetchOrderBook(pair portfolio.Pair, depth int) (*OrderBook, error)
	FetchRecentTrades(pair portfolio.Pair, limit int) ([]Trade, error)
}

// Price represents the market price data for a pair.
type Price struct {
	Pair      portfolio.Pair    `json:"pair"`
	Price     numerical.Decimal `json:"price"`
	BidPrice  numerical.Decimal `json:"bid_price,omitempty"`
	AskPrice  numerical.Decimal `json:"ask_price,omitempty"`
	Volume24h numerical.Decimal `json:"volume_24h,omitempty"`
	Change24h numerical.Decimal `json:"change_24h,omitempty"`
	Source    ExchangeName      `json:"source"`
	Timestamp time.Time         `json:"timestamp"`
}

// Kline represents candlestick data for a trading pair.
type Kline struct {
	Pair        portfolio.Pair `json:"pair"`
	Interval    string         `json:"interval"`
	OpenTime    time.Time      `json:"open_time"`
	Open        float64        `json:"open"`
	High        float64        `json:"high"`
	Low         float64        `json:"low"`
	Close       float64        `json:"close"`
	Volume      float64        `json:"volume"`
	CloseTime   time.Time      `json:"close_time"`
	QuoteVolume float64        `json:"quote_volume,omitempty"`
	TradeCount  int            `json:"trade_count,omitempty"`
	TakerVolume float64        `json:"taker_volume,omitempty"`
}

// OrderBook represents the order book with bids and asks.
type OrderBook struct {
	Pair      portfolio.Pair `json:"pair"`
	Bids      []PriceLevel   `json:"bids"`
	Asks      []PriceLevel   `json:"asks"`
	Timestamp time.Time      `json:"timestamp"`
}

// PriceLevel represents a price level in the order book.
type PriceLevel struct {
	Price    numerical.Decimal `json:"price"`
	Quantity numerical.Decimal `json:"quantity"`
}

// Trade represents a trade executed on the exchange.
type Trade struct {
	ID        string            `json:"id"`
	OrderID   string            `json:"order_id,omitempty"` // Link to the originating order
	Pair      portfolio.Pair    `json:"pair"`
	Exchange  ExchangeName      `json:"exchange"`
	Price     numerical.Decimal `json:"price"`
	Quantity  numerical.Decimal `json:"quantity"`
	Side      OrderSide         `json:"side"`
	IsMaker   bool              `json:"is_maker"`
	Fee       numerical.Decimal `json:"fee,omitempty"`
	Timestamp time.Time         `json:"timestamp"`
}
