package handlers

import (
	"fmt"
	"os"
	"shinkyuShotokan/models"
	"shinkyuShotokan/queries"
	"shinkyuShotokan/structs"
	"shinkyuShotokan/utils"
	"strconv"

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
	instructorsPage := fiber.Map{
		"Page":              structs.Page{PageName: "Instructors", Tabs: utils.CurrentTabs(), Classes: utils.Classes},
		"Instructors":       visible,
		"HiddenInstructors": hidden,
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

	var body struct {
		Name string
		Bio  string
	}
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid form")
	}

	// Handle optional new picture
	if file, err := c.FormFile("NewPicture"); err == nil && file != nil {
		baseDir := fmt.Sprintf("%s/assets/instructors", os.Getenv("UPLOAD_DIR"))
		os.MkdirAll(baseDir, 0700)
		if err := c.SaveFile(file, fmt.Sprintf("%s/%s", baseDir, file.Filename)); err == nil {
			instructor.PictureUrl = fmt.Sprintf("/upload/assets/instructors/%s", file.Filename)
		}
	}

	instructor.Name = body.Name
	instructor.Bio = body.Bio
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
