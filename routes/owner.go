package routes

import (
	"shinkyuShotokan/handlers"
	"shinkyuShotokan/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterOwnerRoutes(app *fiber.App) {
	ownerRoutes := app.Group("owner", middleware.RequireOwnerAuth)

	// Owner-only pages
	ownerRoutes.Get("/users", handlers.AdminUsersPage)
}
