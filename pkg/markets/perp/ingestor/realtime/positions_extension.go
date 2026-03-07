package realtime

import (
	"context"

	"github.com/wisp-trading/sdk/pkg/markets/base/types/ingestors/realtime"
	perpTypes "github.com/wisp-trading/sdk/pkg/markets/perp/types"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	perpConn "github.com/wisp-trading/sdk/pkg/types/connector/perp"
	"github.com/wisp-trading/sdk/pkg/types/logging"
	"github.com/wisp-trading/sdk/pkg/types/portfolio"
)

// positionsExtension handles WebSocket subscriptions for live perp position updates.
// Writes into PerpPositionsStoreExtension — the single source of truth for perp positions.
type positionsExtension struct {
	store  perpTypes.MarketStore
	logger logging.ApplicationLogger
}

func NewPositionsExtension(store perpTypes.MarketStore, logger logging.ApplicationLogger) realtime.WebSocketExtension {
	return &positionsExtension{
		store:  store,
		logger: logger,
	}
}

func (p *positionsExtension) Subscribe(wsConn interface{}, exchangeName connector.ExchangeName, assets []portfolio.Pair) error {
	perpWS, ok := wsConn.(perpConn.WebSocketConnector)
	if !ok {
		p.logger.Debug("WebSocket connector %s does not support perp operations", exchangeName)
		return nil
	}

	for _, asset := range assets {
		if err := perpWS.SubscribePositions(asset); err != nil {
			p.logger.Warn("Failed to subscribe to positions for %s on %s: %v",
				asset.Symbol(), exchangeName, err)
		} else {
			p.logger.Info("Subscribed to position updates for %s on %s", asset.Symbol(), exchangeName)
		}
	}

	return nil
}

func (p *positionsExtension) ProcessChannels(wsConn interface{}, exchangeName connector.ExchangeName, ctx context.Context) {
	perpWS, ok := wsConn.(perpConn.WebSocketConnector)
	if !ok {
		return
	}

	positionChan := perpWS.PositionUpdates()
	p.logger.Info("Starting position channel processor for %s", exchangeName)

	for {
		select {
		case <-ctx.Done():
			p.logger.Debug("Context cancelled, stopping position channel for %s", exchangeName)
			return
		case update, ok := <-positionChan:
			if !ok {
				p.logger.Debug("Position channel closed for %s", exchangeName)
				return
			}
			p.handlePositionUpdate(exchangeName, update)
		}
	}
}

func (p *positionsExtension) handlePositionUpdate(exchangeName connector.ExchangeName, position perpConn.Position) {
	position.Exchange = exchangeName

	if position.Size.IsZero() {
		p.store.RemovePosition(exchangeName, position.Pair)
		p.logger.Debug("Position closed for %s on %s", position.Pair.Symbol(), exchangeName)
	} else {
		p.store.UpsertPosition(position)
		p.logger.Debug("Position updated for %s on %s: size=%s", position.Pair.Symbol(), exchangeName, position.Size.String())
	}
}

func (p *positionsExtension) Unsubscribe(wsConn interface{}, exchangeName connector.ExchangeName) error {
	perpWS, ok := wsConn.(perpConn.WebSocketConnector)
	if !ok {
		return nil
	}
	p.logger.Info("Unsubscribing from position updates for %s", exchangeName)
	_ = perpWS
	return nil
}

var _ realtime.WebSocketExtension = (*positionsExtension)(nil)
