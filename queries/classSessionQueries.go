package queries

import (
	"log"
	"shinkyuShotokan/initializers"
	"shinkyuShotokan/models"
	"time"
)

func GetClassSessionsBetweenDates(startDate time.Time, endDate time.Time) []models.ClassSession {
	endDate = endDate.AddDate(0,0,1)
	var classSessions []models.ClassSession
	result := initializers.DB.Where("start_time >= ? AND start_time <= ?", startDate, endDate).Order("start_time desc").Find(&classSessions)
	if result.Error != nil {
		log.Print(result.Error)
	}
	return classSessions
}

func GetClassSessionsByClassAndBetweenDates(class string, startDate time.Time, endDate time.Time) []models.ClassSession {
	var classSessions []models.ClassSession
	endDate = endDate.AddDate(0,0,1)
	result := initializers.DB.Where("class_name = ? AND start_time >= ? AND start_time <= ?", class, startDate, endDate).Order("start_time desc").Find(&classSessions)
	if result.Error != nil {
		log.Print(result.Error)
	}
	return classSessions
}

func DeleteClassSessionsByClassAndClassPeriod(class string, classPeriod string) {
	result := initializers.DB.Delete(&models.ClassSession{}, "class_name = ? AND period = ?", class, classPeriod)
	if result.Error != nil {
		log.Print(result.Error)
	}
}
