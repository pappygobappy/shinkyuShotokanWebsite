package models

import (
	"time"

	"gorm.io/gorm"
)

type EventTemplate struct {
	gorm.Model
	Name          string
	StartTime     time.Time
	EndTime       time.Time
	CheckInTime   time.Time
	Description   string
	LocationID    string
	EventSubTypes []EventSubType `gorm:"many2many:event_template_subtypes;"`
}

// TableName specifies the table name for EventTemplate
func (EventTemplate) TableName() string {
	return "event_templates"
}
