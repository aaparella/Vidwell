package models

import "github.com/jinzhu/gorm"

type Video struct {
	Title       string
	Uuid        string
	ContentType string

	UserID uint

	gorm.Model
}
