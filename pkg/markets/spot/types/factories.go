package types

import (
	batchTypes "github.com/wisp-trading/sdk/pkg/markets/base/types/ingestors/batch"
	realtimeTypes "github.com/wisp-trading/sdk/pkg/markets/base/types/ingestors/realtime"
)

// SpotBatchIngestorFactory is the spot-domain-typed batch factory.
type SpotBatchIngestorFactory interface {
	batchTypes.BatchIngestorFactory
}

// SpotRealtimeIngestorFactory is the spot-domain-typed realtime factory.
type SpotRealtimeIngestorFactory interface {
	realtimeTypes.RealtimeIngestorFactory
}
