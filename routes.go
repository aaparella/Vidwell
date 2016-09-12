package main

import (
	"net/http"
	"strings"

	"github.com/aaparella/vidwell/controllers"
	"github.com/gorilla/mux"
)

type Controller interface {
	Prefix() string
	Endpoints() map[string]map[string]http.HandlerFunc
}

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
	register := func(c Controller) {
		subrouter := router.PathPrefix(c.Prefix()).Subrouter()

		for path, handlers := range c.Endpoints() {
			for methods, fn := range handlers {
				subrouter.HandleFunc(path, fn).Methods(strings.Split(methods, ", ")...)
			}
		}
	}

	register(controllers.VideoController{})
	register(controllers.UserController{})
}
