package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveSettingsPathExplicit(t *testing.T) {
	got := ResolveSettingsPath("/tmp/custom.yml")
	if got != "/tmp/custom.yml" {
		t.Fatalf("got %q", got)
	}
}

func TestResolveSettingsPathEnv(t *testing.T) {
	t.Setenv(EnvSettingsPath, "/tmp/from-env.yml")
	got := ResolveSettingsPath("")
	if got != "/tmp/from-env.yml" {
		t.Fatalf("got %q", got)
	}
}

func TestResolveSettingsPathCwdMigration(t *testing.T) {
	t.Setenv(EnvSettingsPath, "")
	dir := t.TempDir()
	cwd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(cwd) })
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	// Ensure home default does not exist for this process — still may exist on machine.
	// Create cwd wisp.yml so migration wins when default home file missing.
	// If ~/.wisp/connectors.yml already exists on the machine, it wins — that's OK.
	// Force by setting a fake home.
	fakeHome := t.TempDir()
	t.Setenv("HOME", fakeHome)
	path := filepath.Join(dir, "wisp.yml")
	if err := os.WriteFile(path, []byte("connectors: []\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	got := ResolveSettingsPath("")
	abs, _ := filepath.EvalSymlinks(path)
	gotAbs, _ := filepath.EvalSymlinks(got)
	if abs == "" {
		abs, _ = filepath.Abs(path)
	}
	if gotAbs == "" {
		gotAbs, _ = filepath.Abs(got)
	}
	if gotAbs != abs {
		t.Fatalf("expected cwd wisp.yml %q, got %q", abs, gotAbs)
	}
}
