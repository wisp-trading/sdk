package connectors

import (
	"strings"
	"testing"
)

func TestValidationError_NotConfiguredMessage(t *testing.T) {
	ve := &ValidationError{Exchange: "hyperliquid", NotConfigured: true}
	msg := ve.Error()
	if !strings.Contains(msg, "Settings") || !strings.Contains(msg, "hyperliquid") {
		t.Fatalf("want Settings hint, got %q", msg)
	}
	if strings.Contains(msg, "registry") {
		t.Fatalf("not-configured must not mention registry: %q", msg)
	}
}

func TestStrategyValidationError_NotConfigured(t *testing.T) {
	sve := &StrategyValidationError{
		ExchangeNames: []string{"hyperliquid"},
		ValidationErrors: map[string]*ValidationError{
			"hyperliquid": {Exchange: "hyperliquid", NotConfigured: true},
		},
	}
	msg := sve.Error()
	if !strings.Contains(msg, "no keys") && !strings.Contains(msg, "Settings") {
		t.Fatalf("want keys/Settings hint, got %q", msg)
	}
	if strings.Contains(msg, "not registered") {
		t.Fatalf("misclassified as registry miss: %q", msg)
	}
}

func TestValidationError_ExchangeNotFoundMessage(t *testing.T) {
	ve := &ValidationError{Exchange: "nope", ExchangeNotFound: true}
	msg := ve.Error()
	if !strings.Contains(msg, "not registered") {
		t.Fatalf("got %q", msg)
	}
	if !strings.Contains(msg, "connectors.Module") {
		t.Fatalf("want connectors.Module hint: %q", msg)
	}
}
