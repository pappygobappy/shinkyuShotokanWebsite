package models

import (
	"time"

	"gorm.io/gorm"
)

type Event struct {
	gorm.Model
	Title            string
	Date             time.Time
	Location         string
	Address          string
	GoogleMapsIframe string
	PictureUrl       string
	Alt              string
	Description      string
}
