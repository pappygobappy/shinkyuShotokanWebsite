package models

import (
	"time"

	"gorm.io/gorm"
)

type Event struct {
	gorm.Model
	Title       string
	Date        time.Time
	Location    string
	PictureUrl  string
	Alt         string
	Description string
}
