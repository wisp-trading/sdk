package strategy

import (
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

type TradeAction struct {
	Action   Action                 `json:"action"`
	Asset    portfolio.Pair         `json:"asset"`
	Exchange connector.ExchangeName `json:"exchange"`
	Quantity numerical.Decimal      `json:"quantity"`
	Price    numerical.Decimal      `json:"price"`
}

type Action string

const (
	ActionBuy       Action = "buy"
	ActionSell      Action = "sell"
	ActionSellShort Action = "sell_short"
	ActionCover     Action = "cover"

	ActionHold  Action = "hold"
	ActionClose Action = "close"
)
