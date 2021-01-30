package cmd

import (
	"fmt"
	"os"

	"github.com/kecci/goscription/internal/library"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	// Version project version
	Version = "1.0.0-0"
	rootCmd = &cobra.Command{
		Use:     "goscription",
		Version: Version,
		Short:   "goscription Management CLI",
		Long:    `goscription is skeleton service for golang project`,
		Run: func(cmd *cobra.Command, args []string) {
			httpCmd.Run(cmd, args)
		},
	}

	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Show version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(rootCmd.Version)
		},
	}

	projectCmd = &cobra.Command{
		Use:   "project",
		Short: "Show project name",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(rootCmd.Use)
		},
	}

	mysqlCmd = &cobra.Command{
		Use:   "mysql",
		Short: "Show mysql connection name",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
)

// Execute will run the CLI app of goscription
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(library.InitConfig)
	rootCmd.AddCommand(httpCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(projectCmd)
	rootCmd.AddCommand(mysqlCmd)
}
