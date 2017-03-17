package users

import (
	"net/http"

	"github.com/aaparella/vidwell/models"
	"github.com/aaparella/vidwell/session"
	"github.com/aaparella/vidwell/storage"
)

// GetUser fetches the user with the given id, and an error if there is one.
func GetUser(id uint) (models.User, error) {
	var user models.User
	storage.DB.Find(&user, id)
	return user, storage.DB.Error
}

// GetUsers returns a list of user models.
//
// If fetching a user returns an error, all users that come
// after them will be ignored, and the error that caused
// the problem is returned.
func GetUsers(ids ...uint) ([]models.User, error) {
	var err error
	users := make([]models.User, len(ids))
	get := func(index int, id uint) {
		if err != nil {
			return
		}
		users[index], err = GetUser(id)
	}
	for i, id := range ids {
		get(i, id)
	}
	return users, err
}

// GetUser is a convenience function that returns the user for the request's
// session, or nil if they are not yet logged in.
func GetLoggedInUser(r *http.Request) *models.User {
	val := session.GetSession(r).Values["user"]
	if user, ok := val.(*models.User); !ok {
		return nil
	} else {
		return user
	}
}
