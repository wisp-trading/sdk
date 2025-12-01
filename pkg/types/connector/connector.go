package connector

import (
	"github.com/backtesting-org/kronos-sdk/pkg/types/kronos/numerical"
	"github.com/backtesting-org/kronos-sdk/pkg/types/portfolio"
)

type Exchange struct {
	Name ExchangeName `json:"name"`
}

type ExchangeName string

type Instrument string

const (
	TypeSpot      Instrument = "spot"
	TypePerpetual Instrument = "perpetual"
)

// Connector defines the interface for comprehensive exchange operations including market data, trading, and account management.
type Connector interface {
	FetchRiskFundBalance(symbol string) (*RiskFundBalance, error)
	FetchContracts() ([]ContractInfo, error)
	FetchPrice(symbol string) (*Price, error)
	FetchKlines(symbol, interval string, limit int) ([]Kline, error)
	FetchOrderBook(symbol portfolio.Asset, instrumentType Instrument, depth int) (*OrderBook, error)
	FetchRecentTrades(symbol string, limit int) ([]Trade, error)

	PlaceLimitOrder(symbol string, side OrderSide, quantity, price numerical.Decimal) (*OrderResponse, error)
	PlaceMarketOrder(symbol string, side OrderSide, quantity numerical.Decimal) (*OrderResponse, error)
	CancelOrder(symbol, orderID string) (*CancelResponse, error)
	GetOpenOrders() ([]Order, error)
	GetOrderStatus(orderID string) (*Order, error)

	FetchCurrentFundingRates() (map[portfolio.Asset]FundingRate, error)
	FetchFundingRate(asset portfolio.Asset) (*FundingRate, error)

	FetchHistoricalFundingRates(asset portfolio.Asset, startTime, endTime int64) ([]HistoricalFundingRate, error)
	GetAccountBalance() (*AccountBalance, error)
	GetPositions() ([]Position, error)
	GetTradingHistory(symbol string, limit int) ([]Trade, error)

	FetchAvailableSpotAssets() ([]portfolio.Asset, error)
	FetchAvailablePerpetualAssets() ([]portfolio.Asset, error)

	GetConnectorInfo() *Info
	GetPerpSymbol(symbol portfolio.Asset) string
	SupportsTradingOperations() bool
	SupportsRealTimeData() bool
	SupportsFundingRates() bool
	SupportsPerpetuals() bool
	SupportsSpot() bool

	Initialize(config Config) error
	IsInitialized() bool
}
