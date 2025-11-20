---
sidebar_position: 5
---

# Configuration

Configure Kronos for your exchange, assets, and trading parameters.

## Configuration Structure

Kronos uses two types of configuration:

1. **Project-level config** (`config.yaml` in project root) - Exchanges, logging, risk settings
2. **Strategy-specific config** (`strategies/{name}/config.yaml`) - Strategy settings, backtest parameters

### Project Config (`config.yaml`)

```yaml
# Exchange configuration
exchanges:
  binance:
    api_key: "your-api-key"
    api_secret: "your-api-secret"
    testnet: false
  
  bybit:
    api_key: "your-api-key"
    api_secret: "your-api-secret"
    testnet: false
  
  hyperliquid:
    api_key: "your-api-key"
    api_secret: "your-api-secret"
    testnet: false

# Logging
logging:
  level: "info"            # debug, info, warn, error
  file: "kronos.log"

# Risk management
risk:
  max_position_size: 0.2
  max_total_exposure: 0.8
```

### Strategy Config (`strategies/my-strategy/config.yaml`)

```yaml
# Strategy configuration
strategy:
  name: "my-strategy"
  interval: "1h"           # How often GetSignals() is called
  default_exchange: binance
  assets:
    - "BTC"
    - "ETH"
    - "SOL"

# Backtesting configuration
backtest:
  start_date: "2024-01-01"
  end_date: "2024-12-31"
  initial_capital: 10000
  commission: 0.001        # 0.1% per trade
```

:::danger Never commit API keys
Never commit config files with API keys to version control. Keep sensitive configuration files secure and separate from your codebase.
:::

## Exchange Configuration

### Binance

```yaml
exchanges:
  binance:
    api_key: "your-api-key"
    api_secret: "your-api-secret"
    testnet: false           # Use Binance Testnet
    futures: true            # Trade futures (default: spot)
    margin: false            # Enable margin trading
```

### Bybit

```yaml
exchanges:
  bybit:
    api_key: "your-api-key"
    api_secret: "your-api-secret"
    testnet: false
    futures: true
    margin: false
```

### Hyperliquid

```yaml
exchanges:
  hyperliquid:
    api_key: "your-api-key"
    api_secret: "your-api-secret"
    testnet: false
```

## Strategy Configuration

Strategy configuration lives in `strategies/{name}/config.yaml`.

### Basic Settings

```yaml
strategy:
  name: "my-strategy"
  interval: "1h"              # 1m, 5m, 15m, 1h, 4h, 1d
  default_exchange: binance
  assets:
    - "BTC"
    - "ETH"
```

### Interval Options

How often `GetSignals()` is called:

- `1m` - Every minute
- `5m` - Every 5 minutes
- `15m` - Every 15 minutes
- `30m` - Every 30 minutes
- `1h` - Every hour (default)
- `4h` - Every 4 hours
- `1d` - Every day

### Multiple Exchanges

Trade on multiple exchanges:

```yaml
strategy:
  default_exchange: binance
  exchanges:
    - binance
    - bybit
    - hyperliquid
```

In your strategy:

```go
// Uses default exchange
price := s.k.Market().Price(btc)

// Specify exchange
price := s.k.Market().Price(btc, market.MarketOptions{
    Exchange: connector.Bybit,
})

// Get all prices
prices := s.k.Market().Prices(btc)
```



## Logging Configuration

### Log Level

```yaml
logging:
  level: "info"              # debug, info, warn, error
  file: "kronos.log"
  console: true              # Also log to console
```



## Risk Management

### Position Limits

```yaml
risk:
  max_position_size: 0.2     # Max 20% of capital per position
  max_total_exposure: 0.8    # Max 80% of capital deployed
  max_leverage: 3            # Max 3x leverage
```

### Stop Loss

```yaml
risk:
  default_stop_loss: 0.02    # Default 2% stop loss
  default_take_profit: 0.05  # Default 5% take profit
  trailing_stop: true
  trailing_stop_distance: 0.015  # 1.5%
```



## Example Production Config

```yaml
exchanges:
  binance:
    api_key: "your-api-key"
    api_secret: "your-api-secret"
    testnet: false
    futures: true

strategy:
  name: "production-strategy"
  interval: "1h"
  default_exchange: binance
  assets:
    - "BTC"
    - "ETH"

risk:
  max_position_size: 0.15
  max_total_exposure: 0.6
  default_stop_loss: 0.02
  default_take_profit: 0.05
  trailing_stop: true

logging:
  level: "info"
  file: "logs/production.log"
```

## Multiple Environments

Create different project-level configs for different environments:

```bash
config.yaml           # Default project config
config.dev.yaml       # Development
config.staging.yaml   # Staging
config.prod.yaml      # Production
```

Each strategy keeps its own config in `strategies/{name}/config.yaml`.

Use with:

```bash
kronos live --config config.prod.yaml
```

## Next Steps

- **[Writing Strategies](writing-strategies)** - Build your strategy
- **[Examples](examples)** - See complete implementations
- **[Installation](installation)** - Deploy to production
