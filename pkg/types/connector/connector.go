package connector

type Exchange struct {
	Name       ExchangeName `json:"name"`
	MarketType MarketType   `json:"market_type"`
}

type ExchangeName string

func (e ExchangeName) String() string {
	return string(e)
}

type Instrument string

const (
	TypeSpot       Instrument = "spot"
	TypePerpetual  Instrument = "perpetual"
	TypePrediction Instrument = "prediction"
)

// MarketType identifies the type of market for generic components (ingestors, stores, etc.)
type MarketType string

func (m MarketType) String() string {
	return string(m)
}

const (
	MarketTypeSpot       MarketType = "spot"
	MarketTypePerp       MarketType = "perp"
	MarketTypeFutures    MarketType = "futures"
	MarketTypeOptions    MarketType = "options"
	MarketTypePrediction MarketType = "prediction"
	MarketTypeNFT        MarketType = "nft"
	MarketTypeUnknown    MarketType = "unknown"
)
