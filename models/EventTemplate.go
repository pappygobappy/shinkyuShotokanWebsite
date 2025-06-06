package models

import (
	"time"

	"gorm.io/gorm"
)

type EventTemplate struct {
	gorm.Model
	Name        string
	StartTime   time.Time
	EndTime     time.Time
	CheckInTime time.Time
	Description string
	LocationID  string
}
