package connector

import (
	"time"

	"github.com/shopspring/decimal"
)

type FundingRate struct {
	CurrentRate     decimal.Decimal
	NextFundingTime time.Time
	Timestamp       time.Time
	MarkPrice       decimal.Decimal
	IndexPrice      decimal.Decimal
	Premium         decimal.Decimal
}

type HistoricalFundingRate struct {
	FundingRate decimal.Decimal
	Timestamp   time.Time
}
