package store

import (
	"strconv"
	"strings"
	"sync"

	"github.com/wisp-trading/sdk/pkg/markets/base/store/extensions"
	marketStore "github.com/wisp-trading/sdk/pkg/markets/base/types/stores/market"
	optionsTypes "github.com/wisp-trading/sdk/pkg/markets/options/types"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	"github.com/wisp-trading/sdk/pkg/types/temporal"
)

type optionsStore struct {
	marketStore.TradesStoreExtension
	marketStore.PositionsStoreExtension

	mu sync.RWMutex

	// Market data
	markPrices      map[string]float64
	underlyingPrice map[string]float64
	greeks          map[string]optionsTypes.Greeks
	iv              map[string]float64

	// Positions
	positions map[string]optionsTypes.Position
}

// NewStore creates a new options market store
func NewStore(timeProvider temporal.TimeProvider) optionsTypes.OptionsStore {
	return &optionsStore{
		TradesStoreExtension:    extensions.NewTradesExtension(),
		PositionsStoreExtension: extensions.NewPositionsExtension(),
		markPrices:              make(map[string]float64),
		underlyingPrice:      make(map[string]float64),
		greeks:               make(map[string]optionsTypes.Greeks),
		iv:                   make(map[string]float64),
		positions:            make(map[string]optionsTypes.Position),
	}
}

func (s *optionsStore) contractKey(contract optionsTypes.OptionContract) string {
	return contract.Pair.Symbol() + ":" + contract.Expiration.String() + ":" + strings.ToUpper(contract.OptionType) + ":" + floatToString(contract.Strike)
}

func floatToString(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

// GetPosition returns the position for a contract
func (s *optionsStore) GetPosition(contract optionsTypes.OptionContract) *optionsTypes.Position {
	s.mu.RLock()
	defer s.mu.RUnlock()

	key := s.contractKey(contract)
	if pos, ok := s.positions[key]; ok {
		return &pos
	}
	return nil
}

// SetPosition sets the position for a contract
func (s *optionsStore) SetPosition(contract optionsTypes.OptionContract, position optionsTypes.Position) {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := s.contractKey(contract)
	s.positions[key] = position
}

// GetAllPositions returns all positions in the store
func (s *optionsStore) GetAllPositions() []optionsTypes.Position {
	s.mu.RLock()
	defer s.mu.RUnlock()

	positions := make([]optionsTypes.Position, 0, len(s.positions))
	for _, pos := range s.positions {
		positions = append(positions, pos)
	}
	return positions
}

// QueryPositions returns positions matching the given query (exchange and/or pair filter).
func (s *optionsStore) QueryPositions(q marketStore.ActivityQuery) []optionsTypes.Position {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var out []optionsTypes.Position
	for _, pos := range s.positions {
		if q.Exchange != nil && pos.Exchange != *q.Exchange {
			continue
		}
		if q.Pair != nil && pos.Contract.Pair.Symbol() != q.Pair.Symbol() {
			continue
		}
		out = append(out, pos)
	}
	return out
}

// GetMarkPrice returns the mark price for a contract
func (s *optionsStore) GetMarkPrice(contract optionsTypes.OptionContract) float64 {
	s.mu.RLock()
	defer s.mu.RUnlock()

	key := s.contractKey(contract)
	return s.markPrices[key]
}

// SetMarkPrice sets the mark price for a contract
func (s *optionsStore) SetMarkPrice(contract optionsTypes.OptionContract, price float64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := s.contractKey(contract)
	s.markPrices[key] = price
}

// GetUnderlyingPrice returns the underlying price for a contract
func (s *optionsStore) GetUnderlyingPrice(contract optionsTypes.OptionContract) float64 {
	s.mu.RLock()
	defer s.mu.RUnlock()

	key := s.contractKey(contract)
	return s.underlyingPrice[key]
}

// SetUnderlyingPrice sets the underlying price for a contract
func (s *optionsStore) SetUnderlyingPrice(contract optionsTypes.OptionContract, price float64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := s.contractKey(contract)
	s.underlyingPrice[key] = price
}

// GetGreeks returns the Greeks for a contract
func (s *optionsStore) GetGreeks(contract optionsTypes.OptionContract) optionsTypes.Greeks {
	s.mu.RLock()
	defer s.mu.RUnlock()

	key := s.contractKey(contract)
	return s.greeks[key]
}

// SetGreeks sets the Greeks for a contract
func (s *optionsStore) SetGreeks(contract optionsTypes.OptionContract, greeks optionsTypes.Greeks) {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := s.contractKey(contract)
	s.greeks[key] = greeks
}

// GetIV returns the implied volatility for a contract
func (s *optionsStore) GetIV(contract optionsTypes.OptionContract) float64 {
	s.mu.RLock()
	defer s.mu.RUnlock()

	key := s.contractKey(contract)
	return s.iv[key]
}

// SetIV sets the implied volatility for a contract
func (s *optionsStore) SetIV(contract optionsTypes.OptionContract, iv float64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := s.contractKey(contract)
	s.iv[key] = iv
}

// GetPortfolioGreeks returns aggregated Greeks across all positions
func (s *optionsStore) GetPortfolioGreeks() optionsTypes.Greeks {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := optionsTypes.Greeks{}

	for key, position := range s.positions {
		if greeks, ok := s.greeks[key]; ok {
			result.Delta += greeks.Delta * position.Quantity
			result.Gamma += greeks.Gamma * position.Quantity
			result.Theta += greeks.Theta * position.Quantity
			result.Vega += greeks.Vega * position.Quantity
			result.Rho += greeks.Rho * position.Quantity
		}
	}

	return result
}

// MarketType returns the market type
func (s *optionsStore) MarketType() connector.MarketType {
	return connector.MarketTypeOptions
}

// UpdatePairPrice updates the price for a pair
func (s *optionsStore) UpdatePairPrice(pair portfolio.Pair, exchange connector.ExchangeName, price connector.Price) {
	// Not applicable for options market
}

// UpdatePairPrices updates prices for a pair
func (s *optionsStore) UpdatePairPrices(pair portfolio.Pair, prices marketStore.PriceMap) {
	// Not applicable for options market
}

// GetPairPrice returns the price for a pair
func (s *optionsStore) GetPairPrice(pair portfolio.Pair, exchange connector.ExchangeName) *connector.Price {
	// Not applicable for options market
	return nil
}

// GetPairPrices returns prices for a pair
func (s *optionsStore) GetPairPrices(pair portfolio.Pair) marketStore.PriceMap {
	// Not applicable for options market
	return make(marketStore.PriceMap)
}

// GetLastUpdated returns the last updated map
func (s *optionsStore) GetLastUpdated() marketStore.LastUpdatedMap {
	return make(marketStore.LastUpdatedMap)
}

// UpdateLastUpdated updates the last updated time
func (s *optionsStore) UpdateLastUpdated(key marketStore.UpdateKey) {
	// Not needed for options market
}
