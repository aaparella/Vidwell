package videos

import (
	"log"

	"code.google.com/p/go-uuid/uuid"
	"github.com/aaparella/vidwell/models"
	"github.com/aaparella/vidwell/storage"
)

// StoreVideo is a very badly named function.
func StoreVideo(data []byte, contentType, title, creator string) {
	name := uuid.New()
	if err := UploadVideo(data, contentType, name); err != nil {
		log.Println("Could not store video : ", err)
		return
	}
	if err := CreateVideoRecord(title, name, contentType, 0); err != nil {
		log.Println("Could not create video record : ", err)
	}
}

// UploadVideo uploads video content to content storage, nothing else
func UploadVideo(data []byte, contentType, name string) error {
	return storage.Upload(data, name, "videos", contentType)
}

// CreateVideoRecord creates the database record with metadata about the
// uploaded video.
func CreateVideoRecord(title, uuid, content string, userID uint) error {
	return storage.DB.Create(&models.Video{
		Title:       title,
		Uuid:        uuid,
		ContentType: content,
	}).Error
}
