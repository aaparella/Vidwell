package main

import (
	"github.com/aaparella/vidwell/users"
	"github.com/aaparella/vidwell/videos"
	"github.com/gorilla/mux"
)

// RegisterRoutes prepares the router with all routes needed for the entire
// application.
//
// This router will be used in main as the router for http.ListenAndServe.
// All services in the application must register routes here, so that all
// routes can be seen from a single location. This allows route configuration
// to be kept separate from the rest of the application's logic.
//
// Routes can use variable syntax supported by gorilla mux. Variables in the
// path can then be fetched by using mux.Vars(r) in the handler. For example
// the video viewing page:
//      mux.Register("/videos/{id}", videos.ViewVideo)
//
// Then in the handler, the value of id can be retrieved:
//      keys := mux.Vars(r)
//      id, ok := keys["id"]
func RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/video/{id}", videos.ViewVideo)
	router.HandleFunc("/videos", videos.ViewVideos)

	router.HandleFunc("/user/{id}", users.ViewUser)
	router.HandleFunc("/upload", users.MustBeLoggedIn(videos.UploadVideo))

	router.HandleFunc("/register", users.NewUser).Methods("POST")
	router.HandleFunc("/login", users.Login).Methods("POST")
}
