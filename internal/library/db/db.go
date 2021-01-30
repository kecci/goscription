package db

import (
	"database/sql"
	"fmt"
	"net/url"

	"github.com/ory/viper"
)

// NewMysqlDB generate new database connection
func NewMysqlDB() (*sql.DB, error) {
	// DATABASE
	dbHost := viper.GetString(`database.host`)
	dbPort := viper.GetString(`database.port`)
	dbUser := viper.GetString(`database.user`)
	dbPass := viper.GetString(`database.pass`)
	dbName := viper.GetString(`database.name`)
	connection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	val := url.Values{}
	val.Add("parseTime", "1")
	val.Add("loc", "Asia/Jakarta")
	dsn := fmt.Sprintf("%s?%s", connection, val.Encode())
	return sql.Open(`mysql`, dsn)
}
