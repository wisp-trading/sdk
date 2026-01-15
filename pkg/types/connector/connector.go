package connector

type Exchange struct {
	Name ExchangeName `json:"name"`
}

type ExchangeName string

type Instrument string

const (
	TypeSpot      Instrument = "spot"
	TypePerpetual Instrument = "perpetual"
)

// MarketType identifies the type of market for generic components (ingestors, stores, etc.)
type MarketType string

const (
	MarketTypeSpot       MarketType = "spot"
	MarketTypePerp       MarketType = "perp"
	MarketTypeFutures    MarketType = "futures"
	MarketTypeOptions    MarketType = "options"
	MarketTypePrediction MarketType = "prediction"
	MarketTypeNFT        MarketType = "nft"
)
