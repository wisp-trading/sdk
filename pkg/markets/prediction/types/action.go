package types

import (
	"fmt"

	predictionConnector "github.com/wisp-trading/sdk/pkg/markets/prediction/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/strategy"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

// PredictionAction represents an action for prediction markets
type PredictionAction struct {
	strategy.BaseAction
	Market      predictionConnector.Market  `json:"market"`
	Outcome     predictionConnector.Outcome `json:"outcome"`
	Shares      numerical.Decimal           `json:"shares"`
	MaxPrice    numerical.Decimal           `json:"max_price"`      // Probability (0.0-1.0)
	Expiration  int64                       `json:"expiration"`     // Unix timestamp
	TimeInForce connector.TimeInForce       `json:"time_in_force"`  // GTC (default) | FOK | FAK
}

// GetMarketType returns prediction
func (a *PredictionAction) GetMarketType() connector.MarketType {
	return connector.MarketTypePrediction
}

// Validate checks if the prediction action is valid
func (a *PredictionAction) Validate() error {
	if err := a.ValidateBase(); err != nil {
		return err
	}
	if err := a.Market.Validate(); err != nil {
		return fmt.Errorf("invalid market: %w", err)
	}
	if a.Shares.IsZero() || a.Shares.IsNegative() {
		return fmt.Errorf("shares must be positive")
	}
	if a.MaxPrice.IsZero() || a.MaxPrice.IsNegative() || a.MaxPrice.GreaterThan(numerical.NewFromFloat(1.0)) {
		return fmt.Errorf("max price must be between 0 and 1")
	}
	return nil
}
