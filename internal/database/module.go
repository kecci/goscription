package database

import (
	"github.com/kecci/goscription/internal/database/mysql"
	"go.uber.org/fx"
)

//Module for controller database repository
var Module = fx.Options(
	fx.Provide(
		mysql.NewArticleRepository,
		mysql.NewUserRepository,
	),
)
