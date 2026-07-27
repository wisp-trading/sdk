package strategy

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindStrategiesSurfacesBrokenConfig(t *testing.T) {
	dir := t.TempDir()
	cwd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(cwd) })
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}

	stratDir := filepath.Join("strategies", "broken")
	if err := os.MkdirAll(stratDir, 0o755); err != nil {
		t.Fatal(err)
	}
	// Invalid YAML
	if err := os.WriteFile(filepath.Join(stratDir, "config.yml"), []byte(":\n  - bad"), 0o644); err != nil {
		t.Fatal(err)
	}

	okDir := filepath.Join("strategies", "ok")
	if err := os.MkdirAll(okDir, 0o755); err != nil {
		t.Fatal(err)
	}
	okYAML := "name: ok\nexchanges:\n  - hyperliquid\nassets: {}\n"
	if err := os.WriteFile(filepath.Join(okDir, "config.yml"), []byte(okYAML), 0o644); err != nil {
		t.Fatal(err)
	}

	svc := NewStrategyConfigService()
	list, err := svc.FindStrategies()
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 2 {
		t.Fatalf("want 2 strategies, got %d", len(list))
	}

	var sawBroken, sawOK bool
	for _, s := range list {
		switch s.Name {
		case "broken":
			sawBroken = true
			if s.Error == "" {
				t.Fatal("broken strategy must set Error")
			}
		case "ok":
			sawOK = true
			if s.Error != "" {
				t.Fatalf("ok strategy Error: %s", s.Error)
			}
		}
	}
	if !sawBroken || !sawOK {
		t.Fatalf("missing entries: broken=%v ok=%v list=%+v", sawBroken, sawOK, list)
	}
}

func TestFindStrategiesEmptyNameUsesDir(t *testing.T) {
	dir := t.TempDir()
	cwd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(cwd) })
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	stratDir := filepath.Join("strategies", "from_dir")
	if err := os.MkdirAll(stratDir, 0o755); err != nil {
		t.Fatal(err)
	}
	// no name field
	if err := os.WriteFile(filepath.Join(stratDir, "config.yml"), []byte("exchanges: [hyperliquid]\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	svc := NewStrategyConfigService()
	list, err := svc.FindStrategies()
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 || list[0].Name != "from_dir" {
		t.Fatalf("want name from_dir, got %+v", list)
	}
}
