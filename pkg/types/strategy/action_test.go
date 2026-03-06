package strategy_test

import (
	"testing"

	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
	"github.com/wisp-trading/sdk/pkg/types/strategy"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

func TestBaseActionEmbedding(t *testing.T) {
	// Test that BaseAction methods are inherited by SpotAction
	btc := portfolio.NewAsset("BTC")
	usdt := portfolio.NewAsset("USDT")
	pair := portfolio.NewPair(btc, usdt)

	action := &strategy.SpotAction{
		BaseAction: strategy.BaseAction{
			ActionType: strategy.ActionBuy,
			Exchange:   connector.ExchangeName("binance"),
		},
		Pair:     pair,
		Quantity: numerical.NewFromFloat(0.5),
		Price:    numerical.NewFromFloat(50000),
	}

	// Test inherited methods from BaseAction
	if action.GetType() != strategy.ActionBuy {
		t.Errorf("Expected ActionBuy, got %v", action.GetType())
	}

	if action.GetExchange() != "binance" {
		t.Errorf("Expected binance, got %v", action.GetExchange())
	}

	// Test specific method
	if action.GetMarketType() != "spot" {
		t.Errorf("Expected spot, got %v", action.GetMarketType())
	}

	// Test validation
	if err := action.Validate(); err != nil {
		t.Errorf("Expected valid action, got error: %v", err)
	}
}

func TestPerpActionWithLeverage(t *testing.T) {
	eth := portfolio.NewAsset("ETH")
	usdt := portfolio.NewAsset("USDT")
	pair := portfolio.NewPair(eth, usdt)

	action := &strategy.PerpAction{
		BaseAction: strategy.BaseAction{
			ActionType: strategy.ActionBuy,
			Exchange:   connector.ExchangeName("bybit"),
		},
		Pair:     pair,
		Quantity: numerical.NewFromFloat(10),
		Price:    numerical.NewFromFloat(3000),
		Leverage: numerical.NewFromInt(10),
	}

	if action.GetMarketType() != connector.MarketTypePerp {
		t.Errorf("Expected perpetual, got %v", action.GetMarketType())
	}

	if err := action.Validate(); err != nil {
		t.Errorf("Expected valid action, got error: %v", err)
	}
}

func TestPerpActionInvalidLeverage(t *testing.T) {
	eth := portfolio.NewAsset("ETH")
	usdt := portfolio.NewAsset("USDT")
	pair := portfolio.NewPair(eth, usdt)

	// Test leverage too high
	action := &strategy.PerpAction{
		BaseAction: strategy.BaseAction{
			ActionType: strategy.ActionBuy,
			Exchange:   connector.ExchangeName("bybit"),
		},
		Pair:     pair,
		Quantity: numerical.NewFromFloat(10),
		Price:    numerical.NewFromFloat(3000),
		Leverage: numerical.NewFromInt(200), // Invalid: > 125x
	}

	if err := action.Validate(); err == nil {
		t.Error("Expected validation error for leverage > 125x")
	}
}

func TestPolymorphicActionInterface(t *testing.T) {
	btc := portfolio.NewAsset("BTC")
	usdt := portfolio.NewAsset("USDT")
	pair := portfolio.NewPair(btc, usdt)

	// Create different action types
	actions := []strategy.Action{
		&strategy.SpotAction{
			BaseAction: strategy.BaseAction{
				ActionType: strategy.ActionBuy,
				Exchange:   "binance",
			},
			Pair:     pair,
			Quantity: numerical.NewFromFloat(1),
			Price:    numerical.NewFromFloat(50000),
		},
		&strategy.PerpAction{
			BaseAction: strategy.BaseAction{
				ActionType: strategy.ActionSell,
				Exchange:   "bybit",
			},
			Pair:     pair,
			Quantity: numerical.NewFromFloat(1),
			Price:    numerical.NewFromFloat(50000),
			Leverage: numerical.NewFromInt(5),
		},
	}

	// Test polymorphic interface
	for i, action := range actions {
		if err := action.Validate(); err != nil {
			t.Errorf("Action %d validation failed: %v", i, err)
		}

		if action.GetExchange() == "" {
			t.Errorf("Action %d has empty exchange", i)
		}

		if action.GetMarketType() == "" {
			t.Errorf("Action %d has empty market type", i)
		}
	}

	// Verify different market types
	if actions[0].GetMarketType() != "spot" {
		t.Errorf("Expected spot, got %v", actions[0].GetMarketType())
	}

	if actions[1].GetMarketType() != connector.MarketTypePerp {
		t.Errorf("Expected perpetual, got %v", actions[1].GetMarketType())
	}
}
