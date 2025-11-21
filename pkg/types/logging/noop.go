package logging

// NoOpLogger is a logger that does nothing, useful for testing
type NoOpLogger struct{}

func NewNoOpLogger() *NoOpLogger {
	return &NoOpLogger{}
}

func (l *NoOpLogger) Info(format string, args ...interface{})                               {}
func (l *NoOpLogger) Warn(format string, args ...interface{})                               {}
func (l *NoOpLogger) Error(format string, args ...interface{})                              {}
func (l *NoOpLogger) Debug(format string, args ...interface{})                              {}
func (l *NoOpLogger) Fatal(format string, args ...interface{})                              {}
func (l *NoOpLogger) ErrorWithDebug(format string, rawResponse []byte, args ...interface{}) {}
