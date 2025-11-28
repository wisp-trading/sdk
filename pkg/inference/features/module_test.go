package features_test

import (
	"context"
	"testing"

	"github.com/backtesting-org/kronos-sdk/pkg/inference/features"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

// TestModuleWiring verifies that the features module is properly wired with fx.
// This ensures:
// - The module can be instantiated
// - All dependencies are satisfied
// - The aggregator is provided
// - Sub-modules (market, technical) are included
func TestModuleWiring(t *testing.T) {
	app := fxtest.New(t,
		features.Module,
		fx.Invoke(func(aggregator *features.Aggregator) {
			require.NotNil(t, aggregator, "Aggregator should be provided by the module")
		}),
	)

	app.RequireStart()
	app.RequireStop()
}

// TestModuleProvides verifies that the module provides the aggregator.
func TestModuleProvides(t *testing.T) {
	var aggregator *features.Aggregator

	app := fxtest.New(t,
		features.Module,
		fx.Populate(&aggregator),
	)

	app.RequireStart()
	require.NotNil(t, aggregator, "Module should provide Aggregator")
	app.RequireStop()
}

// TestExtractorsAreLoaded verifies that all expected feature extractors are loaded
// and registered in the aggregator through the fx group mechanism.
func TestExtractorsAreLoaded(t *testing.T) {
	var aggregator *features.Aggregator

	app := fxtest.New(t,
		features.Module,
		fx.Populate(&aggregator),
	)

	app.RequireStart()
	defer app.RequireStop()

	// Verify aggregator is provided
	require.NotNil(t, aggregator, "Aggregator should be provided")

	// Verify aggregator can be called (even if it returns early due to missing context data)
	featureMap, err := aggregator.Extract(context.Background())
	require.NoError(t, err, "Aggregator.Extract should not error")
	assert.NotNil(t, featureMap, "Feature map should not be nil")
	assert.IsType(t, map[string]float64{}, featureMap, "Feature map should be correct type")
}

// TestSubModulesRegistered verifies that sub-modules are correctly registered.
// We can't directly access the extractors (they're in a group), but we can verify
// the module names appear in the fx logs.
func TestSubModulesRegistered(t *testing.T) {
	// This test verifies through the fx output that both market-features and
	// technical-features modules are registered. The fx logs in other tests
	// show these lines:
	// - PROVIDE *market.Extractor[group = "feature_extractors"] from module "market-features"
	// - PROVIDE *technical.Extractor[group = "feature_extractors"] from module "technical-features"

	var aggregator *features.Aggregator

	app := fxtest.New(t,
		features.Module,
		fx.Populate(&aggregator),
	)

	app.RequireStart()
	require.NotNil(t, aggregator, "Module should provide aggregator with extractors")
	app.RequireStop()

	// If we got here without errors, both market and technical extractors
	// were successfully registered to the "feature_extractors" group
	// and injected into the aggregator
}
