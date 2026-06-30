package models

import (
	"time"

	"gorm.io/gorm"
)

type Promotional struct {
	gorm.Model
	Name        string
	Date        time.Time
	StartTime   time.Time
	EndTime     time.Time
	CheckInTime time.Time
	SubType     string
	Location    string
	Description string
}

func (Promotional) TableName() string {
	return "promotionals"
}