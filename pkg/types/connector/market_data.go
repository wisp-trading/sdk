package connector

import (
	"time"

	"github.com/backtesting-org/kronos-sdk/pkg/types/portfolio"
	"github.com/shopspring/decimal"
)

// Price represents the market price data for a symbol.
type Price struct {
	Symbol    string          `json:"symbol"`
	Price     decimal.Decimal `json:"price"`
	BidPrice  decimal.Decimal `json:"bid_price,omitempty"`
	AskPrice  decimal.Decimal `json:"ask_price,omitempty"`
	Volume24h decimal.Decimal `json:"volume_24h,omitempty"`
	Change24h decimal.Decimal `json:"change_24h,omitempty"`
	Source    ExchangeName    `json:"source"`
	Timestamp time.Time       `json:"timestamp"`
}

// Kline represents candlestick data for a trading symbol.
type Kline struct {
	Symbol      string          `json:"symbol"`
	Interval    string          `json:"interval"`
	OpenTime    time.Time       `json:"open_time"`
	Open        decimal.Decimal `json:"open"`
	High        decimal.Decimal `json:"high"`
	Low         decimal.Decimal `json:"low"`
	Close       decimal.Decimal `json:"close"`
	Volume      decimal.Decimal `json:"volume"`
	CloseTime   time.Time       `json:"close_time"`
	QuoteVolume decimal.Decimal `json:"quote_volume,omitempty"`
	TradeCount  int             `json:"trade_count,omitempty"`
	TakerVolume decimal.Decimal `json:"taker_volume,omitempty"`
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
	Price    decimal.Decimal `json:"price"`
	Quantity decimal.Decimal `json:"quantity"`
}

// Trade represents a trade executed on the exchange.
type Trade struct {
	ID        string          `json:"id"`
	Symbol    string          `json:"symbol"`
	Price     decimal.Decimal `json:"price"`
	Quantity  decimal.Decimal `json:"quantity"`
	Side      OrderSide       `json:"side"`
	IsMaker   bool            `json:"is_maker"`
	Fee       decimal.Decimal `json:"fee,omitempty"`
	Timestamp time.Time       `json:"timestamp"`
}
