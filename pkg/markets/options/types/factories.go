package types

import (
	batchTypes "github.com/wisp-trading/sdk/pkg/markets/base/types/ingestors/batch"
	realtimeTypes "github.com/wisp-trading/sdk/pkg/markets/base/types/ingestors/realtime"
)

// OptionsBatchIngestorFactory creates batch ingestors for options
type OptionsBatchIngestorFactory interface {
	CreateIngestors() []batchTypes.BatchIngestor
}

// OptionsRealtimeIngestorFactory creates realtime ingestors for options
type OptionsRealtimeIngestorFactory interface {
	CreateIngestors() []realtimeTypes.RealtimeIngestor
}
