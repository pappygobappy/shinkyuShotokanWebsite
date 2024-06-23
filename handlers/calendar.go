package handlers

import (
	"fmt"
	"log"
	"shinkyuShotokan/initializers"
	"shinkyuShotokan/models"
	"shinkyuShotokan/queries"
	"shinkyuShotokan/structs"
	"shinkyuShotokan/utils"
	"sort"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type CalendarDay struct {
	Day            time.Time
	NotInCurrentMonth bool
	Events         []structs.CalendarItem
}

func Calendar(c *fiber.Ctx) error {

	calendarPage := structs.Page{PageName: "Calendar", Tabs: utils.CurrentTabs(), Classes: utils.Classes}

	hxRequest, err := strconv.ParseBool(c.Get("hx-request"))
	if err != nil {
		hxRequest = false
	}

	//Get current time
	month := time.Now()
	if c.Query("month") != "" {
		month, _ = time.Parse("2006-01", c.Query("month"))
	}

	filteredClass := c.Query("class")

	prevMonth := month.AddDate(0, -1, 0)
	nextMonth := month.AddDate(0, 1, 0)

	currentYear, currentMonth, _ := month.Date()
	currentLocation, _ := time.LoadLocation("America/Los_Angeles")

	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)
	startDate := firstOfMonth
	endDate := lastOfMonth

	if int(startDate.Weekday()) > 0 {
		startDate = firstOfMonth.AddDate(0, 0, -int(startDate.Weekday()))
	}
	if int(endDate.Weekday()) < 6 {
		endDate = endDate.AddDate(0, 0, 6-int(endDate.Weekday()))
	}

	events := queries.GetEventsBetweenDates(startDate, endDate)

	var classSessions []models.ClassSession
	if filteredClass != "" {
		classSessions = queries.GetClassSessionsByClassAndBetweenDates(filteredClass, startDate, endDate)
	} else {
		classSessions = queries.GetClassSessionsBetweenDates(startDate, endDate)
	}

	weeks := []fiber.Map{}
	weeks = append(weeks, emptyWeekMap())
	weekIndex := 0
	for i := startDate; i.Before(endDate) || i.Equal(endDate); i = i.AddDate(0, 0, 1) {
		dayOfWeek := i.Weekday()
		var dayMap CalendarDay
		var eventSlices []structs.CalendarItem
		dayMap.Day = i.In(utils.TZ)
		dayMap.NotInCurrentMonth = i.Month() != firstOfMonth.Month()
		for _, e := range events {
			if e.Date.Equal(i) {
				eventSlices = append(eventSlices, structs.CalendarItem{StartTime: e.StartTime, Title: e.Title, Color: "red", Location: e.Location, Url: fmt.Sprintf("%s%d", "/events/", e.ID)})
			}
		}
		for _, class := range classSessions {
			startYear := class.StartTime.In(utils.TZ).Year()
			startMonth := class.StartTime.In(utils.TZ).Month()
			startDay := class.StartTime.In(utils.TZ).Day()
			if startYear == i.Year() && startMonth == i.Month() && startDay == i.Day() {
				eventSlices = append(eventSlices, structs.CalendarItem{StartTime: class.StartTime, Title: class.ClassName, Color: utils.FindActualClassByName(class.ClassName).Color, Location: class.Location, Url: fmt.Sprintf("%s%d", "/calendar/", class.ID)})
			}
		}
		sort.Slice(eventSlices, func(i, j int) bool {
			if eventSlices[i].StartTime.Equal(eventSlices[j].StartTime) {
				return eventSlices[i].Title[0] < eventSlices[j].Title[0]
			}
			return eventSlices[i].StartTime.Before(eventSlices[j].StartTime)
		})

		dayMap.Events = eventSlices
		weeks[weekIndex][dayOfWeek.String()] = dayMap
		if dayOfWeek == time.Saturday && !i.Equal(endDate) {
			weeks = append(weeks, emptyWeekMap())
			weekIndex += 1
		}
	}
	var periods []models.ClassPeriod
	result := initializers.DB.Find(&periods)
	if result.Error != nil {
		log.Print(result.Error)
	}

	calendarMap := fiber.Map{"Page": calendarPage, "Month": firstOfMonth, "Today": time.Now().In(utils.TZ), "Weeks": weeks, "PrevMonth": prevMonth.Format("2006-01"), "NextMonth": nextMonth.Format("2006-01"), "Locations": queries.GetLocations(), "Classes": utils.ActualClasses, "Periods": periods, "FilteredClass": filteredClass}

	if hxRequest {
		return c.Render("calendarPage2", calendarMap)
	} else {
		return c.Render("calendar", calendarMap)
	}
}

func CalendarItemView(c *fiber.Ctx) error {

	hxRequest, err := strconv.ParseBool(c.Get("hx-request"))
	if err != nil {
		hxRequest = false
	}

	id := c.Params("id")
	var classSession models.ClassSession
	initializers.DB.First(&classSession, id)

	page := structs.Page{PageName: "Calendar", Tabs: utils.CurrentTabs(), Classes: utils.Classes}

	class := queries.FindClassByName(utils.FindActualClassByName(classSession.ClassName).Class)

	calendarItemMap := fiber.Map{
		"Page":         page,
		"ClassSession": classSession,
		"Period":       queries.GetClassPeriodById(classSession.Period),
		"Class":        class,
		"Location":     queries.GetLocationyName(classSession.Location),
	}

	if hxRequest {
		return c.Render("classSessionViewPage", calendarItemMap)
	} else {
		return c.Render("classSessionView", calendarItemMap)
	}
}

func AddClassSessionForm(c *fiber.Ctx) error {
	var periods []models.ClassPeriod
	result := initializers.DB.Find(&periods)
	if result.Error != nil {
		log.Print(result.Error)
	}

	return c.Render("addClassSession", fiber.Map{"Locations": queries.GetLocations(), "Classes": utils.ActualClasses, "Periods": periods})
}

func AddClassSession(c *fiber.Ctx) error {
	var body struct {
		Class     string
		Period    string
		StartTime string
		EndTime   string
		Sunday    bool
		Monday    bool
		Tuesday   bool
		Wednesday bool
		Thursday  bool
		Friday    bool
		Saturday  bool
		Location  string
	}

	if err := c.BodyParser(&body); err != nil {
		log.Print(err)
		return err
	}

	selectedWeekdaysMap := map[string]bool{"Sunday": body.Sunday, "Monday": body.Monday, "Tuesday": body.Tuesday, "Wednesday": body.Wednesday, "Thursday": body.Thursday, "Friday": body.Friday, "Saturday": body.Saturday}

	startTime, error := time.ParseInLocation("15:04", body.StartTime, utils.TZ)
	endTime, error := time.ParseInLocation("15:04", body.EndTime, utils.TZ)

	if error != nil {
		fmt.Println(error)
		return c.Redirect("/")
	}

	//event := models.Event{Title: body.Name, Date: date, StartTime: startTime, EndTime: endTime, Description: body.Description, PictureUrl: photoUrl, Location: body.Location}

	classPeriod := queries.GetClassPeriodById(body.Period)

	for i := classPeriod.StartDate; i.Before(classPeriod.EndDate) || i.Equal(classPeriod.EndDate); i = i.AddDate(0, 0, 1) {
		if selectedWeekdaysMap[i.Weekday().String()] {
			classStartTime := time.Date(i.Year(), i.Month(), i.Day(), startTime.Hour(), startTime.Minute(), 0, 0, utils.TZ)
			classEndTime := time.Date(i.Year(), i.Month(), i.Day(), endTime.Hour(), endTime.Minute(), 0, 0, utils.TZ)
			classSession := models.ClassSession{ClassName: body.Class, Period: body.Period, StartTime: classStartTime, EndTime: classEndTime, Location: body.Location}
			result := initializers.DB.Create(&classSession)

			if result.Error != nil {
				log.Print("Error creating Class Session", result.Error)
				return result.Error
			}
		}
	}

	return c.Redirect("/calendar")
}

func EditClassSessionGet(c *fiber.Ctx) error {
	id := c.Params("id")
	var classSession models.ClassSession
	initializers.DB.First(&classSession, id)

	return c.Render("edit_classSession", fiber.Map{"ClassSession": classSession, "Locations": queries.GetLocations(), "Classes": utils.ActualClasses})
}

func EditClassSessionPut(c *fiber.Ctx) error {
	id := c.Params("id")
	var classSession models.ClassSession
	initializers.DB.First(&classSession, id)

	var body struct {
		Class     string
		Date      string
		StartTime string
		EndTime   string
		Location  string
	}

	if err := c.BodyParser(&body); err != nil {
		log.Print(err)
		return err
	}

	date, error := time.ParseInLocation("2006-01-02", body.Date, utils.TZ)
	startTime, error := time.ParseInLocation("15:04", body.StartTime, utils.TZ)
	endTime, error := time.ParseInLocation("15:04", body.EndTime, utils.TZ)

	if error != nil {
		fmt.Println(error)
		return c.Redirect("/")
	}

	classStartTime := time.Date(date.Year(), date.Month(), date.Day(), startTime.Hour(), startTime.Minute(), 0, 0, utils.TZ)
	classEndTime := time.Date(date.Year(), date.Month(), date.Day(), endTime.Hour(), endTime.Minute(), 0, 0, utils.TZ)

	classSession.ClassName = body.Class
	classSession.StartTime = classStartTime
	classSession.EndTime = classEndTime
	classSession.Location = body.Location

	result := initializers.DB.Save(&classSession)

	if result.Error != nil {
		log.Print("Error saving ClassSession", result.Error)
		return result.Error
	}

	//return c.Redirect("/calendar")
	c.Set("HX-Redirect", "/calendar/"+strconv.FormatUint(uint64(classSession.ID), 10))
	return c.Next()
}

func AddClassPeriodGet(c *fiber.Ctx) error {
	return c.Render("add_classPeriod", nil)
}

func AddClassPeriodPost(c *fiber.Ctx) error {

	var body struct {
		Name      string
		StartDate string
		EndDate   string
	}

	if err := c.BodyParser(&body); err != nil {
		log.Print(err)
		return err
	}

	startDate, error := time.ParseInLocation("2006-01-02", body.StartDate, utils.TZ)
	endDate, error := time.ParseInLocation("2006-01-02", body.EndDate, utils.TZ)

	if error != nil {
		fmt.Println(error)
		return c.Redirect("/")
	}

	classPeriod := models.ClassPeriod{Name: body.Name, StartDate: startDate, EndDate: endDate}

	result := initializers.DB.Create(&classPeriod)

	if result.Error != nil {
		log.Print("Error creating Class Session", result.Error)
		return result.Error
	}

	return c.
		Redirect("/calendar")
}

func EditClassPeriodGet(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Println(id)
	var classPeriod models.ClassPeriod
	initializers.DB.First(&classPeriod, id)
	log.Println(classPeriod)
	return c.Render("edit_classPeriod", classPeriod)
}

func EditClassPeriodPut(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Println(id)
	var classPeriod models.ClassPeriod
	initializers.DB.First(&classPeriod, id)

	var body struct {
		Name      string
		StartDate string
		EndDate   string
	}

	if err := c.BodyParser(&body); err != nil {
		log.Print(err)
		return err
	}

	startDate, error := time.ParseInLocation("2006-01-02", body.StartDate, utils.TZ)
	endDate, error := time.ParseInLocation("2006-01-02", body.EndDate, utils.TZ)

	if error != nil {
		fmt.Println(error)
		return c.Redirect("/")
	}

	classPeriod.Name = body.Name
	classPeriod.StartDate = startDate
	classPeriod.EndDate = endDate

	result := initializers.DB.Save(&classPeriod)

	if result.Error != nil {
		log.Print("Error saving ClassSession", result.Error)
		return result.Error
	}

	c.Set("HX-Redirect", "/calendar")
	return c.Next()
}

func DeleteClassPeriod(c *fiber.Ctx) error {
	id := c.Params("id")
	initializers.DB.Delete(&models.ClassPeriod{}, id)

	c.Set("HX-Redirect", "/calendar")
	return c.Next()
}

func DeleteClassSession(c *fiber.Ctx) error {
	id := c.Params("id")
	initializers.DB.Delete(&models.ClassSession{}, id)

	c.Set("HX-Redirect", "/")
	return c.Next()
}

func emptyWeekMap() fiber.Map {
	weekMap := fiber.Map{}
	weekMap[time.Sunday.String()] = CalendarDay{}
	weekMap[time.Monday.String()] = CalendarDay{}
	weekMap[time.Tuesday.String()] = CalendarDay{}
	weekMap[time.Wednesday.String()] = CalendarDay{}
	weekMap[time.Thursday.String()] = CalendarDay{}
	weekMap[time.Friday.String()] = CalendarDay{}
	weekMap[time.Saturday.String()] = CalendarDay{}
	return weekMap
}
