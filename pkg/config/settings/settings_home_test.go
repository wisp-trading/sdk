package settings

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/wisp-trading/sdk/pkg/types/config"
)

func TestSaveAndLoadConnectorsUnderHome(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv(config.EnvSettingsPath, "")

	path, err := config.DefaultConnectorsPath()
	if err != nil {
		t.Fatal(err)
	}
	// Ensure we write to the resolved home path
	cfg := NewConfiguration(ConfigOptions{SettingsPath: path}).(*settings)

	err = cfg.AddConnector(config.Connector{
		Name:    "hyperliquid",
		Enabled: true,
		Credentials: map[string]string{
			"private_key":     "0xabc",
			"account_address": "0xdef",
		},
	})
	if err != nil {
		t.Fatalf("AddConnector: %v", err)
	}

	// Reload from disk
	cfg2 := NewConfiguration(ConfigOptions{SettingsPath: path}).(*settings)
	list, err := cfg2.GetConnectors()
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 || list[0].Name != "hyperliquid" {
		t.Fatalf("list=%v", list)
	}
	if list[0].Credentials["private_key"] != "0xabc" {
		t.Fatalf("credentials not loaded")
	}

	// File mode should be restrictive
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm()&0o077 != 0 {
		t.Fatalf("expected 0600-ish perms, got %o", info.Mode().Perm())
	}

	// Path lives under ~/.wisp
	if filepath.Base(filepath.Dir(path)) != config.WispHomeDirName {
		t.Fatalf("path not under .wisp: %s", path)
	}
}
