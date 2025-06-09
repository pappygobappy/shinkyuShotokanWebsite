package models

import (
	"gorm.io/gorm"
)

type EventSubType struct {
	gorm.Model
	Name           string
	EventTemplates []EventTemplate `gorm:"many2many:event_template_subtypes;"`
}

// TableName specifies the table name for EventSubType
func (EventSubType) TableName() string {
	return "event_sub_types"
}
