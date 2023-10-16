package handlers

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"shinkyuShotokan/initializers"
	"shinkyuShotokan/models"
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

	var events []models.Event
	result := initializers.DB.Order("date").Find(&events)
	if result.Error != nil {
		log.Print(result.Error)
	}
	homePage := fiber.Map{
		"Page": structs.Page{PageName: "Home", Tabs: utils.Tabs, Classes: utils.Classes},
		"Events": events,
		"ImagePaths": imagePaths,
		"EventPhotos": eventImagePaths,
		"Locations": utils.Locations,
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
