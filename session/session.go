package session

import (
	"encoding/gob"
	"log"
	"net/http"

	"github.com/aaparella/vidwell/config"
	"github.com/aaparella/vidwell/models"
	"github.com/gorilla/sessions"
)

// Store is the CookieStore that all handlers can access session values
// through.
var store *sessions.CookieStore

func init() {
	store = sessions.NewCookieStore([]byte(config.GetSessionConfiguration().Key))
	RegisterModels()
}

// RegisterModels registers any models that we want to be able to store in a
// session, like the User struct.
func RegisterModels() {
	gob.Register(&models.User{})
}

// GetSession returns the vidwell session for the passed request.
// By making this the only way of accessing sessions from the cookie store,
// the session name is guaranteed to be consistent, and it also
// cuts down on boilerplate.
func GetSession(r *http.Request) *sessions.Session {
	s, err := store.Get(r, "vidwell")
	if err != nil {
		log.Println(err)
	}

	return s
}
