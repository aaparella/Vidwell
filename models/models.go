package models

import "github.com/jinzhu/gorm"

// User encapsulates information about a user account.
type User struct {
	// AccountName used to log into Vidwell
	AccountName string
	// DisplayName shown on comments and video pages
	DisplayName string
	// Email used to register account, must be unique.
	Email string `gorm:"unique"`
	// Administrator privilege indicator
	Administrator bool

	// Password is the bcrypt hashed password for this user.
	Password []byte

	// The videos that this user has uploaded
	Uploads []Video

	// ID, CreatedAt, UpdatedAt, DeletedAt
	gorm.Model
}

// Video encapsulates all data about a video. Actual videos
// are stored in object store, referenced by a video's UUID
type Video struct {
	// Title of the video
	Title string
	// Uuid that can be used to fetch the actual video content
	Uuid string
	// ContentType indicates what format the video is in
	ContentType string
	// Views tracks the number of views a video has accrued.
	Views uint
	// UserID of the user that created the video
	UserID uint
	gorm.Model
}

// Subscription represents the subscription of one account
// to the contents of the other
type Subscription struct {
	// Subscriber is the user that IS subscribed
	Subscriber uint
	// Creator is the user that is susbcribed TO
	Creator uint

	gorm.Model
}
