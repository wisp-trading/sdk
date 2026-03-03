package types

import (
	"time"

	"github.com/google/uuid"
	predictionconnector "github.com/wisp-trading/sdk/pkg/markets/prediction/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/strategy"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

// PredictionSignalBuilder provides a fluent API for constructing prediction market trading signals.
type PredictionSignalBuilder interface {
	Buy(market predictionconnector.Market, outcome predictionconnector.Outcome, exchange connector.ExchangeName, shares, maxPrice numerical.Decimal, expiration int64) PredictionSignalBuilder
	Sell(market predictionconnector.Market, outcome predictionconnector.Outcome, exchange connector.ExchangeName, shares, minPrice numerical.Decimal, expiration int64) PredictionSignalBuilder
	Build() PredictionSignal
}

// PredictionSignal is the interface for prediction market signals.
type PredictionSignal interface {
	strategy.Signal
	GetActions() []*PredictionAction
}

type SignalFactory interface {
	NewPrediction(strategyName strategy.StrategyName) PredictionSignalBuilder
}

// predictionSignal is the unexported concrete implementation of PredictionSignal.
type predictionSignal struct {
	id        uuid.UUID
	strategy  strategy.StrategyName
	timestamp time.Time
	actions   []*PredictionAction
}

func (s *predictionSignal) GetID() uuid.UUID                   { return s.id }
func (s *predictionSignal) GetStrategy() strategy.StrategyName { return s.strategy }
func (s *predictionSignal) GetTimestamp() time.Time            { return s.timestamp }
func (s *predictionSignal) GetActions() []*PredictionAction    { return s.actions }

// NewPredictionSignal constructs a PredictionSignal. Used by the builder.
func NewPredictionSignal(id uuid.UUID, name strategy.StrategyName, ts time.Time, actions []*PredictionAction) PredictionSignal {
	return &predictionSignal{id: id, strategy: name, timestamp: ts, actions: actions}
}
