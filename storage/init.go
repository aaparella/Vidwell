package storage

import (
	"fmt"
	"log"

	"github.com/aaparella/vidwell/config"
)

// Create storage client and test connection
func init() {
	log.Printf("Initializing data storage...")
	conf := config.GetStorageConfiguration()
	if err := InitializeObjectStorage(conf); err != nil {
		log.Fatal("Error initializing storage : ", err)
	}
	if err := InitializeDatabase(conf); err != nil {
		log.Fatal("Error initializaing database : ", err)
	}
	log.Println("Storage configured!")
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
