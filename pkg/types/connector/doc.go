// Package connector provides type-safe interfaces for exchange connectors.
//
// The connector package is organized into subpackages:
//   - common: Shared interfaces (BaseConnector, MarketDataReader, OrderExecutor, AccountReader)
//   - spot: Spot market connector interfaces
//   - perp: Perpetual futures connector interfaces
//
// Example usage:
//
//	import (
//		"github.com/backtesting-org/kronos-sdk/pkg/types/connector/spot"
//		"github.com/backtesting-org/kronos-sdk/pkg/types/connector/perp"
//	)
//
//	// Use spot connector
//	var spotConn spot.Connector
//
//	// Use perp connector
//	var perpConn perp.Connector
//
// Legacy note: The old monolithic connector.Connector interface is deprecated.
// Use spot.Connector or perp.Connector for new implementations.
package connector
