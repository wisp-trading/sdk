package monitoring

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/wisp-trading/sdk/pkg/types/monitoring"
)

// server implements the Server interface
type server struct {
	config       monitoring.ServerConfig
	viewRegistry monitoring.ViewRegistry
	socketPath   string
	listener     net.Listener
	httpServer   *http.Server
	shutdownFunc context.CancelFunc // Callback to trigger graceful shutdown
	mu           sync.Mutex
	started      bool
}

// NewServer creates a new monitoring server
func NewServer(
	config monitoring.ServerConfig,
	viewRegistry monitoring.ViewRegistry,
	shutdownFunc context.CancelFunc,
) (monitoring.Server, error) {
	if config.InstanceID == "" {
		return nil, fmt.Errorf("instance ID is required")
	}

	socketDir := config.SocketDir
	if socketDir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get home directory: %w", err)
		}
		socketDir = filepath.Join(homeDir, ".wisp", "sockets")
	}

	if err := os.MkdirAll(socketDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create socket directory: %w", err)
	}

	socketPath := filepath.Join(socketDir, fmt.Sprintf("%s.sock", config.InstanceID))

	return &server{
		config:       config,
		viewRegistry: viewRegistry,
		socketPath:   socketPath,
		shutdownFunc: shutdownFunc,
	}, nil
}

// Start begins listening on the Unix socket
func (s *server) Start() error {
	s.mu.Lock()
	if s.started {
		s.mu.Unlock()
		return fmt.Errorf("server already started")
	}

	// Remove old socket if exists
	if err := os.Remove(s.socketPath); err != nil && !os.IsNotExist(err) {
		s.mu.Unlock()
		return fmt.Errorf("failed to remove old socket: %w", err)
	}

	listener, err := net.Listen("unix", s.socketPath)
	if err != nil {
		s.mu.Unlock()
		return fmt.Errorf("failed to create unix socket: %w", err)
	}

	// Set socket permissions (owner-only access)
	if err := os.Chmod(s.socketPath, 0600); err != nil {
		_ = listener.Close()
		s.mu.Unlock()
		return fmt.Errorf("failed to set socket permissions: %w", err)
	}

	s.listener = listener
	s.started = true

	// Set up HTTP routes
	mux := http.NewServeMux()
	mux.HandleFunc("/health", s.handleHealth)
	mux.HandleFunc("/shutdown", s.handleShutdown)
	mux.HandleFunc("/api/pnl", s.handlePnL)
	mux.HandleFunc("/api/positions", s.handlePositions)
	mux.HandleFunc("/api/trades", s.handleTrades)
	mux.HandleFunc("/api/metrics", s.handleMetrics)
	mux.HandleFunc("/api/markets", s.handleMarkets)
	mux.HandleFunc("/api/orderbook", s.handleOrderbook)
	mux.HandleFunc("/api/orderbook/prediction", s.handlePredictionOrderbook)
	mux.HandleFunc("/api/klines", s.handleKlines)
	mux.HandleFunc("/profiling/stats", s.handleProfilingStats)
	mux.HandleFunc("/profiling/executions", s.handleProfilingExecutions)

	s.httpServer = &http.Server{Handler: mux}
	s.mu.Unlock()

	// Serve blocks until the server is stopped
	if err := s.httpServer.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("server error: %w", err)
	}

	return nil
}

// Stop gracefully shuts down the server
func (s *server) Stop(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.started {
		return nil
	}

	var errs []error

	if s.httpServer != nil {
		if err := s.httpServer.Shutdown(ctx); err != nil {
			errs = append(errs, fmt.Errorf("failed to shutdown http server: %w", err))
		}
	}

	if s.listener != nil {
		if err := s.listener.Close(); err != nil {
			errs = append(errs, fmt.Errorf("failed to close listener: %w", err))
		}
	}

	if err := os.Remove(s.socketPath); err != nil && !os.IsNotExist(err) {
		errs = append(errs, fmt.Errorf("failed to remove socket: %w", err))
	}

	s.started = false

	if len(errs) > 0 {
		return fmt.Errorf("stop errors: %v", errs)
	}

	return nil
}

// SocketPath returns the path to the Unix socket
func (s *server) SocketPath() string {
	return s.socketPath
}

func (s *server) writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, fmt.Sprintf("failed to encode response: %v", err), http.StatusInternalServerError)
	}
}
