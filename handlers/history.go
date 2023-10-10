package handlers

import (
	"shinkyuShotokan/structs"
	"shinkyuShotokan/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func History(c *fiber.Ctx) error {
	historyPage := structs.Page{PageName: "History", Tabs: utils.Tabs, Classes: utils.Classes}
	hxRequest, err := strconv.ParseBool(c.Get("hx-request"))
	if err != nil {
		hxRequest = false
	}

	if hxRequest == true {
		return c.Render("historyPage", fiber.Map{
			"Page": historyPage,
		})
	} else {
		return c.Render("history", fiber.Map{
			"Page": historyPage,
		})
	}
}
