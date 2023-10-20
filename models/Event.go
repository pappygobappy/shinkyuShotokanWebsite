package models

import (
	"html/template"
	"time"

	"gorm.io/gorm"
)

type Event struct {
	gorm.Model
	Title       string
	Date        time.Time
	StartTime   time.Time
	EndTime     time.Time
	Location    string
	PictureUrl  string
	Alt         string
	Description template.HTML
}
