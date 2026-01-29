package connector

import (
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

// Info provides metadata about the exchange's capabilities and limits.
type Info struct {
	Name                ExchangeName      `json:"name"`
	SupportedSymbols    []portfolio.Asset `json:"supported_symbols,omitempty"`
	TradingEnabled      bool              `json:"trading_enabled"`
	WebSocketEnabled    bool              `json:"websocket_enabled"`
	MaxLeverage         numerical.Decimal `json:"max_leverage,omitempty"`
	MinOrderSize        numerical.Decimal `json:"min_order_size,omitempty"`
	MaxOrderSize        numerical.Decimal `json:"max_order_size,omitempty"`
	PricePrecision      int               `json:"price_precision,omitempty"`
	QuantityPrecision   int               `json:"quantity_precision,omitempty"`
	SupportedOrderTypes []OrderType       `json:"supported_order_types,omitempty"`
	QuoteCurrency       string            `json:"quote_currency,omitempty"`
}
