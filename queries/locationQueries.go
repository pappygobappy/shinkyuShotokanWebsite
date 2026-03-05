package queries

import (
	"log"
	"shinkyuShotokan/initializers"
	"shinkyuShotokan/models"
	"strconv"
)

func GetLocations() []models.Location {
	if cached, ok := Cache.Get("locations"); ok {
		return cached.([]models.Location)
	}

	var locations []models.Location
	result := initializers.DB.Find(&locations)
	if result.Error != nil {
		log.Print(result.Error)
	}

	Cache.Set("locations", locations)

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
	locationID, _ := strconv.ParseUint(id, 10, 32)
	var location models.Location
	result := initializers.DB.Where("id = ?", locationID).First(&location)
	if result.Error != nil {
		log.Print(result.Error)
	}
	return location
}
