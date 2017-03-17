package controllers

import (
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/aaparella/vidwell/models"
	"github.com/aaparella/vidwell/render"
	"github.com/aaparella/vidwell/storage"
	"github.com/aaparella/vidwell/users"
	"github.com/aaparella/vidwell/videos"
	"github.com/gorilla/mux"
	"github.com/siddontang/go/log"
)

// VideoController contains all endpoints and webpages regarding viewing
// and editing, uploading, or deleting videos.
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

// ViewVideos displays a list of all public videos. Videos that are private
// will not be included.
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

// ViewVideo displays a single video's page. Performs a check that the user that
// is logged in has access to this video, or if they are not logged in ensures
// that the video is public.
func (vc VideoController) ViewVideo(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Could not find video with id "+mux.Vars(r)["id"], http.StatusBadRequest)
		return
	}

	video, err := videos.GetVideo(uint(id))
	if err != nil {
		http.Error(w, "Could not find video : "+err.Error(), 500)
	}

	video.Views += 1
	storage.DB.Save(&video)

	var creator models.User
	storage.DB.Find(&creator, video.UserID)

	render.Render(w, r, "video", map[string]interface{}{
		"Video":    video,
		"User":     creator,
		"VideoUrl": videos.GetVideoUrl(video),
	})
}

// UploadVideo stores a video in S3, and creates a record for it in the database.
// Only creates the video record in the case of a successful video storage.
// Will redirect to the page where that video can be configured, e.g. title
// and permissions, etc.
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
		users.GetLoggedInUser(r))

	http.Redirect(w, r, "/", http.StatusCreated)
}
