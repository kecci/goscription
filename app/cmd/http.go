package cmd

import (
	"time"

	"github.com/abyanjksatu/goscription/internal/database/mysql"
	delivery "github.com/abyanjksatu/goscription/internal/http"
	"github.com/abyanjksatu/goscription/outbound"
	"github.com/abyanjksatu/goscription/usecase"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/dig"
)

var (
	httpCmd = &cobra.Command{
		Use:   "http",
		Short: "Start HTTP REST API",
		Run:   initHTTP,
	}
)

// NewTimeOutContext is timeout duration
func NewTimeOutContext() time.Duration {
	timeoutContext := time.Duration(viper.GetInt("contextTimeout")) * time.Second
	return timeoutContext
}

// BuildContainer is the container of all dependencies
func BuildContainer() *dig.Container {
	container := dig.New()

	container.Provide(NewTimeOutContext)
	container.Provide(NewDbConn)

	// repository
	container.Provide(mysql.NewArticleRepository)

	// outbound
	container.Provide(outbound.NewGodaddyOutbound)

	// usecase
	container.Provide(usecase.NewArticleUsecase)
	container.Provide(usecase.NewDomainUsecase)

	// handler
	container.Provide(delivery.NewServer)

	return container
}

func initHTTP(cmd *cobra.Command, args []string) {

	container := BuildContainer()

	err := container.Invoke(func(s *delivery.Server) {
		s.Run()
	})

	if err != nil {
		panic(err)
	}
}

func init() {
	rootCmd.AddCommand(httpCmd)
}
