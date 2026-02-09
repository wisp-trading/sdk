package perp

import (
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
)

type AccountReader interface {
	GetMarginBalances() ([]AssetBalance, error)
}

// FundingRateProvider handles funding rate data (perps only)
type FundingRateProvider interface {
	FetchCurrentFundingRates() (map[portfolio.Pair]FundingRate, error)
	FetchFundingRate(pair portfolio.Pair) (*FundingRate, error)
	FetchHistoricalFundingRates(pair portfolio.Pair, startTime, endTime int64) ([]HistoricalFundingRate, error)
}

// PositionManager handles leveraged positions (perps only)
type PositionManager interface {
	GetPositions() ([]Position, error)
}

// ContractProvider handles contract/derivative specifications (perps only)
type ContractProvider interface {
	FetchContracts() ([]connector.ContractInfo, error)
	FetchRiskFundBalance(symbol string) (*RiskFundBalance, error)
}
