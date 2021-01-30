package cmd

import (
	"github.com/kecci/goscription/internal/controller"
	"github.com/kecci/goscription/internal/database"
	"github.com/kecci/goscription/internal/library/db"
	"github.com/kecci/goscription/internal/outbound"
	"github.com/kecci/goscription/internal/server"
	"github.com/kecci/goscription/internal/service"
	"github.com/kecci/goscription/util"
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

func init() {
	rootCmd.AddCommand(httpCmd)
}

func inject() fx.Option {
	return fx.Options(
		fx.Provide(
			util.NewTimeOutContext,
			db.NewMysqlDB,
		),
		server.Module,
		database.Module,
		service.Module,
		outbound.Module,
		controller.Module,
	)
}
