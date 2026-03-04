package auth

import (
	"log"
	"strings"

	"shinkyuShotokan/initializers"
	"shinkyuShotokan/models"
	"shinkyuShotokan/utils"

	"golang.org/x/crypto/bcrypt"
)

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResult struct {
	User      models.User
	Token     string
	ExpiresAt int // seconds until expiry
}

// Login handles user authentication with JWT token generation
func Login(input LoginInput) (*LoginResult, *utils.AppError) {
	email := strings.TrimSpace(input.Email)

	// 1. Look up user (don't leak whether email exists for security)
	var user models.User
	result := initializers.DB.Where("email = ?", email).First(&user)
	if result.Error != nil || user.ID == 0 {
		log.Printf("Login attempt with invalid credentials: %s", email)
		return nil, &utils.AppError{Code: 401, Message: "Email or password is incorrect"}
	}

	// 2. Compare password and saved hash
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password))
	if err != nil {
		log.Printf("Login attempt with invalid credentials: %s", email)
		return nil, &utils.AppError{Code: 401, Message: "Email or password is incorrect"}
	}

	// 3. Generate JWT token
	tokenString, err := generateToken(user.ID)
	if err != nil {
		log.Printf("Failed to generate token for user %d: %v", user.ID, err)
		return nil, &utils.AppError{Code: 500, Message: "Login failed"}
	}

	log.Printf("User %d (%s) logged in successfully", user.ID, user.Email)

	return &LoginResult{
		User:      user,
		Token:     tokenString,
		ExpiresAt: 3600 * 24 * 30, // 30 days in seconds
	}, nil
}
