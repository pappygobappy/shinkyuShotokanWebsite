package auth

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"shinkyuShotokan/initializers"
	"shinkyuShotokan/models"
	"shinkyuShotokan/utils"

	"github.com/golang-jwt/jwt/v5"
	passwordValidator "github.com/wagslane/go-password-validator"
	"golang.org/x/crypto/bcrypt"
)

const minEntropyBits float64 = 60

type SignupInput struct {
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

type SignupResult struct {
	User      models.User
	Token     string
	ExpiresAt int // seconds until expiry
}

// Signup handles user registration with validation and JWT generation
func Signup(input SignupInput) (*SignupResult, *utils.AppError) {
	// 1. Validate required fields
	if input.FirstName == "" || input.LastName == "" || input.Email == "" {
		return nil, &utils.AppError{Code: 422, Message: "All fields are required"}
	}

	// 2. Validate email format (basic check)
	input.Email = strings.TrimSpace(input.Email)
	if !strings.Contains(input.Email, "@") || !strings.Contains(input.Email, ".") {
		return nil, &utils.AppError{Code: 422, Message: "Invalid email format"}
	}

	// 3. Validate password strength
	passValErr := passwordValidator.Validate(input.Password, minEntropyBits)
	if passValErr != nil {
		return nil, &utils.AppError{Code: 422, Message: "Password is not strong enough"}
	}

	// 4. Check if user already exists
	var existingUser models.User
	result := initializers.DB.Where("email = ?", input.Email).First(&existingUser)
	if result.Error == nil && existingUser.ID > 0 {
		return nil, &utils.AppError{Code: 409, Message: "Email already registered"}
	}

	// 5. Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), 10)
	if err != nil {
		log.Printf("Failed to hash password for email %s: %v", input.Email, err)
		return nil, &utils.AppError{Code: 500, Message: "Failed to create account"}
	}

	// 6. Create user
	user := models.User{
		FirstName:    input.FirstName,
		LastName:     input.LastName,
		Email:        input.Email,
		PasswordHash: string(hash),
		Type:         models.AdminUser, // Default to admin for now (first user becomes owner in syncDb)
	}

	if err := initializers.DB.Create(&user).Error; err != nil {
		log.Printf("Failed to create user %s: %v", input.Email, err)
		return nil, &utils.AppError{Code: 500, Message: "Failed to create account"}
	}

	// 7. Generate JWT token
	tokenString, err := generateToken(user.ID)
	if err != nil {
		log.Printf("Failed to generate token for user %d: %v", user.ID, err)
		return nil, &utils.AppError{Code: 500, Message: "Failed to create account"}
	}

	return &SignupResult{
		User:      user,
		Token:     tokenString,
		ExpiresAt: 3600 * 24 * 30, // 30 days in seconds
	}, nil
}

// generateToken creates a JWT token for the given user ID
func generateToken(userID uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	hmacSecret := os.Getenv("HMAC_SECRET")
	if hmacSecret == "" {
		return "", fmt.Errorf("HMAC_SECRET environment variable not set")
	}

	return token.SignedString([]byte(hmacSecret))
}
