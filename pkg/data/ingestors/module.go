package ingestors

import (
	"github.com/wisp-trading/wisp/pkg/data/ingestors/market"
	"github.com/wisp-trading/wisp/pkg/data/ingestors/position"
	"go.uber.org/fx"
)

var Module = fx.Options(
	position.Module,
	market.Module,
)
