package main

import (
	"log"
	"os"
	"shinkyuShotokan/handlers"
	"shinkyuShotokan/initializers"
	"shinkyuShotokan/middleware"
	"shinkyuShotokan/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDb()
	initializers.SyncDb()
	utils.Init()
}

func main() {
	log.Print("hello world")

	//ENV

	//Create App
	engine := html.New("./templates", ".html")
	engine.AddFunc("makeMap", func(values ...interface{}) map[string]interface{}{
		
		makeMap := make(map[string]interface{})
		for i := 0; i < len(values); i+=2 {
			key, ok := values[i].(string)
			if !ok {
				log.Println("error making map")
				return makeMap
			}
			makeMap[key] = values[i +1]
		}
		return makeMap
	})
	app := fiber.New(fiber.Config{
		Views: engine,
		PassLocalsToViews: true,
		BodyLimit: 16 * 1024 * 1024,
	})

	//Register Static Files
	app.Static("/assets", "./assets")

	//Init Data
	

	//Register Routes
	mainRoutes := app.Group("/", middleware.AttachUser)
	mainRoutes.Get("/", handlers.Home)
	mainRoutes.Get("/instructors", handlers.Instructors)
	mainRoutes.Get("/history", handlers.History)
	mainRoutes.Get("/events/:id", handlers.Event)

	mainRoutes.Get("/login", handlers.LoginGet)
	mainRoutes.Get("/signup", handlers.SignupGet)
	mainRoutes.Post("/login", handlers.LoginPost)
	mainRoutes.Post("/signup", handlers.SignupPost)
	for _, class := range utils.Classes {
		mainRoutes.Get(class.GetUrl, handlers.Classes)
	}
	
	adminRoutes := app.Group("admin", middleware.RequireAuth)
	adminRoutes.Get("/", handlers.AdminHome)
	adminRoutes.Post("/events", handlers.AddEvent)
	adminRoutes.Get("/events/:id", handlers.EditEventGet)
	adminRoutes.Put("/events/:id", handlers.EditEventPost)
	adminRoutes.Delete("/events/:id", handlers.DeleteEventPost)
	adminRoutes.Post("/logout", handlers.LogoutPost)
	//app.Get("/pre-karate-class", handlers.PreKarateClasses)

	//Start App
	//log.Fatal(http.ListenAndServe(":8000", nil))
	app.Listen(":" + os.Getenv("PORT"))
}
