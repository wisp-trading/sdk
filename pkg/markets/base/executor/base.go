package executor

import (
	"fmt"

	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/logging"
	"github.com/wisp-trading/sdk/pkg/types/registry"
	"github.com/wisp-trading/sdk/pkg/types/temporal"
)

// Base holds the common dependencies shared by all domain executors.
// Embed it in each domain executor struct to avoid duplicating fields and
// the connector-lookup boilerplate.
type Base struct {
	Connectors   registry.ConnectorRegistry
	Logger       logging.ApplicationLogger
	TimeProvider temporal.TimeProvider
}

// GetOrderExecutor resolves the connector for exchange and asserts it supports
// order execution. Returns a clear error if either check fails.
func (b *Base) GetOrderExecutor(exchange connector.ExchangeName) (connector.OrderExecutor, error) {
	conn, exists := b.Connectors.Connector(exchange)
	if !exists {
		return nil, fmt.Errorf("exchange %s not available", exchange)
	}

	exec, ok := conn.(connector.OrderExecutor)
	if !ok {
		return nil, fmt.Errorf("exchange %s does not support order execution", exchange)
	}

	return exec, nil
}
