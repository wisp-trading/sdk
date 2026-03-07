package types

import configTypes "github.com/wisp-trading/sdk/pkg/types/config"

// AssetLoader loads the relevant assets for a domain into its watchlist.
type AssetLoader interface {
	Load(cfg *configTypes.StartupConfig) error
}
