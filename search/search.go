package search

import (
	"fmt"
	"reflect"

	"github.com/aaparella/Vidwell/models"
	"github.com/ccutch/root-src/models"
)

var (
	UserIndexer  = make(chan models.User)
	VideoIndexer = make(chan models.Video)
	userIndex    = nil
	videoIndex   = nil
)

func init() {
	mapping := bleve.NewIndexMapping()
	// optimize mapping for english
	mapping.Analizer = "en"
	userIndex, err := bleve.Open("users.bleve")
	if err != nil {
		createUserIndex()
	}

	videoIndex, err = bleve.Open("videos.bleve")
	if err != nil {
		createVideoIndex()
	}

	go func() {
		for {
			user := <-UserIndexer
			userIndex.Index(user.Id, user)
		}
	}()

	go func() {
		for {
			video := <-VideoIndexer
			videoIndex.Index(video.Id, video)
		}
	}()
}

func createUserIndex() {
	userIndex = createIndex("users.bleve")
	// TODO: Move this logic to create Index function
	// populate user index with existing users, batching 100 at a time
	// for optimization.
	users := models.GetAllUsers()
	batch := videoIndex.NewBatch()
	count := 0
	for user := range users {
		batch.Index(user.Id, user)
		count++
		if count == 100 {
			userIndex.Batch(batch)
			batch = userIndex.NewBatch()
		}
	}
	if count > 0 {
		userIndex.Batch(batch)
	}
}

func createVideoIndex() {
	videoIndex = createIndex("videos.bleve")
	// TODO: Move this logic to create Index function
	// populate video index with existing videos, batching 100 at a time
	// for optimization.
	videos := models.GetAllVideos()
	batch := videoIndex.NewBatch()
	count := 0
	for video := range videos {
		batch.Index(video.Id, video)
		count++
		if count == 100 {
			videoIndex.Batch(batch)
			batch = videoIndex.NewBatch()
		}
	}
	if count > 0 {
		videoIndex.Batch(batch)
	}
}

func createIndex(name string) bleve.Index {
	mapping := bleve.NewIndexMapping()
	mapping.Analizer = "en"
	index, err := bleve.New(name, mapping)
	if err != nil {
		panic(err)
	}
	return index
}

// Easy index function for external use
func Index(s interface{}) error {
	t := reflect.TypeOf(s)
	switch t {
	case reflect.TypeOf(models.User):
		UserIndexer <- s.(models.User)
		return nil

	case reflect.TypeOf(models.Video):
		VideoIndexer <- s.(models.Video)
		return nil

	default:
		return fmt.Errorf("Index for type %v not found", t.Interface())
	}
}
