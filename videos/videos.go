package videos

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"code.google.com/p/go-uuid/uuid"
	"github.com/aaparella/vidwell/models"
	"github.com/aaparella/vidwell/render"
	"github.com/aaparella/vidwell/storage"
	"github.com/aaparella/vidwell/users"
	"github.com/gorilla/mux"
)

func UploadVideo(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(2 << 32)
	file, handler, err := r.FormFile("fileupload")
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	go StoreVideo(data,
		handler.Header.Get("Content-Type"),
		r.FormValue("title"),
		users.GetUser(r))

	fmt.Fprintf(w, "Thank you for the video!")
}

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

type VideoPageData struct {
	Video    models.Video
	VideoURL string
}

func ViewVideo(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var video models.Video
	if err := storage.DB.Find(&video, id).Error; err != nil {
		fmt.Fprintf(w, "Could not find video with ID: %s", id)
		return
	}
	url := storage.GetVideoUrl(video.Uuid)

	data := VideoPageData{
		Video:    video,
		VideoURL: url.String(),
	}

	render.Render(w, "video", data)
}

func ViewVideos(w http.ResponseWriter, r *http.Request) {
	var videos []models.Video
	if err := storage.DB.Find(&videos).Error; err != nil {
		fmt.Fprintf(w, "Could not find videos : %s", err.Error())
		return
	}
	render.Render(w, "videos", videos)
}
