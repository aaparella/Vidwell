package models

import "github.com/jinzhu/gorm"

// User encapsulates information about a user account.
type User struct {
	AccountName   string
	DisplayName   string
	Email         string
	Administrator bool

	// The videos that this user has uploaded
	Uploads []Video

	// ID, CreatedAt, UpdatedAt, DeletedAt
	gorm.Model
}
