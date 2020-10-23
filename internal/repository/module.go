package repository

import (
	"go.uber.org/fx"
)

//Module for controller database repository
var Module = fx.Options(
	fx.Provide(
		NewArticleRepository,
		NewUserRepository,
	),
)
