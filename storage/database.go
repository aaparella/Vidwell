package storage

import (
	"database/sql"
	"fmt"

	"github.com/aaparella/vidwell/config"
	_ "github.com/lib/pq"
)

var db *sql.DB

func InitializeDatabase(conf config.StorageConfiguration) error {
	var err error
	db, err = sql.Open("postgres", conf.Database+"&dbname=vidwell")
	if err != nil {
		return fmt.Errorf("Could not connect to database %s:%s",
			conf.Database, err.Error())
	}
	return nil
}
