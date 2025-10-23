# Simple SMA Crossover Strategy

A complete example of a trading strategy using the Kronos SDK.

## What This Strategy Does

This is a **Simple Moving Average (SMA) Crossover** strategy that:
- Calculates two SMAs: a short-period and a long-period
- Generates a **BUY signal** when the short SMA crosses **above** the long SMA (golden cross)
- Generates a **SELL signal** when the short SMA crosses **below** the long SMA (death cross)

## Configuration

```yaml
strategies:
  - name: sma-crossover
    enabled: true
    config:
      short_period: 10      # Fast SMA period
      long_period: 50       # Slow SMA period
      asset: BTC           # Asset to trade
      exchange: Paradex    # Exchange to use
      interval: 5m         # Kline interval
      quantity: "1.0"      # Amount to trade
```

## Key Features Demonstrated

1. **Strategy Registration** - Uses `init()` to register with SDK
2. **BaseStrategy** - Embeds BaseStrategy for common functionality
3. **Context Usage** - Accesses market data and portfolio through Context
4. **Signal Generation** - Creates typed signals with metadata
5. **Validation** - Implements Validator interface for config validation
6. **Lifecycle Hooks** - Implements OnEnable/OnDisable for initialization
7. **Testing** - Comprehensive tests using MockContext

## Code Structure

```
simple_sma/
├── strategy.go       # Main strategy implementation
├── strategy_test.go  # Comprehensive tests
└── README.md         # This file
```

## How It Works

### 1. Registration

```go
func init() {
    strategy.Register("sma-crossover", NewFromConfig)
}
```

The strategy registers itself when imported. Kronos can then load it by name.

### 2. Implementation

```go
type SMAStrategy struct {
    *strategy.BaseStrategy  // Embed for common functionality
    config Config           // Strategy-specific config
}

func (s *SMAStrategy) GetSignals() ([]*types.Signal, error) {
    // 1. Check if enabled
    if !s.IsEnabled() {
        return nil, nil
    }
    
    // 2. Get market data
    klines := s.Context().Market().GetKlines(...)
    
    // 3. Calculate indicators
    shortSMA := s.calculateSMA(klines, s.config.ShortPeriod)
    longSMA := s.calculateSMA(klines, s.config.LongPeriod)
    
    // 4. Detect crossover and generate signal
    if goldenCross {
        return []*types.Signal{buySignal}, nil
    }
    
    return nil, nil
}
```

### 3. Testing

```go
func TestSMAStrategy_GoldenCross(t *testing.T) {
    // Create mock context
    ctx := testing.NewMockContext()
    
    // Setup test data
    ctx.Market().SetKlines("BTC", "Paradex", "5m", mockKlines)
    
    // Create strategy
    strat := New(ctx, config)
    strat.Enable()
    
    // Test
    signals, _ := strat.GetSignals()
    assert.Equal(t, types.ActionBuy, signals[0].Actions[0].Action)
}
```

## Running Tests

```bash
cd sdk/examples/simple_sma
go test -v
```

## Using This Strategy

### Option 1: Direct Import (for learning)

```go
import "chatapi/sdk/examples/simple_sma"

strat := simple_sma.New(ctx, simple_sma.Config{
    ShortPeriod: 10,
    LongPeriod: 50,
    Asset: "BTC",
    Exchange: "Paradex",
    Interval: "5m",
    Quantity: "1.0",
})
```

### Option 2: Via Registry (production)

```go
import (
    "chatapi/sdk/strategy"
    _ "chatapi/sdk/examples/simple_sma"  // Side-effect import
)

strat, err := strategy.Create("sma-crossover", ctx, config)
```

### Option 3: In Kronos (via config)

```yaml
# config/strategies/default.yml
strategies:
  - name: sma-crossover
    enabled: true
    config:
      short_period: 10
      long_period: 50
      asset: BTC
      exchange: Paradex
      interval: 5m
      quantity: "1.0"
```

## Extending This Example

To create your own strategy based on this:

1. Copy this directory
2. Change the package name
3. Update the registration name in `init()`
4. Modify the `GetSignals()` logic
5. Update tests
6. Register in your Kronos config

## What You Learn Here

- ✅ How to implement the Strategy interface
- ✅ How to use BaseStrategy for common functionality
- ✅ How to access market data through Context
- ✅ How to check portfolio state before trading
- ✅ How to generate typed signals with metadata
- ✅ How to validate configuration
- ✅ How to implement lifecycle hooks
- ✅ How to write comprehensive tests with MockContext
- ✅ How to register strategies for config-driven loading

## Performance Notes

This is a **teaching example** optimized for clarity, not performance. For production:
- Cache SMA calculations
- Use more efficient data structures
- Add risk management
- Add position sizing logic
- Add stop-loss/take-profit levels

## Related Examples

- `rsi_momentum/` - RSI-based momentum strategy
- `bollinger_bands/` - Mean reversion strategy
- `grid_trading/` - Grid trading strategy
- `multi_asset/` - Multi-asset portfolio strategy

## Questions?

See the main [SDK documentation](../../README.md) or [PLAN.md](../../PLAN.md).

