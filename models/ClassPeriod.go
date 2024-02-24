package models

import (
	"time"

	"gorm.io/gorm"
)

type ClassPeriod struct {
	gorm.Model
	Name      string
	StartDate time.Time
	EndDate   time.Time
}
