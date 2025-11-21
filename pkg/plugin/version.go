package plugin

import (
	"fmt"
	"plugin"

	"github.com/backtesting-org/kronos-sdk/pkg/version"
)

// extractSDKVersion extracts the SDK version from a loaded plugin
// Plugins must export: var SDKVersion = version.SDKVersion
func extractSDKVersion(p *plugin.Plugin) (string, error) {
	// Look up the SDK version variable that plugins export
	// Plugins should include: var SDKVersion = version.SDKVersion
	sym, err := p.Lookup("SDKVersion")
	if err != nil {
		return "", fmt.Errorf("plugin does not export SDKVersion variable. Plugins must include: var SDKVersion = version.SDKVersion")
	}

	// Type assert to string pointer
	versionPtr, ok := sym.(*string)
	if !ok {
		return "", fmt.Errorf("SDKVersion symbol has unexpected type: %T", sym)
	}

	return *versionPtr, nil
}

// validateSDKVersion enforces strict version matching
// Only exact version matches are allowed - no backward compatibility
func validateSDKVersion(pluginSDKVersion string) error {
	currentVersion := version.SDKVersion

	if pluginSDKVersion != currentVersion {
		return fmt.Errorf(
			"SDK version mismatch: plugin requires %s, but running SDK is %s. Please rebuild plugin with SDK %s",
			pluginSDKVersion,
			currentVersion,
			currentVersion,
		)
	}

	return nil
}
