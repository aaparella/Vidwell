package notifier

import "github.com/aaparella/vidwell/models"

// Notifier defines the interface for any struct or method
// that notifies a user of an event.
type Notifier interface {
	Notify(*models.User) error
}
