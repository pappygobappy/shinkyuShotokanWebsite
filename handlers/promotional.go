package handlers

import (
	"shinkyuShotokan/structs"
	"shinkyuShotokan/utils"
	"github.com/gofiber/fiber/v2"
)

func AdminPromotionalsPage(c *fiber.Ctx) error {
	page := fiber.Map{
		"Page":  structs.Page{PageName: "Promotionals", 
		Tabs: utils.CurrentTabs(), Classes: utils.Classes},
	}
	return c.Render("adminPage", page)
}