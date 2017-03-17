package main

import (
	"net/http"
	"strings"

	"github.com/aaparella/vidwell/controllers"
	"github.com/gorilla/mux"
)

// Controller defines a web controller that can be registered with a mux
// to serve webpages and endpoints.
type Controller interface {
	// Prefix is the root path that all endpoints will be accessible through
	// for this controller.
	Prefix() string
	// Endpoints defines the web pages and APIs that this controller offers.
	// The first key is the path for the handler, and the second maps each
	// method to the appropriate handler.
	//
	// map[string]map[string]http.HandlerFunc{
	//   "/users": {
	//   	"GET": c.ViewUsers,
	//   },
	//   "/users/new": {
	//	 	"GET": c.NewUserPage,
	//   	"POST": c.MakeNewUser,
	//   }
	// }
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
	// register adds the routes defined by the controller to the router.
	register := func(c Controller) {
		registerController(router, c)
	}

	// Register a controller here in order to register it's endpoints
	register(controllers.VideoController{})
	register(controllers.UserController{})
}

func registerController(router *mux.Router, c Controller) {
	subrouter := router.PathPrefix(c.Prefix()).Subrouter()

	for path, handlers := range c.Endpoints() {
		for methods, fn := range handlers {
			subrouter.HandleFunc(path, fn).Methods(strings.Split(methods, ", ")...)
		}
	}
}
