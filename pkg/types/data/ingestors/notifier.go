package ingestors

// DataUpdateNotifier signals when market data has been updated
type DataUpdateNotifier interface {
	// Notify signals that data has been updated (non-blocking)
	Notify()

	// Updates returns a read-only channel that receives update notifications
	Updates() <-chan struct{}
}
