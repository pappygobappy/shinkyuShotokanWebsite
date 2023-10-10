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
	"text/template"
	"time"

	"github.com/gofiber/fiber/v2"
)

func getEventFilePaths(event models.Event) map[string]string {
	paths, _ := os.Getwd()
	basePath := fmt.Sprintf("%s/assets/event/%s/files/", paths, strconv.FormatUint(uint64(event.ID), 10))
	var files = make(map[string]string)
	err := filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {

		if err != nil {
			fmt.Println(err)
			return nil
		}
		if !info.IsDir() {
			files[strings.Replace(path, basePath, "", 1)] = strings.Replace(path, paths, "", 1)
			//files = append(files, strings.Replace(path, paths, "", 1))
		}
		return nil
	})
	if err != nil {
		log.Println(err)
	}
	return files
}

func Event(c *fiber.Ctx) error {
	id := c.Params("id")
	var event models.Event
	initializers.DB.First(&event, id)
	page := structs.Page{PageName: "Event", Tabs: utils.Tabs, Classes: utils.Classes}
	files := getEventFilePaths(event)

	return c.Render("event", fiber.Map{
		"Page":        page,
		"Event":       event,
		"Description": strings.Replace(template.HTMLEscapeString(event.Description), "\n", "<br/>", -1),
		"Files":       files,
	})
}

func EditEventGet(c *fiber.Ctx) error {
	id := c.Params("id")
	var event models.Event
	initializers.DB.First(&event, id)
	page := structs.Page{PageName: "Event", Tabs: utils.Tabs, Classes: utils.Classes}
	paths, _ := os.Getwd()
	files := getEventFilePaths(event)

	var eventImagePaths []string
	err := filepath.Walk(paths+"/assets/events/", func(path string, info os.FileInfo, err error) error {

		if err != nil {
			fmt.Println(err)
			return nil
		}
		if !info.IsDir() {
			eventImagePaths = append(eventImagePaths, strings.Replace(path, paths, "", 1))
		}
		return nil
	})

	if err != nil {
		log.Println("Failed to get files")
	}

	var events []models.Event
	result := initializers.DB.Find(&events)
	if result.Error != nil {
		log.Print(result.Error)
	}
	return c.Render("edit_event", fiber.Map{
		"Page":        page,
		"Event":       event,
		"EventPhotos": eventImagePaths,
		"Description": strings.Replace(template.HTMLEscapeString(event.Description), "\n", "<br/>", -1),
		"Files":       files,
	})
}

func EditEventPost(c *fiber.Ctx) error {
	id := c.Params("id")
	var event models.Event
	initializers.DB.First(&event, id)
	page := structs.Page{PageName: "Event", Tabs: utils.Tabs, Classes: utils.Classes}

	var body struct {
		Name               string
		Date               string
		Location           string
		Address            string
		GoogleMapsIframe   string
		Description        string
		ExistingCoverPhoto string
		DeletedFiles       string
	}

	if err := c.BodyParser(&body); err != nil {
		log.Print(err)
		return err
	}

	date, error := time.ParseInLocation("2006-01-02", body.Date, time.Local)

	filesToDelete := strings.Split(body.DeletedFiles, ",")

	for _, file := range filesToDelete {
		file = "."+file
		os.Remove(file)
	}

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

	//event = models.Event{Title: body.Name, Date: date, Description: body.Description, PictureUrl: photoUrl, Location: body.Location, GoogleMapsIframe: body.GoogleMapsIframe, Address: body.Address}
	event.Title = body.Name
	event.Date = date
	event.Description = body.Description
	event.PictureUrl = photoUrl
	event.Location = body.Location
	event.GoogleMapsIframe = body.GoogleMapsIframe
	event.Address = body.Address

	result := initializers.DB.Save(&event)

	log.Println("Saved")

	if result.Error != nil {
		log.Print("Error creating Event", result.Error)
		return result.Error
	}

	//Handle Files
	if form, err := c.MultipartForm(); err == nil {
		// => *multipart.Form

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

	files := getEventFilePaths(event)
	return c.Render("event", fiber.Map{
		"Page":        page,
		"Event":       event,
		"Description": strings.Replace(template.HTMLEscapeString(event.Description), "\n", "<br/>", -1),
		"Files":       files,
	})
}

func DeleteEventPost(c *fiber.Ctx) error {
	id := c.Params("id")
	initializers.DB.Delete(&models.Event{}, id)
	os.RemoveAll(fmt.Sprintf("./assets/event/%s/files/", id))

	c.Set("HX-Redirect", "/")
	return c.Next()
}
