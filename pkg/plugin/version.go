package plugin

import (
	"fmt"
	"plugin"

	"github.com/backtesting-org/kronos-sdk/pkg/version"
)

// extractSDKVersion extracts the SDK version from a loaded plugin
// The plugin's compiled binary contains the version.SDKVersion constant
// since the plugin imports SDK packages
func extractSDKVersion(p *plugin.Plugin) (string, error) {
	// Look up the SDK version constant that's embedded in the plugin
	// When the plugin is built, it imports SDK packages which embed this constant
	sym, err := p.Lookup("github.com/backtesting-org/kronos-sdk/pkg/version.SDKVersion")
	if err != nil {
		return "", fmt.Errorf("plugin does not contain SDK version information: %w", err)
	}

	// Type assert to string pointer (constants are stored as pointers in plugin symbols)
	versionPtr, ok := sym.(*string)
	if !ok {
		return "", fmt.Errorf("SDK version symbol has unexpected type")
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
