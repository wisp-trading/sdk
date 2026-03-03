package data

import (
	"github.com/wisp-trading/sdk/pkg/data/ingestors"
	"github.com/wisp-trading/sdk/pkg/data/stores"
	"go.uber.org/fx"
)

var Module = fx.Options(
	ingestors.Module,
	stores.Module,

	fx.Provide(
		NewMarketWatchlist,
	),
)
