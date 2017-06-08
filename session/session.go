package session

import (
	"encoding/gob"
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

// GetSessionValue returns the value for the provided key, for the provided request.
func GetSessionValue(key string, r *http.Request) interface{} {
	return getSession(r).Values[key]
}

// StoreSessionValue stores the passed value as the value for the passed key for the
// passed request. The underlying GetSession method will create a session if one does not
// already exist, and store the value. Error propogated back up so that callers
// can indicate failure to the user if necessary
func StoreSessionValue(key string, val interface{}, r *http.Request, w http.ResponseWriter) error {
	s := getSession(r)
	s.Values[key] = val
	return s.Save(r, w)
}

// getSession returns the vidwell session for the passed request.
// By making this the only way of accessing sessions from the cookie store,
// the session name is guaranteed to be consistent, and it also
// cuts down on boilerplate.
//
// Private as all other modules should interact only with get/set session value
// methods
func getSession(r *http.Request) *sessions.Session {
	s, err := store.Get(r, "vidwell")
	if err != nil {
		s, _ = store.New(r, "vidwell")
	}
	return s
}
