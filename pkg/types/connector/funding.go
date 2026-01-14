package connector

import (
	"time"

	"github.com/backtesting-org/kronos-sdk/pkg/types/kronos/numerical"
	"github.com/backtesting-org/kronos-sdk/pkg/types/portfolio"
)

type FundingRate struct {
	Asset           portfolio.Asset
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
