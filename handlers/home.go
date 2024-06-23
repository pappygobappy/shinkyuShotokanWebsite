package handlers

import (
	"fmt"
	"os"
	"path/filepath"
	"shinkyuShotokan/queries"
	"shinkyuShotokan/structs"
	"shinkyuShotokan/utils"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func Home(c *fiber.Ctx) error {
	paths, _ := os.Getwd()
	var imagePaths []string
	err := filepath.Walk(paths+"/public/image_carousel/", func(path string, info os.FileInfo, err error) error {

		if err != nil {
			fmt.Println(err)
			return nil
		}
		if !info.IsDir() {
			imagePaths = append(imagePaths, strings.Replace(path, paths, "", 1))
		}
		return nil
	})

	eventImagePaths := getExistingEventCoverPhotos()

	events := queries.GetUpcomingEvents()
	homePage := fiber.Map{
		"Page":        structs.Page{PageName: "Home", Tabs: utils.CurrentTabs(), Classes: utils.Classes},
		"Events":      events,
		"ImagePaths":  imagePaths,
		"EventPhotos": eventImagePaths,
		"Locations":   queries.GetLocations(),
		"message": c.Locals("message"),
	}

	hxRequest, err := strconv.ParseBool(c.Get("hx-request"))
	if err != nil {
		hxRequest = false
	}

	if hxRequest == true {
		return c.Render("homePage", homePage)
	} else {
		return c.Render("home", homePage)
	}
}
