package plugin

import (
	"github.com/backtesting-org/kronos-sdk/pkg/types/plugin"
	"github.com/backtesting-org/kronos-sdk/pkg/types/strategy"
)

// extractMetadata extracts metadata from a strategy instance
func extractMetadata(strat strategy.Strategy) *plugin.Metadata {
	metadata := &plugin.Metadata{
		Name:        string(strat.GetName()),
		Description: strat.GetDescription(),
		RiskLevel:   string(strat.GetRiskLevel()),
		Type:        string(strat.GetStrategyType()),
		Version:     "1.0.0", // Default version
		Parameters:  make(map[string]plugin.ParameterDef),
	}

	// Try to extract parameters if strategy implements ParameterProvider interface
	if paramProvider, ok := strat.(plugin.ParameterProvider); ok {
		params := paramProvider.GetParameters()
		for _, p := range params {
			metadata.Parameters[p.Name] = p
		}
	}

	return metadata
}
