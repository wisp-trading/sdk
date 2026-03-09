package monitoring

import (
	"net/http"
	"strings"

	prediction "github.com/wisp-trading/sdk/pkg/markets/prediction/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
)

func (s *server) handleOrderbook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	exchange := r.URL.Query().Get("exchange")
	if exchange == "" {
		http.Error(w, "exchange parameter required", http.StatusBadRequest)
		return
	}

	pair, ok := parsePairParam(w, r)
	if !ok {
		return
	}

	orderbook := s.viewRegistry.GetOrderbook(connector.ExchangeName(exchange), pair)
	if orderbook == nil {
		http.Error(w, "orderbook not found", http.StatusNotFound)
		return
	}

	s.writeJSON(w, orderbook)
}

func (s *server) handlePredictionOrderbook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	exchange := r.URL.Query().Get("exchange")
	marketID := r.URL.Query().Get("market_id")
	outcomeID := r.URL.Query().Get("outcome_id")

	if exchange == "" || marketID == "" || outcomeID == "" {
		http.Error(w, "exchange, market_id and outcome_id parameters required", http.StatusBadRequest)
		return
	}

	ob := s.viewRegistry.GetPredictionOrderbookView(
		connector.ExchangeName(exchange),
		prediction.MarketIDFromString(marketID),
		prediction.OutcomeIDFromString(outcomeID),
	)
	if ob == nil {
		http.Error(w, "orderbook not found", http.StatusNotFound)
		return
	}

	s.writeJSON(w, ob)
}

// parsePairParam parses the "pair" query param (e.g. "BTC-USDT") into a portfolio.Pair.
// Writes an error response and returns false if the param is missing or malformed.
func parsePairParam(w http.ResponseWriter, r *http.Request) (portfolio.Pair, bool) {
	raw := r.URL.Query().Get("pair")
	if raw == "" {
		http.Error(w, "pair parameter required", http.StatusBadRequest)
		return portfolio.Pair{}, false
	}

	parts := strings.Split(raw, "-")
	if len(parts) != 2 {
		http.Error(w, "invalid pair format, expected BASE-QUOTE", http.StatusBadRequest)
		return portfolio.Pair{}, false
	}

	return portfolio.NewPair(
		portfolio.NewAsset(parts[0]),
		portfolio.NewAsset(parts[1]),
	), true
}
