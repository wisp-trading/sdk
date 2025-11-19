---
sidebar_position: 1
---

# Installation

Get up and running with Kronos in minutes.

## Install Kronos CLI

Install via Homebrew:

```bash
brew install backtesting-org/tap/kronos
```

Verify the installation:

```bash
kronos --version
```

## Prerequisites

You'll need Go installed to write strategies:

- **Go 1.21+** - [Install Go](https://go.dev/doc/install)

Check your Go version:

```bash
go version
```

## Create a New Strategy

Initialize a new Kronos project:

```bash
kronos init my-strategy
cd my-strategy
```

This creates:
- `strategy.go` - Your strategy implementation
- `go.mod` - Go dependencies
- `config.yaml` - Kronos configuration
- `README.md` - Project documentation

## Project Structure

```
my-strategy/
├── strategy.go      # Your strategy code
├── go.mod           # Go module
├── config.yaml      # Kronos config
└── README.md        # Documentation
```

## Backtest Your Strategy

Test your strategy with historical data:

```bash
kronos backtest
```

This runs your strategy against past market data and shows performance metrics:
- Total return
- Win rate
- Maximum drawdown
- Sharpe ratio
- Number of trades

### Backtest Options

```bash
# Backtest specific date range
kronos backtest --start 2024-01-01 --end 2024-12-31

# Use specific exchange
kronos backtest --exchange binance

# Specify starting capital
kronos backtest --capital 10000

# Detailed output
kronos backtest --verbose
```

## Deploy to Live Trading

Once you're satisfied with backtest results, deploy your strategy:

```bash
kronos live
```

This starts your strategy in live trading mode. Kronos will:
1. Connect to your configured exchange(s)
2. Start monitoring markets in real-time
3. Execute trades based on your strategy signals
4. Log all activity

### Live Trading Options

```bash
# Dry run (paper trading - no real orders)
kronos live --dry-run

# Specific exchange
kronos live --exchange binance

# Enable verbose logging
kronos live --verbose

# Custom config file
kronos live --config production.yaml
```

## Configuration

Before running live, configure your exchange API keys in `config.yaml`:

```yaml
exchanges:
  binance:
    api_key: "your-api-key"
    api_secret: "your-api-secret"
    testnet: false
  
  bybit:
    api_key: "your-api-key"
    api_secret: "your-api-secret"
    testnet: false

strategy:
  name: "my-strategy"
  interval: "1h"
  assets:
    - "BTC"
    - "ETH"
```

:::warning
Never commit sensitive config files to version control. Keep config files with API keys out of your repository.
:::

## Quick Commands Reference

```bash
# Create new strategy
kronos init <name>

# Backtest strategy
kronos backtest

# Run live (paper trading)
kronos live --dry-run

# Run live (real trading)
kronos live

# View logs
kronos logs

# Check status
kronos status

# Stop running strategy
kronos stop
```

## Next Steps

- **[Quick Reference](quick-reference)** - Learn the Kronos API basics
- **[Writing Strategies](writing-strategies)** - Deep dive into strategy development
- **[Configuration](configuration)** - Detailed config options
