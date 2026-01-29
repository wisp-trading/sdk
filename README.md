# Wisp SDK

Build algorithmic trading strategies in Go with a clean, intuitive API.

## Quick Start

```go
package main

import (
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/wisp"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
	"github.com/wisp-trading/sdk/pkg/types/strategy"
)

type MyStrategy struct {
	k wisp.Wisp
}

func (s *MyStrategy) GetSignals() ([]*strategy.Signal, error) {
	btc := s.k.Asset("BTC")
	rsi, _ := s.k.Indicators().RSI(btc, 14)

	if rsi.LessThan(numerical.NewFromInt(30)) {
		signal := s.k.Signal(s.GetName()).
			Buy(btc, connector.Binance, numerical.NewFromFloat(0.1)).
			Build()
		return []*strategy.Signal{signal}, nil
	}

	return nil, nil
}

func (s *MyStrategy) GetName() strategy.StrategyName         { return "MyStrategy" }
func (s *MyStrategy) GetDescription() string                 { return "RSI momentum" }
func (s *MyStrategy) GetRiskLevel() strategy.RiskLevel       { return strategy.RiskLevelMedium }
func (s *MyStrategy) GetStrategyType() strategy.StrategyType { return strategy.StrategyTypeTechnical }
```

## Installation

```bash
# Clone the repository
git clone https://github.com/wisp-trading/sdk

# Get dependencies
go mod download
```

## Features

- **Built-in Indicators** - RSI, MACD, Moving Averages, Bollinger Bands, Stochastic, ATR
- **Multi-Exchange Support** - Binance, Bybit, Hyperliquid, Paradex
- **Automatic Data Management** - No manual data fetching or caching
- **Type-Safe API** - Compile-time safety for assets, exchanges, and decimals
- **Backtest & Live** - Same code runs in both modes

## Documentation

**📚 [Full Documentation](https://wisp-docs.vercel.app)** (coming soon)

- [Getting Started](docs/getting-started/)
- [Writing Strategies](docs/getting-started/writing-strategies.md)
- [Examples](docs/examples/)
- [API Reference](docs/api/)

## Examples

- [Simple RSI](examples/rsi/) - Basic momentum strategy
- [Moving Average Crossover](examples/ma-crossover/) - Trend following
- [Mean Reversion](examples/mean_reversion/) - Bollinger Bands strategy
- [Multi-Indicator](examples/multi-indicator/) - Confirmation signals
- [Arbitrage](examples/arbitrage/) - Cross-exchange trading
- [Portfolio](examples/portfolio/) - Multi-asset strategies

## Project Structure

```
wisp/
├── pkg/
│   ├── wisp/          # Main SDK (Indicators, Market, Analytics, Signal)
│   └── types/           # Public interfaces (Strategy, Connector, Portfolio)
├── internal/            # Internal implementations
├── examples/            # Example strategies
└── docs/                # Documentation site
```

## CLI Usage

```bash
# Backtest a strategy
wisp backtest --strategy MyStrategy --start 2024-01-01 --end 2024-12-31

# Run live
wisp live --strategy MyStrategy --exchange binance

# List available strategies
wisp list
```

## Contributing

Contributions welcome! See [CONTRIBUTING.md](CONTRIBUTING.md)

## License

MIT License - see [LICENSE](LICENSE)
