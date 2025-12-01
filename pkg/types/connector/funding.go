package connector

import (
	"time"

	"github.com/backtesting-org/kronos-sdk/pkg/types/kronos/numerical"
)

type FundingRate struct {
	CurrentRate     numerical.Decimal
	NextFundingTime time.Time
	Timestamp       time.Time
	MarkPrice       numerical.Decimal
	IndexPrice      numerical.Decimal
	Premium         numerical.Decimal
}

type HistoricalFundingRate struct {
	FundingRate numerical.Decimal
	Timestamp   time.Time
}
