package connector

// WebSocketCapable provides WebSocket lifecycle management
type WebSocketCapable interface {
	StartWebSocket() error
	StopWebSocket() error
	IsWebSocketConnected() bool
	ErrorChannel() <-chan error
}
