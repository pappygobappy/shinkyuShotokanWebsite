package handlers

import (
	"log"
	"strconv"
	"shinkyuShotokan/models"
	"shinkyuShotokan/queries"
	"shinkyuShotokan/structs"
	"shinkyuShotokan/utils"

	"github.com/gofiber/fiber/v2"
)

func GearPage(c *fiber.Ctx) error {
	items := queries.GetGearItems()
	var user *models.User
	if u := c.Locals("user"); u != nil {
		u := u.(models.User)
		if u.Type != "" {
			user = &u
		}
	}
	page := fiber.Map{
		"Page":  structs.Page{PageName: "Gear", Tabs: utils.CurrentTabs(), Classes: utils.Classes},
		"Items": items,
		"user":  user,
	}
	return c.Render("gear", page)
}

func AdminGearPage(c *fiber.Ctx) error {
	items := queries.GetGearItems()
	page := fiber.Map{
		"Page":  structs.Page{PageName: "Gear", Tabs: utils.CurrentTabs(), Classes: utils.Classes},
		"Items": items,
	}
	return c.Render("adminPage", page)
}

func AddGearItem(c *fiber.Ctx) error {
	var body struct {
		Name        string
		Link        string
		ShopName   string
		Description string
	}

	if err := c.BodyParser(&body); err != nil {
		log.Print(err)
		return err
	}

	item := models.GearItem{
		Name:        body.Name,
		Link:        body.Link,
		ShopName:   body.ShopName,
		Description: body.Description,
	}
	result := queries.CreateGearItem(&item)
	if result != nil {
		log.Print("Error creating GearItem", result)
		return result
	}

	return c.Redirect("/gear")
}

func EditGearItemGet(c *fiber.Ctx) error {
	id := c.Params("id")
	item := queries.GetGearItemById(parseUint(id))
	return c.Render("edit_gear_form", fiber.Map{"Item": item})
}

func GearItemRow(c *fiber.Ctx) error {
	id := c.Params("id")
	item := queries.GetGearItemById(parseUint(id))
	var user *models.User
	if u := c.Locals("user"); u != nil {
		u := u.(models.User)
		if u.Type != "" {
			user = &u
		}
	}
return c.Render("gear_item_id", fiber.Map{"Item": item, "user": user})
}

func GearItemDisplay(c *fiber.Ctx) error {
	id := c.Params("id")
	item := queries.GetGearItemById(parseUint(id))
	return c.Render("gear_item_id", fiber.Map{"Item": item})
}

func MoveGearItem(c *fiber.Ctx) error {
	idParam := c.Params("id")
	idUint64, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid id")
	}
	direction := c.Query("direction")
	if direction != "up" && direction != "down" {
		return fiber.NewError(fiber.StatusBadRequest, "invalid direction")
	}
	if err := queries.MoveGearItem(uint(idUint64), direction); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return GearPage(c)
}

func EditGearItemPut(c *fiber.Ctx) error {
	id := c.Params("id")
	item := queries.GetGearItemById(parseUint(id))

var body struct {
		Name        string
		Link        string
		ShopName   string
		Description string
	}

	if err := c.BodyParser(&body); err != nil {
		log.Print(err)
		return err
	}

	item.Name = body.Name
	item.Link = body.Link
	item.ShopName = body.ShopName
	item.Description = body.Description

	result := queries.UpdateGearItem(&item)
	if result != nil {
		log.Print("Error updating GearItem", result)
		return result
	}

	if c.Get("hx-request") != "" {
		c.Set("HX-Redirect", "/gear")
		return c.SendStatus(200)
	}
	return c.Redirect("/gear")
}

func DeleteGearItem(c *fiber.Ctx) error {
	id := c.Params("id")
	result := queries.DeleteGearItem(parseUint(id))
	if result != nil {
		log.Print("Error deleting GearItem", result)
		return result
	}
	return c.Redirect("/gear")
}

func parseUint(id string) uint {
	var n uint
	for _, c := range id {
		if c >= '0' && c <= '9' {
			n = n*10 + uint(c-'0')
		}
	}
	return n
}
