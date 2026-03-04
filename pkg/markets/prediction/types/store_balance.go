package types

import (
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/data/stores/market"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

// BalanceStoreExtension stores asset balances per exchange for prediction markets.
type BalanceStoreExtension interface {
	market.StoreExtension

	// UpdateBalance sets the current balance for an asset on an exchange.
	UpdateBalance(exchange connector.ExchangeName, asset portfolio.Asset, balance numerical.Decimal)

	// GetBalance returns the stored balance for an asset on an exchange.
	// Returns zero and false if no balance has been recorded yet.
	GetBalance(exchange connector.ExchangeName, asset portfolio.Asset) (numerical.Decimal, bool)
}
