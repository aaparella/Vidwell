package users

import (
	"net/http"

	"github.com/aaparella/vidwell/models"
	"github.com/aaparella/vidwell/session"
)

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
