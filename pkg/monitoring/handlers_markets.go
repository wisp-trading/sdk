package monitoring

import (
	"net/http"
)

func (s *server) handleMarkets(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	s.writeJSON(w, s.viewRegistry.GetMarketViews())
}
