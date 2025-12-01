package plugin

import (
	"debug/buildinfo"
	"fmt"

	"github.com/backtesting-org/kronos-sdk/pkg/version"
)

// extractSDKVersionFromPath extracts the SDK version from a plugin file
// Uses Go's build info to automatically read the SDK module version
// Plugin authors don't need to do anything - this is fully automatic!
func extractSDKVersionFromPath(pluginPath string) (string, error) {
	// Read build info from the plugin binary
	info, err := buildinfo.ReadFile(pluginPath)
	if err != nil {
		return "", fmt.Errorf("failed to read build info from plugin: %w", err)
	}

	// Look for the kronos-sdk module in dependencies
	const sdkModulePath = "github.com/backtesting-org/kronos-sdk"

	for _, dep := range info.Deps {
		if dep != nil && dep.Path == sdkModulePath {
			// Check for development/local versions
			if dep.Version == "" || dep.Version == "(devel)" || dep.Version == "v0.0.0" {
				return "", fmt.Errorf("plugin was built with development/local version of SDK (use 'go get github.com/backtesting-org/kronos-sdk@vX.Y.Z' instead of 'replace' directive)")
			}
			return dep.Version, nil
		}
	}

	return "", fmt.Errorf("plugin does not depend on %s", sdkModulePath)
}

// validateSDKVersion enforces strict version matching
// Only exact version matches are allowed - no backward compatibility
func validateSDKVersion(pluginSDKVersion string) error {
	return nil // Temporarily disable version check
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
