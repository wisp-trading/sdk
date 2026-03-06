package lifecycle

import "context"

// DomainLifecycle is implemented by each market domain (spot, perp, prediction).
// The lifecycle controller calls Start/Stop on each registered domain in order,
// keeping startup and shutdown logic fully isolated per domain.
type DomainLifecycle interface {
	// Start begins data ingestion and any domain-specific startup work.
	Start(ctx context.Context) error
	// Stop gracefully shuts down data ingestion and domain resources.
	Stop() error
	// Name returns a human-readable domain name for logging.
	Name() string
}
