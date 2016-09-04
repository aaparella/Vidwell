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
	return CreateSchema()
}

func CreateSchema() error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS videos(
		name varchar(255),
		creator varchar(255),
		uuid varchar(255),
		created timestamp default current_timestamp
	);`)
	return err
}

func CreateVideoRecord(title, creator, uuid string) error {
	_, err := db.Exec(`INSERT INTO videos VALUES(
		$1, $2, $3
	);`, title, creator, uuid)
	return err
}
