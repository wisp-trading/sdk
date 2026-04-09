package price_feeds

import (
	"github.com/wisp-trading/sdk/pkg/markets/base/ingestor"
	"github.com/wisp-trading/sdk/pkg/markets/price_feeds/ingestor/batch"
	"github.com/wisp-trading/sdk/pkg/markets/price_feeds/ingestor/realtime"
	"github.com/wisp-trading/sdk/pkg/markets/price_feeds/store"
	priceFeedTypes "github.com/wisp-trading/sdk/pkg/markets/price_feeds/types"
	"github.com/wisp-trading/sdk/pkg/types"
	lifecycleTypes "github.com/wisp-trading/sdk/pkg/types/lifecycle"
	"github.com/wisp-trading/sdk/pkg/types/logging"
	"go.uber.org/fx"
)

var Module = fx.Module("price_feeds",
	fx.Provide(
		store.NewStore,
		batch.NewFactory,
		realtime.NewFactory,
		providePriceFeedsAdapter,
		fx.Annotate(
			newPriceFeedsDomainLifecycle,
			fx.ResultTags(`group:"domain_lifecycles"`),
		),
	),
)

// providePriceFeedsAdapter wraps the internal PriceFeedsStore as the public types.PriceFeeds interface.
func providePriceFeedsAdapter(store priceFeedTypes.PriceFeedsStore) types.PriceFeeds {
	return NewPriceFeedsAdapter(store)
}

func newPriceFeedsDomainLifecycle(
	batchFactory priceFeedTypes.PriceFeedsBatchIngestorFactory,
	realtimeFactory priceFeedTypes.PriceFeedsRealtimeIngestorFactory,
	logger logging.ApplicationLogger,
) lifecycleTypes.DomainLifecycle {
	return ingestor.NewDomainCoordinator("price_feeds", nil, batchFactory, realtimeFactory, logger)
}
