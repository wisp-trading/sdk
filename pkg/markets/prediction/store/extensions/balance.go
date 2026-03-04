package extensions

import (
	"sync"

	predictionTypes "github.com/wisp-trading/sdk/pkg/markets/prediction/types"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

type balanceKey struct {
	exchange connector.ExchangeName
	asset    portfolio.Asset
}

type predictionBalanceExtension struct {
	mu       sync.RWMutex
	balances map[balanceKey]numerical.Decimal
}

func NewPredictionBalanceExtension() predictionTypes.BalanceStoreExtension {
	return &predictionBalanceExtension{
		balances: make(map[balanceKey]numerical.Decimal),
	}
}

func (e *predictionBalanceExtension) UpdateBalance(exchange connector.ExchangeName, asset portfolio.Asset, balance numerical.Decimal) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.balances[balanceKey{exchange, asset}] = balance
}

func (e *predictionBalanceExtension) GetBalance(exchange connector.ExchangeName, asset portfolio.Asset) (numerical.Decimal, bool) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	balance, ok := e.balances[balanceKey{exchange, asset}]
	return balance, ok
}
