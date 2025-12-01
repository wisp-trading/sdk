package strategy

import (
	"github.com/backtesting-org/kronos-sdk/pkg/types/connector"
	"github.com/backtesting-org/kronos-sdk/pkg/types/kronos/numerical"
	"github.com/backtesting-org/kronos-sdk/pkg/types/portfolio"
)

type TradeAction struct {
	Action   Action                 `json:"action"`
	Asset    portfolio.Asset        `json:"asset"`
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
