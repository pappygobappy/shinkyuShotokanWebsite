package handlers

import (
	"fmt"
	"log"
	"shinkyuShotokan/initializers"
	"shinkyuShotokan/models"
	"shinkyuShotokan/queries"
	"shinkyuShotokan/structs"
	"shinkyuShotokan/utils"
	"time"

	"github.com/gofiber/fiber/v2"
)

func AdminEventTemplatesPage(c *fiber.Ctx) error {
	persistedLocations := queries.GetLocations()
	eventTemplates := queries.GetEventTemplates()
	fmt.Println(eventTemplates)
	adminPage := fiber.Map{
		"Page":           structs.Page{PageName: "Event Templates", Tabs: utils.CurrentTabs(), Classes: utils.Classes},
		"Locations":      persistedLocations,
		"EventTemplates": eventTemplates,
	}
	fmt.Println(adminPage["Page"].(structs.Page).PageName)
	return c.Render("adminPage", adminPage)
}

func AddEventTemplate(c *fiber.Ctx) error {
	var body struct {
		Name        string
		StartTime   string
		EndTime     string
		CheckInTime string
		Description string
		LocationID  string
	}

	if err := c.BodyParser(&body); err != nil {
		log.Print(err)
		return err
	}

	// Parse times using today's date
	today := time.Now().Format("2006-01-02")
	startTime, err := time.ParseInLocation("2006-01-02 15:04", today+" "+body.StartTime, utils.TZ)
	if err != nil {
		return err
	}
	endTime, err := time.ParseInLocation("2006-01-02 15:04", today+" "+body.EndTime, utils.TZ)
	if err != nil {
		return err
	}
	checkInTime, err := time.ParseInLocation("2006-01-02 15:04", today+" "+body.CheckInTime, utils.TZ)
	if err != nil {
		return err
	}

	template := models.EventTemplate{
		Name:        body.Name,
		StartTime:   startTime,
		EndTime:     endTime,
		CheckInTime: checkInTime,
		Description: body.Description,
		LocationID:  body.LocationID,
	}

	result := initializers.DB.Create(&template)
	if result.Error != nil {
		log.Print("Error creating Event Template", result.Error)
		return result.Error
	}

	return c.Redirect("/admin/event-templates")
}

func EditEventTemplateGet(c *fiber.Ctx) error {
	id := c.Params("id")
	template := queries.GetEventTemplateByID(id)

	return c.Render("edit_event_template_form", fiber.Map{
		"Template":  template,
		"Locations": queries.GetLocations(),
	})
}

func EditEventTemplatePut(c *fiber.Ctx) error {
	id := c.Params("id")
	var template models.EventTemplate
	initializers.DB.First(&template, id)

	var body struct {
		Name        string
		StartTime   string
		EndTime     string
		CheckInTime string
		Description string
		LocationID  string
	}

	if err := c.BodyParser(&body); err != nil {
		log.Print(err)
		return err
	}

	// Parse times using today's date
	today := time.Now().Format("2006-01-02")
	startTime, err := time.ParseInLocation("2006-01-02 15:04", today+" "+body.StartTime, utils.TZ)
	if err != nil {
		return err
	}
	endTime, err := time.ParseInLocation("2006-01-02 15:04", today+" "+body.EndTime, utils.TZ)
	if err != nil {
		return err
	}
	checkInTime, err := time.ParseInLocation("2006-01-02 15:04", today+" "+body.CheckInTime, utils.TZ)
	if err != nil {
		return err
	}

	template.Name = body.Name
	template.Description = body.Description
	template.StartTime = startTime
	template.EndTime = endTime
	template.CheckInTime = checkInTime
	template.LocationID = body.LocationID

	result := initializers.DB.Save(&template)
	if result.Error != nil {
		log.Print("Error updating Event Template", result.Error)
		return result.Error
	}

	return c.Render("event_template", template)
}
