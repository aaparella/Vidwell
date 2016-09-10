package users

import (
	"fmt"
	"net/http"

	"github.com/aaparella/vidwell/models"
	"github.com/aaparella/vidwell/session"
	"github.com/aaparella/vidwell/storage"
	"golang.org/x/crypto/bcrypt"
)

// Login logs in a user by verifying the passed credentials, and setting
// the user value in the request session if they are successful.
//
// Performs the bcrypt comparison even if we already know the user's email
// does not exist so that timing attacks cannot be used to identify user
// emails.
func Login(w http.ResponseWriter, r *http.Request) {
	if user := GetUser(r); user != nil {
		fmt.Fprintf(w, "You are already logged in!")
		return
	}

	email, pass := r.FormValue("email"), r.FormValue("password")
	if !validEmail(email) || !validPassword(pass) {
		fmt.Fprintf(w, "Please enter a valid email and password")
		return
	}

	u := &models.User{}
	storage.DB.Where(&models.User{Email: email}).First(u)
	if u.Email == "" {
		u.Password = []byte("   ")
	}

	if err := bcrypt.CompareHashAndPassword(u.Password, []byte(pass)); err != nil {
		fmt.Fprintf(w, "Incorrect username / password combination")
	} else {
		sess := session.GetSession(r)
		sess.Values["user"] = u
		sess.Save(r, w)
		fmt.Fprintf(w, "Welcome back")
	}
}

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

// Checks that an email is valid. Clearly needs work.
func validEmail(email string) bool {
	return true
}

// Checks that a password is valid. Right now only enforces a length
// requirement.
func validPassword(pass string) bool {
	return len(pass) > 6
}
