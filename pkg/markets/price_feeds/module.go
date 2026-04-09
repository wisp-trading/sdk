package price_feeds

import (
	"github.com/wisp-trading/sdk/pkg/markets/base/ingestor"
	"github.com/wisp-trading/sdk/pkg/markets/price_feeds/ingestor/realtime"
	"github.com/wisp-trading/sdk/pkg/markets/price_feeds/store"
	priceFeedTypes "github.com/wisp-trading/sdk/pkg/markets/price_feeds/types"
	lifecycleTypes "github.com/wisp-trading/sdk/pkg/types/lifecycle"
	"github.com/wisp-trading/sdk/pkg/types/logging"
	"go.uber.org/fx"
)

var Module = fx.Module("price_feeds",
	fx.Provide(
		store.NewStore,
		realtime.NewFactory,
		fx.Annotate(
			newPriceFeedsDomainLifecycle,
			fx.ResultTags(`group:"domain_lifecycles"`),
		),
	),
)

func newPriceFeedsDomainLifecycle(
	batchFactory priceFeedTypes.PriceFeedsBatchIngestorFactory,
	logger logging.ApplicationLogger,
) lifecycleTypes.DomainLifecycle {
	return ingestor.NewDomainCoordinator("price_feeds", nil, batchFactory, nil, logger)
}
