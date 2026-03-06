package types

import (
	batchTypes "github.com/wisp-trading/sdk/pkg/markets/base/types/ingestors/batch"
	realtimeTypes "github.com/wisp-trading/sdk/pkg/markets/base/types/ingestors/realtime"
)

// PerpBatchIngestorFactory is the perp-domain-typed batch factory.
type PerpBatchIngestorFactory interface {
	batchTypes.BatchIngestorFactory
}

// PerpRealtimeIngestorFactory is the perp-domain-typed realtime factory.
type PerpRealtimeIngestorFactory interface {
	realtimeTypes.RealtimeIngestorFactory
}
