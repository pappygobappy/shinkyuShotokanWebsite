package handlers

import (
	"log"
	"shinkyuShotokan/models"
	"shinkyuShotokan/queries"
	"shinkyuShotokan/structs"
	"shinkyuShotokan/utils"
	"time"

	"github.com/gofiber/fiber/v2"
)

func AdminPromotionalsPage(c *fiber.Ctx) error {
	promotionals := queries.GetAllPromotionals()

	counts := make(map[uint]int)
	for _, p := range promotionals {
		count, err := queries.GetCheckedInCountByPromotionalID(p.ID)
		if err != nil {
			return err
		}
		counts[p.ID] = count
	}

	var user *models.User
	if u := c.Locals("user"); u != nil {
		currentUser, ok := u.(models.User)
		if ok && currentUser.Type != "" {
			user = &currentUser
		}
	}

	tabs := utils.CurrentTabs()

	page := fiber.Map{
		"Page": structs.Page{
			PageName: "Promotionals",
			Tabs:     tabs,
			Classes:  utils.Classes,
		},
		"Tabs":           tabs,
		"Promotionals":   promotionals,
		"CheckedInCount": counts,
		"user":           user,
	}

	return c.Render("adminPage", page)
}

func AddPromotional(c *fiber.Ctx) error {
	var body struct {
		Name        string
		Date        string
		StartTime   string
		EndTime     string
		CheckInTime string
		SubType     string
		Location    string
		Description string
	}

	if err := c.BodyParser(&body); err != nil {
		log.Print(err)
		return err
	}

	p := models.Promotional{
		Name:        body.Name,
		SubType:     body.SubType,
		Location:    body.Location,
		Description: body.Description,
	}

	// Parse date and time strings
	if body.Date != "" {
		p.Date, _ = parseDate(body.Date)
	}
	if body.StartTime != "" {
		p.StartTime, _ = parseTime(body.StartTime)
	}
	if body.EndTime != "" {
		p.EndTime, _ = parseTime(body.EndTime)
	}
	if body.CheckInTime != "" {
		p.CheckInTime, _ = parseTime(body.CheckInTime)
	}

	if err := queries.CreatePromotional(&p); err != nil {
		log.Print("Error creating Promotional", err)
		return err
	}

	return c.Redirect("/admin/promotionals")
}

func EditPromotionalGet(c *fiber.Ctx) error {
	id := c.Params("id")
	p := queries.GetPromotionalByID(parseUint(id))
	return c.Render("admin_promotional_edit", fiber.Map{
		"Promotional": p,
	})
}

func EditPromotionalPut(c *fiber.Ctx) error {
	id := c.Params("id")
	p := queries.GetPromotionalByID(parseUint(id))

	var body struct {
		Name        string
		Date        string
		StartTime   string
		EndTime     string
		CheckInTime string
		SubType     string
		Location    string
		Description string
	}

	if err := c.BodyParser(&body); err != nil {
		log.Print(err)
		return err
	}

	p.Name = body.Name
	p.SubType = body.SubType
	p.Location = body.Location
	p.Description = body.Description

	if body.Date != "" {
		p.Date, _ = parseDate(body.Date)
	}
	if body.StartTime != "" {
		p.StartTime, _ = parseTime(body.StartTime)
	}
	if body.EndTime != "" {
		p.EndTime, _ = parseTime(body.EndTime)
	}
	if body.CheckInTime != "" {
		p.CheckInTime, _ = parseTime(body.CheckInTime)
	}

	if err := queries.UpdatePromotional(&p); err != nil {
		log.Print("Error updating Promotional", err)
		return err
	}

	if c.Get("hx-request") != "" {
		c.Set("HX-Redirect", "/admin/promotionals")
		return c.SendStatus(200)
	}
	return c.Redirect("/admin/promotionals")
}

func DeletePromotional(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := queries.SoftDeletePromotional(parseUint(id)); err != nil {
		log.Print("Error deleting Promotional", err)
		return err
	}
	return c.Redirect("/admin/promotionals")
}

// parseUint is defined in handlers/gear.go — available in the same package.
// Keep it here only if gear.go's version differs. Check gear.go first.

func parseDate(s string) (time.Time, error) {
	return time.Parse("2006-01-02", s)
}

func parseTime(s string) (time.Time, error) {
	return time.Parse("15:04", s)
}