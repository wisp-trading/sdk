package signal

import (
	"github.com/google/uuid"
	"github.com/wisp-trading/sdk/pkg/markets/prediction/types"
	predictionconnector "github.com/wisp-trading/sdk/pkg/markets/prediction/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/strategy"
	"github.com/wisp-trading/sdk/pkg/types/temporal"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

// predictionBuilder is the concrete implementation of strategy.PredictionSignalBuilder.
type predictionBuilder struct {
	strategyName strategy.StrategyName
	actions      []types.PredictionAction
	timeProvider temporal.TimeProvider
}

// Buy adds a buy action for a prediction market outcome.
// maxPrice is the maximum probability to pay (0.0–1.0).
// expiration is a Unix timestamp after which the order is cancelled.
func (b *predictionBuilder) Buy(market predictionconnector.Market, outcome predictionconnector.Outcome, exchange connector.ExchangeName, shares, maxPrice numerical.Decimal, expiration int64) types.PredictionSignalBuilder {
	b.actions = append(b.actions, types.PredictionAction{
		BaseAction: strategy.BaseAction{ActionType: strategy.ActionBuy, Exchange: exchange},
		Market:     market,
		Outcome:    outcome,
		Shares:     shares,
		MaxPrice:   maxPrice,
		Expiration: expiration,
	})
	return b
}

// Sell adds a sell action to exit a prediction market position.
// minPrice is the minimum probability to accept (0.0–1.0).
// expiration is a Unix timestamp after which the order is cancelled.
func (b *predictionBuilder) Sell(market predictionconnector.Market, outcome predictionconnector.Outcome, exchange connector.ExchangeName, shares, minPrice numerical.Decimal, expiration int64) types.PredictionSignalBuilder {
	b.actions = append(b.actions, types.PredictionAction{
		BaseAction: strategy.BaseAction{ActionType: strategy.ActionSell, Exchange: exchange},
		Market:     market,
		Outcome:    outcome,
		Shares:     shares,
		MaxPrice:   minPrice,
		Expiration: expiration,
	})
	return b
}

// Build constructs and returns the PredictionSignal.
func (b *predictionBuilder) Build() types.PredictionSignal {
	return types.NewPredictionSignal(uuid.New(), b.strategyName, b.timeProvider.Now(), b.actions)
}
