package users

import (
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/aaparella/vidwell/models"
	"github.com/aaparella/vidwell/render"
	"github.com/aaparella/vidwell/session"
	"github.com/aaparella/vidwell/storage"
	"github.com/gorilla/mux"
)

func ViewUser(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var user models.User

	if err := storage.DB.Find(&user, id); err != nil {
		fmt.Fprintf(w, "Could not find user with ID: %s", id)
		return
	}

	render.Render(w, "user", user)
}

// NewUser creates a new user from the given form values. Creates database
// record and automatically logs in the user.
func NewUser(w http.ResponseWriter, r *http.Request) {
	u := &models.User{
		Email:       r.FormValue("email"),
		Password:    []byte(r.FormValue("password")),
		AccountName: r.FormValue("username"),
	}

	if !validEmail(u.Email) || !validPassword(string(u.Password)) {
		fmt.Fprintf(w, "Username, email or password invalid")
		return
	}

	password, err := bcrypt.GenerateFromPassword(u.Password, 0)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	u.Password = password

	if err := storage.DB.Create(u).Error; err != nil {
		fmt.Fprintf(w, "Could not create new user: %s", err.Error())
	} else {
		LoginUser(w, r, u)
		fmt.Fprintf(w, "Thank you for registering!")
	}
}

// GetUser is a convenience function that returns the user for the request's
// session, or nil if they are not yet logged in.
func GetUser(r *http.Request) *models.User {
	val := session.GetSession(r).Values["user"]
	if user, ok := val.(*models.User); !ok {
		return nil
	} else {
		return user
	}
}

// LoginUser writes the passed user to the session for the passed request,
// and saves the session.
func LoginUser(w http.ResponseWriter, r *http.Request, user *models.User) {
	sess := session.GetSession(r)
	sess.Values["user"] = user
	sess.Save(r, w)
}
