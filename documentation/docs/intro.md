---
sidebar_position: 1
---

# Welcome to Kronos SDK

Kronos SDK is a Go framework that makes building trading strategies simple. Write your strategy logic once, and Kronos handles all the complexity—market data, exchange APIs, indicators, and execution.

## The Kronos Way

Instead of managing exchange connections, API calls, and data pipelines, you write this:

```go
func (s *Strategy) GetSignals() ([]*strategy.Signal, error) {
    btc := s.k.Asset("BTC")
    
    // Just ask for what you need - Kronos figures out the rest
    rsi, _ := s.k.Indicators().RSI(btc, 14)
    price, _ := s.k.Market().Price(btc)
    
    if rsi.LessThan(decimal.NewFromInt(30)) {
        signal := s.k.Signal(s.GetName()).
            Buy(btc, connector.Binance, decimal.NewFromFloat(0.1)).
            Build()
        return []*strategy.Signal{signal}, nil
    }
    
    return nil, nil
}
```

That's it. No exchange client setup. No data fetching. No indicator calculations from scratch.

## What Kronos Does for You

### Automatic Data Management

```go
// You write:
rsi, _ := s.k.Indicators().RSI(btc, 14)

// Kronos handles:
// - Fetching price data from the exchange
// - Caching and managing historical data
// - Calculating the RSI indicator
// - Returning the current value
```

### Exchange Abstraction

```go
// Same code works on any exchange
price, _ := s.k.Market().Price(btc)  // Automatically uses your default exchange

// Or specify explicitly
price, _ := s.k.Market().Price(btc, indicators.IndicatorOptions{
    Exchange: connector.Binance,
})
```

### Built-in Indicators

All technical indicators work the same way:

```go
rsi, _ := s.k.Indicators().RSI(btc, 14)
sma, _ := s.k.Indicators().SMA(btc, 20)
ema, _ := s.k.Indicators().EMA(btc, 50)
macd, _ := s.k.Indicators().MACD(btc, 12, 26, 9)
bb, _ := s.k.Indicators().BollingerBands(btc, 20, 2.0)
stoch, _ := s.k.Indicators().Stochastic(btc, 14, 3)
atr, _ := s.k.Indicators().ATR(btc, 14)
```

Clean, consistent API. Pass the asset and parameters. Kronos does the rest.

### Market Data

```go
// Current price
price, _ := s.k.Market().Price(btc)

// Funding rates (for perpetuals)
funding, _ := s.k.Market().FundingRate(btc)

// Order book
book, _ := s.k.Market().OrderBook(btc)

// Prices across all exchanges
prices := s.k.Market().Prices(btc)
```

### Cross-Exchange Operations

```go
// Find arbitrage opportunities
arbOpps := s.k.Market().FindArbitrage(btc, decimal.NewFromInt(10)) // 10 bps minimum

for _, opp := range arbOpps {
    fmt.Printf("Buy on %s @ %s, Sell on %s @ %s\n",
        opp.BuyExchange, opp.BuyPrice,
        opp.SellExchange, opp.SellPrice)
}
```

## A Complete Strategy

Here's a full working strategy:

```go
package main

import (
    sdk "github.com/backtesting-org/kronos-sdk/pkg/kronos"
    "github.com/backtesting-org/kronos-sdk/pkg/types/strategy"
    "github.com/shopspring/decimal"
)

type MomentumStrategy struct {
    k *sdk.Kronos
}

func NewMomentum(k *sdk.Kronos) *MomentumStrategy {
    return &MomentumStrategy{k: k}
}

func (s *MomentumStrategy) GetSignals() ([]*strategy.Signal, error) {
    btc := s.k.Asset("BTC")
    eth := s.k.Asset("ETH")
    
    // Get indicators (Kronos fetches data automatically)
    btcRSI, _ := s.k.Indicators().RSI(btc, 14)
    ethRSI, _ := s.k.Indicators().RSI(eth, 14)
    
    var signals []*strategy.Signal
    
    // BTC oversold
    if btcRSI.LessThan(decimal.NewFromInt(30)) {
        signals = append(signals, 
            s.k.Signal(s.GetName()).
                Buy(btc, connector.Binance, decimal.NewFromFloat(0.1)).
                Build())
        s.k.Log().Opportunity("Momentum", "BTC", "RSI oversold")
    }
    
    // ETH overbought
    if ethRSI.GreaterThan(decimal.NewFromInt(70)) {
        signals = append(signals, 
            s.k.Signal(s.GetName()).
                Sell(eth, connector.Binance, decimal.NewFromFloat(1.0)).
                Build())
        s.k.Log().Opportunity("Momentum", "ETH", "RSI overbought")
    }
    
    return signals, nil
}

// Interface methods
func (s *MomentumStrategy) GetName() strategy.StrategyName {
    return "Momentum"
}

func (s *MomentumStrategy) GetDescription() string {
    return "RSI-based momentum strategy"
}

func (s *MomentumStrategy) GetRiskLevel() strategy.RiskLevel {
    return strategy.RiskLevelMedium
}

func (s *MomentumStrategy) GetStrategyType() strategy.StrategyType {
    return strategy.StrategyTypeTechnical
}
```

## Key Features

### Type Safety

Full compile-time checking with Go's type system:

```go
// This won't compile - caught before runtime
rsi, _ := s.k.Indicators().RSI(btc, "14")  // ❌ string instead of int

// This is correct
rsi, _ := s.k.Indicators().RSI(btc, 14)    // ✅
```

### Decimal Precision

No floating-point errors:

```go
import "github.com/shopspring/decimal"

// Financial-grade precision
price := decimal.NewFromFloat(50000.50)
quantity := decimal.NewFromFloat(0.1)
total := price.Mul(quantity)  // Exact calculation
```

### Write Once, Run Anywhere

The same strategy code works in:
- **Backtesting** - Test against historical data
- **Paper Trading** - Practice with simulated funds
- **Live Trading** - Deploy to real exchanges

No code changes. No environment checks. Just works.

### Exchange Support

Currently supported:
- Binance
- Bybit
- Hyperliquid

Your strategy code doesn't care which exchange you use. Kronos abstracts it all.

## How It Works

When you call an indicator or market data function, Kronos:

1. **Identifies the data source** - Uses your configured exchange or specified exchange
2. **Fetches required data** - Gets price history, order book, funding rates, etc.
3. **Caches intelligently** - Stores data to avoid redundant API calls
4. **Calculates/returns** - Computes indicators or returns market data
5. **Handles errors** - Manages rate limits, retries, and error cases

You just write `s.k.Indicators().RSI(btc, 14)` and Kronos does the rest.

## Multiple Timeframes

Use different intervals for different purposes:

```go
// Long-term trend (4-hour chart)
sma200, _ := s.k.Indicators().SMA(btc, 200, indicators.IndicatorOptions{
    Interval: "4h",
})

// Short-term signal (1-hour chart)
rsi, _ := s.k.Indicators().RSI(btc, 14, indicators.IndicatorOptions{
    Interval: "1h",
})

price, _ := s.k.Market().Price(btc)

// Only buy if:
// - Price is above 4h SMA200 (uptrend)
// - 1h RSI is oversold
if price.GreaterThan(sma200) && rsi.LessThan(decimal.NewFromInt(30)) {
    signal := s.k.Signal(s.GetName()).
        Buy(btc, connector.Binance, decimal.NewFromFloat(0.1)).
        Build()
    return []*strategy.Signal{signal}, nil
}
```

## Multiple Assets

Trade multiple assets in one strategy:

```go
func (s *Strategy) GetSignals() ([]*strategy.Signal, error) {
    btc := s.k.Asset("BTC")
    eth := s.k.Asset("ETH")
    sol := s.k.Asset("SOL")
    
    var signals []*strategy.Signal
    
    // Check each asset
    for _, asset := range []portfolio.Asset{btc, eth, sol} {
        rsi := s.k.Indicators().RSI(asset, 14)
        
        if rsi.LessThan(decimal.NewFromInt(30)) {
            signals = append(signals, 
                s.k.Signal(s.GetName()).
                    Buy(asset, connector.Binance, decimal.NewFromFloat(0.1)).
                    Build())
        }
    }
    
    return signals, nil
}
```

## Analytics

Beyond indicators, Kronos provides market analytics:

```go
// Volatility analysis
vol, _ := s.k.Analytics().Volatility(btc, 24)

// Trend analysis
trend, _ := s.k.Analytics().Trend(btc, 50)
s.k.Log().Info("Trend: %s (Strength: %s%%)", 
    trend.Direction, trend.Strength)

// Volume analysis
volAnalysis, _ := s.k.Analytics().VolumeAnalysis(btc, 24)
if volAnalysis.IsVolumeSpike {
    s.k.Log().MarketCondition("Volume spike detected")
}

// Price change
change, _ := s.k.Analytics().GetPriceChange(btc, 24)
s.k.Log().Info("24h change: %s%%", change.ChangePercent)
```

## Architecture

Kronos consists of three components:

```
┌─────────────────────────────────────────────────────────┐
│  kronos-cli                                             │
│  • Single CLI for all operations                        │
│  • Manages backtesting and live runtimes               │
└──────────────────┬──────────────────────────────────────┘
                   │
         ┌─────────┴─────────┐
         │                   │
         ▼                   ▼
┌─────────────────┐  ┌─────────────────┐
│ kronos-backtest │  │  kronos-live    │
│ • Simulated     │  │ • Real exchange │
│   exchange      │  │   connectors    │
│ • Historical    │  │ • Live data     │
│   data replay   │  │ • Real orders   │
└────────┬────────┘  └────────┬────────┘
         │                    │
         └─────────┬──────────┘
                   │
                   ▼
         ┌─────────────────────┐
         │   kronos-sdk        │
         │   (this package)    │
         │                     │
         │ • Strategy API      │
         │ • Indicators        │
         │ • Market services   │
         └─────────────────────┘
```

As a strategy developer, you only interact with `kronos-sdk`. The CLI and runtimes are handled automatically.

## Getting Started

Ready to build your first strategy? Head to [Getting Started](getting-started).

## API Reference

Explore the complete API:

- **[Indicators](api/indicators/stochastic)** - Technical analysis indicators
- **Market Data** - Prices, order books, funding rates
- **Analytics** - Volatility, trend, volume analysis
- **Strategy Framework** - Building and testing strategies

## Philosophy

Kronos is built on these principles:

1. **Simplicity** - Hide complexity, expose clean APIs
2. **Type Safety** - Catch errors at compile time
3. **Precision** - Use decimal arithmetic for financial calculations
4. **Write Once** - Same code works everywhere
5. **Focus** - You write strategy logic, we handle everything else

## Community

- **GitHub**: [backtesting-org/kronos-sdk](https://github.com/backtesting-org/kronos-sdk)
- **Issues**: [Report bugs or request features](https://github.com/backtesting-org/kronos-sdk/issues)
- **Go Docs**: [pkg.go.dev](https://pkg.go.dev/github.com/backtesting-org/kronos-sdk)

## License

Kronos SDK is open source under the MIT License.
