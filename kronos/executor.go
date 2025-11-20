package kronos

import (
	"github.com/backtesting-org/kronos-sdk/kronos/trade"
)

// KronosExecutor embeds the base Kronos and adds trade execution capabilities.
// This is only used by the orchestrator during the execution phase.
// Strategies in GetSignals methods receive the base Kronos without trade access.
type KronosExecutor struct {
	kronos // Embed base Kronos for all read operations

	// Trade service for execution - only available in executor
	Trade *trade.TradeService
}

// NewKronosExecutor creates a new KronosExecutor with full capabilities.
// This should only be used by the orchestrator, not injected into strategies.
func NewKronosExecutor(baseKronos kronos, tradeService *trade.TradeService) *KronosExecutor {
	return &KronosExecutor{
		kronos: baseKronos,
		Trade:  tradeService,
	}
}
