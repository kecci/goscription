package cmd

import (
	"github.com/abyanjksatu/goscription/internal/controller"
	"github.com/abyanjksatu/goscription/internal/database"
	"github.com/abyanjksatu/goscription/internal/outbound"
	"github.com/abyanjksatu/goscription/internal/server"
	"github.com/abyanjksatu/goscription/internal/service"
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
			NewTimeOutContext,
			NewDbConn,
		),
		server.Module,
		database.Module,
		service.Module,
		outbound.Module,
		controller.Module,
	)
}
