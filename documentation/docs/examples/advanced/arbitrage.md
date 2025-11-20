---
sidebar_position: 2
---

# Cross-Exchange Arbitrage

Find and exploit price differences across exchanges.

## Strategy Overview

- **Type**: Arbitrage
- **Indicators**: None (price-based)
- **Risk Level**: Low
- **Assets**: Single asset, multiple exchanges

## Complete Code

```go
package main

import (
    sdk "github.com/backtesting-org/kronos-sdk/pkg/kronos"
    "github.com/backtesting-org/kronos-sdk/pkg/types/connector"
    "github.com/backtesting-org/kronos-sdk/pkg/types/strategy"
    "github.com/shopspring/decimal"
)

type Arbitrage struct {
    k *sdk.Kronos
}

func NewArbitrage(k *sdk.Kronos) *Arbitrage {
    return &Arbitrage{k: k}
}

func (s *Arbitrage) GetSignals() ([]*strategy.Signal, error) {
    btc := s.k.Asset("BTC")
    
    // Get prices from all exchanges
    prices := s.k.Market().Prices(btc)
    
    // Find min and max
    var minPrice, maxPrice decimal.Decimal
    var minExchange, maxExchange connector.ExchangeName
    
    for exchange, price := range prices {
        if minPrice.IsZero() || price.LessThan(minPrice) {
            minPrice = price
            minExchange = exchange
        }
        if maxPrice.IsZero() || price.GreaterThan(maxPrice) {
            maxPrice = price
            maxExchange = exchange
        }
    }
    
    // Calculate spread
    spread := maxPrice.Sub(minPrice).Div(minPrice).Mul(decimal.NewFromInt(100))
    
    // If spread > 0.5%, arbitrage
    if spread.GreaterThan(decimal.NewFromFloat(0.5)) {
        s.k.Log().Opportunity("Arbitrage", btc.Symbol(),
            "Buy %s @ %s, Sell %s @ %s, Spread: %s%%",
            minExchange, minPrice, maxExchange, maxPrice, spread)
        
        qty := decimal.NewFromFloat(0.1)
        signal := s.k.Signal(s.GetName()).
            Buy(btc, minExchange, qty).
            Sell(btc, maxExchange, qty).
            Build()
        return []*strategy.Signal{signal}, nil
    }
    
    return nil, nil
}

func (s *Arbitrage) GetName() strategy.StrategyName { return "Arbitrage" }
func (s *Arbitrage) GetDescription() string { return "Cross-exchange arbitrage" }
func (s *Arbitrage) GetRiskLevel() strategy.RiskLevel { return strategy.RiskLevelLow }
func (s *Arbitrage) GetStrategyType() strategy.StrategyType { return strategy.StrategyTypeArbitrage }
```

## How It Works

1. **Fetch All Prices**: Get current price from all configured exchanges
2. **Find Extremes**: Identify lowest and highest prices
3. **Calculate Spread**: Compute percentage difference
4. **Execute Arbitrage**: Buy low exchange, sell high exchange simultaneously

## Key Concepts

- **Price Inefficiency**: Exchanges sometimes have different prices
- **Simultaneous Execution**: Buy and sell at the same time
- **Market Neutral**: No directional exposure
- **Low Risk**: Profit from spread, not price movement
- **Minimum Spread**: Only trade when spread > 0.5% (covers fees + slippage)

## Backtesting

Run with:

```bash
kronos backtest
```

Expected characteristics:
- Variable frequency (depends on market conditions)
- Small profits per trade
- High win rate
- Requires fast execution
- Sensitive to fees and slippage

## Important Considerations

### Fees
Make sure spread covers:
- Trading fees on both exchanges
- Withdrawal fees (if moving funds)
- Network fees (blockchain)

### Execution Risk
- Prices can change between signal and execution
- One leg may fill while other doesn't
- Requires sufficient liquidity on both exchanges

### Capital Requirements
- Need funds on both exchanges
- Or fast transfer between exchanges
- Consider rebalancing costs

## Improvements

Consider adding:
- Minimum profit threshold (not just spread)
- Liquidity checking (order book depth)
- Historical spread analysis
- Automatic rebalancing
- Multiple asset pairs

## Related Strategies

- [Portfolio](portfolio) - Multi-asset management
- [Best Price Execution](best-price) - Always use best exchange
