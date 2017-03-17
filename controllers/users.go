package controllers

import (
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/aaparella/vidwell/models"
	"github.com/aaparella/vidwell/render"
	"github.com/aaparella/vidwell/storage"
	"github.com/aaparella/vidwell/users"
	"github.com/gorilla/mux"
)

// UserController consolidates behavior regarding the viewing and creation
// of new user accounts, as well as logging in.
type UserController struct {
}

func (uc UserController) Prefix() string {
	return "/users"
}

// Endpoints defines the endpoints for user operations. Makes heavy use of
// methods defined in the User module for creating / managing user data.
//
// ViewUser displays a user's profile page.
// Login allows a user to log in to their account.
// NewUser allows a user to sign up and create a new account.
func (uc UserController) Endpoints() map[string]map[string]http.HandlerFunc {
	return map[string]map[string]http.HandlerFunc{
		"/{id}": {
			"GET": uc.ViewUser,
		},
		"/login": {
			"POST": uc.Login,
		},
		"/new": {
			"POST": uc.NewUser,
		},
	}
}

// ViewUser displays the user account page for the specified user.
// The user id is retrieved from the url path.
func (uc UserController) ViewUser(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var user models.User

	if err := storage.DB.Find(&user, id).Error; err != nil {
		http.Error(w, "Could not find user with that ID", http.StatusNotFound)
		return
	}

	render.Render(w, r, "user", user)
}

// NewUser creates a new user from the given form values. Creates database
// record and automatically logs in the user.
func (uc UserController) NewUser(w http.ResponseWriter, r *http.Request) {
	u := &models.User{
		Email:       r.FormValue("email"),
		Password:    []byte(r.FormValue("password")),
		AccountName: r.FormValue("username"),
	}

	if !users.ValidEmail(u.Email) || !users.ValidPassword(string(u.Password)) {
		http.Error(w, "Username, password, or email invalid", http.StatusBadRequest)
		return
	}

	u.Password, _ = bcrypt.GenerateFromPassword(u.Password, 0)
	if err := storage.DB.Create(u).Error; err != nil {
		log.Println(w, "Could not create new user: %s", err.Error())
		http.Error(w, "Could not create new user", http.StatusInternalServerError)
	} else {
		users.LoginUser(w, r, u)
		http.Redirect(w, r, "/users/", http.StatusCreated)
	}
}

// Login logs in a user by verifying the passed credentials, and setting
// the user value in the request session if they are successful.
//
// Performs the bcrypt comparison even if we already know the user's email
// does not exist so that timing attacks cannot be used to identify user
// emails.
func (uc UserController) Login(w http.ResponseWriter, r *http.Request) {
	if user := users.GetLoggedInUser(r); user != nil {
		http.Redirect(w, r, "/", http.StatusOK)
		return
	}

	email, pass := r.FormValue("email"), r.FormValue("password")
	user := users.CheckLoginInformation(email, pass)

	if user == nil {
		http.Error(w, "Incorrect username and password combination", http.StatusUnauthorized)
	} else {
		users.LoginUser(w, r, user)
		http.Redirect(w, r, "/", http.StatusOK)
	}
}
