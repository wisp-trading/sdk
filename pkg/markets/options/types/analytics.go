package types

// OptionsAnalytics provides Greeks-based risk analytics and insights
type OptionsAnalytics interface {
	// Greeks exposure analysis
	GetDeltaExposure() float64
	GetGammaExposure() float64
	GetThetaExposure() float64
	GetVegaExposure() float64

	// Portfolio Greeks
	GetPortfolioGreeks() Greeks

	// Risk metrics
	GetDailyThetaDecay() float64
	GetIVSensitivity() float64
}
