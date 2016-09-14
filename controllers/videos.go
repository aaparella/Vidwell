package controllers

import (
	"io/ioutil"
	"net/http"

	"github.com/aaparella/vidwell/models"
	"github.com/aaparella/vidwell/render"
	"github.com/aaparella/vidwell/storage"
	"github.com/aaparella/vidwell/users"
	"github.com/aaparella/vidwell/videos"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/siddontang/go/log"
)

type VideoController struct {
}

func (vc VideoController) Prefix() string {
	return "/videos"
}

func (vc VideoController) Endpoints() map[string]map[string]http.HandlerFunc {
	return map[string]map[string]http.HandlerFunc{
		"": {
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
		log.Error("Could not find videos : ", err.Error())
		http.Error(w, "Error finding videos", http.StatusInternalServerError)
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
		if err != gorm.ErrRecordNotFound {
			http.Error(w, "Error accessing video", http.StatusInternalServerError)
			return
		}
		http.Error(w, "No video with that ID", http.StatusNotFound)
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
		http.Error(w, "Invalid form", http.StatusBadRequest)
		return
	}

	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		http.Error(w, "Could not read file", http.StatusInternalServerError)
		return
	}

	go videos.StoreVideo(data,
		handler.Header.Get("Content-Type"),
		r.FormValue("title"),
		users.GetUser(r))

	http.Redirect(w, r, "/", http.StatusCreated)
}
