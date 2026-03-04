package auth

import (
	"log"
	"strings"
	"time"

	"shinkyuShotokan/initializers"
	"shinkyuShotokan/models"
	"shinkyuShotokan/queries"
	"shinkyuShotokan/utils"

	passwordValidator "github.com/wagslane/go-password-validator"
	"golang.org/x/crypto/bcrypt"
)

const passwordResetMinEntropyBits float64 = 60

type ForgotPasswordInput struct {
	Email string `json:"email"`
}

type ResetPasswordInput struct {
	Token              string `json:"token"`
	NewPassword        string `json:"new_password"`
	ConfirmNewPassword string `json:"confirm_new_password"`
}

// ForgotPassword validates email and returns user if found (doesn't send email - handler does that)
func ForgotPassword(input ForgotPasswordInput) (*models.User, *utils.AppError) {
	email := strings.TrimSpace(input.Email)
	if email == "" {
		return nil, &utils.AppError{Code: 422, Message: "Email is required"}
	}

	user := queries.GetUserByEmail(email)

	// Always return success to prevent email enumeration attacks
	if user.ID != 0 {
		log.Printf("Password reset request for: %s", email)
	}

	return &user, nil
}

// ResetPassword handles password reset with token validation and new password setting
func ResetPassword(input ResetPasswordInput) (*models.User, *utils.AppError) {
	// 1. Validate token exists
	reset := queries.GetPasswordResetToken(input.Token)
	if reset.ID == 0 || reset.UsedAt != nil || reset.ExpiresAt.Before(time.Now()) {
		return nil, &utils.AppError{Code: 422, Message: "The password reset link is invalid or has expired"}
	}

	// 2. Validate passwords match
	if input.NewPassword != input.ConfirmNewPassword {
		return nil, &utils.AppError{Code: 422, Message: "Passwords don't match"}
	}

	// 3. Validate password strength
	passValErr := passwordValidator.Validate(input.NewPassword, passwordResetMinEntropyBits)
	if passValErr != nil {
		return nil, &utils.AppError{Code: 422, Message: "Password is not strong enough"}
	}

	// 4. Get user and hash new password
	user := queries.GetUserById(reset.UserID)
	if user.ID == 0 {
		return nil, &utils.AppError{Code: 404, Message: "User not found"}
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), 10)
	if err != nil {
		log.Printf("Failed to hash password for user %d: %v", user.ID, err)
		return nil, &utils.AppError{Code: 500, Message: "Unable to reset password at this time"}
	}

	// 5. Update user password
	user.PasswordHash = string(hash)
	if err := initializers.DB.Save(&user).Error; err != nil {
		log.Printf("Failed to save user %d during password reset: %v", user.ID, err)
		return nil, &utils.AppError{Code: 500, Message: "Unable to reset password at this time"}
	}

	// 6. Mark token as used and clean up old tokens
	now := time.Now()
	reset.UsedAt = &now
	if err := initializers.DB.Save(&reset).Error; err != nil {
		log.Printf("Failed to mark reset token used for user %d: %v", user.ID, err)
	}

	initializers.DB.Where("user_id = ? AND token <> ?", user.ID, input.Token).Delete(&models.PasswordResetToken{})

	log.Printf("Password successfully reset for user %d", user.ID)

	return &user, nil
}
