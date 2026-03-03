package routes

import (
	"shinkyuShotokan/handlers"
	"shinkyuShotokan/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterAdminRoutes(app *fiber.App) {
	adminRoutes := app.Group("admin", middleware.RequireAuth)

	// Admin dashboard
	adminRoutes.Get("/", handlers.AdminPage)

	// Location management
	adminRoutes.Get("/locations", handlers.AdminLocationPage)
	adminRoutes.Post("/locations", handlers.AddLocation)
	adminRoutes.Get("/locations/:id", handlers.LocationGet)
	adminRoutes.Get("/locations/:id/edit", handlers.EditLocationGet)
	adminRoutes.Put("/locations/:id", handlers.EditLocationPut)

	// Carousel images
	adminRoutes.Get("/upload-carousel-image", handlers.UploadCarouselImagePage)
	adminRoutes.Post("/upload-carousel-image", handlers.UploadCarouselImage)

	// Class management
	adminRoutes.Get("/classes/:id", handlers.EditClassGet)
	adminRoutes.Put("/classes/:id", handlers.EditClassPut)
	adminRoutes.Post("/classSession", handlers.AddClassSession)
	adminRoutes.Delete("/classSession", handlers.DeleteClassPeriodSessions)
	adminRoutes.Get("/calendar/:id", handlers.EditClassSessionGet)
	adminRoutes.Put("/calendar/:id", handlers.EditClassSessionPut)
	adminRoutes.Delete("/calendar/:id", handlers.DeleteClassSession)
	adminRoutes.Get("/classSessionForm", handlers.AddClassSessionForm)
	adminRoutes.Get("/deleteClassSessionsForm", handlers.GetDeleteClassSessionsForm)

	// Class periods
	adminRoutes.Get("/classPeriod", handlers.AddClassPeriodGet)
	adminRoutes.Post("/classPeriod", handlers.AddClassPeriodPost)
	adminRoutes.Get("/classPeriod/:id/edit", handlers.EditClassPeriodGet)
	adminRoutes.Put("/classPeriod/:id/edit", handlers.EditClassPeriodPut)
	adminRoutes.Delete("/classPeriod/:id", handlers.DeleteClassPeriod)

	// Event management
	adminRoutes.Post("/events", handlers.AddEvent)
	adminRoutes.Get("/events/:id", handlers.EditEventGet)
	adminRoutes.Put("/events/:id", handlers.EditEventPost)
	adminRoutes.Delete("/events/:id", handlers.DeleteEventPost)
	adminRoutes.Post("/logout", handlers.LogoutPost)

	// Password management
	adminRoutes.Get("/reset-password", handlers.ChangePasswordGet)
	adminRoutes.Post("/reset-password", handlers.ChangePasswordPost)
	adminRoutes.Post("/start_add_event", handlers.StartAddEvent)

	// Event templates
	adminRoutes.Get("/event-templates", handlers.AdminEventTemplatesPage)
	adminRoutes.Post("/event-templates", handlers.AddEventTemplate)
	adminRoutes.Get("/event-templates/:id", handlers.EditEventTemplateGet)
	adminRoutes.Put("/event-templates/:id", handlers.EditEventTemplatePut)
	adminRoutes.Delete("/event-templates/:id", handlers.DeleteEventTemplatePost)

	// Instructor management
	adminRoutes.Get("/instructors/:id/edit", handlers.EditInstructorGet)
	adminRoutes.Put("/instructors/:id", handlers.EditInstructorPut)
	adminRoutes.Post("/instructors/:id/move", handlers.MoveInstructor)
	adminRoutes.Post("/instructors/:id/toggle-hidden", handlers.ToggleInstructorHidden)
	adminRoutes.Get("/instructors/new", handlers.AddInstructorGet)
	adminRoutes.Post("/instructors", handlers.AddInstructorPost)
	adminRoutes.Get("/instructors/upload-page-image", handlers.UploadCurrentInstructorsImagePage)
	adminRoutes.Post("/instructors/upload-page-image", handlers.UploadCurrentInstructorsImage)

	// Carousel images admin
	adminRoutes.Get("/carousel-images", handlers.AdminCarouselImagesPage)
	adminRoutes.Post("/carousel-images/:id/move", handlers.MoveCarouselImage)
	adminRoutes.Post("/carousel-images/:id/remove", handlers.SoftDeleteCarouselImage)
	adminRoutes.Post("/carousel-images/:id/restore", handlers.RestoreCarouselImage)
	adminRoutes.Post("/carousel-images/:id/hard-delete", handlers.HardDeleteCarouselImage)

	// User profile
	adminRoutes.Get("/userProfile", handlers.AdminUserProfilePage)
	adminRoutes.Get("/userProfile/edit", handlers.GetUserProfilePageEdit)
	adminRoutes.Post("/userProfile", handlers.PostUserProfilePageEdit)
}
