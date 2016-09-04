package storage

import (
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
