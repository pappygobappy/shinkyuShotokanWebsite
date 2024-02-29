package handlers

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"shinkyuShotokan/initializers"
	"shinkyuShotokan/models"
	"shinkyuShotokan/structs"
	"shinkyuShotokan/utils"
	"strconv"
	"strings"
	"time"

	ics "github.com/arran4/golang-ical"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func getExistingEventCoverPhotos() []string {
	workingDir, _ := os.Getwd()
	var eventImagePaths []string
	err := filepath.Walk(workingDir+"/public/events/", func(path string, info os.FileInfo, err error) error {

		if err != nil {
			fmt.Println(err)
			return nil
		}
		if !info.IsDir() {
			eventImagePaths = append(eventImagePaths, strings.Replace(path, workingDir, "", 1))
		}
		return nil
	})

	if err != nil {
		log.Println("Failed to get existing cover photos")
	}

	uploadedCoverPhotosPath := fmt.Sprintf("%s/assets/events/", os.Getenv("UPLOAD_DIR"))

	err = filepath.Walk(uploadedCoverPhotosPath, func(path string, info os.FileInfo, err error) error {

		if err != nil {
			fmt.Println(err)
			return nil
		}
		if !info.IsDir() {
			eventImagePaths = append(eventImagePaths, strings.Replace(path, os.Getenv("UPLOAD_DIR"), "/upload", 1))
		}
		return nil
	})

	if err != nil {
		log.Println("Failed to get uploaded cover photos")
	}

	return eventImagePaths
}

func uploadEventFiles(event models.Event, c *fiber.Ctx) {
	if form, err := c.MultipartForm(); err == nil {
		// Get all files from "Files" key:
		files := form.File["Files"]
		// => []*multipart.FileHeader
		os.MkdirAll(fmt.Sprintf("%s/assets/event/%s/files", os.Getenv("UPLOAD_DIR"), strconv.FormatUint(uint64(event.ID), 10)), 0700)
		// Loop through files:
		for _, file := range files {
			// Save the files to disk:
			if err := c.SaveFile(file, fmt.Sprintf("%s/assets/event/%s/files/%s", os.Getenv("UPLOAD_DIR"), strconv.FormatUint(uint64(event.ID), 10), file.Filename)); err != nil {
				log.Println(err)
			}
		}
	}
}

func getEventFilePaths(event models.Event) map[string]string {
	basePath := fmt.Sprintf("%s/assets/event/%s/files/", os.Getenv("UPLOAD_DIR"), strconv.FormatUint(uint64(event.ID), 10))
	//key filename, value path
	var files = make(map[string]string)
	err := filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {

		if err != nil {
			fmt.Println(err)
			return nil
		}
		if !info.IsDir() {
			files[info.Name()] = strings.Replace(path, os.Getenv("UPLOAD_DIR"), "/upload", 1)
			//files = append(files, strings.Replace(path, paths, "", 1))
		}
		return nil
	})
	if err != nil {
		log.Println(err)
	}
	return files
}

func getCoverPhotoUrl(existingCoverPhoto string, c *fiber.Ctx) string {
	newCoverPhoto, err := c.FormFile("NewCoverPhoto")
	if err != nil {
		fmt.Println("No new cover photo")
		return existingCoverPhoto
	} else {
		os.MkdirAll(fmt.Sprintf("%s/assets/events/", os.Getenv("UPLOAD_DIR")), 0700)
		c.SaveFile(newCoverPhoto, fmt.Sprintf("%s/assets/events/%s", os.Getenv("UPLOAD_DIR"), newCoverPhoto.Filename))
		return fmt.Sprintf("/upload/assets/events/%s", newCoverPhoto.Filename)
	}
}

func Event(c *fiber.Ctx) error {
	id := c.Params("id")
	var event models.Event
	initializers.DB.First(&event, id)
	pageName := "Upcoming Events"
	if event.Date.Before(time.Now()) {
		pageName = "Past Events"
	}
	page := structs.Page{PageName: pageName, Tabs: utils.CurrentTabs(), Classes: utils.Classes}
	files := getEventFilePaths(event)

	return c.Render("event", fiber.Map{
		"Page":        page,
		"Event":       event,
		"Description": event.Description,
		"Files":       files,
		"Location":    utils.Locations[event.Location],
	})
}

func EditEventGet(c *fiber.Ctx) error {
	id := c.Params("id")
	var event models.Event
	initializers.DB.First(&event, id)
	page := structs.Page{PageName: "Event", Tabs: utils.CurrentTabs(), Classes: utils.Classes}
	files := getEventFilePaths(event)

	eventImagePaths := getExistingEventCoverPhotos()

	return c.Render("edit_event", fiber.Map{
		"Page":        page,
		"Event":       event,
		"EventPhotos": eventImagePaths,
		"Description": event.Description,
		"Files":       files,
		"Locations":   utils.Locations,
	})
}

func AddEvent(c *fiber.Ctx) error {
	var body struct {
		Name               string
		Date               string
		StartTime          string
		EndTime            string
		Location           string
		Description        template.HTML
		ExistingCoverPhoto string
	}

	if err := c.BodyParser(&body); err != nil {
		log.Print(err)
		return err
	}

	date, error := time.ParseInLocation("2006-01-02", body.Date, utils.TZ)
	startTime, error := time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", body.Date, body.StartTime), utils.TZ)
	endTime, error := time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", body.Date, body.EndTime), utils.TZ)

	if error != nil {
		fmt.Println(error)
		return c.Redirect("/")
	}

	photoUrl := getCoverPhotoUrl(body.ExistingCoverPhoto, c)

	event := models.Event{Title: body.Name, Date: date, StartTime: startTime, EndTime: endTime, Description: body.Description, PictureUrl: photoUrl, Location: body.Location}

	result := initializers.DB.Create(&event)

	if result.Error != nil {
		log.Print("Error creating Event", result.Error)
		return result.Error
	}

	//Handle Files
	uploadEventFiles(event, c)
	createEventIcs(event, utils.Locations[event.Location])

	return c.Redirect("/")
}

func EditEventPost(c *fiber.Ctx) error {
	id := c.Params("id")
	var event models.Event
	initializers.DB.First(&event, id)
	//page := structs.Page{PageName: "Event", Tabs: utils.CurrentTabs(), Classes: utils.Classes}

	var body struct {
		Name               string
		Date               string
		StartTime          string
		EndTime            string
		Location           string
		Description        template.HTML
		ExistingCoverPhoto string
		DeletedFiles       string
	}

	if err := c.BodyParser(&body); err != nil {
		log.Print(err)
		return err
	}

	date, error := time.ParseInLocation("2006-01-02", body.Date, utils.TZ)
	startTime, error := time.ParseInLocation("2006-01-02 15:04:00", fmt.Sprintf("%s %s", body.Date, body.StartTime), utils.TZ)
	endTime, error := time.ParseInLocation("2006-01-02 15:04:00", fmt.Sprintf("%s %s", body.Date, body.EndTime), utils.TZ)

	filesToDelete := strings.Split(body.DeletedFiles, ",")

	for _, file := range filesToDelete {
		os.Remove(fmt.Sprintf("%s/assets/event/%s/files/%s", os.Getenv("UPLOAD_DIR"), id, file))
	}

	if error != nil {
		fmt.Println(error)
		return c.Redirect("/")
	}

	photoUrl := getCoverPhotoUrl(body.ExistingCoverPhoto, c)

	event.Title = body.Name
	event.Date = date
	event.StartTime = startTime
	event.EndTime = endTime
	event.Description = body.Description
	event.PictureUrl = photoUrl
	event.Location = body.Location

	result := initializers.DB.Save(&event)

	if result.Error != nil {
		log.Print("Error creating Event", result.Error)
		return result.Error
	}

	//Handle Files
	uploadEventFiles(event, c)
	createEventIcs(event, utils.Locations[event.Location])

	//files := getEventFilePaths(event)
	c.Set("HX-Redirect", "/events/"+strconv.FormatUint(uint64(event.ID), 10))
	return c.Next()
}

func DeleteEventPost(c *fiber.Ctx) error {
	id := c.Params("id")
	initializers.DB.Delete(&models.Event{}, id)
	os.RemoveAll(fmt.Sprintf("%s/assets/event/%s/files/", os.Getenv("UPLOAD_DIR"), id))

	c.Set("HX-Redirect", "/")
	return c.Next()
}

func createEventIcs(e models.Event, l models.Location) {
	cal := ics.NewCalendar()
	cal.SetMethod(ics.MethodRequest)
	event := cal.AddEvent(uuid.New().String())
	event.SetCreatedTime(time.Now())
	event.SetDtStampTime(time.Now())
	event.SetModifiedAt(time.Now())
	event.SetStartAt(e.StartTime)
	event.SetEndAt(e.EndTime)
	event.SetSummary(e.Title)
	event.SetLocation(l.Name)
	re := regexp.MustCompile(`\r?\n`)
	htmlDesc := re.ReplaceAllString(string(e.Description), "\n")
	desc := strings.Replace(htmlDesc, "<b>", "*", -1)
	desc = strings.Replace(desc, "</b>", "*", -1)
	for strings.Contains(desc, "<a") {
		start := strings.Index(desc, "<a")
		end := strings.Index(desc, "</a>") + 4
		newDesc := strings.TrimSpace(desc[0:start])
		desc = newDesc + strings.TrimSpace(desc[end:])
	}
	desc = strings.Replace(desc, "<a>*</a>", "", -1)
	event.SetDescription(desc)
	event.SetProperty("X-ALT-DESC;FMTTYPE=text/html", fmt.Sprintf("<!DOCTYPE HTML PUBLIC \"\"-//W3C//DTD HTML 3.2//EN\"\"><HTML><BODY>%s</BODY></HTML>", htmlDesc))
	ics := cal.Serialize()
	f, err := os.Create(fmt.Sprintf("%s/assets/event/%s/%s.ics", os.Getenv("UPLOAD_DIR"), strconv.FormatUint(uint64(e.ID), 10), e.Title))
	if err != nil {
		fmt.Println(err)
		return
	}
	f.WriteString(ics)
}
