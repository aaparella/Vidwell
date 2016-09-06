package videos

import (
	"fmt"
	"log"
	"net/http"

	"code.google.com/p/go-uuid/uuid"
	"github.com/aaparella/vidwell/models"
	"github.com/aaparella/vidwell/storage"
	"github.com/aaparella/vidwell/views"
	"github.com/gorilla/mux"
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

func ViewVideo(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var video models.Video
	if err := storage.DB.Find(&video, id).Error; err != nil {
		fmt.Fprintf(w, "Could not find video with ID: %s", id)
		return
	}
	views.Render(w, "video", video)
}

func ViewVideos(w http.ResponseWriter, r *http.Request) {
	var videos []models.Video
	if err := storage.DB.Find(&videos).Error; err != nil {
		fmt.Fprintf(w, "Could not find videos : %s", err.Error())
		return
	}
	views.Render(w, "videos", videos)
}
