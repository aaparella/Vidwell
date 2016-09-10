package models

import "github.com/jinzhu/gorm"

// User encapsulates information about a user account.
type User struct {
	AccountName   string
	DisplayName   string
	Email         string `gorm:"unique"`
	Administrator bool

	// Password is the bcrypt hashed password for this user.
	Password []byte

	// The videos that this user has uploaded
	Uploads []Video

	// ID, CreatedAt, UpdatedAt, DeletedAt
	gorm.Model
}
