# Executor Package

The executor package provides the execution layer for Wisp trading signals. It includes a default executor implementation with support for custom hooks, enabling users to extend execution behavior through plugins.

## Architecture

```
┌────────────────────────────────────┐
│  DefaultExecutor                   │
│  - Places orders                   │
│  - Tracks positions                │
│  - Manages hooks                   │
└────────────┬───────────────────────┘
             │
             │ Calls hooks at lifecycle points
             ▼
┌────────────────────────────────────┐
│  ExecutionHooks                    │
│  - BeforeExecute                   │
│  - AfterExecute                    │
│  - OnError                         │
└────────────────────────────────────┘
```

## Core Interfaces

### Executor

The main interface for executing trading signals:

```go
type Executor interface {
    ExecuteSignal(signal *strategy.Signal) error
    HandleTradeExecution(trade connector.Trade) error
    RegisterHook(hook ExecutionHook)
}
```

### ExecutionHook

Interface for customizing execution behavior:

```go
type ExecutionHook interface {
    BeforeExecute(ctx *ExecutionContext) error
    AfterExecute(ctx *ExecutionContext, result *ExecutionResult) error
    OnError(ctx *ExecutionContext, err error) error
}
```

## Usage

### Basic Usage (Default Executor)

```go
import "github.com/wisp-trading/wisp/pkg/executor"

// Create executor with dependencies
exec := executor.NewDefaultExecutor(
    connectorRegistry,
    positionStore,
    logger,
    timeProvider,
)

// Execute signals
signal := strategy.GenerateSignal()
err := exec.ExecuteSignal(signal)
```

### With Built-in Hooks

```go
import (
    "github.com/wisp-trading/wisp/pkg/executor"
    "github.com/wisp-trading/wisp/pkg/executor/hooks"
)

exec := executor.NewDefaultExecutor(...)

// Add risk management
exec.RegisterHook(hooks.NewBasicRiskHook(
    decimal.NewFromFloat(1000), // max position size
    50,                          // max daily trades
))

// Add logging
exec.RegisterHook(hooks.NewLoggingHook(logger))

// Add metrics
exec.RegisterHook(hooks.NewMetricsHook())
```

### With Custom Hooks (.so plugin)

```go
// Load user's custom hooks from plugin
hookPlugin, err := executor.LoadHookPlugin("./user-hooks.so")
if err == nil {
    for _, hook := range hookPlugin.CreateHooks() {
        exec.RegisterHook(hook)
    }
}
```

## Built-in Hooks

### BasicRiskHook

Provides basic risk management:
- Position size limits
- Daily trade limits

```go
hook := hooks.NewBasicRiskHook(
    decimal.NewFromFloat(1000), // max position size
    50,                          // max daily trades
)
```

### LoggingHook

Logs execution events:
- Before execution
- After execution
- On errors

```go
hook := hooks.NewLoggingHook(logger)
```

### MetricsHook

Tracks execution statistics:
- Total executions
- Success rate
- Order count

```go
hook := hooks.NewMetricsHook()

// Later, get stats
stats := hook.GetStats()
// "Executions: 100 | Success: 95 (95.0%) | Failures: 5 | Orders: 142"
```

## Creating Custom Hooks

Users can create custom hooks as `.so` plugins:

### 1. Create Hook Plugin

```go
// hooks.go
package main

import "github.com/wisp-trading/wisp/pkg/executor"

var HookPlugin hookPlugin

type hookPlugin struct{}

func (p hookPlugin) Name() string { return "my-hooks" }
func (p hookPlugin) Version() string { return "1.0.0" }

func (p hookPlugin) CreateHooks() []executor.ExecutionHook {
    return []executor.ExecutionHook{
        &MyRiskHook{},
        &MyNotificationHook{},
    }
}

// MyRiskHook implements custom risk logic
type MyRiskHook struct{}

func (h *MyRiskHook) BeforeExecute(ctx *executor.ExecutionContext) error {
    // Custom risk checks
    // - Call ML models
    // - Check market conditions
    // - Validate against portfolio
    return nil
}

func (h *MyRiskHook) AfterExecute(ctx *executor.ExecutionContext, result *executor.ExecutionResult) error {
    // Custom post-execution logic
    // - Update analytics
    // - Send notifications
    return nil
}

func (h *MyRiskHook) OnError(ctx *executor.ExecutionContext, err error) error {
    // Custom error handling
    return err
}
```

### 2. Build Plugin

```bash
go build -buildmode=plugin -o hooks.so
```

### 3. Load in Application

```go
hookPlugin, err := executor.LoadHookPlugin("./hooks.so")
if err == nil {
    for _, hook := range hookPlugin.CreateHooks() {
        exec.RegisterHook(hook)
    }
}
```

## Hook Execution Order

Hooks are executed in the order they are registered:

```go
exec.RegisterHook(hook1) // Executes first
exec.RegisterHook(hook2) // Executes second
exec.RegisterHook(hook3) // Executes third

// Execution flow:
// 1. hook1.BeforeExecute()
// 2. hook2.BeforeExecute()
// 3. hook3.BeforeExecute()
// 4. [Core Execution]
// 5. hook1.AfterExecute()
// 6. hook2.AfterExecute()
// 7. hook3.AfterExecute()
```

**Important:** If any `BeforeExecute` returns an error, execution is cancelled and all hooks' `OnError` methods are called.

## ExecutionContext

The context passed to hooks:

```go
type ExecutionContext struct {
    Signal    *strategy.Signal         // The signal being executed
    Timestamp time.Time                 // Execution timestamp
    Metadata  map[string]interface{}   // Custom metadata
}
```

Hooks can store data in `Metadata` for communication:

```go
func (h *Hook1) BeforeExecute(ctx *ExecutionContext) error {
    ctx.Metadata["risk_score"] = 0.75
    return nil
}

func (h *Hook2) BeforeExecute(ctx *ExecutionContext) error {
    riskScore := ctx.Metadata["risk_score"].(float64)
    if riskScore > 0.8 {
        return fmt.Errorf("risk too high")
    }
    return nil
}
```

## Best Practices

### 1. Keep Hooks Focused

Each hook should do one thing well:

```go
// Good
exec.RegisterHook(NewRiskHook())
exec.RegisterHook(NewNotificationHook())
exec.RegisterHook(NewMetricsHook())

// Avoid - one hook doing everything
exec.RegisterHook(NewMegaHook()) // risk + notifications + metrics
```

### 2. Use Metadata for Communication

```go
func (h *RiskHook) BeforeExecute(ctx *ExecutionContext) error {
    ctx.Metadata["risk_passed"] = true
    return nil
}

func (h *LoggingHook) AfterExecute(ctx *ExecutionContext, result *ExecutionResult) error {
    if ctx.Metadata["risk_passed"] == true {
        // Log that risk checks passed
    }
    return nil
}
```

### 3. Handle Errors Gracefully

```go
func (h *Hook) OnError(ctx *ExecutionContext, err error) error {
    // Log error, send alerts, etc.
    // But don't panic or modify the error unnecessarily
    return err
}
```

### 4. Make Hooks Stateless When Possible

```go
// Good - stateless
type ValidationHook struct {
    maxSize decimal.Decimal
}

// Be careful with state
type CounterHook struct {
    count int // This state persists across executions
}
```

## Testing Hooks

```go
func TestMyHook(t *testing.T) {
    hook := NewMyHook()
    
    ctx := &executor.ExecutionContext{
        Signal: &strategy.Signal{
            Strategy: "test",
            Actions: []strategy.TradeAction{
                // ... test actions
            },
        },
        Timestamp: time.Now(),
        Metadata:  make(map[string]interface{}),
    }
    
    err := hook.BeforeExecute(ctx)
    if err != nil {
        t.Errorf("Expected no error, got: %v", err)
    }
}
```

## Complete Example

See `examples/custom_hooks/` for a complete working example of custom hooks.

## Integration with Live Trading App

In your live trading application:

```go
import "github.com/wisp-trading/wisp/pkg/executor"

func main() {
    // Create executor
    exec := executor.NewDefaultExecutor(
        config.Connectors(),
        config.PositionStore(),
        config.Logger(),
        config.TimeProvider(),
    )
    
    // Add built-in hooks
    exec.RegisterHook(hooks.NewBasicRiskHook(...))
    
    // Load user's strategy plugin
    strategyPlugin := loadStrategyPlugin("./user-strategy.so")
    
    // Load user's custom hooks (optional)
    if fileExists("./user-hooks.so") {
        hookPlugin, err := executor.LoadHookPlugin("./user-hooks.so")
        if err == nil {
            for _, hook := range hookPlugin.CreateHooks() {
                exec.RegisterHook(hook)
            }
        }
    }
    
    // Start orchestrator with executor
    orch := orchestrator.New(k, exec)
    orch.Start()
}
```

## Dependencies

The executor requires these dependencies (injected via constructor):

- `ConnectorRegistry`: Access to exchange connectors
- `PositionStore`: Position and order tracking
- `Logger`: Logging functionality
- `TimeProvider`: Time-related operations

All dependencies are defined as interfaces in this package for maximum flexibility.

