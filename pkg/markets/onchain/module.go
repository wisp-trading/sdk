package onchain

import (
	"github.com/wisp-trading/sdk/pkg/markets/onchain/activity"
	"github.com/wisp-trading/sdk/pkg/markets/onchain/executor"
	"github.com/wisp-trading/sdk/pkg/markets/onchain/facade"
	"github.com/wisp-trading/sdk/pkg/markets/onchain/store"
	lifecycleTypes "github.com/wisp-trading/sdk/pkg/types/lifecycle"
	"go.uber.org/fx"
)

// Module wires the onchain market domain (AMM / UniV3 path).
var Module = fx.Module("onchain",
	fx.Provide(
		store.NewStore,
		facade.NewOnchain,
		executor.NewExecutor,
		NewOnchainWatchlist,
		activity.NewOnchainPNL,
		fx.Annotate(
			newOnchainDomainLifecycle,
			fx.ResultTags(`group:"domain_lifecycles"`),
		),
	),
)

// ensure domain lifecycle type is referenced for docs
var _ = lifecycleTypes.DomainLifecycle(nil)
