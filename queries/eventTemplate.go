package queries

import (
	"shinkyuShotokan/initializers"
	"shinkyuShotokan/models"
)

func GetEventTemplates() []models.EventTemplate {
	var templates []models.EventTemplate
	initializers.DB.Preload("Location").Find(&templates)
	return templates
}

func GetEventTemplateByID(id string) models.EventTemplate {
	var template models.EventTemplate
	initializers.DB.Preload("Location").First(&template, id)
	return template
}
