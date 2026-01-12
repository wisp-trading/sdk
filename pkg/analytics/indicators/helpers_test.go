package indicators_test

// Helper function to create float64 slices for indicator tests
func makeDecimals(values ...float64) []float64 {
	result := make([]float64, len(values))
	copy(result, values)
	return result
}
