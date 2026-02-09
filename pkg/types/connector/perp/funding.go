package perp

import (
	"time"

	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

type FundingRate struct {
	Pair            portfolio.Pair
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
