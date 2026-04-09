package realtime

import (
	"go.uber.org/fx"

	priceFeedTypes "github.com/wisp-trading/sdk/pkg/markets/price_feeds/types"
	"github.com/wisp-trading/sdk/pkg/types/logging"
)

// FactoryParams contains dependencies for creating an ingestor.
type FactoryParams struct {
	fx.In

	Connector Connector
	Store     priceFeedTypes.PriceFeedsStore
	Logger    logging.ApplicationLogger
}

// NewIngestor is an fx provider that creates a Pyth price feed ingestor.
func NewIngestor(params FactoryParams) Ingestor {
	return New(params.Connector, params.Store, params.Logger)
}
