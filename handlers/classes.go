package handlers

import (
	"log"
	"shinkyuShotokan/initializers"
	"shinkyuShotokan/models"
	"shinkyuShotokan/queries"
	"shinkyuShotokan/structs"
	"shinkyuShotokan/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func Classes(c *fiber.Ctx) error {

	//Get selected Class
	class := queries.FindClassByPath(c.Path())

	//Build Page Data
	classesPage := fiber.Map{
		"Page":  structs.Page{PageName: "Classes", Tabs: utils.CurrentTabs(), Classes: utils.Classes},
		"Class": class,
	}

	//Render Page
	hxRequest, err := strconv.ParseBool(c.Get("hx-request"))
	if err != nil {
		hxRequest = false
	}

	if hxRequest {
		return c.Render("classPage", classesPage)
	} else {
		return c.Render("class", classesPage)
	}
}

func EditClassGet(c *fiber.Ctx) error {
	log.Println("Hello")
	class := queries.FindClassByID(c.Params("id"))

	return c.Render("edit_class_form", fiber.Map{
		"Class":     class,
		"Locations": queries.GetLocations(),
	})
}

func EditClassPut(c *fiber.Ctx) error {
	id := c.Params("id")
	var class models.Class
	initializers.DB.First(&class, id)

	var body struct {
		MinAge       int
		MaxAge    int
		Schedule string
		Description string
		Location string
		Annotations []string
	}

	if err := c.BodyParser(&body); err != nil {
		log.Print(err)
		return err
	}

	class.StartAge = body.MinAge
	class.EndAge = body.MaxAge
	class.Schedule = body.Schedule
	class.Description = body.Description
	class.LocationID = body.Location

	initializers.DB.Unscoped().Model(&class).Association("Annotations").Unscoped().Clear()

	var classAnnotations []models.ClassAnnotation

	for _, annotation := range body.Annotations {
		classAnnotation := models.ClassAnnotation{Annotation: annotation, ClassID: class.ID}
		// res := initializers.DB.Save(&classAnnotation)

		// if res.Error != nil {
		// 	log.Print("Error updating Class", res.Error)
		// 	return res.Error
		// }
		classAnnotations = append(classAnnotations, classAnnotation)
	}

	class.Annotations = classAnnotations;

	result := initializers.DB.Save(&class)

	if result.Error != nil {
		log.Print("Error updating Class", result.Error)
		return result.Error
	}

	class = queries.FindClassByID(id)

	classesPage := fiber.Map{
		"Page":  structs.Page{PageName: "Classes", Tabs: utils.CurrentTabs(), Classes: utils.Classes},
		"Class": class,
	}

	return c.Render("class", classesPage)
}

// func findClassByPath(path string) models.Class {
// 	for _, class := range utils.Classes {
// 		if class.GetUrl == path {
// 			return class
// 		}
// 	}
// 	return models.Class{}
// }
