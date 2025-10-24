package logging

import (
	"fmt"
	"strings"

	"github.com/backtesting-org/kronos-sdk/pkg/types/logging"
	"go.uber.org/zap"
)

// tradingLogger internal implementation
type tradingLogger struct {
	sugared *zap.SugaredLogger
	env     string
}

func NewTradingLogger(sugared *zap.SugaredLogger, env string) logging.TradingLogger {
	return &tradingLogger{
		sugared: sugared,
		env:     env,
	}
}

// Ensure tradingLogger implements the interface
var _ logging.TradingLogger = (*tradingLogger)(nil)

func (l *tradingLogger) Opportunity(strategy, asset, msg string, args ...interface{}) {
	formatted := fmt.Sprintf("💰 [%s][%s] %s", strategy, asset, fmt.Sprintf(msg, args...))
	l.sugared.Info(formatted)
}

func (l *tradingLogger) Success(strategy, asset, msg string, args ...interface{}) {
	formatted := fmt.Sprintf("✅ [%s][%s] %s", strategy, asset, fmt.Sprintf(msg, args...))
	l.sugared.Info(formatted)
}

func (l *tradingLogger) Failed(strategy, asset, msg string, args ...interface{}) {
	formatted := fmt.Sprintf("❌ [%s][%s] %s", strategy, asset, fmt.Sprintf(msg, args...))
	l.sugared.Warn(formatted)
}

func (l *tradingLogger) MarketCondition(msg string, args ...interface{}) {
	formatted := fmt.Sprintf("📊 [MARKET] %s", fmt.Sprintf(msg, args...))
	l.sugared.Info(formatted)
}

func (l *tradingLogger) DataCollection(exchange, msg string, args ...interface{}) {
	formatted := fmt.Sprintf("📡 [%s] %s", exchange, fmt.Sprintf(msg, args...))
	l.sugared.Debug(formatted)
}

func (l *tradingLogger) OrderLifecycle(msg, asset string, args ...interface{}) {
	formatted := fmt.Sprintf("🔄 [%s] %s", asset, fmt.Sprintf(msg, args...))
	l.sugared.Info(formatted)
}

func (l *tradingLogger) Info(msg string, args ...interface{}) {
	formatted := fmt.Sprintf("ℹ️ %s", fmt.Sprintf(msg, args...))
	l.sugared.Info(formatted)
}

func (l *tradingLogger) Debug(strategy, asset, msg string, args ...interface{}) {
	// Only log debug in development environments
	if strings.ToLower(l.env) == "debug" || strings.ToLower(l.env) == "development" {
		formatted := fmt.Sprintf("🔍 [%s][%s] %s", strategy, asset, fmt.Sprintf(msg, args...))
		l.sugared.Debug(formatted)
	}
}
