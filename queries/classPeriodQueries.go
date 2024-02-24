package queries

import (
	"log"
	"shinkyuShotokan/initializers"
	"shinkyuShotokan/models"
)

func GetClassPeriodByName(period string) models.ClassPeriod {
	var classPeriod models.ClassPeriod
	result := initializers.DB.Where("name = ?", period).First(&classPeriod)
	if result.Error != nil {
		log.Print(result.Error)
	}
	return classPeriod
}
