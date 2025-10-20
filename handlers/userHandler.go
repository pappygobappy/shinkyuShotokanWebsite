package handlers

import (
	"log"
	"os"
	"strings"
	"time"

	"shinkyuShotokan/initializers"
	"shinkyuShotokan/models"
	"shinkyuShotokan/queries"
	"shinkyuShotokan/structs"
	"shinkyuShotokan/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	passwordValidator "github.com/wagslane/go-password-validator"
	"golang.org/x/crypto/bcrypt"
)

var minEntropyBits float64 = 60

func SignupGet(c *fiber.Ctx) error {
	return c.Render("signup", fiber.Map{})
}

func SignupPost(c *fiber.Ctx) error {
	//Get the Username/password
	var body struct {
		FirstName       string
		LastName        string
		Email           string
		Password        string
		ConfirmPassword string
	}

	if err := c.BodyParser(&body); err != nil {
		log.Panic(err)
		return err
	}

	if body.FirstName == "" || body.LastName == "" || body.Email == "" {
		return c.Render("signup", fiber.Map{
			"error": "All fields are required",
		})
	}

	//Verify password pattern
	passValErr := passwordValidator.Validate(body.Password, minEntropyBits)
	if passValErr != nil {
		return c.Render("signup", fiber.Map{
			"error": "Password is not strong enough",
		})
	}

	//Hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	if err != nil {
		log.Panic("Failed to hash password", err)
		return err
	}

	//Create the user
	user := models.User{FirstName: body.FirstName, LastName: body.LastName, Email: body.Email, PasswordHash: string(hash), Type: models.AdminUser}

	result := initializers.DB.Create(&user)

	if result.Error != nil {
		log.Panic("Error creating User", result.Error)
		return err
	}

	//Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	hmacSampleSecret := (os.Getenv("HMAC_SECRET"))

	tokenString, err := token.SignedString([]byte(hmacSampleSecret))

	if err != nil {
		log.Print(err)
		return c.Render("login", fiber.Map{})
	}

	//send it back
	//c.Append("token", tokenString)
	c.Cookie(&fiber.Cookie{
		Name:     "Authorization",
		Value:    tokenString,
		MaxAge:   3600 * 24 * 30,
		HTTPOnly: true,
	})

	//Respond

	return c.Redirect("/")
}

func LoginGet(c *fiber.Ctx) error {
	resetMessage := ""
	if c.Query("reset") == "1" {
		resetMessage = "Your password has been reset. Please log in."
	}

	return c.Render("login", fiber.Map{
		"success": resetMessage,
	})
}

func LoginPost(c *fiber.Ctx) error {
	//Get Email and Pass

	var body struct {
		Email    string
		Password string
	}

	if err := c.BodyParser(&body); err != nil {
		log.Panic(err)
		return err
	}

	//Look up requested user
	user := queries.GetUserByEmail(body.Email)

	if user.ID == 0 {
		log.Println("Invalid email or password")
		return c.Render("login", fiber.Map{
			"error": "Email or password is incorrect",
		})
	}

	//Compare password and saved hash
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(body.Password))

	if err != nil {
		log.Println("Invalid email or password")
		return c.Render("login", fiber.Map{
			"error": "Email or password is incorrect",
		})
	}

	//Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	hmacSampleSecret := (os.Getenv("HMAC_SECRET"))

	tokenString, err := token.SignedString([]byte(hmacSampleSecret))

	if err != nil {
		log.Print(err)
		return c.Render("login", fiber.Map{})
	}

	//send it back
	//c.Append("token", tokenString)
	c.Cookie(&fiber.Cookie{
		Name:     "Authorization",
		Value:    tokenString,
		MaxAge:   3600 * 24 * 30,
		HTTPOnly: true,
	})
	return c.Redirect("/")
}

func ForgotPasswordGet(c *fiber.Ctx) error {
	return c.Render("forgot_password", fiber.Map{})
}

func ForgotPasswordPost(c *fiber.Ctx) error {
	var body struct {
		Email string
	}

	if err := c.BodyParser(&body); err != nil {
		log.Println("failed to parse forgot password request", err)
		return c.Render("forgot_password", fiber.Map{
			"error": "Unable to process your request. Please try again.",
		})
	}

	email := strings.TrimSpace(body.Email)
	if email != "" {
		user := queries.GetUserByEmail(email)
		if user.ID != 0 {
			queries.CreatePasswordResetToken(user, c)
		}
	}

	return c.Render("forgot_password", fiber.Map{
		"success": "If an account exists for that email, a password reset link has been sent.",
	})
}

func ResetPasswordTokenGet(c *fiber.Ctx) error {
	token := c.Params("token")
	if token == "" {
		return c.Render("forgot_password", fiber.Map{
			"error": "The password reset link is invalid or has expired. Please request a new one.",
		})
	}

	reset := queries.GetPasswordResetToken(token)
	if reset.ID == 0 || reset.UsedAt != nil || reset.ExpiresAt.Before(time.Now()) {
		return c.Render("forgot_password", fiber.Map{
			"error": "The password reset link is invalid or has expired. Please request a new one.",
		})
	}

	return c.Render("reset_password_token", fiber.Map{
		"token": token,
	})
}

func ResetPasswordTokenPost(c *fiber.Ctx) error {
	token := c.Params("token")
	if token == "" {
		return c.Render("forgot_password", fiber.Map{
			"error": "The password reset link is invalid or has expired. Please request a new one.",
		})
	}

	var body struct {
		NewPassword        string
		ConfirmNewPassword string
	}

	if err := c.BodyParser(&body); err != nil {
		log.Println("failed to parse password reset token request", err)
		return c.Render("forgot_password", fiber.Map{
			"error": "Unable to process your request. Please try again.",
		})
	}

	if body.NewPassword != body.ConfirmNewPassword {
		return c.Render("reset_password_token", fiber.Map{
			"token": token,
			"error": "Passwords don't match",
		})
	}

	if err := passwordValidator.Validate(body.NewPassword, minEntropyBits); err != nil {
		return c.Render("reset_password_token", fiber.Map{
			"token": token,
			"error": "Password is not strong enough.",
		})
	}

	reset := queries.GetPasswordResetToken(token)
	if reset.ID == 0 || reset.UsedAt != nil || reset.ExpiresAt.Before(time.Now()) {
		return c.Render("forgot_password", fiber.Map{
			"error": "The password reset link is invalid or has expired. Please request a new one.",
		})
	}

	user := queries.GetUserById(reset.UserID)
	if user.ID == 0 {
		return c.Render("forgot_password", fiber.Map{
			"error": "The password reset link is invalid or has expired. Please request a new one.",
		})
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.NewPassword), 10)
	if err != nil {
		log.Println("failed to hash password during reset", err)
		return c.Render("reset_password_token", fiber.Map{
			"token": token,
			"error": "Unable to reset password at this time.",
		})
	}

	user.PasswordHash = string(hash)
	if err := initializers.DB.Save(&user).Error; err != nil {
		log.Println("failed to save user during password reset", err)
		return c.Render("reset_password_token", fiber.Map{
			"token": token,
			"error": "Unable to reset password at this time.",
		})
	}

	now := time.Now()
	reset.UsedAt = &now
	if err := initializers.DB.Save(&reset).Error; err != nil {
		log.Println("failed to mark reset token used", err)
	}

	initializers.DB.Where("user_id = ? AND token <> ?", user.ID, reset.Token).Delete(&models.PasswordResetToken{})

	return c.Redirect("/login?reset=1")
}

func ChangePasswordGet(c *fiber.Ctx) error {
	return c.Render("change_password", fiber.Map{})
}

func ChangePasswordPost(c *fiber.Ctx) error {

	u := c.Locals("user")
	user, _ := u.(models.User)

	// var user models.User
	// initializers.DB.First(&user, "email = ?", u)

	if user.ID == 0 {
		log.Println("Invalid email or password")
		return c.Render("change_password", fiber.Map{
			"error": "Email or password is incorrect",
		})
	}

	//Get Current and NewPass
	var body struct {
		CurrentPassword    string
		NewPassword        string
		ConfirmNewPassword string
	}

	if err := c.BodyParser(&body); err != nil {
		log.Panic(err)
		return err
	}

	//Verify passwords match
	if body.ConfirmNewPassword != body.NewPassword {
		return c.Render("change_password", fiber.Map{
			"error": "Passwords don't match",
		})
	}

	//Verify password pattern
	passValErr := passwordValidator.Validate(body.NewPassword, minEntropyBits)
	if passValErr != nil {
		return c.Render("change_password", fiber.Map{
			"error": "Password is not strong enough.",
		})
	}

	//Compare current password
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(body.CurrentPassword))

	if err != nil {
		log.Println("Invalid email or password")
		return c.Render("change_password", fiber.Map{
			"error": "Email or password is incorrect",
		})
	}

	//Hash the new password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.NewPassword), 10)

	if err != nil {
		log.Panic("Failed to hash password", err)
		return err
	}

	//Update the user
	user.PasswordHash = string(hash)
	initializers.DB.Save(&user)

	c.Locals("message", "Successfully reset Password!")
	c.Set("HX-Redirect", "/admin/userProfile")
	return c.Next()
}

func LogoutPost(c *fiber.Ctx) error {

	//Clear JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(-(time.Hour * 24 * 30)).Unix(),
	})

	hmacSampleSecret := (os.Getenv("HMAC_SECRET"))

	tokenString, err := token.SignedString([]byte(hmacSampleSecret))

	if err != nil {
		log.Print(err)
		return c.Render("login", fiber.Map{})
	}

	//send it back
	//c.Append("token", tokenString)
	c.Cookie(&fiber.Cookie{
		Name:     "Authorization",
		Value:    tokenString,
		MaxAge:   3600 * 24 * 30,
		HTTPOnly: true,
	})
	c.Set("HX-Redirect", "/")
	return c.Next()
}

func AdminUserProfilePage(c *fiber.Ctx) error {
	return c.Render("adminPage", fiber.Map{
		"Page": structs.Page{PageName: "User Profile", Tabs: utils.CurrentTabs(), Classes: utils.Classes},
	})
}

func GetUserProfilePageEdit(c *fiber.Ctx) error {
	return c.Render("userProfileEdit", fiber.Map{})
}

func PostUserProfilePageEdit(c *fiber.Ctx) error {
	log.Println("PostUserProfilePageEdit")
	user := c.Locals("user").(models.User)

	var body struct {
		FirstName string
		LastName  string
		Email     string
	}

	if err := c.BodyParser(&body); err != nil {
		log.Print(err)
		return err
	}

	user.FirstName = body.FirstName
	user.LastName = body.LastName
	user.Email = body.Email

	initializers.DB.Save(&user)

	c.Set("HX-Redirect", "/admin/userProfile")
	return c.Next()
}
