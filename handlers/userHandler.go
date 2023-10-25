package handlers

import (
	"log"
	"os"
	"shinkyuShotokan/initializers"
	"shinkyuShotokan/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	passwordValidator "github.com/wagslane/go-password-validator"
)

var minEntropyBits float64 = 60

func SignupGet(c *fiber.Ctx) error {
	return c.Render("signup", fiber.Map{})
}

func SignupPost(c *fiber.Ctx) error {
	//Get the Username/password
	var body struct {
		Email           string
		Password        string
		ConfirmPassword string
	}

	if err := c.BodyParser(&body); err != nil {
		log.Panic(err)
		return err
	}

	//Verify passwords match
	if body.ConfirmPassword != body.Password {
		return c.Render("signup", fiber.Map{
			"error": "Passwords don't match",
		})
	}

	//Verify password pattern
	passValErr := passwordValidator.Validate(body.Password, minEntropyBits)
	if passValErr != nil {
		return c.Render("reset_password", fiber.Map{
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
	user := models.User{Email: body.Email, PasswordHash: string(hash)}

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
	return c.Render("login", fiber.Map{})
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
	var user models.User
	initializers.DB.First(&user, "email = ?", body.Email)

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

func ResetPasswordGet(c *fiber.Ctx) error {
	return c.Render("reset_password", fiber.Map{})
}

func ResetPasswordPost(c *fiber.Ctx) error {

	u := c.Locals("user")
	user, _ := u.(models.User)

	// var user models.User
	// initializers.DB.First(&user, "email = ?", u)

	if user.ID == 0 {
		log.Println("Invalid email or password")
		return c.Render("reset_password", fiber.Map{
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
		return c.Render("reset_password", fiber.Map{
			"error": "Passwords don't match",
		})
	}

	//Verify password pattern
	passValErr := passwordValidator.Validate(body.NewPassword, minEntropyBits)
	if passValErr != nil {
		return c.Render("reset_password", fiber.Map{
			"error": "Password is not strong enough.",
		})
	}

	//Compare current password
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(body.CurrentPassword))

	if err != nil {
		log.Println("Invalid email or password")
		return c.Render("reset_password", fiber.Map{
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
	return Home(c)
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
	return c.Redirect("/")
}
