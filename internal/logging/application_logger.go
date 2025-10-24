package logging

import (
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/backtesting-org/kronos-sdk/pkg/types/logging"
	"go.uber.org/zap"
)

// applicationLogger internal implementation
type applicationLogger struct {
	sugared *zap.SugaredLogger
	env     string
}

func NewApplicationLogger(sugared *zap.SugaredLogger, env string) logging.ApplicationLogger {
	return &applicationLogger{
		sugared: sugared,
		env:     env,
	}
}

// Ensure applicationLogger implements the interface
var _ logging.ApplicationLogger = (*applicationLogger)(nil)

func (l *applicationLogger) Error(msg string, args ...interface{}) {
	formatted := fmt.Sprintf(msg, args...)

	// Get caller info for clean file:line display
	_, file, line, ok := runtime.Caller(1)
	if ok {
		filename := filepath.Base(file)
		formatted = fmt.Sprintf("%s (at %s:%d)", formatted, filename, line)
	}

	l.sugared.Error(formatted)
}

func (l *applicationLogger) Fatal(msg string, args ...interface{}) {
	formatted := fmt.Sprintf(msg, args...)
	l.sugared.Fatal(formatted)

}

func (l *applicationLogger) Warn(msg string, args ...interface{}) {
	formatted := fmt.Sprintf(msg, args...)
	l.sugared.Warn(formatted)
}

func (l *applicationLogger) Info(msg string, args ...interface{}) {
	l.sugared.Infof(msg, args...)
}

func (l *applicationLogger) Debug(msg string, args ...interface{}) {
	l.sugared.Debugf(msg, args...)
}

func (l *applicationLogger) ErrorWithDebug(msg string, rawResponse []byte, args ...interface{}) {
	formatted := fmt.Sprintf(msg, args...)

	// Add response debugging info for logs
	responsePreview := string(rawResponse)
	if len(responsePreview) > 200 {
		responsePreview = responsePreview[:200] + "..."
	}

	debugMsg := fmt.Sprintf("%s | Response: %s", formatted, responsePreview)

	// Get clean caller info
	_, file, line, ok := runtime.Caller(1)
	if ok {
		filename := filepath.Base(file)
		debugMsg = fmt.Sprintf("%s (at %s:%d)", debugMsg, filename, line)
	}

	l.sugared.Error(debugMsg)
}
