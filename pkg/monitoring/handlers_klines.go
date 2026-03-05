package monitoring

import (
	"net/http"
	"strconv"

	"github.com/wisp-trading/sdk/pkg/types/connector"
)

func (s *server) handleSpotKlines(w http.ResponseWriter, r *http.Request) {
	s.handleKlines(w, r, false)
}

func (s *server) handlePerpKlines(w http.ResponseWriter, r *http.Request) {
	s.handleKlines(w, r, true)
}

// handleKlines is the shared klines handler. isPerp selects spot vs perp.
// Query params: pair (e.g. "BTC-USDT"), exchange, interval (e.g. "1m"), limit (optional, default 100).
func (s *server) handleKlines(w http.ResponseWriter, r *http.Request, isPerp bool) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	pair, ok := parsePairParam(w, r)
	if !ok {
		return
	}

	exchange := r.URL.Query().Get("exchange")
	if exchange == "" {
		http.Error(w, "exchange parameter required", http.StatusBadRequest)
		return
	}

	interval := r.URL.Query().Get("interval")
	if interval == "" {
		interval = "1m"
	}

	limit := 100
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	var klines []connector.Kline
	if isPerp {
		klines = s.viewRegistry.GetPerpKlines(pair, exchange, interval, limit)
	} else {
		klines = s.viewRegistry.GetSpotKlines(pair, exchange, interval, limit)
	}

	if klines == nil {
		klines = []connector.Kline{}
	}

	s.writeJSON(w, klines)
}
