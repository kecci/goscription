package cmd

import (
	"github.com/kecci/goscription/internal/controller"
	"github.com/kecci/goscription/internal/http"
	"github.com/kecci/goscription/internal/library"
	"github.com/kecci/goscription/internal/library/db"
	"github.com/kecci/goscription/internal/repository"
	"github.com/kecci/goscription/internal/service"
	"github.com/kecci/goscription/utility"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

var (
	httpCmd = &cobra.Command{
		Use:   "http",
		Short: "Start HTTP REST API",
		Run:   initHTTP,
	}
)

func initHTTP(cmd *cobra.Command, args []string) {
	fx.New(inject()).Run()
}

func inject() fx.Option {
	return fx.Options(
		fx.Provide(
			library.NewConfig,
			utility.NewTimeOutContext,
			db.NewMysqlDB,
		),
		repository.Module,
		service.Module,
		controller.Module,
		http.Module,
	)
}
