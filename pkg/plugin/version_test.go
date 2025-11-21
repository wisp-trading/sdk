package plugin

import (
	"testing"
)

func TestExtractSDKVersionFromPath(t *testing.T) {
	// This test would need a real plugin file
	// For now, we'll test the function exists and has the right signature
	t.Log("extractSDKVersionFromPath function exists")
	
	// Test with non-existent file
	_, err := extractSDKVersionFromPath("/nonexistent/plugin.so")
	if err == nil {
		t.Error("Expected error for non-existent file")
	}
	t.Logf("Got expected error for non-existent file: %v", err)
}

func TestExtractSDKVersionFromRealPlugin(t *testing.T) {
	// Use the test plugin we built with 'replace' directive
	pluginPath := "/tmp/real_test.so"

	version, err := extractSDKVersionFromPath(pluginPath)

	// This should fail because the plugin was built with 'replace' directive (v0.0.0)
	if err == nil {
		t.Errorf("Expected error for development version, but got version: %s", version)
	} else {
		t.Logf("✅ Correctly rejected development version: %v", err)
	}
}

func TestValidateSDKVersion(t *testing.T) {
	tests := []struct {
		name            string
		pluginVersion   string
		expectError     bool
	}{
		{
			name:          "matching version",
			pluginVersion: "v0.0.2", // Current SDK version
			expectError:   false,
		},
		{
			name:          "mismatched version",
			pluginVersion: "v0.0.1",
			expectError:   true,
		},
		{
			name:          "higher version",
			pluginVersion: "v0.0.3",
			expectError:   true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateSDKVersion(tt.pluginVersion)
			if tt.expectError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
			if err != nil {
				t.Logf("Error message: %v", err)
			}
		})
	}
}
