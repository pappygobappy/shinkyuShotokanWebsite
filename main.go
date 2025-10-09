package main

import (
	"html/template"
	"log"
	"net/url"
	"os"
	"shinkyuShotokan/handlers"
	"shinkyuShotokan/initializers"
	"shinkyuShotokan/middleware"
	"shinkyuShotokan/models"
	"shinkyuShotokan/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDb()
	initializers.SyncDb()
	utils.Init()
}

func addEngineFuncs(engine *html.Engine) {
	engine.AddFunc("makeMap", func(values ...interface{}) map[string]interface{} {

		makeMap := make(map[string]interface{})
		for i := 0; i < len(values); i += 2 {
			key, ok := values[i].(string)
			if !ok {
				log.Println("error making map")
				return makeMap
			}
			makeMap[key] = values[i+1]
		}
		return makeMap
	})

	engine.AddFunc("htmlRender", func(s string) template.HTML {
		return template.HTML(s)
	})

	engine.AddFunc("gmtRfc5545", func(t time.Time) string {
		return t.In(time.FixedZone("GMT", 0)).Format("20060102T150405Z")
	})

	engine.AddFunc("yahooDateFormat", func(t time.Time) string {
		return t.In(utils.TZ).Format("20060102T150405Z")
	})

	engine.AddFunc("outlookCalInvite", func(event models.Event) string {
		baseURL := "https://outlook.live.com/calendar/0/deeplink/compose?"
		params := url.Values{}
		params.Add("subject", event.Title)
		params.Add("startdt", event.StartTime.In(utils.TZ).Format("2006-01-02T15:04:05"))
		params.Add("enddt", event.EndTime.In(utils.TZ).Format("2006-01-02T15:04:05"))
		params.Add("location", event.Location)
		params.Add("body", string(event.OutlookDescription()))

		return baseURL + params.Encode()
	})

	engine.AddFunc("startTimePSTString", func(t time.Time) string {
		return t.In(utils.TZ).Format("15:04:05")
	})

	engine.AddFunc("formatTimePST", func(t time.Time, format string) string {
		return t.In(utils.TZ).Format(format)
	})

	engine.AddFunc("isToday", func(today time.Time, day time.Time) bool {
		return today.Year() == day.Year() && today.Month() == day.Month() && today.Day() == day.Day()
	})

	engine.AddFunc("minus", func(target int, subtract int) int {
		return target - subtract
	})
}

func main() {
	log.Print("hello world")

	//ENV

	//Create App
	engine := html.New("./templates", ".html")
	addEngineFuncs(engine)

	app := fiber.New(fiber.Config{
		Views:             engine,
		PassLocalsToViews: true,
		BodyLimit:         16 * 1024 * 1024,
	})

	//Register Static Files
	app.Static("/public", "./public")
	app.Static("/upload", os.Getenv("UPLOAD_DIR"))

	//Init Data

	//Register Routes
	mainRoutes := app.Group("/", middleware.AttachUser)
	mainRoutes.Get("/", handlers.Home)
	mainRoutes.Get("/instructors", handlers.Instructors)
	mainRoutes.Get("/history", handlers.History)
	mainRoutes.Get("/events/:id", handlers.Event)
	mainRoutes.Get("/requirements/:rank", handlers.Requirements)
	mainRoutes.Get("/contact-us", handlers.ContactUs)
	mainRoutes.Get("/calendar", handlers.Calendar)
	mainRoutes.Get("/calendar/:id", handlers.CalendarItemView)
	mainRoutes.Get("/login", handlers.LoginGet)
	mainRoutes.Get("/signup", handlers.SignupGet)
	mainRoutes.Post("/login", handlers.LoginPost)
	mainRoutes.Post("/signup", handlers.SignupPost)
	for _, class := range utils.Classes {
		mainRoutes.Get(class.GetUrl, handlers.Classes)
	}

	adminRoutes := app.Group("admin", middleware.RequireAuth)
	adminRoutes.Get("/", handlers.AdminPage)
	adminRoutes.Get("/locations", handlers.AdminLocationPage)
	adminRoutes.Get("/users", handlers.AdminUsersPage)
	adminRoutes.Get("/upload-carousel-image", handlers.UploadCarouselImagePage)
	adminRoutes.Post("/upload-carousel-image", handlers.UploadCarouselImage)
	adminRoutes.Post("/locations", handlers.AddLocation)
	adminRoutes.Get("/locations/:id", handlers.LocationGet)
	adminRoutes.Get("/locations/:id/edit", handlers.EditLocationGet)
	adminRoutes.Put("/locations/:id", handlers.EditLocationPut)
	adminRoutes.Get("/classes/:id", handlers.EditClassGet)
	adminRoutes.Put("/classes/:id", handlers.EditClassPut)
	adminRoutes.Post("/classSession", handlers.AddClassSession)
	adminRoutes.Delete("/classSession", handlers.DeleteClassPeriodSessions)
	adminRoutes.Get("/calendar/:id", handlers.EditClassSessionGet)
	adminRoutes.Put("/calendar/:id", handlers.EditClassSessionPut)
	adminRoutes.Delete("/calendar/:id", handlers.DeleteClassSession)
	adminRoutes.Get("/classSessionForm", handlers.AddClassSessionForm)
	adminRoutes.Get("/deleteClassSessionsForm", handlers.GetDeleteClassSessionsForm)
	adminRoutes.Get("/classPeriod", handlers.AddClassPeriodGet)
	adminRoutes.Post("/classPeriod", handlers.AddClassPeriodPost)
	adminRoutes.Get("/classPeriod/:id/edit", handlers.EditClassPeriodGet)
	adminRoutes.Put("/classPeriod/:id/edit", handlers.EditClassPeriodPut)
	adminRoutes.Delete("/classPeriod/:id", handlers.DeleteClassPeriod)
	adminRoutes.Post("/events", handlers.AddEvent)
	adminRoutes.Get("/events/:id", handlers.EditEventGet)
	adminRoutes.Put("/events/:id", handlers.EditEventPost)
	adminRoutes.Delete("/events/:id", handlers.DeleteEventPost)
	adminRoutes.Post("/logout", handlers.LogoutPost)
	adminRoutes.Get("/reset-password", handlers.ResetPasswordGet)
	adminRoutes.Post("/reset-password", handlers.ResetPasswordPost)
	adminRoutes.Post("/start_add_event", handlers.StartAddEvent)

	// Event Template Admin Routes
	adminRoutes.Get("/event-templates", handlers.AdminEventTemplatesPage)
	adminRoutes.Post("/event-templates", handlers.AddEventTemplate)
	adminRoutes.Get("/event-templates/:id", handlers.EditEventTemplateGet)
	adminRoutes.Put("/event-templates/:id", handlers.EditEventTemplatePut)
	adminRoutes.Delete("/event-templates/:id", handlers.DeleteEventTemplatePost)

	// Instructor Admin Routes
	adminRoutes.Get("/instructors/:id/edit", handlers.EditInstructorGet)
	adminRoutes.Put("/instructors/:id", handlers.EditInstructorPut)
	adminRoutes.Post("/instructors/:id/move", handlers.MoveInstructor)
	adminRoutes.Post("/instructors/:id/toggle-hidden", handlers.ToggleInstructorHidden)
	adminRoutes.Get("/instructors/new", handlers.AddInstructorGet)
	adminRoutes.Post("/instructors", handlers.AddInstructorPost)

	// // Carousel Images Admin Routes
	adminRoutes.Get("/carousel-images", handlers.AdminCarouselImagesPage)
	adminRoutes.Post("/carousel-images/:id/move", handlers.MoveCarouselImage)
	adminRoutes.Post("/carousel-images/:id/remove", handlers.SoftDeleteCarouselImage)
	adminRoutes.Post("/carousel-images/:id/restore", handlers.RestoreCarouselImage)

	//Start App
	//log.Fatal(http.ListenAndServe(":8000", nil))
	app.Listen(":" + os.Getenv("PORT"))
}
