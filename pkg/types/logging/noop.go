package logging

// NoOpLogger is a logger that does nothing, useful for testing
type NoOpLogger struct{}

func NewNoOpLogger() ApplicationLogger {
	return &NoOpLogger{}
}

// Printf-style
func (l *NoOpLogger) Info(format string, args ...interface{})                               {}
func (l *NoOpLogger) Warn(format string, args ...interface{})                               {}
func (l *NoOpLogger) Error(format string, args ...interface{})                              {}
func (l *NoOpLogger) Debug(format string, args ...interface{})                              {}
func (l *NoOpLogger) Fatal(format string, args ...interface{})                              {}
func (l *NoOpLogger) ErrorWithDebug(format string, rawResponse []byte, args ...interface{}) {}

// Structured key-value
func (l *NoOpLogger) Infof(msg string, args ...interface{})  {}
func (l *NoOpLogger) Debugf(msg string, args ...interface{}) {}
func (l *NoOpLogger) Warnf(msg string, args ...interface{})  {}
func (l *NoOpLogger) Errorf(msg string, args ...interface{}) {}

var _ ApplicationLogger = (*NoOpLogger)(nil)
