package monitoring

import (
	"net/http"
	"strconv"

	"github.com/wisp-trading/sdk/pkg/types/connector"
)

// handleKlines handles kline requests for both spot and perp.
// The exchange param is used to automatically determine the market type.
// Query params: exchange, pair (e.g. "BTC-USDT"), interval (e.g. "1m"), limit (optional, default 100).
func (s *server) handleKlines(w http.ResponseWriter, r *http.Request) {
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

	klines := s.viewRegistry.GetKlines(connector.ExchangeName(exchange), pair, interval, limit)
	if klines == nil {
		klines = []connector.Kline{}
	}

	s.writeJSON(w, klines)
}
