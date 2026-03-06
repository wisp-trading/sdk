package activity

import (
	"context"

	storeTypes "github.com/wisp-trading/sdk/pkg/markets/base/types/stores/activity"
	perpTypes "github.com/wisp-trading/sdk/pkg/markets/perp/types"
	"github.com/wisp-trading/sdk/pkg/types/connector"
	perpConn "github.com/wisp-trading/sdk/pkg/types/connector/perp"
	"github.com/wisp-trading/sdk/pkg/types/registry"
	"github.com/wisp-trading/sdk/pkg/types/wisp/numerical"
)

// perpPNL sources PNL directly from live connector positions.
// The exchange owns realized/unrealized — we don't recalculate.
type perpPNL struct {
	connectors registry.ConnectorRegistry
	trades     storeTypes.Trades
}

func NewPerpPNL(connectors registry.ConnectorRegistry, trades storeTypes.Trades) perpTypes.PerpPNL {
	return &perpPNL{connectors: connectors, trades: trades}
}

func (p *perpPNL) Positions(_ context.Context) []perpTypes.PerpPositionPNL {
	positions := p.livePositions()
	results := make([]perpTypes.PerpPositionPNL, 0, len(positions))
	for _, pos := range positions {
		results = append(results, perpTypes.PerpPositionPNL{
			Position:   pos,
			Realized:   pos.RealizedPnL,
			Unrealized: pos.UnrealizedPnL,
		})
	}
	return results
}

func (p *perpPNL) Realized(_ context.Context) numerical.Decimal {
	total := numerical.Zero()
	for _, pos := range p.livePositions() {
		total = total.Add(pos.RealizedPnL)
	}
	return total
}

func (p *perpPNL) Unrealized(_ context.Context) numerical.Decimal {
	total := numerical.Zero()
	for _, pos := range p.livePositions() {
		total = total.Add(pos.UnrealizedPnL)
	}
	return total
}

func (p *perpPNL) Fees(_ context.Context) numerical.Decimal {
	return sumTradeFees(p.trades.GetAllTrades())
}

// livePositions fetches current open positions across all ready perp connectors.
func (p *perpPNL) livePositions() []perpConn.Position {
	perpConnectors := p.connectors.FilterPerp(
		registry.NewFilter().ReadyOnly().Build(),
	)
	var all []perpConn.Position
	for _, conn := range perpConnectors {
		pm, ok := conn.(perpConn.PositionManager)
		if !ok {
			continue
		}
		positions, err := pm.GetPositions()
		if err != nil {
			continue
		}
		all = append(all, positions...)
	}
	return all
}

var _ perpTypes.PerpPNL = (*perpPNL)(nil)

func sumTradeFees(trades []connector.Trade) numerical.Decimal {
	total := numerical.Zero()
	for _, t := range trades {
		total = total.Add(t.Fee)
	}
	return total
}
