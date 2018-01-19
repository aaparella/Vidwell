package videos

import (
	"net/url"

	"github.com/aaparella/vidwell/models"
	"github.com/aaparella/vidwell/storage"
	"github.com/pborman/uuid"
	"github.com/sirupsen/logrus"
)

// StoreVideo is a very badly named function.
func StoreVideo(data []byte, contentType, title string, user *models.User) {
	name := uuid.New()
	if err := StoreVideoObject(data, contentType, name); err != nil {
		logrus.Error("Could not store video : ", err)
		return
	}
	if err := CreateVideoRecord(title, name, contentType, user.ID); err != nil {
		logrus.Error("Could not create video record : ", err)
	}
}

// GetVideo fetches the video model for a given id.
func GetVideo(id uint) (models.Video, error) {
	var video models.Video
	storage.DB.Find(&video, id)
	return video, storage.DB.Error
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

// GetVideoUrl gets a publicly accessible URL for the video with the specified
// uuid.
func GetVideoUrl(video models.Video) *url.URL {
	return storage.GetObjectUrl(video.Uuid, "vidwell.videos")
}
