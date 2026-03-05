package logging

import (
	"github.com/wisp-trading/sdk/pkg/types/logging"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ZapApplicationLogger adapts zap.Logger to ApplicationLogger interface
type ZapApplicationLogger struct {
	logger *zap.Logger
}

// NewDefaultZapLogger creates a production-ready zap logger with proper encoding
func NewDefaultZapLogger() (*zap.Logger, error) {
	config := zap.NewProductionConfig()
	config.Encoding = "json"
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return config.Build()
}

// NewZapApplicationLogger creates an application logger from a zap logger
func NewZapApplicationLogger(logger *zap.Logger) logging.ApplicationLogger {
	return &ZapApplicationLogger{logger: logger}
}

// Printf-style — interpolates args into the format string.
func (z *ZapApplicationLogger) Info(format string, args ...interface{}) {
	z.logger.Sugar().Infof(format, args...)
}

func (z *ZapApplicationLogger) Debug(format string, args ...interface{}) {
	z.logger.Sugar().Debugf(format, args...)
}

func (z *ZapApplicationLogger) Warn(format string, args ...interface{}) {
	z.logger.Sugar().Warnf(format, args...)
}

func (z *ZapApplicationLogger) Error(format string, args ...interface{}) {
	z.logger.Sugar().Errorf(format, args...)
}

func (z *ZapApplicationLogger) Fatal(format string, args ...interface{}) {
	z.logger.Sugar().Fatalf(format, args...)
}

func (z *ZapApplicationLogger) ErrorWithDebug(msg string, rawResponse []byte, args ...interface{}) {
	fields := []interface{}{"raw_response", string(rawResponse)}
	fields = append(fields, args...)
	z.logger.Sugar().Errorw(msg, fields...)
}

// Structured key-value — emits args as discrete JSON fields.
func (z *ZapApplicationLogger) Infof(msg string, args ...interface{}) {
	z.logger.Sugar().Infow(msg, args...)
}

func (z *ZapApplicationLogger) Debugf(msg string, args ...interface{}) {
	z.logger.Sugar().Debugw(msg, args...)
}

func (z *ZapApplicationLogger) Warnf(msg string, args ...interface{}) {
	z.logger.Sugar().Warnw(msg, args...)
}

func (z *ZapApplicationLogger) Errorf(msg string, args ...interface{}) {
	z.logger.Sugar().Errorw(msg, args...)
}

// ZapTradingLogger adapts zap.Logger to TradingLogger interface
type ZapTradingLogger struct {
	logger *zap.Logger
}

// NewZapTradingLogger creates a trading logger from a zap logger
func NewZapTradingLogger(logger *zap.Logger) logging.TradingLogger {
	return &ZapTradingLogger{logger: logger}
}

// Printf-style — interpolates args into the format string.
func (z *ZapTradingLogger) Info(format string, args ...interface{}) {
	z.logger.Sugar().Infof(format, args...)
}

// Structured key-value — emits args as discrete JSON fields.
func (z *ZapTradingLogger) Infof(msg string, args ...interface{}) {
	z.logger.Sugar().Infow(msg, args...)
}

func (z *ZapTradingLogger) Opportunity(strategy, asset, msg string, args ...interface{}) {
	fields := []interface{}{"strategy", strategy, "asset", asset}
	fields = append(fields, args...)
	z.logger.Sugar().Infow(msg, fields...)
}

func (z *ZapTradingLogger) Success(strategy, asset, msg string, args ...interface{}) {
	fields := []interface{}{"strategy", strategy, "asset", asset}
	fields = append(fields, args...)
	z.logger.Sugar().Infow(msg, fields...)
}

func (z *ZapTradingLogger) Failed(strategy, asset, msg string, args ...interface{}) {
	fields := []interface{}{"strategy", strategy, "asset", asset}
	fields = append(fields, args...)
	z.logger.Sugar().Errorw(msg, fields...)
}

func (z *ZapTradingLogger) MarketCondition(msg string, args ...interface{}) {
	z.logger.Sugar().Infow(msg, args...)
}

func (z *ZapTradingLogger) OrderLifecycle(msg, asset string, args ...interface{}) {
	fields := []interface{}{"asset", asset}
	fields = append(fields, args...)
	z.logger.Sugar().Infow(msg, fields...)
}

func (z *ZapTradingLogger) Debug(strategy, asset, msg string, args ...interface{}) {
	fields := []interface{}{"strategy", strategy, "asset", asset}
	fields = append(fields, args...)
	z.logger.Sugar().Debugw(msg, fields...)
}

func (z *ZapTradingLogger) DataCollection(exchange, msg string, args ...interface{}) {
	fields := []interface{}{"exchange", exchange}
	fields = append(fields, args...)
	z.logger.Sugar().Infow(msg, fields...)
}
