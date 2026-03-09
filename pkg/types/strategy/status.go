package strategy

import "time"

// StatusPhase is the coarse operational state the strategy reports.
// It is an open string type — strategies define their own phase values.
// The CLI uses this to colour-code the status header.
// Common conventions: "idle", "scanning", "in_trade", "cooling", "error".
type StatusPhase string

// StrategyStatus is a self-describing snapshot of a strategy's operational state.
// Strategies push updates via BaseStrategy.EmitStatus(). The view registry caches
// the latest snapshot per strategy and serves it over the monitoring socket.
//
// Fields is intentionally map[string]string — strategies describe their own state
// without needing interface changes. The CLI renders it as a two-column table.
//
// Example:
//
//	s.EmitStatus(strategy.StrategyStatus{
//	    Phase:   strategy.PhaseInTrade,
//	    Summary: "short BTC-USDT @ 65420, 4m elapsed, +0.02% toward TP",
//	    Fields: map[string]string{
//	        "entry_price": "65420.00",
//	        "entry_time":  active.EntryTime.Format("15:04:05"),
//	        "elapsed":     "4m12s",
//	        "tp_progress": "0.02%",
//	        "tp_required": "0.05%",
//	        "watchdog":    "active",
//	    },
//	})
type StrategyStatus struct {
	Phase   StatusPhase       `json:"phase"`
	Summary string            `json:"summary"`
	Fields  map[string]string `json:"fields,omitempty"`
	At      time.Time         `json:"at"`
}
