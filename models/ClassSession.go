package models

import (
	"time"

	"gorm.io/gorm"
)

type ClassSession struct {
	gorm.Model
	ClassName   string
	Period      string
	StartTime   time.Time
	EndTime     time.Time
	Location    string
	IsCancelled bool
}
