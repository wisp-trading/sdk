package prediction

import (
	"fmt"
	"strings"

	"github.com/wisp-trading/sdk/pkg/types/portfolio"
)

// PredictionPair represents a tradeable prediction market outcome.
// Format: {market-slug}:{outcome-name}:{quote-symbol} (e.g., "trump-2024:YES:USDC")
type PredictionPair struct {
	portfolio.Pair
	market  string
	outcome string
}

// NewPredictionPair creates a new prediction pair for a specific outcome.
// Uses the outcome name as the base asset and quote asset from portfolio.
func NewPredictionPair(market, outcome string, quote portfolio.Asset, separator ...string) PredictionPair {
	// Create base asset from outcome name
	outcomeName := fmt.Sprintf("%s:%s", market, outcome)

	baseAsset := portfolio.NewAsset(outcomeName)

	// Create embedded Pair with market slug as base symbol prefix
	basePair := portfolio.NewPair(baseAsset, quote, separator...)

	return PredictionPair{
		Pair:    basePair,
		market:  market,
		outcome: outcome,
	}
}

func PredictionPairFromPair(pair portfolio.Pair) (PredictionPair, error) {
	market, outcome := SplitAssetSymbol(pair.Base().Symbol(), ":")

	if market == "" || outcome == "" {
		return PredictionPair{}, fmt.Errorf("invalid pair symbol format: expected 'market:outcome', got '%s'", pair.Base().Symbol())
	}


	return PredictionPair{
		Pair:    pair,
		market:  market,
		outcome: outcome,
	}, nil
}

func SplitAssetSymbol(symbol, separator string) (market, outcome string) {
	parts := strings.SplitN(symbol, separator, 2)
	if len(parts) != 2 {
		return "", ""
	}

	return parts[0], parts[1]
}

func (p *PredictionPair) Market() string {
	return p.market
}

func (p *PredictionPair) Outcome() string {
	return p.outcome
}

func (p *PredictionPair) Validate() error {
	if p.market == "" {
		return fmt.Errorf("prediction pair must have a market set")
	}

	if p.outcome == "" {
		return fmt.Errorf("prediction pair must have an outcome set")
	}

	return nil
}
