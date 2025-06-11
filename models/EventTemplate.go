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

type EventType struct {
	Name      string
	Templates []FormattedTemplate
}

type FormattedTemplate struct {
	Type        string         `json:"Type"`
	Name        string         `json:"Name"`
	StartTime   string         `json:"StartTime"`
	EndTime     string         `json:"EndTime"`
	CheckInTime string         `json:"CheckInTime"`
	Location    string         `json:"Location"`
	SubTypes    []EventSubType `json:"SubTypes"`
}

// GetSubTypes returns all unique EventSubTypes from all templates
func (et *EventType) GetSubTypes() []EventSubType {
	subTypeMap := make(map[uint]EventSubType)

	for _, template := range et.Templates {
		for _, subType := range template.SubTypes {
			subTypeMap[subType.ID] = subType
		}
	}

	subTypes := make([]EventSubType, 0, len(subTypeMap))
	for _, subType := range subTypeMap {
		subTypes = append(subTypes, subType)
	}

	return subTypes
}
