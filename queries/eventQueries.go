package queries

import (
	"log"
	"shinkyuShotokan/initializers"
	"shinkyuShotokan/models"
)

func GetUpcomingEvents() []models.Event {
	var events []models.Event
	result := initializers.DB.Where("date >= date(now())").Order("date").Find(&events)
	if result.Error != nil {
		log.Print(result.Error)
	}
	return events
}

func GetPastEventsForTheYear() []models.Event {
	var events []models.Event
	result := initializers.DB.Where("date < date(now()) AND date > now() - interval '1 year'").Order("date desc").Find(&events)
	if result.Error != nil {
		log.Print(result.Error)
	}
	return events
}
