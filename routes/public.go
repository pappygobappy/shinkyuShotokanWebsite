package routes

import (
	"shinkyuShotokan/handlers"
	"shinkyuShotokan/middleware"
	"shinkyuShotokan/utils"

	"github.com/gofiber/fiber/v2"
)

func RegisterPublicRoutes(app *fiber.App) {
	mainRoutes := app.Group("/", middleware.AttachUser)

	// Public pages
	mainRoutes.Get("/", handlers.Home)
	mainRoutes.Get("/instructors", handlers.Instructors)
	mainRoutes.Get("/instructors/sensei_sue", handlers.SenseiSue)
	mainRoutes.Get("/history", handlers.History)
	mainRoutes.Get("/events/:id", handlers.Event)
	mainRoutes.Get("/requirements/:rank", handlers.Requirements)
	mainRoutes.Get("/contact-us", handlers.ContactUs)
	mainRoutes.Get("/calendar", handlers.Calendar)
	mainRoutes.Get("/calendar/:id", handlers.CalendarItemView)

	// Classes (dynamic routes from utils.Classes)
	for _, class := range utils.Classes {
		mainRoutes.Get(class.GetUrl, handlers.Classes)
	}

	// Auth pages
	mainRoutes.Get("/login", handlers.LoginGet)
	mainRoutes.Get("/signup", handlers.SignupGet)
	mainRoutes.Post("/login", handlers.LoginPost)
	mainRoutes.Post("/signup", handlers.SignupPost)
	mainRoutes.Get("/forgot-password", handlers.ForgotPasswordGet)
	mainRoutes.Post("/forgot-password", handlers.ForgotPasswordPost)
	mainRoutes.Get("/reset-password/:token", handlers.ResetPasswordTokenGet)
	mainRoutes.Post("/reset-password/:token", handlers.ResetPasswordTokenPost)

	// Gear page
	mainRoutes.Get("/gear", handlers.GearPage)
}
