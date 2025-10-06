package queries

import (
	"log"
	"shinkyuShotokan/initializers"
	"shinkyuShotokan/models"
)

func GetUsers() []models.User {
	var users []models.User
	result := initializers.DB.Find(&users)
	if result.Error != nil {
		log.Print(result.Error)
	}
	return users
}

func GetUserById(id string) models.User {
	var user models.User
	result := initializers.DB.Where("id = ?", id).First(&user)
	if result.Error != nil {
		log.Print(result.Error)
	}
	return user
}

func GetUserByEmail(email string) models.User {
	var user models.User
	result := initializers.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		log.Print(result.Error)
	}
	return user
}
