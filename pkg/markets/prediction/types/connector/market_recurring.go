package connector

import (
	"time"
)

// RecurringMarket represents an individual market instance that belongs to a recurring series
type RecurringMarket struct {
	RecurrenceInterval RecurrenceInterval `json:"recurrence_interval"`
}

// RecurrenceInterval represents how often new markets are created in a series
type RecurrenceInterval string

const (
	Recurrence1Min   RecurrenceInterval = "1m"
	Recurrence5Min   RecurrenceInterval = "5m"
	Recurrence15Min  RecurrenceInterval = "15m"
	Recurrence30Min  RecurrenceInterval = "30m"
	Recurrence1Hour  RecurrenceInterval = "1h"
	Recurrence4Hour  RecurrenceInterval = "4h"
	Recurrence1Day   RecurrenceInterval = "1d"
	Recurrence1Week  RecurrenceInterval = "1w"
	Recurrence1Month RecurrenceInterval = "1mo"
)

// Duration returns the time.Duration represented by this recurrence interval
func (r RecurrenceInterval) Duration() (time.Duration, bool) {
	switch r {
	case Recurrence1Min:
		return time.Minute, true
	case Recurrence5Min:
		return 5 * time.Minute, true
	case Recurrence15Min:
		return 15 * time.Minute, true
	case Recurrence30Min:
		return 30 * time.Minute, true
	case Recurrence1Hour:
		return time.Hour, true
	case Recurrence4Hour:
		return 4 * time.Hour, true
	case Recurrence1Day:
		return 24 * time.Hour, true
	case Recurrence1Week:
		return 7 * 24 * time.Hour, true
	case Recurrence1Month:
		// Approximation - 30 days
		return 30 * 24 * time.Hour, true
	default:
		return 0, false
	}
}

// String returns the string representation of the recurrence interval
func (r RecurrenceInterval) String() string {
	return string(r)
}
