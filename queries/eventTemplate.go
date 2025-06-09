package queries

import (
	"shinkyuShotokan/initializers"
	"shinkyuShotokan/models"
)

func GetEventTemplates() []models.EventTemplate {
	var templates []models.EventTemplate
	initializers.DB.Preload("EventSubTypes").Find(&templates)
	return templates
}

func GetEventTemplateByID(id string) models.EventTemplate {
	var template models.EventTemplate
	initializers.DB.Preload("EventSubTypes").First(&template, id)
	return template
}
