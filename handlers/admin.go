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
		"Page":        structs.Page{PageName: "Home", Tabs: utils.CurrentTabs(), Classes: utils.Classes},
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

func EditImageCarousel(c *fiber.Ctx) error {

	if form, err := c.MultipartForm(); err == nil {
		// => *multipart.Form
		os.MkdirAll(fmt.Sprintf("%s/assets/image_carousel/", os.Getenv("UPLOAD_DIR")), 0700)
	
		if token := form.Value["token"]; len(token) > 0 {
		  // Get key value:
		  fmt.Println(token[0])
		}
	
		// Get all files from "documents" key:
		files := form.File["NewImages"]
		// => []*multipart.FileHeader
	
		// Loop through files:
		for _, file := range files {
		  fmt.Println(file.Filename, file.Size, file.Header["Content-Type"][0])
		  // => "tutorial.pdf" 360641 "application/pdf"
	
		  // Save the files to disk:
		  if err := c.SaveFile(file, fmt.Sprintf("%s/assets/image_carousel/%s", os.Getenv("UPLOAD_DIR"), file.Filename)); err != nil {
			return err
		  }
		}
	  }
	log.Println("about to create image_carousel")
	log.Println("created image_carousel")
	// c.SaveFile(newCoverPhoto, fmt.Sprintf("%s/assets/image_carousel/%s", os.Getenv("UPLOAD_DIR"), newCoverPhoto.Filename))
	//return fmt.Sprintf("%s/assets/image_carousel/%s", os.Getenv("UPLOAD_DIR"), newCoverPhoto.Filename)
	c.Set("HX-Redirect", "/")
	return c.Next()
}
