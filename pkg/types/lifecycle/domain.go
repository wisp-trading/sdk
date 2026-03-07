package lifecycle

import (
	"context"

	"github.com/wisp-trading/sdk/pkg/types/config"
)

// DomainLifecycle is implemented by each market domain (spot, perp, prediction).
// The lifecycle controller calls Start/Stop on each registered domain in order,
// keeping startup and shutdown logic fully isolated per domain.
type DomainLifecycle interface {
	// Start begins data ingestion and any domain-specific startup work.
	// StartupConfig is passed through from the runtime so the domain can seed
	// its own watchlist without any graph or cross-domain coupling.
	Start(ctx context.Context, cfg *config.StartupConfig) error
	// Stop gracefully shuts down data ingestion and domain resources.
	Stop() error
	// Name returns a human-readable domain name for logging.
	Name() string
}
