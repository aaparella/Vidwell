package controllers

import (
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/aaparella/vidwell/models"
	"github.com/aaparella/vidwell/render"
	"github.com/aaparella/vidwell/session"
	"github.com/aaparella/vidwell/storage"
	"github.com/aaparella/vidwell/users"
	"github.com/gorilla/mux"
)

type UserController struct {
}

func (uc UserController) Prefix() string {
	return "/users"
}

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

func (uc UserController) ViewUser(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var user models.User

	if err := storage.DB.Find(&user, id); err != nil {
		fmt.Fprintf(w, "Could not find user with ID: %s", id)
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
		users.LoginUser(w, r, u)
		fmt.Fprintf(w, "Thank you for registering!")
	}
}

// Login logs in a user by verifying the passed credentials, and setting
// the user value in the request session if they are successful.
//
// Performs the bcrypt comparison even if we already know the user's email
// does not exist so that timing attacks cannot be used to identify user
// emails.
func (uc UserController) Login(w http.ResponseWriter, r *http.Request) {
	if user := users.GetUser(r); user != nil {
		fmt.Fprintf(w, "You are already logged in!")
		return
	}

	email, pass := r.FormValue("email"), r.FormValue("password")
	if !users.ValidEmail(email) || !users.ValidPassword(pass) {
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
