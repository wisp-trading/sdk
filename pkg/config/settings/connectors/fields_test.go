package connectors

import (
	"testing"

	"github.com/wisp-trading/sdk/pkg/types/connector"
)

// sampleConfig mirrors exchange configs: required keys have no omitempty;
// optional URL/flags use omitempty so zero-value JSON omits them.
type sampleConfig struct {
	APIKey    string  `json:"api_key"`
	APISecret string  `json:"api_secret"`
	BaseURL   string  `json:"base_url,omitempty"`
	UseTest   bool    `json:"use_testnet,omitempty"`
	Slippage  float64 `json:"default_slippage,omitempty"`
}

func (sampleConfig) ExchangeName() connector.ExchangeName { return "sample" }
func (sampleConfig) Validate() error                      { return nil }

func TestCredentialFieldsFromConfig_RequiredOnly(t *testing.T) {
	fields := credentialFieldsFromConfig(sampleConfig{})
	want := []string{"api_key", "api_secret"}
	if len(fields) != len(want) {
		t.Fatalf("fields=%v want %v", fields, want)
	}
	for i := range want {
		if fields[i] != want[i] {
			t.Fatalf("fields=%v want %v", fields, want)
		}
	}
}

type noisyConfig struct {
	Key     string `json:"private_key"`
	Network string `json:"network"`
	Base    string `json:"base_url"`
}

func (noisyConfig) ExchangeName() connector.ExchangeName { return "noisy" }
func (noisyConfig) Validate() error                      { return nil }

func TestCredentialFieldsFromConfig_FiltersNonCredentials(t *testing.T) {
	got := credentialFieldsFromConfig(noisyConfig{})
	if len(got) != 1 || got[0] != "private_key" {
		t.Fatalf("got %v, want [private_key]", got)
	}
}

type sortedConfig struct {
	B string `json:"z_last"`
	A string `json:"a_first"`
}

func (sortedConfig) ExchangeName() connector.ExchangeName { return "sorted" }
func (sortedConfig) Validate() error                      { return nil }

func TestCredentialFieldsFromConfig_Sorted(t *testing.T) {
	got := credentialFieldsFromConfig(sortedConfig{})
	if len(got) != 2 || got[0] != "a_first" || got[1] != "z_last" {
		t.Fatalf("unsorted or wrong: %v", got)
	}
}
