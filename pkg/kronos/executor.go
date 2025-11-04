package kronos

import (
	"github.com/backtesting-org/kronos-sdk/pkg/types/logging"
	"github.com/backtesting-org/kronos-sdk/pkg/types/portfolio/store"
)

// KronosExecutor embeds the base Kronos and adds trade execution capabilities.
// This is only used by the orchestrator during the execution phase.
// Strategies in GetSignals methods receive the base Kronos without trade access.
type KronosExecutor struct {
	*Kronos // Embed base Kronos for all read operations

	// Trade service for execution - only available in executor
	Trade *TradeService
}

// NewKronosExecutor creates a new KronosExecutor with full capabilities.
// This should only be used by the orchestrator, not injected into strategies.
func NewKronosExecutor(store store.Store, logger logging.ApplicationLogger) *KronosExecutor {
	// Create base Kronos
	baseKronos := NewKronos(store, logger)

	// Create executor with trade service
	executor := &KronosExecutor{
		Kronos: baseKronos,
		Trade: &TradeService{
			logger: logger,
		},
	}

	return executor
}

// Example usage pattern:
//
// In strategy GetSignals (receives base Kronos):
//   func (s *MyStrategy) GetSignals() ([]*Signal, error) {
//       sma := s.k.Indicators.SMA(BTC, 20)
//       price := s.k.Market.Price(BTC)
//       // Cannot access s.k.Trade - not available
//       return signals, nil
//   }
//
// In orchestrator execution phase:
//   executor := kronos.NewKronosExecutor(store, logger)
//   result, err := executor.Trade.Buy(asset, exchange, quantity)
