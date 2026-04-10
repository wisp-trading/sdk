package signal

import "go.uber.org/fx"

// Module is kept as a placeholder.
// Signal factories are now provided directly by domain packages:
//   - spot:    pkg/markets/spot/signal
//   - perp:    pkg/markets/perp/signal
//   - options: pkg/markets/options/signal
//   - predict: pkg/markets/prediction/signal
var Module = fx.Module("signal")
