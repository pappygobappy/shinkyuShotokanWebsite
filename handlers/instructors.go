package handlers

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"shinkyuShotokan/models"
	"shinkyuShotokan/queries"
	"shinkyuShotokan/structs"
	"shinkyuShotokan/utils"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func ToggleInstructorHidden(c *fiber.Ctx) error {
	idParam := c.Params("id")
	idUint64, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid id")
	}
	hidden := c.FormValue("hidden") == "true"
	if err := queries.SetInstructorHidden(uint(idUint64), hidden); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return Instructors(c)
}

func Instructors(c *fiber.Ctx) error {
	user := c.Locals("user")
	var visible, hidden []models.Instructor
	visible = queries.GetVisibleInstructors()
	if user != nil {
		hidden = queries.GetHiddenInstructors()
	}
	currentInstructorsPagePhoto := queries.GetCurrentInstructorsPagePhoto()
	instructorsPage := fiber.Map{
		"Page":                        structs.Page{PageName: "Instructors", Tabs: utils.CurrentTabs(), Classes: utils.Classes},
		"Instructors":                 visible,
		"HiddenInstructors":           hidden,
		"CurrentInstructorsPagePhoto": currentInstructorsPagePhoto,
	}
	hxRequest, err := strconv.ParseBool(c.Get("hx-request"))
	if err != nil {
		hxRequest = false
	}
	if hxRequest {
		return c.Render("instructorsPage", instructorsPage)
	} else {
		return c.Render("instructors", instructorsPage)
	}
}

func SenseiSue(c *fiber.Ctx) error {
	paths, _ := os.Getwd()
	var imagePaths []string
	// Add images from public/image_carousel (legacy)
	err := filepath.Walk(paths+"/public/instructors/sensei_sue", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return nil
		}
		if !info.IsDir() {
			if paths != "/" {
				path = strings.Replace(path, paths, "", 1)
			}
			imagePaths = append(imagePaths, path)
		}
		return nil
	})
	if err != nil {
		log.Println("error making map")
	}

	instructorsPage := fiber.Map{
		"Page":       structs.Page{PageName: "Instructors", Tabs: utils.CurrentTabs(), Classes: utils.Classes},
		"ImagePaths": imagePaths,
	}

	return c.Render("sensei_sue", instructorsPage)
}

func MoveInstructor(c *fiber.Ctx) error {
	idParam := c.Params("id")
	idUint64, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid id")
	}
	direction := c.FormValue("direction")
	if direction != "up" && direction != "down" {
		return fiber.NewError(fiber.StatusBadRequest, "invalid direction")
	}
	if err := queries.MoveInstructor(uint(idUint64), direction); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return Instructors(c)
}

func EditInstructorGet(c *fiber.Ctx) error {
	id := c.Params("id")
	idUint64, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid id")
	}
	instructor := queries.GetInstructorByID(uint(idUint64))
	page := structs.Page{PageName: "Instructor", Tabs: utils.CurrentTabs(), Classes: utils.Classes}
	return c.Render("edit_instructor", fiber.Map{
		"Page":       page,
		"Instructor": instructor,
	})
}

func EditInstructorPut(c *fiber.Ctx) error {
	id := c.Params("id")
	idUint64, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid id")
	}
	instructor := queries.GetInstructorByID(uint(idUint64))

	// Get form values
	name := c.FormValue("Name")
	bio := c.FormValue("Bio")
	zoomLevel, _ := strconv.ParseInt(c.FormValue("ZoomLevel"), 10, 64)
	offsetX, _ := strconv.ParseInt(c.FormValue("OffsetX"), 10, 64)
	offsetY, _ := strconv.ParseInt(c.FormValue("OffsetY"), 10, 64)

	// Handle optional new picture
	if file, err := c.FormFile("NewPicture"); err == nil && file != nil {
		baseDir := fmt.Sprintf("%s/assets/instructors", os.Getenv("UPLOAD_DIR"))
		os.MkdirAll(baseDir, 0700)
		if err := c.SaveFile(file, fmt.Sprintf("%s/%s", baseDir, file.Filename)); err == nil {
			instructor.PictureUrl = fmt.Sprintf("/upload/assets/instructors/%s", file.Filename)
		}
	}

	instructor.Name = name
	instructor.Bio = bio
	instructor.ZoomLevel = int(zoomLevel)
	instructor.OffsetX = int(offsetX)
	instructor.OffsetY = int(offsetY)
	queries.UpdateInstructor(instructor)

	// After update, go back to instructors list
	return Instructors(c)
}

func AddInstructorGet(c *fiber.Ctx) error {
	page := structs.Page{PageName: "Add Instructor", Tabs: utils.CurrentTabs(), Classes: utils.Classes}
	return c.Render("edit_instructor", fiber.Map{
		"Page":       page,
		"Instructor": fiber.Map{"ID": 0, "Name": "", "Bio": "", "PictureUrl": ""},
	})
}

func AddInstructorPost(c *fiber.Ctx) error {
	var body struct {
		Name string
		Bio  string
	}
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid form")
	}

	pictureUrl := ""
	if file, err := c.FormFile("NewPicture"); err == nil && file != nil {
		baseDir := fmt.Sprintf("%s/assets/instructors", os.Getenv("UPLOAD_DIR"))
		os.MkdirAll(baseDir, 0700)
		if err := c.SaveFile(file, fmt.Sprintf("%s/%s", baseDir, file.Filename)); err == nil {
			pictureUrl = fmt.Sprintf("/upload/assets/instructors/%s", file.Filename)
		}
	}

	order := queries.GetNextInstructorDisplayOrder()
	instructor := models.Instructor{Name: body.Name, Bio: body.Bio, PictureUrl: pictureUrl, DisplayOrder: order}
	queries.CreateInstructor(instructor)
	return Instructors(c)
}

func UploadCurrentInstructorsImagePage(c *fiber.Ctx) error {
	return c.Render("upload_image_form", fiber.Map{
		"Page":       structs.Page{PageName: "Upload Current Instructors Page Image", Tabs: utils.CurrentTabs(), Classes: utils.Classes},
		"user":       c.Locals("user"),
		"title":      "Upload Current Instructors Page Image",
		"formAction": "/admin/instructors/upload-page-image",
	})
}

func UploadCurrentInstructorsImage(c *fiber.Ctx) error {
	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(400).SendString("No file uploaded")
	}
	uploadDir := os.Getenv("UPLOAD_DIR") + "/assets/instructors/"
	os.MkdirAll(uploadDir, 0700)
	destination := uploadDir + file.Filename
	if err := c.SaveFile(file, destination); err != nil {
		return c.Status(500).SendString("Failed to save file")
	}
	queries.SetCurrentInstructorsPagePhoto(strings.Replace(destination, os.Getenv("UPLOAD_DIR"), "/upload", 1))
	if c.Get("HX-Request") != "" {
		c.Set("HX-Redirect", "/instructors")
		return c.SendStatus(200)
	}
	return c.Redirect("/instructors")
}
