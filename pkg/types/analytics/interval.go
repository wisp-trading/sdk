package analytics

// Interval represents time intervals for kline/candlestick data.
// These constants are used across the SDK for fetching market data,
// calculating indicators, and analyzing price movements.
const (
	Interval1Minute  = "1m"
	Interval5Minute  = "5m"
	Interval15Minute = "15m"
	Interval30Minute = "30m"
	Interval1Hour    = "1h"
	Interval2Hour    = "2h"
	Interval4Hour    = "4h"
	Interval6Hour    = "6h"
	Interval12Hour   = "12h"
	Interval1Day     = "1d"
	Interval3Day     = "3d"
	Interval1Week    = "1w"
	Interval1Month   = "1M"
)

// DefaultInterval is the default interval used when none is specified
const DefaultInterval = Interval1Hour

// PeriodsPerYear maps each interval to the number of periods in a year.
// Used for annualizing volatility and other time-based calculations.
var PeriodsPerYear = map[string]float64{
	Interval1Minute:  525600.0, // 60 * 24 * 365
	Interval5Minute:  105120.0, // 12 * 24 * 365
	Interval15Minute: 35040.0,  // 4 * 24 * 365
	Interval30Minute: 17520.0,  // 2 * 24 * 365
	Interval1Hour:    8760.0,   // 24 * 365
	Interval2Hour:    4380.0,   // 12 * 365
	Interval4Hour:    2190.0,   // 6 * 365
	Interval6Hour:    1460.0,   // 4 * 365
	Interval12Hour:   730.0,    // 2 * 365
	Interval1Day:     365.0,    // 365
	Interval3Day:     121.67,   // 365 / 3
	Interval1Week:    52.14,    // 365 / 7
	Interval1Month:   12.0,     // 12 months
}
