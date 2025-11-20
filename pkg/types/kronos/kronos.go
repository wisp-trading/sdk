package kronos

import (
	"github.com/backtesting-org/kronos-sdk/kronos/signal"
	"github.com/backtesting-org/kronos-sdk/pkg/types/logging"
	"github.com/backtesting-org/kronos-sdk/pkg/types/portfolio"
	"github.com/backtesting-org/kronos-sdk/pkg/types/stores/market"
	"github.com/backtesting-org/kronos-sdk/pkg/types/strategy"
)

type Kronos interface {
	Log() logging.TradingLogger
	Store() market.MarketData
	Asset(symbol string) portfolio.Asset
	Signal(strategyName strategy.StrategyName) *signal.Builder
}
