package packages

import (
	"github.com/backtesting-org/kronos-sdk/pkg/adapters"
	"github.com/backtesting-org/kronos-sdk/pkg/analytics"
	"github.com/backtesting-org/kronos-sdk/pkg/events"
	"github.com/backtesting-org/kronos-sdk/pkg/ingestors"
	"github.com/backtesting-org/kronos-sdk/pkg/registry"
	"github.com/backtesting-org/kronos-sdk/pkg/runtime"
	"github.com/backtesting-org/kronos-sdk/pkg/signal"
	"github.com/backtesting-org/kronos-sdk/pkg/stores"
	"go.uber.org/fx"
)

var Module = fx.Options(
	adapters.Module,
	analytics.Module,
	events.Module,
	ingestors.Module,
	registry.Module,
	runtime.Module,
	signal.Module,
	stores.Module,
)
