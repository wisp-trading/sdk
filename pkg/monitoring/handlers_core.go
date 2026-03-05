package monitoring

import (
	"net/http"
	"strconv"
	"time"

	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/monitoring"
	"github.com/wisp-trading/sdk/pkg/types/monitoring/health"
)

func (s *server) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	healthReport := s.viewRegistry.GetHealth()
	if healthReport == nil {
		healthReport = &health.SystemHealthReport{OverallState: health.StateConnected}
	}

	s.writeJSON(w, healthReport)
}

func (s *server) handleShutdown(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	s.writeJSON(w, map[string]string{"status": "shutting down"})

	if s.shutdownFunc != nil {
		go func() {
			time.Sleep(100 * time.Millisecond)
			s.shutdownFunc()
		}()
	}
}

func (s *server) handlePnL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	pnl := s.viewRegistry.GetPnLView()
	if pnl == nil {
		s.writeJSON(w, struct{}{})
		return
	}

	s.writeJSON(w, pnl)
}

func (s *server) handlePositions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	positions := s.viewRegistry.GetPositionsView()
	if positions == nil {
		s.writeJSON(w, struct{}{})
		return
	}

	s.writeJSON(w, positions)
}

func (s *server) handleTrades(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	limit := 50
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	trades := s.viewRegistry.GetRecentTrades(limit)
	if trades == nil {
		trades = []connector.Trade{}
	}

	s.writeJSON(w, trades)
}

func (s *server) handleMetrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	metrics := s.viewRegistry.GetMetrics()
	if metrics == nil {
		s.writeJSON(w, &monitoring.StrategyMetrics{})
		return
	}

	s.writeJSON(w, metrics)
}
