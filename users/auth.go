package users

import (
	"fmt"
	"net/http"

	"github.com/aaparella/vidwell/models"
	"github.com/aaparella/vidwell/session"
)

// MustBeLoggedIn wraps another http.HandlerFunc that requires a user to be
// logged in in order to access that endpoint (e.g video uploading, and
// user editing)
func MustBeLoggedIn(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if user := GetUser(r); user != nil {
			f(w, r)
		} else {
			fmt.Fprintf(w, "You must be logged in for this!")
		}
	}
}

// LoginUser writes the passed user to the session for the passed request,
// and saves the session.
func LoginUser(w http.ResponseWriter, r *http.Request, user *models.User) {
	sess := session.GetSession(r)
	sess.Values["user"] = user
	sess.Save(r, w)
}

// Checks that an email is valid. Clearly needs work.
func ValidEmail(email string) bool {
	return true
}

// Checks that a password is valid. Right now only enforces a length
// requirement.
func ValidPassword(pass string) bool {
	return len(pass) > 6
}
