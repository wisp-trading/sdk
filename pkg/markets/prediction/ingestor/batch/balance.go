package batch

import (
	"github.com/wisp-trading/sdk/pkg/markets/prediction/types"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/logging"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
)

// balanceCollectionExtension polls the connector for account balances and writes them to the store.
type balanceCollectionExtension struct {
	store  types.BalanceStoreExtension
	logger logging.ApplicationLogger
}

func NewBalanceCollectionExtension(
	store types.BalanceStoreExtension,
	logger logging.ApplicationLogger,
) types.PredictionCollectionExtension {
	return &balanceCollectionExtension{
		store:  store,
		logger: logger,
	}
}

func (e *balanceCollectionExtension) Collect(conn interface{}, exchangeName connector.ExchangeName) {
	accountReader, ok := conn.(connector.AccountReader)
	if !ok {
		e.logger.Debug("Connector %s does not implement AccountReader, skipping balance collection", exchangeName)
		return
	}

	asset := portfolio.NewAsset("USD")

	balance, err := accountReader.GetBalance(asset)
	if err != nil {
		e.logger.Error("Failed to fetch balances from %s: %v", exchangeName, err)
		return
	}

	e.store.UpdateBalance(exchangeName, balance.Asset, balance.Free)
	e.logger.Info("Updated balance for %s on %s: %s", balance.Asset.Symbol(), exchangeName, balance.Free.StringFixed(4))
}
