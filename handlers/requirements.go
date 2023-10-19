package handlers

import (
	"fmt"
	"shinkyuShotokan/structs"
	"shinkyuShotokan/utils"

	"github.com/gofiber/fiber/v2"
)

func Requirements(c *fiber.Ctx) error {
	rank := c.Params("rank")
	return c.Render(fmt.Sprintf("%s_requirements", rank), fiber.Map{"Page": structs.Page{PageName: "Requirements", Tabs: utils.CurrentTabs(), Classes: utils.Classes}})
}
