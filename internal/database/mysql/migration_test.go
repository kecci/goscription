package mysql_test

import (
	"database/sql"
	"strings"

	"github.com/golang-migrate/migrate"
	_mysql "github.com/golang-migrate/migrate/database/mysql"
)

type migration struct {
	Migrate *migrate.Migrate
}

func (this *migration) Up() (error, bool) {
	err := this.Migrate.Up()
	if err != nil {
		if err == migrate.ErrNoChange {
			return nil, true
		}
		return err, false
	}
	return nil, true
}

func (this *migration) Down() (error, bool) {
	err := this.Migrate.Down()
	if err != nil {
		return err, false
	}
	return nil, true
}

func runMigration(dbConn *sql.DB, migrationsFolderLocation string) (*migration, error) {
	dataPath := []string{}
	dataPath = append(dataPath, "file://")
	dataPath = append(dataPath, migrationsFolderLocation)

	pathToMigrate := strings.Join(dataPath, "")

	driver, err := _mysql.WithInstance(dbConn, &_mysql.Config{})
	if err != nil {
		return nil, err
	}

	m, err := migrate.NewWithDatabaseInstance(pathToMigrate, mysql, driver)
	if err != nil {
		return nil, err
	}
	return &migration{Migrate: m}, nil
}
