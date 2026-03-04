package handlers

import (
	"log"

	"shinkyuShotokan/initializers"
	"shinkyuShotokan/models"

	"github.com/gofiber/fiber/v2"
	passwordValidator "github.com/wagslane/go-password-validator"
	"golang.org/x/crypto/bcrypt"
)

const minEntropyBits float64 = 60

// DEPRECATED: All authentication handlers moved to auth.go
// This file now only contains ChangePassword and UserProfileEdit functions

func ChangePasswordGet(c *fiber.Ctx) error {
	return c.Render("change_password", fiber.Map{})
}

func ChangePasswordPost(c *fiber.Ctx) error {
	u := c.Locals("user")
	user, ok := u.(models.User)

	if !ok || user.ID == 0 {
		log.Println("Invalid user session")
		return c.Render("change_password", fiber.Map{
			"error": "Please log in again",
		})
	}

	var body struct {
		CurrentPassword    string
		NewPassword        string
		ConfirmNewPassword string
	}

	if err := c.BodyParser(&body); err != nil {
		log.Panic(err)
		return err
	}

	if body.ConfirmNewPassword != body.NewPassword {
		return c.Render("change_password", fiber.Map{
			"error": "Passwords don't match",
		})
	}

	passValErr := passwordValidator.Validate(body.NewPassword, minEntropyBits)
	if passValErr != nil {
		return c.Render("change_password", fiber.Map{
			"error": "Password is not strong enough.",
		})
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(body.CurrentPassword))
	if err != nil {
		log.Println("Invalid current password")
		return c.Render("change_password", fiber.Map{
			"error": "Email or password is incorrect",
		})
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.NewPassword), 10)
	if err != nil {
		log.Panic("Failed to hash password", err)
		return err
	}

	user.PasswordHash = string(hash)
	initializers.DB.Save(&user)

	c.Locals("message", "Successfully reset Password!")
	c.Set("HX-Redirect", "/admin/userProfile")
	return c.Next()
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
