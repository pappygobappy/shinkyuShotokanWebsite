package handlers

import (
	"fmt"
	"log"
	"os"
	"time"

	"shinkyuShotokan/initializers"
	"shinkyuShotokan/models"
	"shinkyuShotokan/queries"
	"shinkyuShotokan/services/auth"
	"shinkyuShotokan/structs"
	"shinkyuShotokan/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	passwordValidator "github.com/wagslane/go-password-validator"
	"golang.org/x/crypto/bcrypt"
)

// SignupGet renders the signup page
func SignupGet(c *fiber.Ctx) error {
	if c.Locals("user") != nil {
		return c.Redirect("/")
	}
	return c.Render("signup", fiber.Map{})
}

// SignupPost handles user registration
func SignupPost(c *fiber.Ctx) error {
	var input auth.SignupInput

	if err := c.BodyParser(&input); err != nil {
		log.Printf("Failed to parse signup request: %v", err)
		return handleHTMXError(c, 422, "Invalid request body")
	}

	result, appErr := auth.Signup(input)
	if appErr != nil {
		return handleHTMXError(c, appErr.Code, appErr.Message)
	}

	setAuthCookie(c, result.Token, result.ExpiresAt)

	hxRequest := c.Locals("hxRequest") == true
	if hxRequest {
		return c.Redirect("/")
	}

	return c.Redirect("/")
}

// LoginGet renders the login page
func LoginGet(c *fiber.Ctx) error {
	if c.Locals("user") != nil {
		return c.Redirect("/")
	}

	resetMessage := ""
	if c.Query("reset") == "1" {
		resetMessage = "Your password has been reset. Please log in."
	}

	return c.Render("login", fiber.Map{
		"success": resetMessage,
	})
}

// LoginPost handles user authentication
func LoginPost(c *fiber.Ctx) error {
	var input auth.LoginInput

	if err := c.BodyParser(&input); err != nil {
		log.Printf("Failed to parse login request: %v", err)
		return handleHTMXError(c, 422, "Invalid request body")
	}

	result, appErr := auth.Login(input)
	if appErr != nil {
		return handleHTMXError(c, appErr.Code, appErr.Message)
	}

	setAuthCookie(c, result.Token, result.ExpiresAt)

	hxRequest := c.Locals("hxRequest") == true
	if hxRequest {
		return c.Redirect("/")
	}

	return c.Redirect("/")
}

// ForgotPasswordGet renders the forgot password page
func ForgotPasswordGet(c *fiber.Ctx) error {
	return c.Render("forgot_password", fiber.Map{})
}

// ForgotPasswordPost handles password reset email sending
func ForgotPasswordPost(c *fiber.Ctx) error {
	var input auth.ForgotPasswordInput

	if err := c.BodyParser(&input); err != nil {
		log.Printf("Failed to parse forgot password request: %v", err)
		return handleHTMXError(c, 422, "Unable to process your request. Please try again.")
	}

	user, appErr := auth.ForgotPassword(input)
	if appErr != nil {
		return handleHTMXError(c, appErr.Code, appErr.Message)
	}

	// Only create token and send email if user exists (prevent email enumeration)
	if user.ID != 0 {
		queries.CreatePasswordResetToken(*user, c)
		log.Printf("Password reset email sent for: %s", user.Email)
	}

	return c.Render("forgot_password", fiber.Map{
		"success": "If an account exists for that email, a password reset link has been sent.",
	})
}

// ResetPasswordTokenGet renders the reset password form
func ResetPasswordTokenGet(c *fiber.Ctx) error {
	token := c.Params("token")
	if token == "" {
		return handleHTMXError(c, 422, "The password reset link is invalid or has expired. Please request a new one.")
	}

	reset := queries.GetPasswordResetToken(token)
	if reset.ID == 0 || reset.UsedAt != nil || reset.ExpiresAt.Before(time.Now()) {
		return handleHTMXError(c, 422, "The password reset link is invalid or has expired. Please request a new one.")
	}

	return c.Render("reset_password_token", fiber.Map{
		"token": token,
	})
}

// ResetPasswordTokenPost handles password reset with token validation
func ResetPasswordTokenPost(c *fiber.Ctx) error {
	token := c.Params("token")
	if token == "" {
		return handleHTMXError(c, 422, "The password reset link is invalid or has expired. Please request a new one.")
	}

	var body struct {
		NewPassword        string `json:"new_password"`
		ConfirmNewPassword string `json:"confirm_new_password"`
	}

	if err := c.BodyParser(&body); err != nil {
		log.Printf("Failed to parse password reset request: %v", err)
		return handleHTMXError(c, 422, "Unable to process your request. Please try again.")
	}

	// Validate token exists and is valid
	reset := queries.GetPasswordResetToken(token)
	if reset.ID == 0 || reset.UsedAt != nil || reset.ExpiresAt.Before(time.Now()) {
		return c.Render("reset_password_token", fiber.Map{
			"token": token,
			"error": "The password reset link is invalid or has expired.",
		})
	}

	// Validate passwords match
	if body.NewPassword != body.ConfirmNewPassword {
		return c.Render("reset_password_token", fiber.Map{
			"token": token,
			"error": "Passwords don't match",
		})
	}

	// Validate password strength (reuse existing validation)
	minEntropyBits := 60.0
	passValErr := passwordValidator.Validate(body.NewPassword, minEntropyBits)
	if passValErr != nil {
		return c.Render("reset_password_token", fiber.Map{
			"token": token,
			"error": "Password is not strong enough.",
		})
	}

	// Get user and hash new password
	user := queries.GetUserById(reset.UserID)
	if user.ID == 0 {
		return c.Render("reset_password_token", fiber.Map{
			"token": token,
			"error": "User not found.",
		})
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.NewPassword), 10)
	if err != nil {
		log.Printf("Failed to hash password for user %d: %v", user.ID, err)
		return c.Render("reset_password_token", fiber.Map{
			"token": token,
			"error": "Unable to reset password at this time.",
		})
	}

	// Update user password
	user.PasswordHash = string(hash)
	if err := initializers.DB.Save(&user).Error; err != nil {
		log.Printf("Failed to save user %d during password reset: %v", user.ID, err)
		return c.Render("reset_password_token", fiber.Map{
			"token": token,
			"error": "Unable to reset password at this time.",
		})
	}

	// Mark token as used and clean up old tokens
	now := time.Now()
	reset.UsedAt = &now
	if err := initializers.DB.Save(&reset).Error; err != nil {
		log.Printf("Failed to mark reset token used for user %d: %v", user.ID, err)
	}

	initializers.DB.Where("user_id = ? AND token <> ?", user.ID, token).Delete(&models.PasswordResetToken{})

	log.Printf("Password successfully reset for user %d", user.ID)
	return c.Redirect("/login?reset=1")
}

// LogoutPost handles user logout by clearing the auth cookie
func LogoutPost(c *fiber.Ctx) error {
	u := c.Locals("user")
	user, ok := u.(models.User)

	if !ok || user.ID == 0 {
		log.Println("Invalid user session during logout")
		return c.Redirect("/")
	}

	expiredToken, _ := generateExpiredCookie(user.ID)
	setAuthCookie(c, expiredToken, -1) // -1 = immediately expire

	c.Set("HX-Redirect", "/")
	return c.Next()
}

// AdminUserProfilePage renders the admin user profile page
func AdminUserProfilePage(c *fiber.Ctx) error {
	return c.Render("adminPage", fiber.Map{
		"Page": structs.Page{PageName: "User Profile", Tabs: utils.CurrentTabs(), Classes: utils.Classes},
		"user": c.Locals("user"),
	})
}

// GetUserProfilePageEdit renders the user profile edit page
func GetUserProfilePageEdit(c *fiber.Ctx) error {
	return c.Render("userProfile_edit", fiber.Map{})
}

// Helper functions

func setAuthCookie(c *fiber.Ctx, token string, maxAge int) {
	c.Cookie(&fiber.Cookie{
		Name:     "Authorization",
		Value:    token,
		Path:     "/",
		MaxAge:   maxAge,
		HTTPOnly: true,
	})
}

func handleHTMXError(c *fiber.Ctx, code int, message string) error {
	hxRequest := c.Locals("hxRequest") == true

	if hxRequest {
		return c.Status(code).Render("partials/error", fiber.Map{
			"error": message,
		})
	}

	return c.Status(code).JSON(fiber.Map{
		"error": message,
	})
}

func generateExpiredCookie(userID uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(-(time.Hour * 24 * 30)).Unix(),
		"sub": userID,
	})

	hmacSecret := os.Getenv("HMAC_SECRET")
	if hmacSecret == "" {
		return "", fmt.Errorf("HMAC_SECRET environment variable not set")
	}

	return token.SignedString([]byte(hmacSecret))
}
