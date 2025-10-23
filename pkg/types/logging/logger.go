package logging

// ApplicationLogger interface for system/code errors that go to Sentry
type ApplicationLogger interface {
	Info(msg string, args ...interface{})
	Debug(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Fatal(msg string, args ...interface{})
	ErrorWithDebug(msg string, rawResponse []byte, args ...interface{})
}

// TradingLogger interface for business events that don't go to Sentry
type TradingLogger interface {
	// Market condition logging
	MarketCondition(msg string, args ...interface{})

	// Trading operations
	Opportunity(strategy, asset, msg string, args ...interface{})
	Success(strategy, asset, msg string, args ...interface{})
	Failed(strategy, asset, msg string, args ...interface{})
	OrderLifecycle(msg, asset string, args ...interface{})

	// Data collection events
	DataCollection(exchange, msg string, args ...interface{})

	// Debug logging for trading events
	Debug(strategy, asset, msg string, args ...interface{})

	Info(msg string, args ...interface{})
}
