package types

// OptionsAssetLoader discovers available options expirations for underlyings
type OptionsAssetLoader interface {
	LoadAssets(connectorRegistry interface{}) error
}

// UniverseProvider builds the current trading universe
type OptionsUniverseProvider interface {
	Universe() OptionsUniverse
}
