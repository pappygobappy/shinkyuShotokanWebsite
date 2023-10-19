package handlers

import (
	"shinkyuShotokan/structs"
	"shinkyuShotokan/utils"

	"github.com/gofiber/fiber/v2"
)

func ContactUs(c *fiber.Ctx) error {
	return c.Render("contact_us", fiber.Map{"Page": structs.Page{PageName: "Contact Us", Tabs: utils.CurrentTabs(), Classes: utils.Classes}})
}
