# Wisp SDK

The SDK is the Go library that powers the Wisp trading runtime. It defines the interfaces strategies are written against, the market domain packages (spot, perp, prediction), the signal execution pipeline, the monitoring server, and the fx dependency-injection wiring that holds everything together.

## For Strategy Developers

If you're writing trading strategies with Wisp SDK, start at **[usewisp.dev/docs](https://usewisp.dev/docs)**:

- **[Getting Started](https://usewisp.dev/docs/getting-started)** — Installation, configuration, your first strategy
- **[Writing Strategies](https://usewisp.dev/docs/getting-started/writing-strategies)** — 13 strategy patterns with complete examples
- **[Strategy Examples](https://usewisp.dev/docs/examples)** — Real strategies from basic to advanced
- **[API Reference](https://usewisp.dev/docs/api/indicators/rsi)** — Indicators, market data, and more

---

## For Contributors

This document explains the SDK codebase structure and is aimed at developers contributing to the Wisp SDK itself.

---

## Repository Layout

```
pkg/
├── types/              # Public interfaces — the contract surface of the SDK
│   ├── wisp/           # wisp.Wisp — the handle injected into every strategy
│   │   ├── analytics/  # Indicators and Analytics interfaces
│   │   ├── numerical/  # Decimal wrapper (shopspring/decimal)
│   │   └── activity/   # Activity / PNL read interfaces
│   ├── strategy/       # Strategy, Signal, SignalBuilder, BaseStrategy
│   ├── connector/      # Exchange-facing types: Order, Trade, Kline, OrderBook
│   ├── execution/      # Executor interface
│   ├── lifecycle/      # Orchestrator and DomainLifecycle interfaces
│   ├── registry/       # StrategyRegistry and Hooks interfaces
│   ├── monitoring/     # Monitoring server and view interfaces
│   └── plugin/         # Plugin manager interface
│
├── markets/            # Domain implementations (spot, perp, prediction)
│   ├── base/           # Shared ingestor and store primitives
│   ├── spot/           # Spot market: watchlist, store, ingestor, executor, views
│   ├── perp/           # Perp market: watchlist, store, ingestor, executor, views
│   └── prediction/     # Prediction market: watchlist, store, ingestor, executor
│
├── analytics/          # Indicator and analytics implementations
│   └── indicators/     # RSI, MACD, EMA, SMA, ATR, Bollinger, Stochastic
│
├── signal/             # Concrete SpotSignalBuilder / PerpSignalBuilder
├── executor/           # Top-level signal router (dispatches to domain executors)
├── registry/           # StrategyRegistry and Hooks implementations
├── lifecycle/          # Orchestrator and monitoring lifecycle
├── plugin/             # Go plugin loader (.so → strategy.Strategy)
├── monitoring/         # Unix-socket HTTP monitoring server
├── activity/           # Cross-domain activity aggregation
├── config/             # Config loading (viper-backed)
├── runtime/            # TimeProvider and runtime utilities
├── adapters/           # Logging adapters (zap)
└── modules.go          # Root fx.Module — wires all packages together
```

---

## Core Concepts

### `wisp.Wisp` — the strategy handle

`pkg/types/wisp/wisp.go` defines the `Wisp` interface. This is the only dependency injected into a strategy. It exposes:

- `Spot() / Perp() / Predict()` — domain-scoped objects for market data and signal creation
- `Indicators()` — pure-function technical indicator calculations
- `Analytics()` — higher-level analytics (trend, volatility, volume)
- `Activity()` — read access to positions, trades, and PNL
- `Emit(signal)` — routes a built signal directly to the executor
- `Asset(symbol) / Pair(base, quote)` — portfolio type constructors
- `Log()` — strategy-scoped logger

### `strategy.Strategy` — the strategy interface

`pkg/types/strategy/strategy.go` defines what a strategy must implement:

```go
type Strategy interface {
    GetName() StrategyName
    Start(ctx context.Context) error
    Stop(ctx context.Context) error
    Signals() <-chan Signal
    LatestStatus() StrategyStatus
    StatusLog() []StrategyStatus
}
```

Strategies are **self-directed** — they own their run loop. The orchestrator only calls `Start`/`Stop`. Signals are pushed asynchronously: `wisp.Emit(signal)` routes to the executor; `base.Emit(signal)` publishes to the observable `Signals()` channel.

### `strategy.BaseStrategy` — embed this

`pkg/types/strategy/base.go` provides `BaseStrategy`, which should be embedded in every concrete strategy. It implements `Stop`, `Signals`, `LatestStatus`, `StatusLog`, and the `Emit`/`Mark`/`EmitStatus` helpers. Concrete strategies call `StartWithRunner(ctx, s.run)` from their own `Start` method to launch the run goroutine under the base lifecycle.

```go
type momentumStrategy struct {
    strategy.BaseStrategy
    k wisp.Wisp
}

func NewStrategy(k wisp.Wisp) strategy.Strategy {
    s := &momentumStrategy{k: k}
    s.BaseStrategy = *strategy.NewBaseStrategy(strategy.BaseStrategyConfig{Name: "momentum"})
    return s
}

func (s *momentumStrategy) Start(ctx context.Context) error {
    return s.StartWithRunner(ctx, s.run)
}

func (s *momentumStrategy) run(ctx context.Context) {
    // your tick loop here
}
```

See [Strategy Examples](https://usewisp.dev/docs/examples) for complete working strategies using this pattern, from RSI momentum to multi-exchange arbitrage.

### Signal builders

Signal builders live on the domain objects, not on `wisp.Wisp` directly:

```go
// Perp market order
signal, err := s.k.Perp().Signal(s.GetName()).
    Buy(pair, connector.Hyperliquid, numerical.NewFromFloat(0.1)).
    Build()

// Spot limit order
signal, err := s.k.Spot().Signal(s.GetName()).
    BuyLimit(pair, connector.Hyperliquid, qty, price).
    Build()
```

`Build()` returns `(Signal, error)`. The concrete builders are in `pkg/signal/builder.go`; the interfaces are in `pkg/types/strategy/builder.go`.

For a complete walkthrough of signal patterns and strategy structure, see [Writing Strategies](https://usewisp.dev/docs/getting-started/writing-strategies).

### Indicators

Indicators are pure functions — they take a `[]connector.Kline` slice, not an asset or exchange name. Fetch klines from the domain object first, then pass them in:

```go
klines := s.k.Perp().Klines(connector.Hyperliquid, pair, "1h", 60)

rsi, _   := s.k.Indicators().RSI(klines, 14)
macd, _  := s.k.Indicators().MACD(klines, 12, 26, 9)
bb, _    := s.k.Indicators().BollingerBands(klines, 20, 2.0)
ema, _   := s.k.Indicators().EMA(klines, 50)
sma, _   := s.k.Indicators().SMA(klines, 200)
atr, _   := s.k.Indicators().ATR(klines, 14)
stoch, _ := s.k.Indicators().Stochastic(klines, 14, 3)
```

Fetching klines once and reusing them across multiple indicator calls avoids redundant store reads. Implementations live in `pkg/analytics/indicators/`.

For strategy-focused indicator usage and multiple indicator combinations, see [API Reference](https://usewisp.dev/docs/api/indicators/rsi) and the pattern guide on [Writing Strategies](https://usewisp.dev/docs/getting-started/writing-strategies).

---

## Market Domains

Each market type (spot, perp, prediction) follows the same internal structure:

```
markets/{domain}/
├── module.go       # fx wiring for this domain
├── watchlist.go    # tracks which pairs the strategy has subscribed to
├── universe.go     # available pairs / markets
├── store/          # in-memory store for klines, orderbooks, positions
├── ingestor/       # pulls data from connectors into the store
│   ├── batch/      # initial snapshot ingestor
│   └── realtime/   # live WebSocket ingestor
├── executor/       # converts signals into exchange orders
├── activity/       # PNL calculation for this domain
├── views/          # read-only projections exposed to strategies
└── types/          # domain-specific interfaces (Spot, Perp, Predict SDKs)
```

The domain's public interface (`types/spot_sdk.go`, `types/perp_sdk.go`) is what strategies interact with. Everything underneath is internal to the domain package.

---

## Signal Flow

```
strategy.run() → wisp.Emit(signal)
    → executor.Execute(signal)         pkg/executor/default.go
        → spotExecutor / perpExecutor  pkg/markets/{domain}/executor/
            → connector.PlaceOrder     (external connectors package)
```

The top-level executor (`pkg/executor/`) type-switches on the signal (`SpotSignal` vs `PerpSignal`) and delegates to the appropriate domain executor. Domain executors translate `strategy.Signal` → `connector.Order` and call out to the exchange connector.

---

## Dependency Injection

The SDK uses [uber-go/fx](https://github.com/uber-go/fx). Every package exposes a `var Module = fx.Module(...)` that declares its `fx.Provide` entries. The root module in `pkg/modules.go` composes them all:

```go
var Module = fx.Options(
    activity.Module,
    adapters.Module,
    analytics.Module,
    config.Module,
    monitoring.Module,
    lifecycle.Module,
    plugin.Module,
    registry.Module,
    runtime.Module,
    signal.Module,
    executor.Module,
    prediction.Module,
    perp.Module,
    spot.Module,
)
```

When adding a new package, create a `module.go` with `var Module = fx.Module(...)` and add it here.

---

## Monitoring

The monitoring server (`pkg/monitoring/`) runs a Unix-socket HTTP server inside the strategy process. The TUI connects to this socket to read live data. Views are registered via `ViewRegistry` and queried by the TUI over JSON. Handlers live in `handlers_*.go` files, grouped by concern: core, klines, markets, orderbook, profiling.

---

## Plugin Loading

Strategies are compiled as Go plugins (`.so` files) and loaded at runtime by `pkg/plugin/manager.go`. The plugin must export a `NewStrategy(wisp.Wisp) strategy.Strategy` symbol. The manager loads the `.so`, resolves the symbol, calls it with the injected `wisp.Wisp`, and registers the returned strategy with the `StrategyRegistry`.

---

## Testing

- Unit tests sit alongside the files they test (`_test.go`).
- Integration/spec tests use [Ginkgo](https://github.com/onsi/ginkgo) (`spec_test.go`, `suite_test.go`).
- Mocks are generated with [mockery](https://github.com/vektra/mockery) and live under `mocks/`.

Run all tests:

```bash
make test
```

---

## Resources

- **[Wisp SDK Documentation](https://usewisp.dev/docs)** — User guide, examples, and API reference
- **[GitHub Repository](https://github.com/wisp-trading/sdk)** — Source code and issue tracking
- **[Getting Started](https://usewisp.dev/docs/getting-started)** — Installation and first steps

---

## Adding a New Indicator

1. Add an implementation file in `pkg/analytics/indicators/` (e.g. `mfi.go`).
2. Add the method signature to `pkg/types/wisp/analytics/indicators.go`.
3. Wire it into the concrete `analyticsIndicators` struct in `pkg/analytics/`.
4. Add tests in `pkg/analytics/indicators/mfi_test.go`.
