package db

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"

	"github.com/ory/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Database struct
type Database struct {
	Mysql    *sql.DB
	Postgres *gorm.DB
}

// NewDB initiate database
func NewDB() Database {
	return Database{
		Mysql:    newMysqlDB(),
		Postgres: newPostgresDB(),
	}
}

// NewMysqlDB generate mysql database
func newMysqlDB() *sql.DB {
	dbHost := viper.GetString(`database.mysql.host`)
	dbPort := viper.GetString(`database.mysql.port`)
	dbUser := viper.GetString(`database.mysql.user`)
	dbPass := viper.GetString(`database.mysql.pass`)
	dbName := viper.GetString(`database.mysql.name`)
	connection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	val := url.Values{}
	val.Add("parseTime", "1")
	val.Add("loc", "Asia/Jakarta")
	dsn := fmt.Sprintf("%s?%s", connection, val.Encode())
	db, err := sql.Open(`mysql`, dsn)
	if err != nil {
		fmt.Println(err.Error())
	}
	return db
}

// NewPostgresDB generate postgres db
func newPostgresDB() *gorm.DB {
	dbHost := viper.GetString(`database.postgres.host`)
	dbPort := viper.GetString(`database.postgres.port`)
	dbUser := viper.GetString(`database.postgres.user`)
	dbPass := viper.GetString(`database.postgres.pass`)
	dbName := viper.GetString(`database.postgres.name`)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPass, dbHost, dbPort, dbName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println(err.Error())
	}
	return db
}
