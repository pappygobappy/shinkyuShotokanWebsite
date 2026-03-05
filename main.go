package main

import (
	"html/template"
	"log"
	"net/url"
	"os"
	"shinkyuShotokan/initializers"
	"shinkyuShotokan/models"
	"shinkyuShotokan/packages/cache"
	"shinkyuShotokan/queries"
	"shinkyuShotokan/routes"
	"shinkyuShotokan/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

var MainCache *cache.MemoryCache

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDb()
	initializers.SyncDb()
	MainCache = cache.New(5 * time.Minute)
	queries.InitCache(MainCache)
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
	routes.RegisterPublicRoutes(app)
	routes.RegisterAdminRoutes(app)
	routes.RegisterOwnerRoutes(app)

	//Start App
	//log.Fatal(http.ListenAndServe(":8000", nil))
	app.Listen(":" + os.Getenv("PORT"))
}
