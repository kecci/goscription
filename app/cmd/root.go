package cmd

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"os"

	// imported for Mysql driver
	"github.com/abyanjksatu/goscription/internal/database/mysql"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	rootCmd = &cobra.Command{
		Use:   "goscription",
		Short: "Article Management CLI",
	}
	articleRepository mysql.ArticleRepository
	dbConn            *sql.DB
)

// Execute will run the CLI app of goscription
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig, initHTTPServiceDependecies)
}

func initConfig() {
	viper.SetConfigType("toml")
	viper.SetConfigFile("config.toml")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}

	logrus.Info("Using Config file: ", viper.ConfigFileUsed())

	if viper.GetBool("debug") {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Warn("Comment service is Running in Debug Mode")
		return
	}
	logrus.SetLevel(logrus.InfoLevel)
	logrus.Warn("Comment service is Running in Production Mode")
	logrus.SetFormatter(&logrus.JSONFormatter{})
}

// NewDbConn generate new database connection
func NewDbConn() (*sql.DB, error) {
	// DATABASE
	dbHost := viper.GetString(`database.host`)
	dbPort := viper.GetString(`database.port`)
	dbUser := viper.GetString(`database.user`)
	dbPass := viper.GetString(`database.pass`)
	dbName := viper.GetString(`database.name`)
	connection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	fmt.Println(connection)
	val := url.Values{}
	val.Add("parseTime", "1")
	val.Add("loc", "Asia/Jakarta")
	dsn := fmt.Sprintf("%s?%s", connection, val.Encode())
	return sql.Open(`mysql`, dsn)
}

func initHTTPServiceDependecies() {
	dbConn, err := NewDbConn()

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	err = dbConn.Ping()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
