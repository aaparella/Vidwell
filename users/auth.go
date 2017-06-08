package users

import (
	"fmt"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/aaparella/vidwell/models"
	"github.com/aaparella/vidwell/session"
	"github.com/aaparella/vidwell/storage"
)

// MustBeLoggedIn wraps another http.HandlerFunc that requires a user to be
// logged in in order to access that endpoint (e.g video uploading, and
// user editing)
func MustBeLoggedIn(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if user := GetLoggedInUser(r); user != nil {
			f(w, r)
		} else {
			fmt.Fprintf(w, "You must be logged in for this!")
		}
	}
}

// CheckLoginInformation gets the user that the given email and
// password match for. Returns nil if they do not match for any
// users.
func CheckLoginInformation(email, password string) *models.User {
	u := &models.User{}
	if err := storage.DB.Where(&models.User{Email: email}).First(u); err != nil {
		u = &models.User{Password: []byte{}}
	}
	if err := bcrypt.CompareHashAndPassword(u.Password, []byte(password)); err != nil {
		return nil
	}
	return u
}

// LoginUser writes the passed user to the session for the passed request,
// and saves the session.
func LoginUser(w http.ResponseWriter, r *http.Request, user *models.User) {
	if err := session.StoreSessionValue("user", user, r, w); err != nil {
		//TODO handle this properly
		log.Println("Could not log in user : ", err)
	}
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
