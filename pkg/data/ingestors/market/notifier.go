package market

import "github.com/wisp-trading/wisp/pkg/types/data/ingestors"

// channelNotifier implements DataUpdateNotifier using a buffered channel
type channelNotifier struct {
	ch chan struct{}
}

// NewDataUpdateNotifier creates a new channel-based notifier
func NewDataUpdateNotifier(bufferSize int) ingestors.DataUpdateNotifier {
	return &channelNotifier{
		ch: make(chan struct{}, bufferSize),
	}
}

// Notify sends a signal (non-blocking)
func (n *channelNotifier) Notify() {
	select {
	case n.ch <- struct{}{}:
	default:
		// Channel full, skip notification
	}
}

// Updates returns the channel for receiving notifications
func (n *channelNotifier) Updates() <-chan struct{} {
	return n.ch
}
