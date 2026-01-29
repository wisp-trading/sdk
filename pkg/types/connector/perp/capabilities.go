package perp

import (
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
)

// FundingRateProvider handles funding rate data (perps only)
type FundingRateProvider interface {
	FetchCurrentFundingRates() (map[portfolio.Asset]FundingRate, error)
	FetchFundingRate(asset portfolio.Asset) (*FundingRate, error)
	FetchHistoricalFundingRates(asset portfolio.Asset, startTime, endTime int64) ([]HistoricalFundingRate, error)
}

// PositionManager handles leveraged positions (perps only)
type PositionManager interface {
	GetPositions() ([]connector.Position, error)
}

// ContractProvider handles contract/derivative specifications (perps only)
type ContractProvider interface {
	FetchContracts() ([]connector.ContractInfo, error)
	FetchRiskFundBalance(symbol string) (*connector.RiskFundBalance, error)
}
