package queries

import (
	"log"
	"shinkyuShotokan/initializers"
	"shinkyuShotokan/models"
)

func GetLocations() []models.Location {
	var locations []models.Location
	result := initializers.DB.Find(&locations)
	if result.Error != nil {
		log.Print(result.Error)
	}
	return locations
}

func GetLocationByName(name string) models.Location {
	var location models.Location
	result := initializers.DB.Where("name = ?", name).First(&location)
	if result.Error != nil {
		log.Print(result.Error)
	}
	return location
}

func GetLocationById(id string) models.Location {
	var location models.Location
	result := initializers.DB.Where("id = ?", id).First(&location)
	if result.Error != nil {
		log.Print(result.Error)
	}
	return location
}
