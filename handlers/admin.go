package handlers

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"shinkyuShotokan/initializers"
	"shinkyuShotokan/models"
	"shinkyuShotokan/queries"
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
		"Page":        structs.Page{PageName: "Home", Tabs: utils.CurrentTabs(), Classes: utils.Classes},
		"Events":      events,
		"ImagePaths":  imagePaths,
		"EventPhotos": eventImagePaths,
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

func AdminPage(c *fiber.Ctx) error {
	return c.Redirect("/admin/locations")
}

func AdminLocationPage(c *fiber.Ctx) error {
	persistedLocations := queries.GetLocations()
	adminPage := fiber.Map{
		"Page":      structs.Page{PageName: "Admin", Tabs: utils.CurrentTabs(), Classes: utils.Classes},
		"Locations": persistedLocations,
	}

	return c.Render("adminPage", adminPage)
}

func AddLocation(c *fiber.Ctx) error {
	var body struct {
		Name       string
		Address    string
		GoogleMaps string
	}

	if err := c.BodyParser(&body); err != nil {
		log.Print(err)
		return err
	}

	location := models.Location{Name: body.Name, Address: body.Address, GoogleMapsIframe: body.GoogleMaps}
	result := initializers.DB.Create(&location)

	if result.Error != nil {
		log.Print("Error creating Location", result.Error)
		return result.Error
	}

	return c.Redirect("/admin")
}

func LocationGet(c *fiber.Ctx) error {
	id := c.Params("id")
	var location models.Location
	initializers.DB.First(&location, id)

	return c.Render("location", location)
}

func EditLocationGet(c *fiber.Ctx) error {
	id := c.Params("id")
	var location models.Location
	initializers.DB.First(&location, id)

	return c.Render("edit_location_form", fiber.Map{
		"Location": location,
	})
}

func EditLocationPut(c *fiber.Ctx) error {
	id := c.Params("id")
	var location models.Location
	initializers.DB.First(&location, id)

	var body struct {
		Name       string
		Address    string
		GoogleMaps string
	}

	if err := c.BodyParser(&body); err != nil {
		log.Print(err)
		return err
	}

	location.Name = body.Name
	location.Address = body.Address
	location.GoogleMapsIframe = body.GoogleMaps

	result := initializers.DB.Save(&location)

	if result.Error != nil {
		log.Print("Error updating Location", result.Error)
		return result.Error
	}

	return c.Render("location", location)
}
