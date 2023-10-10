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
	"time"

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

func AddEvent(c *fiber.Ctx) error {
	var body struct {
		Name               string
		Date               string
		Location           string
		Address            string
		GoogleMapsIframe   string
		Description        string
		ExistingCoverPhoto string
	}

	if err := c.BodyParser(&body); err != nil {
		log.Print(err)
		return err
	}

	date, error := time.ParseInLocation("2006-01-02", body.Date, time.Local)

	if error != nil {
		fmt.Println(error)
		return c.Redirect("/")
	}
	
	newCoverPhoto, err := c.FormFile("NewCoverPhoto")

	var photoUrl string
	if err != nil {
		fmt.Println("No new cover photo")
		photoUrl = body.ExistingCoverPhoto
	} else {
		photoUrl = fmt.Sprintf("/assets/events/%s", newCoverPhoto.Filename)
		c.SaveFile(newCoverPhoto, fmt.Sprintf(".%s", photoUrl))
	}

	event := models.Event{Title: body.Name, Date: date, Description: body.Description, PictureUrl: photoUrl, Location: body.Location, GoogleMapsIframe: body.GoogleMapsIframe, Address: body.Address}

	result := initializers.DB.Create(&event)

	if result.Error != nil {
		log.Print("Error creating Event", result.Error)
		return result.Error
	}

	//Handle Files
	if form, err := c.MultipartForm(); err == nil {
		// => *multipart.Form

		if token := form.Value["token"]; len(token) > 0 {
			// Get key value:
			fmt.Println(token[0])
		}

		// Get all files from "Files" key:
		files := form.File["Files"]
		// => []*multipart.FileHeader
		os.MkdirAll(fmt.Sprintf("./assets/event/%s/files", strconv.FormatUint(uint64(event.ID), 10)), 0700)
		// Loop through files:
		for _, file := range files {
			// => "tutorial.pdf" 360641 "application/pdf"

			// Save the files to disk:
			if err := c.SaveFile(file, fmt.Sprintf("./assets/event/%s/files/%s", strconv.FormatUint(uint64(event.ID), 10), file.Filename)); err != nil {
				log.Println(err)
			}
		}
	}

	return c.Redirect("/")
}
