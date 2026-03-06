package market

import (
	"time"

	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
)

// KlineMap defines the structure for storing kline data, organized by exchange and pair
type KlineMap map[connector.ExchangeName]map[string][]connector.Kline

// KlineStoreExtension provides kline data storage
type KlineStoreExtension interface {
	StoreExtension
	UpdateKline(pair portfolio.Pair, exchange connector.ExchangeName, kline connector.Kline)
	GetKlines(pair portfolio.Pair, exchange connector.ExchangeName, interval string, limit int) []connector.Kline
	GetKlinesSince(pair portfolio.Pair, exchange connector.ExchangeName, interval string, since time.Time) []connector.Kline
}

// KlineWriter is a narrower interface for components that only write klines
type KlineWriter interface {
	UpdateKline(pair portfolio.Pair, exchange connector.ExchangeName, kline connector.Kline)
}
