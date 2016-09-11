package storage

import (
	"fmt"
	"sync"

	"github.com/aaparella/vidwell/config"
	"github.com/aaparella/vidwell/models"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

var Mutex *sync.Mutex
var DB *gorm.DB

func InitializeDatabase(conf config.StorageConfiguration) error {
	var err error
	DB, err = gorm.Open("postgres", conf.Database)
	if err != nil {
		return fmt.Errorf("Could not connect to database %s:%s",
			conf.Database, err.Error())
	}
	DB.LogMode(conf.DatabaseLog)
	return MigrateTables()
}

// TeardownDatabase is called from main on application exit. This is to allow
// for easy cleanup. This is the *only* place that db.Close() will be called
// from in the entire application.
func TeardownDatabase() error {
	return DB.Close()
}

// MigrateTables creates all tables needed if they do not already exist.
// Important to note that it will not delete rows that are no longer necessary,
// if a field is removed from a model. The table must be dropped and recreated,
// or altered directly through a SQL interface.
func MigrateTables() error {
	m := []interface{}{
		&models.User{},
		&models.Video{},
	}

	if err := DB.AutoMigrate(m...).Error; err != nil {
		return fmt.Errorf("Could not migrate schema : %s", err.Error())
	}
	return nil
}
