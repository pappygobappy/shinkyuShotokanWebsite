package handlers

import (
	"shinkyuShotokan/structs"
	"shinkyuShotokan/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func Instructors(c *fiber.Ctx) error {
	instructorsPage := fiber.Map{
		"Page": structs.Page{PageName: "Instructors", Tabs: utils.Tabs, Classes: utils.Classes},
		"Instructors": utils.Instructors,
	}
	hxRequest, err := strconv.ParseBool(c.Get("hx-request"))
	if err != nil {
		hxRequest = false
	}

	if hxRequest == true {
		return c.Render("instructorsPage", instructorsPage)
	} else {
		return c.Render("instructors", instructorsPage)
	}
}
