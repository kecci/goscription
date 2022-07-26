package repository

import (
	"github.com/kecci/goscription/internal/repository/mysql"
	"github.com/kecci/goscription/internal/repository/postgres"
	"go.uber.org/fx"
)

// Module for controller database repository
var Module = fx.Options(
	fx.Provide(
		mysql.NewArticleRepository,
		mysql.NewUserRepository,
		postgres.NewAddressRepository,
	),
)
