package config

import (
	"os"
	"path/filepath"
)

const (
	// WispConfigurationFileName is the historical project-local basename (wisp.yml).
	// Prefer ConnectorsFileName under ~/.wisp for credentials.
	WispConfigurationFileName = "wisp"

	// ConnectorsFileName is the global connector credentials file under WispHomeDir.
	ConnectorsFileName = "connectors.yml"

	// WispHomeDirName is the directory under the user home for CLI state + credentials.
	WispHomeDirName = ".wisp"

	// EnvSettingsPath overrides the connectors settings path when set.
	EnvSettingsPath = "WISP_SETTINGS"

	StrategiesDirectory = "strategies"
)

// WispHome returns ~/.wisp, creating it if needed.
func WispHome() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(home, WispHomeDirName)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return "", err
	}
	return dir, nil
}

// DefaultConnectorsPath is ~/.wisp/connectors.yml (canonical global credentials).
func DefaultConnectorsPath() (string, error) {
	dir, err := WispHome()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, ConnectorsFileName), nil
}

// ResolveSettingsPath picks where connector credentials live.
//
// Order:
//  1. explicit non-empty path (caller / --wisp flag)
//  2. $WISP_SETTINGS
//  3. ~/.wisp/connectors.yml if present
//  4. cwd wisp.yml or exchanges.yml (migration from project-local files)
//  5. default ~/.wisp/connectors.yml (may not exist yet — create via CLI Settings)
func ResolveSettingsPath(explicit string) string {
	if explicit != "" {
		return explicit
	}
	if env := os.Getenv(EnvSettingsPath); env != "" {
		return env
	}
	if def, err := DefaultConnectorsPath(); err == nil {
		if fileExists(def) {
			return def
		}
	}
	if cwd, err := os.Getwd(); err == nil {
		for _, name := range []string{
			WispConfigurationFileName + ".yml",
			"exchanges.yml",
		} {
			p := filepath.Join(cwd, name)
			if fileExists(p) {
				return p
			}
		}
	}
	if def, err := DefaultConnectorsPath(); err == nil {
		return def
	}
	return ConnectorsFileName
}

func fileExists(path string) bool {
	st, err := os.Stat(path)
	return err == nil && !st.IsDir()
}
