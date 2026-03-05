package monitoring

import (
	"net/http"
	"strconv"

	"github.com/wisp-trading/sdk/pkg/types/monitoring"
)

func (s *server) handleProfilingStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	stats := s.viewRegistry.GetProfilingStats()
	if stats == nil {
		s.writeJSON(w, &monitoring.ProfilingStats{})
		return
	}

	s.writeJSON(w, stats)
}

func (s *server) handleProfilingExecutions(w http.ResponseWriter, r *http.Request) {
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

	executions := s.viewRegistry.GetRecentExecutions(limit)
	if executions == nil {
		executions = []monitoring.ProfilingMetrics{}
	}

	s.writeJSON(w, executions)
}
