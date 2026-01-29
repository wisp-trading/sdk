package stores

import (
	"github.com/wisp-trading/wisp/pkg/data/stores/activity"
	"github.com/wisp-trading/wisp/pkg/data/stores/market"
	"go.uber.org/fx"
)

var Module = fx.Options(
	activity.Module,
	market.Module,
)
