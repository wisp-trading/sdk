package logging

// ApplicationLogger interface for system/code errors that go to Sentry.
//
// Two logging styles are supported:
//
//  1. Printf-style (plain methods) — for simple messages with interpolated values.
//     Existing call sites use this style and it remains unchanged.
//     e.g. logger.Info("Updated balance for %s on %s: %s", asset, exchange, value)
//
//  2. Structured key-value (f-suffixed methods) — preferred for new code.
//     Args are alternating key/value pairs, emitted as discrete JSON fields.
//     This makes logs machine-queryable in aggregators like Datadog / Loki.
//     e.g. logger.Infof("Balance updated", "asset", asset, "exchange", exchange, "free", value)
type ApplicationLogger interface {
	// Printf-style — format string followed by format operands.
	Info(format string, args ...interface{})
	Debug(format string, args ...interface{})
	Warn(format string, args ...interface{})
	Error(format string, args ...interface{})
	Fatal(format string, args ...interface{})
	ErrorWithDebug(msg string, rawResponse []byte, args ...interface{})

	// Structured key-value — message followed by alternating key, value pairs.
	Infof(msg string, args ...interface{})
	Debugf(msg string, args ...interface{})
	Warnf(msg string, args ...interface{})
	Errorf(msg string, args ...interface{})
}

// TradingLogger interface for business events that don't go to Sentry.
//
// The same two logging styles apply as ApplicationLogger.
// Trading-specific methods (Opportunity, Success, etc.) always use structured
// logging — strategy and asset are automatically included as fields.
type TradingLogger interface {
	// Printf-style — format string followed by format operands.
	Info(format string, args ...interface{})

	// Structured key-value — message followed by alternating key, value pairs.
	Infof(msg string, args ...interface{})

	// Market condition logging (structured).
	MarketCondition(msg string, args ...interface{})

	// Trading operations (structured — strategy & asset added automatically).
	Opportunity(strategy, asset, msg string, args ...interface{})
	Success(strategy, asset, msg string, args ...interface{})
	Failed(strategy, asset, msg string, args ...interface{})
	OrderLifecycle(msg, asset string, args ...interface{})

	// Data collection events (structured — exchange added automatically).
	DataCollection(exchange, msg string, args ...interface{})

	// Debug logging (structured — strategy & asset added automatically).
	Debug(strategy, asset, msg string, args ...interface{})
}
