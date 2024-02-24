package initializers

import "shinkyuShotokan/models"

func SyncDb() {
	DB.AutoMigrate(&models.User{}, &models.Event{}, &models.ClassSession{}, &models.ClassPeriod{})
}
