# Kronos Strategy SDK

A type-safe SDK for building trading strategies on the Kronos trading platform. This SDK provides comprehensive interfaces for market data access, exchange connectivity, portfolio management, and strategy development.

## 🚀 Quick Start

### Installation

```bash
go get kronos/sdk
```

### Your First Strategy

```go
package mystrategy

import (
    "kronos/sdk/pkg/types/connector"
    "kronos/sdk/pkg/types/logging"
    "kronos/sdk/pkg/types/portfolio"
    "kronos/sdk/pkg/types/portfolio/store"
    "kronos/sdk/pkg/types/strategy"
    "time"

    "github.com/google/uuid"
    "github.com/shopspring/decimal"
)

type MyStrategy struct {
    *strategy.BaseStrategy
    assetStore store.Store
    logger     logging.ApplicationLogger
}

func NewMyStrategy(
    assetStore store.Store,
    logger logging.ApplicationLogger,
) *MyStrategy {
    base := strategy.NewBaseStrategy(
        strategy.StrategyName("My Strategy"),
        "A simple example strategy",
        strategy.RiskLevelLow,
        strategy.StrategyTypeTechnical,
    )
    
    return &MyStrategy{
        BaseStrategy: base,
        assetStore:   assetStore,
        logger:       logger,
    }
}

func (s *MyStrategy) GetSignals() ([]*strategy.Signal, error) {
    if !s.IsEnabled() {
        return nil, nil
    }
    
    // Get market data from the asset store
    btc := portfolio.NewAsset("BTC")
    price := s.assetStore.GetAssetPrice(btc, connector.Paradex)
    
    if price == nil {
        return nil, nil // No data yet
    }
    
    // Simple buy signal example
    signal := &strategy.Signal{
        ID:       uuid.New(),
        Strategy: s.GetName(),
        Actions: []strategy.TradeAction{
            {
                Action:   strategy.ActionBuy,
                Asset:    btc,
                Exchange: connector.Paradex,
                Quantity: decimal.NewFromFloat(0.01),
                Price:    price.Price,
            },
        },
        Timestamp: time.Now(),
    }
    
    return []*strategy.Signal{signal}, nil
}
```

## 📚 Core Concepts

### Strategy Interface

All strategies must implement the `strategy.Strategy` interface:

```go
type Strategy interface {
    GetSignals() ([]*Signal, error)     // Generate trading signals
    
    GetName() StrategyName              // Strategy name
    GetDescription() string              // Strategy description
    GetRiskLevel() RiskLevel            // Risk classification
    GetStrategyType() StrategyType      // Strategy type
    
    Enable() error                       // Enable the strategy
    Disable() error                      // Disable the strategy
    IsEnabled() bool                     // Check if enabled
}
```

### BaseStrategy

Use `BaseStrategy` for common functionality:

```go
base := strategy.NewBaseStrategy(
    strategy.CashCarry,
    "Cash and carry arbitrage strategy",
    strategy.RiskLevelLow,
    strategy.StrategyTypeCashCarry,
)

// Access base methods
base.Enable()
base.Disable()
base.IsEnabled()
base.GetName()
base.GetRiskLevel()
```

### Asset Store

The `store.Store` interface provides access to all market data:

```go
// Funding rates
fundingRates := assetStore.GetFundingRatesForAsset(asset)
fundingRate := assetStore.GetFundingRate(asset, exchange)
assetsWithFunding := assetStore.GetAllAssetsWithFundingRates()

// Order books
orderBook := assetStore.GetOrderBook(asset, exchange, connector.TypePerpetual)
orderBooks := assetStore.GetOrderBooks(asset)
assetsWithOrderBooks := assetStore.GetAllAssetsWithOrderBooks()

// Prices
price := assetStore.GetAssetPrice(asset, exchange)
prices := assetStore.GetAssetPrices(asset)
```

## 🔌 Connector Interface

The `Connector` interface defines exchange integration:

### Market Data
```go
// Price data
price, err := connector.FetchPrice("BTC")
klines, err := connector.FetchKlines("BTC", types.Interval5Minute, 100)
orderBook, err := connector.FetchOrderBook(asset, connector.TypePerpetual, 10)

// Funding rates
fundingRate, err := connector.FetchFundingRate(asset)
fundingRates, err := connector.FetchCurrentFundingRates()
historicalRates, err := connector.FetchHistoricalFundingRates(asset, startTime, endTime)
```

### Trading Operations
```go
// Place orders
response, err := connector.PlaceLimitOrder("BTC", connector.OrderSideBuy, quantity, price)
response, err := connector.PlaceMarketOrder("BTC", connector.OrderSideBuy, quantity)

// Manage orders
cancelResp, err := connector.CancelOrder("BTC", orderID)
orders, err := connector.GetOpenOrders()
order, err := connector.GetOrderStatus(orderID)
```

### Account Management
```go
// Account data
balance, err := connector.GetAccountBalance()
positions, err := connector.GetPositions()
trades, err := connector.GetTradingHistory("BTC", 100)

// Available assets
spotAssets, err := connector.FetchAvailableSpotAssets()
perpAssets, err := connector.FetchAvailablePerpetualAssets()
```

### WebSocket Support

For real-time data, use `WebSocketConnector`:

```go
// Lifecycle
err := wsConnector.StartWebSocket(ctx)
defer wsConnector.StopWebSocket()

// Subscriptions
wsConnector.SubscribeOrderBook(asset, connector.TypePerpetual)
wsConnector.SubscribeFundingRate(asset)
wsConnector.SubscribeKlines(asset, types.Interval5Minute)

// Data channels
orderBooks := wsConnector.OrderBookUpdates()
klines := wsConnector.KlineUpdates()
errors := wsConnector.ErrorChannel()
```

## 📊 Data Types

### Market Data Types

```go
// Price
type Price struct {
    Symbol    string
    Price     decimal.Decimal
    BidPrice  decimal.Decimal
    AskPrice  decimal.Decimal
    Volume24h decimal.Decimal
    Timestamp time.Time
}

// Kline (candlestick)
type Kline struct {
    Symbol    string
    Interval  string
    OpenTime  time.Time
    Open      decimal.Decimal
    High      decimal.Decimal
    Low       decimal.Decimal
    Close     decimal.Decimal
    Volume    decimal.Decimal
}

// OrderBook
type OrderBook struct {
    Asset     portfolio.Asset
    Bids      []PriceLevel
    Asks      []PriceLevel
    Timestamp time.Time
}

// Funding Rate
type FundingRate struct {
    CurrentRate     decimal.Decimal
    NextFundingTime time.Time
    MarkPrice       decimal.Decimal
    IndexPrice      decimal.Decimal
}
```

### Trading Types

```go
// Trade Action
type TradeAction struct {
    Action   Action                 // buy, sell, sell_short, cover
    Asset    portfolio.Asset
    Exchange connector.ExchangeName
    Quantity decimal.Decimal
    Price    decimal.Decimal
}

// Signal
type Signal struct {
    ID        uuid.UUID
    Strategy  StrategyName
    Actions   []TradeAction
    Timestamp time.Time
}

// Order
type Order struct {
    ID           string
    Symbol       string
    Side         OrderSide       // BUY, SELL
    Type         OrderType       // LIMIT, MARKET
    Status       OrderStatus     // NEW, FILLED, CANCELED
    Quantity     decimal.Decimal
    Price        decimal.Decimal
    FilledQty    decimal.Decimal
}

// Position
type Position struct {
    Symbol        portfolio.Asset
    Exchange      ExchangeName
    Side          OrderSide
    Size          decimal.Decimal
    EntryPrice    decimal.Decimal
    UnrealizedPnL decimal.Decimal
}
```

## 🎯 Strategy Types & Actions

### Strategy Types
- `StrategyTypeVolumeMaximizer` - Volume maximization
- `StrategyTypeCashCarry` - Cash and carry arbitrage
- `StrategyTypeArbitrage` - General arbitrage
- `StrategyTypeTechnical` - Technical analysis
- `StrategyTypeMomentum` - Momentum trading
- `StrategyTypeMeanReversion` - Mean reversion

### Actions
- `ActionBuy` - Buy (long position)
- `ActionSell` - Sell (close long)
- `ActionSellShort` - Short sell (open short)
- `ActionCover` - Cover (close short)
- `ActionHold` - Hold position
- `ActionClose` - Close position

### Risk Levels
- `RiskLevelLow` - Conservative strategies
- `RiskLevelMedium` - Moderate risk
- `RiskLevelHigh` - Aggressive strategies

## 🏢 Supported Exchanges

- `Hyperliquid` - Hyperliquid DEX
- `Paradex` - Paradex exchange
- `Binance` - Binance
- `Bybit` - Bybit

## 📝 Logging

Two logging interfaces for different purposes:

### ApplicationLogger
For system errors and debugging

```go
logger.Info("Starting strategy")
logger.Debug("Current state: %v", state)
logger.Warn("Unusual condition detected")
logger.Error("Failed to fetch data: %v", err)
logger.ErrorWithDebug("API error", rawResponse)
```

### TradingLogger
For business events

```go
tradingLogger.Opportunity("CashCarry", "BTC", "Funding rate: %v", rate)
tradingLogger.Success("CashCarry", "BTC", "Position opened")
tradingLogger.Failed("CashCarry", "BTC", "Order rejected")
tradingLogger.MarketCondition("High volatility detected")
tradingLogger.OrderLifecycle("Order filled", "BTC")
```

## 🔧 Strategy Registry

Register and manage strategies:

```go
registry := registry.NewStrategyRegistry()

// Register strategies
registry.Register(cashCarryStrategy)
registry.Register(momentumStrategy)

// Get strategies
strategy, exists := registry.GetStrategy(strategy.CashCarry)
allStrategies := registry.GetAllStrategies()
enabledOnly := registry.GetEnabledStrategies()

// Manage state
registry.EnableStrategy(strategy.CashCarry)
registry.DisableStrategy(strategy.CashCarry)
```

## ⏰ Temporal Interface

For time operations that work in both live and simulation modes:

```go
type TimeProvider interface {
    Now() time.Time
    After(d time.Duration) <-chan time.Time
    NewTimer(d time.Duration) Timer
    Since(t time.Time) time.Duration
    NewTicker(d time.Duration) Ticker
    Sleep(d time.Duration)
}
```

## 📖 Example: Cash & Carry Strategy

See `examples/cash_carry/` for a complete working example demonstrating:
- Funding rate monitoring
- Multi-exchange signal generation
- Asset store integration
- Strategy interface implementation

```go
type CashCarryStrategy struct {
    assetStore     store.Store
    logger         logging.ApplicationLogger
    minFundingRate decimal.Decimal
}

func (s *CashCarryStrategy) GetSignals() ([]*strategy.Signal, error) {
    assets := s.assetStore.GetAllAssetsWithFundingRates()
    
    var signals []*strategy.Signal
    for _, asset := range assets {
        fundingRates := s.assetStore.GetFundingRatesForAsset(asset)
        
        for exchange, rate := range fundingRates {
            if rate.CurrentRate.GreaterThan(s.minFundingRate) {
                signal := s.createSignal(asset, exchange, rate.CurrentRate)
                signals = append(signals, signal)
            }
        }
    }
    
    return signals, nil
}
```

## 🏗️ Project Structure

```
kronos/sdk/
├── pkg/types/              # Public type definitions
│   ├── connector/          # Exchange connector interfaces
│   │   ├── connector.go    # Main Connector interface
│   │   ├── market_data.go  # Market data types
│   │   ├── trading.go      # Trading types
│   │   ├── account.go      # Account types
│   │   └── funding.go      # Funding rate types
│   ├── strategy/           # Strategy interfaces
│   │   ├── strategy.go     # Strategy interface
│   │   ├── base.go         # BaseStrategy implementation
│   │   ├── signal.go       # Signal types
│   │   └── action.go       # Trade action types
│   ├── portfolio/          # Portfolio types
│   │   ├── asset.go        # Asset type
│   │   └── store/          # Asset store interface
│   ├── logging/            # Logger interfaces
│   └── temporal/           # Time abstraction
├── internal/               # Internal implementations
│   ├── registry/           # Strategy registry
│   ├── logging/            # Logger implementations
│   └── time/               # Time provider implementations
└── examples/               # Example strategies
    └── cash_carry/         # Cash & carry example
```

## 💡 Best Practices

1. **Use decimal.Decimal**: Always use `decimal.Decimal` for prices and quantities
2. **Check enabled state**: Start `GetSignals()` with enabled check
3. **Handle nil data**: Always check for nil when fetching market data
4. **Proper error handling**: Return errors, don't panic
5. **Use appropriate logger**: ApplicationLogger for errors, TradingLogger for business events
6. **Type safety**: Use the provided asset and exchange types
7. **Resource cleanup**: Close channels and stop websockets properly

## 🔐 Type Safety

The SDK uses strong typing throughout:

```go
// Type-safe assets
asset := portfolio.NewAsset("BTC")

// Type-safe exchanges
exchange := connector.Paradex

// Type-safe actions
action := strategy.ActionBuy

// Type-safe order sides
side := connector.OrderSideBuy
```

## 📦 Dependencies

- `github.com/shopspring/decimal` - Precise decimal arithmetic
- `github.com/google/uuid` - UUID generation
- `go.uber.org/zap` - Structured logging

## 📄 License

See LICENSE file for details.


