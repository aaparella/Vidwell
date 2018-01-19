package notifier

import (
	"github.com/aaparella/vidwell/models"
	"github.com/aaparella/vidwell/storage"
	"github.com/sirupsen/logrus"
)

// Notifier defines the interface for any struct or method
// that notifies a user of an event.
type Notifier interface {
	Notify(*models.User) error
}

// GetSubscribers returns a list of users that are subscribed to the
// user with the provided id.
func GetSubscribers(creator uint) ([]*models.User, error) {
	var subscriptions []models.Subscription
	if err := storage.DB.Where(&models.Subscription{Creator: creator}).Find(&subscriptions).Error; err != nil {
		logrus.Println("Could not fetch subscriptions : ", storage.DB.Error)
	}
	return nil, nil
}

// NotifySubscribers notifies all those who are subscribed to
// the creator of a video by sending them a personalized email
// with some basic information about the video.
func NotifySubscribers(creator *models.User, video *models.Video, n Notifier) {
	if subscribers, err := GetSubscribers(creator.ID); err == nil {
		for _, subscriber := range subscribers {
			n.Notify(subscriber)
		}
	} else {
		logrus.Println("Could not fetch subscribers : ", err)
	}
}
