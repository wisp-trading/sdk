package onchain

import (
	"context"

	"github.com/wisp-trading/sdk/pkg/types/config"
	lifecycleTypes "github.com/wisp-trading/sdk/pkg/types/lifecycle"
	"github.com/wisp-trading/sdk/pkg/types/logging"
)

// noopLifecycle is used when the domain has no batch/realtime ingestors
// (Dex / strategy supply market data; connector is swap-only).
type noopLifecycle struct {
	name   string
	logger logging.ApplicationLogger
}

func newOnchainDomainLifecycle(logger logging.ApplicationLogger) lifecycleTypes.DomainLifecycle {
	return &noopLifecycle{name: "onchain", logger: logger}
}

func (n *noopLifecycle) Start(_ context.Context, _ *config.StartupConfig) error {
	n.logger.Info("onchain domain started (no market-data ingestors)")
	return nil
}

func (n *noopLifecycle) Stop() error { return nil }

func (n *noopLifecycle) Name() string { return n.name }
