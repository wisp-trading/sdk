package extensions

import (
	"sync"

	"github.com/wisp-trading/sdk/pkg/markets/base/types/stores/market"
	perpTypes "github.com/wisp-trading/sdk/pkg/markets/perp/types"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	perpConn "github.com/wisp-trading/sdk/pkg/types/connector/perp"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
)

type perpPositionsExtension struct {
	mu        sync.RWMutex
	positions map[positionKey]perpConn.Position
}

type positionKey struct {
	exchange connector.ExchangeName
	pair     string // Pair.Symbol()
}

func NewPerpPositionsExtension() perpTypes.PerpPositionsStoreExtension {
	return &perpPositionsExtension{
		positions: make(map[positionKey]perpConn.Position),
	}
}

func (e *perpPositionsExtension) UpsertPosition(position perpConn.Position) {
	e.mu.Lock()
	defer e.mu.Unlock()
	key := positionKey{exchange: position.Exchange, pair: position.Pair.Symbol()}
	e.positions[key] = position
}

func (e *perpPositionsExtension) RemovePosition(exchange connector.ExchangeName, pair portfolio.Pair) {
	e.mu.Lock()
	defer e.mu.Unlock()
	delete(e.positions, positionKey{exchange: exchange, pair: pair.Symbol()})
}

func (e *perpPositionsExtension) GetPositions() []perpConn.Position {
	e.mu.RLock()
	defer e.mu.RUnlock()
	out := make([]perpConn.Position, 0, len(e.positions))
	for _, p := range e.positions {
		out = append(out, p)
	}
	return out
}

func (e *perpPositionsExtension) GetPosition(exchange connector.ExchangeName, pair portfolio.Pair) *perpConn.Position {
	e.mu.RLock()
	defer e.mu.RUnlock()
	p, ok := e.positions[positionKey{exchange: exchange, pair: pair.Symbol()}]
	if !ok {
		return nil
	}
	return &p
}

func (e *perpPositionsExtension) QueryPositions(q market.ActivityQuery) []perpConn.Position {
	e.mu.RLock()
	defer e.mu.RUnlock()
	var out []perpConn.Position
	for _, p := range e.positions {
		if q.Exchange != nil && p.Exchange != *q.Exchange {
			continue
		}
		if q.Pair != nil && p.Pair.Symbol() != q.Pair.Symbol() {
			continue
		}
		out = append(out, p)
	}
	return out
}

var _ perpTypes.PerpPositionsStoreExtension = (*perpPositionsExtension)(nil)
