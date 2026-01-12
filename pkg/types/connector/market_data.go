package connector

import (
	"time"

	"github.com/backtesting-org/kronos-sdk/pkg/types/kronos/numerical"
	"github.com/backtesting-org/kronos-sdk/pkg/types/portfolio"
)

// Price represents the market price data for a symbol.
type Price struct {
	Symbol    string            `json:"symbol"`
	Price     numerical.Decimal `json:"price"`
	BidPrice  numerical.Decimal `json:"bid_price,omitempty"`
	AskPrice  numerical.Decimal `json:"ask_price,omitempty"`
	Volume24h numerical.Decimal `json:"volume_24h,omitempty"`
	Change24h numerical.Decimal `json:"change_24h,omitempty"`
	Source    ExchangeName      `json:"source"`
	Timestamp time.Time         `json:"timestamp"`
}

// Kline represents candlestick data for a trading symbol.
type Kline struct {
	Symbol      string    `json:"symbol"`
	Interval    string    `json:"interval"`
	OpenTime    time.Time `json:"open_time"`
	Open        float64   `json:"open"`
	High        float64   `json:"high"`
	Low         float64   `json:"low"`
	Close       float64   `json:"close"`
	Volume      float64   `json:"volume"`
	CloseTime   time.Time `json:"close_time"`
	QuoteVolume float64   `json:"quote_volume,omitempty"`
	TradeCount  int       `json:"trade_count,omitempty"`
	TakerVolume float64   `json:"taker_volume,omitempty"`
}

// OrderBook represents the order book with bids and asks.
type OrderBook struct {
	Asset     portfolio.Asset `json:"symbol"`
	Bids      []PriceLevel    `json:"bids"`
	Asks      []PriceLevel    `json:"asks"`
	Timestamp time.Time       `json:"timestamp"`
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
	Symbol    string            `json:"symbol"`
	Exchange  ExchangeName      `json:"exchange"`
	Price     numerical.Decimal `json:"price"`
	Quantity  numerical.Decimal `json:"quantity"`
	Side      OrderSide         `json:"side"`
	IsMaker   bool              `json:"is_maker"`
	Fee       numerical.Decimal `json:"fee,omitempty"`
	Timestamp time.Time         `json:"timestamp"`
}
