package outbound

import (
	"go.uber.org/fx"
)

//Module for outbound module
var Module = fx.Provide(
	NewGodaddyOutbound,
)
