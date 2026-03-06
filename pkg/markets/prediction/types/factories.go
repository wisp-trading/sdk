package types

import (
	batchTypes "github.com/wisp-trading/sdk/pkg/markets/base/types/ingestors/batch"
	realtimeTypes "github.com/wisp-trading/sdk/pkg/markets/base/types/ingestors/realtime"
)

// PredictionBatchIngestorFactory is the prediction-domain-typed batch factory.
type PredictionBatchIngestorFactory interface {
	batchTypes.BatchIngestorFactory
}

// PredictionRealtimeIngestorFactory is the prediction-domain-typed realtime factory.
type PredictionRealtimeIngestorFactory interface {
	realtimeTypes.RealtimeIngestorFactory
}
