package handlers

import (
	"shinkyuShotokan/structs"
	"shinkyuShotokan/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func Classes(c *fiber.Ctx) error {

	//Get selected Class
	class := findClass(c.Path())

	//Build Page Data
	classesPage := fiber.Map{
		"Page": structs.Page{PageName: "Classes", Tabs: utils.Tabs, Classes: utils.Classes},
		"Class": class,
	}

	//Render Page
	hxRequest, err := strconv.ParseBool(c.Get("hx-request"))
	if err != nil {
		hxRequest = false
	}

	if hxRequest == true {
		return c.Render("classPage", classesPage)
	} else {
		return c.Render("class", classesPage)
	}
}

func findClass(path string) structs.Class {
	for _, class := range utils.Classes {
		if class.GetUrl == path {
			return class
		}
	}
	return structs.Class{};
}
