package main

import (
	"os"

	sdk "github.com/backtesting-org/kronos-sdk/pkg/kronos"
	"github.com/backtesting-org/kronos-sdk/pkg/types/connector"
	"github.com/backtesting-org/kronos-sdk/pkg/types/kronos/numerical"
	"github.com/backtesting-org/kronos-sdk/pkg/types/strategy"
	"gopkg.in/yaml.v3"
)

// ExchangesConfig holds the global exchanges configuration
type ExchangesConfig struct {
	Exchanges []struct {
		Name    string   `yaml:"name"`
		Enabled bool     `yaml:"enabled"`
		Assets  []string `yaml:"assets"`
	} `yaml:"exchanges"`
}

// MeanReversionConfig holds configuration for mean reversion strategy
type MeanReversionConfig struct {
	Exchange   string `yaml:"exchange"`
	Indicators struct {
		BollingerBands struct {
			Period int     `yaml:"period"`
			StdDev float64 `yaml:"std_dev"`
		} `yaml:"bollinger_bands"`
		RSI struct {
			Period     int `yaml:"period"`
			Oversold   int `yaml:"oversold"`
			Overbought int `yaml:"overbought"`
		} `yaml:"rsi"`
	} `yaml:"indicators"`
	Parameters struct {
		PositionSize float64 `yaml:"position_size"`
	} `yaml:"parameters"`
}

// MeanReversionStrategy implements a Bollinger Bands mean reversion strategy
type meanReversionStrategy struct {
	strategy.BaseStrategy
	k        *sdk.Kronos
	config   MeanReversionConfig
	exConfig ExchangesConfig
}

// NewMeanReversion creates a new mean reversion strategy instance
func NewMeanReversion(k *sdk.Kronos) strategy.Strategy {
	// Load configuration
	var config MeanReversionConfig
	data, err := os.ReadFile("config.yml")
	if err == nil {
		yaml.Unmarshal(data, &config)
	}

	// Load exchanges configuration
	var exConfig ExchangesConfig
	exData, err := os.ReadFile("exchanges.yml")
	if err == nil {
		yaml.Unmarshal(exData, &exConfig)
	}

	// Set defaults if not loaded
	if config.Indicators.BollingerBands.Period == 0 {
		config.Indicators.BollingerBands.Period = 20
		config.Indicators.BollingerBands.StdDev = 2.0
	}
	if config.Indicators.RSI.Period == 0 {
		config.Indicators.RSI.Period = 14
		config.Indicators.RSI.Oversold = 35
		config.Indicators.RSI.Overbought = 65
	}
	if config.Exchange == "" {
		config.Exchange = "binance"
	}
	if config.Parameters.PositionSize == 0 {
		config.Parameters.PositionSize = 0.1
	}

	return &meanReversionStrategy{k: k, config: config, exConfig: exConfig}
}

// GetSignals generates trading signals based on Bollinger Bands mean reversion
func (s *meanReversionStrategy) GetSignals() ([]*strategy.Signal, error) {
	// Determine asset symbol and exchange from exchanges.yml based on config
	exchangeStr := s.config.Exchange
	assetSymbol := "BTC" // default
	exchange := connector.ExchangeName(exchangeStr)

	// Check if exchange is enabled and get assets
	exFound := false
	for _, ex := range s.exConfig.Exchanges {
		if ex.Name == exchangeStr && ex.Enabled {
			exFound = true
			if len(ex.Assets) > 0 {
				// Parse "BTC/USDT" to get "BTC"
				assetSymbol = ex.Assets[0]
			}
			break
		}
	}

	if !exFound {
		s.k.Log().Info("MeanReversion", "Exchange %s not enabled in exchanges.yml, using defaults", exchangeStr)
		exchange = connector.Binance
	}

	asset := s.k.Asset(assetSymbol)
	quantity := numerical.NewFromFloat(s.config.Parameters.PositionSize)

	// Get indicators with config values
	bb, err := s.k.Indicators.BollingerBands(asset, s.config.Indicators.BollingerBands.Period, s.config.Indicators.BollingerBands.StdDev)
	if err != nil {
		s.k.Log().Debug("MeanReversion", assetSymbol, "Failed to get Bollinger Bands: %v", err)
		return nil, nil
	}

	price, err := s.k.Market.Price(asset)
	if err != nil {
		s.k.Log().Debug("MeanReversion", assetSymbol, "Failed to get price: %v", err)
		return nil, nil
	}

	rsi, err := s.k.Indicators.RSI(asset, s.config.Indicators.RSI.Period)
	if err != nil {
		s.k.Log().Debug("MeanReversion", assetSymbol, "Failed to get RSI: %v", err)
		return nil, nil
	}

	oversoldThreshold := numerical.NewFromInt(int64(s.config.Indicators.RSI.Oversold))
	overboughtThreshold := numerical.NewFromInt(int64(s.config.Indicators.RSI.Overbought))

	// Buy at lower band with RSI confirmation
	if price.LessThan(bb.Lower) && rsi.LessThan(oversoldThreshold) {
		s.k.Log().Opportunity("MeanReversion", assetSymbol,
			"Price below lower band (%.2f), RSI oversold (%.2f), targeting middle band (%.2f)",
			price, rsi, bb.Middle)

		signal := s.k.Signal(s.GetName()).
			Buy(asset, exchange, quantity).
			Build()
		return []*strategy.Signal{signal}, nil
	}

	// Sell at upper band with RSI confirmation
	if price.GreaterThan(bb.Upper) && rsi.GreaterThan(overboughtThreshold) {
		s.k.Log().Opportunity("MeanReversion", assetSymbol,
			"Price above upper band (%.2f), RSI overbought (%.2f), targeting middle band (%.2f)",
			price, rsi, bb.Middle)

		signal := s.k.Signal(s.GetName()).
			Sell(asset, exchange, quantity).
			Build()
		return []*strategy.Signal{signal}, nil
	}

	return nil, nil
}

// Interface implementation
func (s *meanReversionStrategy) GetName() strategy.StrategyName {
	return "Mean Reversion"
}

func (s *meanReversionStrategy) GetDescription() string {
	return "Bollinger Bands mean reversion with RSI confirmation"
}

func (s *meanReversionStrategy) GetRiskLevel() strategy.RiskLevel {
	return strategy.RiskLevelMedium
}

func (s *meanReversionStrategy) GetStrategyType() strategy.StrategyType {
	return strategy.StrategyTypeMeanReversion
}
