package activity

import (
	"context"

	predTypes "github.com/wisp-trading/sdk/pkg/markets/prediction/types"
	predConn "github.com/wisp-trading/sdk/pkg/markets/prediction/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

type predictionPNL struct {
	store predTypes.MarketStore
}

func NewPredictionPNL(store predTypes.MarketStore) predTypes.PredictionPNL {
	return &predictionPNL{store: store}
}

func (p *predictionPNL) positions() predTypes.PositionsStoreExtension {
	return p.store.(predTypes.PositionsStoreExtension)
}

func (p *predictionPNL) orderbook() predTypes.OrderBookStoreExtension {
	return p.store.(predTypes.OrderBookStoreExtension)
}

func (p *predictionPNL) Positions(_ context.Context) []predTypes.PredictionPositionPNL {
	orders := p.positions().GetOrders()
	results := make([]predTypes.PredictionPositionPNL, 0, len(orders))
	for _, order := range orders {
		unrealized := p.impliedValue(order)
		realized := numerical.Zero()
		if order.Status == connector.OrderStatusFilled {
			realized = order.RealizedPnL
		}
		results = append(results, predTypes.PredictionPositionPNL{
			Order:      order,
			Realized:   realized,
			Unrealized: unrealized,
		})
	}
	return results
}

func (p *predictionPNL) Realized(_ context.Context) numerical.Decimal {
	total := numerical.Zero()
	for _, order := range p.positions().GetOrders() {
		if order.Status == connector.OrderStatusFilled {
			total = total.Add(order.RealizedPnL)
		}
	}
	return total
}

func (p *predictionPNL) Unrealized(_ context.Context) numerical.Decimal {
	total := numerical.Zero()
	for _, order := range p.positions().GetOrders() {
		total = total.Add(p.impliedValue(order))
	}
	return total
}

func (p *predictionPNL) Fees(_ context.Context) numerical.Decimal {
	total := numerical.Zero()
	for _, order := range p.positions().GetOrders() {
		total = total.Add(order.Fee)
	}
	return total
}

func (p *predictionPNL) impliedValue(order predTypes.PredictionOrder) numerical.Decimal {
	if order.Status != connector.OrderStatusFilled {
		return numerical.Zero()
	}

	ob := p.orderbook().GetOrderBook(order.Exchange, predConn.MarketID(order.MarketSlug), order.OutcomeID)
	if ob == nil || len(ob.Bids) == 0 || len(ob.Asks) == 0 {
		return numerical.Zero()
	}

	mid := ob.Bids[0].Price.Add(ob.Asks[0].Price).Div(numerical.NewFromInt(2))
	return order.Shares.Mul(mid).Sub(order.Shares.Mul(order.Price))
}

var _ predTypes.PredictionPNL = (*predictionPNL)(nil)
