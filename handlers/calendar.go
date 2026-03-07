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
	Day               time.Time
	NotInCurrentMonth bool
	Events            []structs.CalendarItem
}

func Calendar(c *fiber.Ctx) error {
	calendarPage := structs.Page{PageName: "Calendar", Tabs: utils.CurrentTabs(), Classes: utils.Classes}

	hxRequest, _ := strconv.ParseBool(c.Get("hx-request"))

	monthStr := c.Query("month")
	if monthStr == "" {
		monthStr = time.Now().Format("2006-01")
	}
	month, _ := time.Parse("2006-01", monthStr)

	filteredClass := c.Query("class")

	result := buildCalendarView(month, filteredClass)

	calendarMap := fiber.Map{
		"Page":          calendarPage,
		"Month":         result.Month,
		"Today":         result.Today,
		"Weeks":         result.Weeks,
		"PrevMonth":     result.PrevMonth,
		"NextMonth":     result.NextMonth,
		"Locations":     result.Locations,
		"Classes":       utils.ActualClasses,
		"Periods":       result.Periods,
		"FilteredClass": result.FilteredClass,
	}

	if hxRequest {
		return c.Render("calendarPage2", calendarMap)
	}
	return c.Render("calendar", calendarMap)
}

func buildCalendarView(month time.Time, filteredClass string) structs.CalendarViewResult {
	currentLocation, _ := time.LoadLocation("America/Los_Angeles")

	currentYear, currentMonth, _ := month.Date()
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

	weeks := buildWeeks(events, classSessions, firstOfMonth, startDate, endDate)

	var periods []models.ClassPeriod
	result := initializers.DB.Find(&periods)
	if result.Error != nil {
		log.Print(result.Error)
	}

	return structs.CalendarViewResult{
		Weeks:         weeks,
		Month:         firstOfMonth,
		Today:         time.Now().In(utils.TZ),
		PrevMonth:     month.AddDate(0, -1, 0).Format("2006-01"),
		NextMonth:     month.AddDate(0, 1, 0).Format("2006-01"),
		Locations:     queries.GetLocations(),
		Classes:       utils.ActualClasses,
		Periods:       periods,
		FilteredClass: filteredClass,
	}
}

func buildWeeks(events []models.Event, classSessions []models.ClassSession, firstOfMonth, startDate, endDate time.Time) []structs.CalendarWeek {
	var weeks []structs.CalendarWeek
	weeks = append(weeks, structs.CalendarWeek{})

	for i := startDate; i.Before(endDate) || i.Equal(endDate); i = i.AddDate(0, 0, 1) {
		dayOfWeek := i.Weekday().String()
		var dayMap structs.CalendarDay
		var eventSlices []structs.CalendarItem
		dayMap.Day = i.In(utils.TZ)
		dayMap.NotInCurrentMonth = i.Month() != firstOfMonth.Month()

		for _, e := range events {
			if e.Date.Equal(i) {
				eventSlices = append(eventSlices, structs.CalendarItem{
					StartTime: e.StartTime, Title: e.Title, Color: "red", Location: e.Location, Url: fmt.Sprintf("/events/%d", e.ID),
				})
			}
		}

		for _, class := range classSessions {
			startYear := class.StartTime.In(utils.TZ).Year()
			startMonth := class.StartTime.In(utils.TZ).Month()
			startDay := class.StartTime.In(utils.TZ).Day()
			if startYear == i.Year() && startMonth == i.Month() && startDay == i.Day() {
				actualClass := utils.FindActualClassByName(class.ClassName)
				eventSlices = append(eventSlices, structs.CalendarItem{
					StartTime: class.StartTime, Title: class.ClassName, Color: actualClass.Color, Location: class.Location, Url: fmt.Sprintf("/calendar/%d", class.ID), IsCancelled: class.IsCancelled,
				})
			}
		}

		sort.Slice(eventSlices, func(i, j int) bool {
			if eventSlices[i].StartTime.Equal(eventSlices[j].StartTime) {
				return len(eventSlices[i].Title) > 0 && len(eventSlices[j].Title) > 0 && eventSlices[i].Title[0] < eventSlices[j].Title[0]
			}
			return eventSlices[i].StartTime.Before(eventSlices[j].StartTime)
		})

		dayMap.Events = eventSlices
		weeks[len(weeks)-1][dayOfWeek] = dayMap

		if i.Weekday() == time.Saturday && !i.Equal(endDate) {
			weeks = append(weeks, structs.CalendarWeek{})
		}
	}

	return weeks
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
		"Location":     queries.GetLocationByName(classSession.Location),
	}

	if hxRequest {
		return c.Render("classSessionViewPage", calendarItemMap)
	} else {
		return c.Render("classSessionView", calendarItemMap)
	}
}

func AddClassSessionForm(c *fiber.Ctx) error {
	var periods []models.ClassPeriod
	result := initializers.DB.Where("end_date > ?", time.Now()).Order("start_date").Limit(4).Find(&periods)
	if result.Error != nil {
		log.Print(result.Error)
	}

	return c.Render("addClassSession", fiber.Map{"Locations": queries.GetLocations(), "Classes": utils.ActualClasses, "Periods": periods})
}

func GetDeleteClassSessionsForm(c *fiber.Ctx) error {
	var periods []models.ClassPeriod
	result := initializers.DB.Order("start_date desc").Find(&periods)
	if result.Error != nil {
		log.Print(result.Error)
	}

	return c.Render("deleteClassSessions", fiber.Map{"Classes": utils.ActualClasses, "Periods": periods})
}

func AddClassSession(c *fiber.Ctx) error {
	var body struct {
		Class     string
		Period    string
		StartDate string
		StartTime string
		EndTime   string
		EndDate   string
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
	startDate, error := time.ParseInLocation("2006-01-02", body.StartDate, utils.TZ)
	endDate, error := time.ParseInLocation("2006-01-02", body.EndDate, utils.TZ)
	startTime, error := time.ParseInLocation("15:04", body.StartTime, utils.TZ)
	endTime, error := time.ParseInLocation("15:04", body.EndTime, utils.TZ)

	if error != nil {
		fmt.Println(error)
		return c.Redirect("/")
	}

	//event := models.Event{Title: body.Name, Date: date, StartTime: startTime, EndTime: endTime, Description: body.Description, PictureUrl: photoUrl, Location: body.Location}

	//classPeriod := queries.GetClassPeriodById(body.Period)

	for i := startDate; i.Before(endDate) || i.Equal(endDate); i = i.AddDate(0, 0, 1) {
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

func DeleteClassPeriodSessions(c *fiber.Ctx) error {
	var body struct {
		Class  string
		Period string
	}

	if err := c.BodyParser(&body); err != nil {
		log.Print(err)
		return err
	}

	queries.DeleteClassSessionsByClassAndClassPeriod(body.Class, body.Period)
	var period models.ClassPeriod = queries.GetClassPeriodById(body.Period)
	log.Print("Deleted ", body.Class, " ", period.Name)
	c.Set("HX-Redirect", "/calendar")
	return c.Next()
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
		Class       string
		Date        string
		StartTime   string
		EndTime     string
		Location    string
		IsCancelled bool
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
	classSession.IsCancelled = body.IsCancelled

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
	var session models.ClassSession
	initializers.DB.Delete(&session, id)

	month := session.StartTime.Format("2006-01")
	c.Set("HX-Redirect", fmt.Sprintf("/calendar?month=%s", month))
	return c.Next()
}
