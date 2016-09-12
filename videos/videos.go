package videos

import (
	"log"
	"net/url"

	"code.google.com/p/go-uuid/uuid"
	"github.com/aaparella/vidwell/models"
	"github.com/aaparella/vidwell/storage"
)

// StoreVideo is a very badly named function.
func StoreVideo(data []byte, contentType, title string, user *models.User) {
	name := uuid.New()
	if err := StoreVideoObject(data, contentType, name); err != nil {
		log.Println("Could not store video : ", err)
		return
	}
	if err := CreateVideoRecord(title, name, contentType, user.ID); err != nil {
		log.Println("Could not create video record : ", err)
	}
}

// UploadVideo uploads video content to content storage, nothing else
func StoreVideoObject(data []byte, contentType, name string) error {
	return storage.Upload(data, name, "vidwell.videos", contentType)
}

// CreateVideoRecord creates the database record with metadata about the
// uploaded video.
func CreateVideoRecord(title, uuid, content string, userID uint) error {
	return storage.DB.Create(&models.Video{
		Title:       title,
		Uuid:        uuid,
		ContentType: content,
		UserID:      userID,
	}).Error
}

func GetVideoUrl(uuid string) *url.URL {
	return storage.GetObjectUrl(uuid, "vidwell.videos")
}
