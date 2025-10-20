package queries

import (
	"fmt"
	"log"
	"shinkyuShotokan/initializers"
	"shinkyuShotokan/models"
	"shinkyuShotokan/services"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetUsers() []models.User {
	var users []models.User
	result := initializers.DB.Find(&users)
	if result.Error != nil {
		log.Print(result.Error)
	}
	return users
}

func GetUserById(id uint) models.User {
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

func CreatePasswordResetToken(user models.User, c *fiber.Ctx) models.PasswordResetToken {
	token := uuid.NewString()
	initializers.DB.Where("user_id = ?", user.ID).Delete(&models.PasswordResetToken{})
	reset := models.PasswordResetToken{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(time.Hour),
	}
	if err := initializers.DB.Create(&reset).Error; err != nil {
		log.Println("failed to create password reset token", err)
	} else {
		resetURL := fmt.Sprintf("%s/reset-password/%s", c.BaseURL(), token)
		if err := services.SendPasswordResetEmail(user.Email, resetURL); err != nil {
			log.Println("failed to send password reset email", err)
		}
	}
	return reset
}

func GetPasswordResetToken(token string) models.PasswordResetToken {
	var reset models.PasswordResetToken
	initializers.DB.Where("token = ?", token).First(&reset)
	return reset
}
