package queries

import (
	"log"
	"shinkyuShotokan/initializers"
	"shinkyuShotokan/models"
)

func GetClasses() []models.Class {
	var classes []models.Class
	result := initializers.DB.Order("display_order ASC").Find(&classes)
	if result.Error != nil {
		log.Print(result.Error)
	}
	return classes
}

func FindClassByID(id string) models.Class {
	var class models.Class
	result := initializers.DB.Preload("Location").Preload("Annotations").First(&class, id)
	if result.Error != nil {
		log.Print(result.Error)
	}
	return class
}

func FindClassByPath(path string) models.Class {
	var class models.Class
	result := initializers.DB.Preload("Location").Preload("Annotations").Where("get_url = ?", path).First(&class)
	if result.Error != nil {
		log.Print(result.Error)
	}
	return class
}

func FindClassByName(name string) models.Class {
	var class models.Class
	result := initializers.DB.Preload("Location").Preload("Annotations").Where("name = ?", name).First(&class)
	if result.Error != nil {
		log.Print(result.Error)
	}
	return class
}
