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
	if !strings.Contains(msg, "No keys") {
		t.Fatalf("got %q", msg)
	}
	if strings.Contains(strings.Split(msg, "hyperliquid")[1], "SDK registry") {
		// only fail if registry appears for this exchange line — coarse check
	}
	if strings.Contains(msg, "Not found in SDK registry") {
		t.Fatalf("misclassified as registry miss: %q", msg)
	}
}
