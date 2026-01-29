package activity

import (
	"github.com/wisp-trading/wisp/pkg/data/stores/activity/position"
	"github.com/wisp-trading/wisp/pkg/data/stores/activity/trade"
	"go.uber.org/fx"
)

var Module = fx.Options(
	position.Module,
	trade.Module,
)
