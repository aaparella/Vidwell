package storage

import (
	"fmt"

	"github.com/aaparella/vidwell/config"
	"github.com/sirupsen/logrus"
)

// Create storage client and test connection
func init() {
	logrus.Printf("Initializing data storage...")
	conf := config.GetStorageConfiguration()
	if err := InitializeObjectStorage(conf); err != nil {
		logrus.Fatal("Error initializing storage : ", err)
	}
	if err := InitializeDatabase(conf); err != nil {
		logrus.Fatal("Error initializaing database : ", err)
	}
	logrus.Println("Storage configured!")
}

// Teardown is called in main when the application is closed, and closes any
// database connections or other needed steps. S3 doesn't require any
// additional work, so we only worry about the database here.
func Teardown() error {
	if err := TeardownDatabase(); err != nil {
		return fmt.Errorf("Error closing database connections: %s", err.Error())
	}
	return nil
}
