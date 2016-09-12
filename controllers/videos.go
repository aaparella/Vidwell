package controllers

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/aaparella/vidwell/models"
	"github.com/aaparella/vidwell/render"
	"github.com/aaparella/vidwell/storage"
	"github.com/aaparella/vidwell/users"
	"github.com/aaparella/vidwell/videos"
	"github.com/gorilla/mux"
)

type VideoController struct {
}

func (vc VideoController) Prefix() string {
	return "/videos"
}

func (vc VideoController) Endpoints() map[string]map[string]http.HandlerFunc {
	return map[string]map[string]http.HandlerFunc{
		"/": {
			"GET": vc.ViewVideos,
		},
		"/{id}": {
			"GET": vc.ViewVideo,
		},
		"/upload": {
			"POST": users.MustBeLoggedIn(vc.UploadVideo),
		},
	}
}

func (vc VideoController) ViewVideos(w http.ResponseWriter, r *http.Request) {
	var videos []models.Video
	if err := storage.DB.Find(&videos).Error; err != nil {
		fmt.Fprintf(w, "Could not find videos : %s", err.Error())
		return
	}
	render.Render(w, r, "videos", map[string]interface{}{
		"Videos": videos,
	})
}

func (vc VideoController) ViewVideo(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var video models.Video
	var user models.User
	if err := storage.DB.Find(&video, id).Error; err != nil {
		fmt.Fprintf(w, "Could not find video with ID: %s", id)
		return
	}
	url := videos.GetVideoUrl(video.Uuid)
	video.Views += 1
	storage.DB.Save(&video)

	storage.DB.Find(&user, video.UserID)

	render.Render(w, r, "video", map[string]interface{}{
		"Video":    video,
		"User":     user,
		"VideoUrl": url.String(),
	})
}

func (vc VideoController) UploadVideo(w http.ResponseWriter, r *http.Request) {
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

	go videos.StoreVideo(data,
		handler.Header.Get("Content-Type"),
		r.FormValue("title"),
		users.GetUser(r))

	fmt.Fprintf(w, "Thank you for the video!")
}
