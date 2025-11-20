package plugin

import (
	"fmt"
	"plugin"
)

// ValidatePluginFile validates that a file is a valid Go plugin
func ValidatePluginFile(filePath string) error {
	// Try to open the plugin
	p, err := plugin.Open(filePath)
	if err != nil {
		return fmt.Errorf("invalid plugin file: %w", err)
	}

	// Check for required symbols
	_, err = p.Lookup("NewStrategy")
	if err != nil {
		// Try alternative symbol
		_, err = p.Lookup("Strategy")
		if err != nil {
			return fmt.Errorf("plugin must export NewStrategy function or Strategy variable")
		}
	}

	return nil
}
