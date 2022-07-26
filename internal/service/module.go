package service

import (
	"go.uber.org/fx"
)

// Module for service module
var Module = fx.Provide(
	NewArticleService,
	NewUserService,
	NewHealthService,
)
