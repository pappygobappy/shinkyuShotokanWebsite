package queries

import (
	"fmt"
	"shinkyuShotokan/initializers"
	"shinkyuShotokan/models"
)

func GetEventTemplates() []models.EventTemplate {
	var templates []models.EventTemplate
	initializers.DB.Preload("EventSubTypes").Find(&templates)
	return templates
}

func GetEventTypes() []models.EventType {
	templates := GetEventTemplates()
	typeMap := make(map[string]models.EventType)

	for _, template := range templates {
		// Create formatted template
		formattedTemplate := models.FormattedTemplate{
			Type:        template.Name,
			Name:        template.EventSubTypes[0].Name, // Assuming first subtype is the main one
			StartTime:   template.StartTime.Format("15:04:05"),
			EndTime:     template.EndTime.Format("15:04"),
			CheckInTime: template.CheckInTime.Format("15:04"),
			Location:    template.LocationID,
			SubTypes:    template.EventSubTypes,
		}

		if eventType, exists := typeMap[template.Name]; exists {
			eventType.Templates = append(eventType.Templates, formattedTemplate)
			typeMap[template.Name] = eventType
		} else {
			typeMap[template.Name] = models.EventType{
				Name:      template.Name,
				Templates: []models.FormattedTemplate{formattedTemplate},
			}
		}
	}

	eventTypes := make([]models.EventType, 0, len(typeMap))
	for _, eventType := range typeMap {
		eventTypes = append(eventTypes, eventType)
	}
	fmt.Println(eventTypes)
	return eventTypes
}

func GetEventSubTypes() []models.EventSubType {
	var subTypes []models.EventSubType
	initializers.DB.Find(&subTypes)
	return subTypes
}

func GetEventSubTypesByIDs(ids []uint) []models.EventSubType {
	var subTypes []models.EventSubType
	initializers.DB.Where("id IN ?", ids).Find(&subTypes)
	return subTypes
}

func GetEventTemplateByID(id string) models.EventTemplate {
	var template models.EventTemplate
	initializers.DB.Preload("EventSubTypes").First(&template, id)
	return template
}

func GetEventTemplatesByName(name string) []models.EventTemplate {
	var templates []models.EventTemplate
	initializers.DB.Preload("EventSubTypes").Where("name LIKE ?", "%"+name+"%").Find(&templates)
	return templates
}

func GetEventTemplatesByNameAndSubType(name string, subType string) models.EventTemplate {
	var template models.EventTemplate
	initializers.DB.Preload("EventSubTypes").
		Joins("JOIN event_template_subtypes ON event_template_subtypes.event_template_id = event_templates.id").
		Joins("JOIN event_sub_types ON event_sub_types.id = event_template_subtypes.event_sub_type_id").
		Where("event_templates.name = ? AND event_sub_types.name = ?", name, subType).
		First(&template)
	return template
}
