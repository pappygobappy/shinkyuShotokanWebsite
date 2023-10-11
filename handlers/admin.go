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

func AdminHome(c *fiber.Ctx) error {
	paths, _ := os.Getwd()
	var imagePaths []string
	var eventImagePaths []string
	err := filepath.Walk(paths+"/assets/image_carousel/", func(path string, info os.FileInfo, err error) error {

		if err != nil {
			fmt.Println(err)
			return nil
		}
		if !info.IsDir() {
			imagePaths = append(imagePaths, strings.Replace(path, paths, "", 1))
		}
		return nil
	})
	err = filepath.Walk(paths+"/assets/events/", func(path string, info os.FileInfo, err error) error {

		if err != nil {
			fmt.Println(err)
			return nil
		}
		if !info.IsDir() {
			eventImagePaths = append(eventImagePaths, strings.Replace(path, paths, "", 1))
		}
		return nil
	})

	var events []models.Event
	result := initializers.DB.Find(&events)
	if result.Error != nil {
		log.Print(result.Error)
	}
	homePage := fiber.Map{
		"Page":        structs.Page{PageName: "Home", Tabs: utils.Tabs, Classes: utils.Classes},
		"Events":      events,
		"ImagePaths":  imagePaths,
		"EventPhotos": eventImagePaths,
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
