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
	// Add images from public/image_carousel (legacy)
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
	// Add images from UPLOAD_DIR/assets/image_carousel (uploaded)
	uploadedCarouselPath := os.Getenv("UPLOAD_DIR") + "/assets/image_carousel/"
	filepath.Walk(uploadedCarouselPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return nil
		}
		if !info.IsDir() {
			imagePaths = append(imagePaths, strings.Replace(path, os.Getenv("UPLOAD_DIR"), "/upload", 1))
		}
		return nil
	})
	eventImagePaths := getExistingEventCoverPhotos()
	eventCardImagePaths := getExistingEventCardPhotos()

	events := queries.GetUpcomingEvents()
	eventTypes := queries.GetEventTypes()

	homePage := fiber.Map{
		"Page":            structs.Page{PageName: "Home", Tabs: utils.CurrentTabs(), Classes: utils.Classes},
		"Events":          events,
		"ImagePaths":      imagePaths,
		"EventPhotos":     eventImagePaths,
		"EventCardPhotos": eventCardImagePaths,
		"EventTypes":      eventTypes,
		"Locations":       queries.GetLocations(),
		"message":         c.Locals("message"),
	}

	hxRequest, err := strconv.ParseBool(c.Get("hx-request"))
	if err != nil {
		hxRequest = false
	}

	if hxRequest {
		return c.Render("homePage", homePage)
	} else {
		return c.Render("home", homePage)
	}
}
