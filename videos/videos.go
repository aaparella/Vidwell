package videos

import (
	"log"
	"time"

	"code.google.com/p/go-uuid/uuid"
	"github.com/aaparella/vidwell/storage"
)

type Video struct {
	ID        int
	Title     string
	Creator   string
	Timestamp time.Time
}

// UploadVideo uploads video content to content storage, nothing else
func UploadVideo(data []byte, contentType, name string) error {
	return storage.Upload(data, name, "videos", contentType)
}

// StoreVideo is a very badly named function
func StoreVideo(data []byte, contentType, title, creator string) {
	name := uuid.New()
	if err := UploadVideo(data, contentType, name); err != nil {
		log.Println("Could not store video : ", err)
		return
	}
	if err := CreateVideoRecord(title, creator, name); err != nil {
		log.Println("Could not create video record : ", err)
	}
}

func CreateVideoRecord(title, creator, uuid string) error {
	return storage.CreateVideoRecord(title, creator, uuid)
}
