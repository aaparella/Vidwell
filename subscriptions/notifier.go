package service

import (
	"log"

	"github.com/aaparella/vidwell/models"
	"github.com/aaparella/vidwell/storage"
)

type Notifier interface {
	Notify(*models.User)
}

func GetSubscribers(creator uint) ([]*models.User, error) {
	var subscriptions []models.Subscription
	if err := storage.DB.Where(&models.Subscription{Creator: creator.ID}).Find(&subscriptions).Error; err != nil {
		log.Println("Could not fetch subscriptions : ", storage.DB.Error)
	}

}

// NotifySubscribers notifies all those who are subscribed to
// the creator of a video by sending them a personalized email
// with some basic information about the video.
func NotifySubscribers(creator *models.User, video *models.Video, notifier service.Notifier) {
	if subscribers, err := GetSubscribers(creator.ID); err == nil {
		for _, subscription := range subscriptions {
			notifier.Notify(subscription.Subscriber)
		}
	} else {
		log.Println("Could not fetch subscribers : ", err)
	}
}
