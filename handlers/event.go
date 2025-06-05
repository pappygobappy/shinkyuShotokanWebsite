package handlers

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"shinkyuShotokan/initializers"
	"shinkyuShotokan/models"
	"shinkyuShotokan/queries"
	"shinkyuShotokan/structs"
	"shinkyuShotokan/utils"
	"strconv"
	"strings"
	"time"

	ics "github.com/arran4/golang-ical"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func getExistingEventCoverPhotos() []string {
	workingDir, _ := os.Getwd()
	var eventImagePaths []string
	err := filepath.Walk(workingDir+"/public/events/", func(path string, info os.FileInfo, err error) error {

		if err != nil {
			fmt.Println(err)
			return nil
		}
		if !info.IsDir() {
			eventImagePaths = append(eventImagePaths, strings.Replace(path, workingDir, "", 1))
		}
		return nil
	})

	if err != nil {
		log.Println("Failed to get existing cover photos")
	}

	uploadedCoverPhotosPath := fmt.Sprintf("%s/assets/events/", os.Getenv("UPLOAD_DIR"))

	err = filepath.Walk(uploadedCoverPhotosPath, func(path string, info os.FileInfo, err error) error {

		if err != nil {
			fmt.Println(err)
			return nil
		}
		if !info.IsDir() {
			eventImagePaths = append(eventImagePaths, strings.Replace(path, os.Getenv("UPLOAD_DIR"), "/upload", 1))
		}
		return nil
	})

	if err != nil {
		log.Println("Failed to get uploaded cover photos")
	}

	return eventImagePaths
}

func getExistingEventCardPhotos() []string {
	workingDir, _ := os.Getwd()
	var eventImagePaths []string

	err := filepath.Walk(workingDir+"/public/events/cards/", func(path string, info os.FileInfo, err error) error {

		if err != nil {
			fmt.Println(err)
			return nil
		}
		if !info.IsDir() {
			eventImagePaths = append(eventImagePaths, strings.Replace(path, workingDir, "", 1))
		}
		return nil
	})

	if err != nil {
		log.Println("Failed to get existing card photos")
	}

	uploadedCoverPhotosPath := fmt.Sprintf("%s/assets/events/cards", os.Getenv("UPLOAD_DIR"))

	err = filepath.Walk(uploadedCoverPhotosPath, func(path string, info os.FileInfo, err error) error {

		if err != nil {
			fmt.Println(err)
			return nil
		}
		if !info.IsDir() {
			eventImagePaths = append(eventImagePaths, strings.Replace(path, os.Getenv("UPLOAD_DIR"), "/upload", 1))
		}
		return nil
	})

	if err != nil {
		log.Println("Failed to get uploaded card photos")
	}

	return eventImagePaths
}

func uploadEventFiles(event models.Event, c *fiber.Ctx) {
	if form, err := c.MultipartForm(); err == nil {
		// Get all files from "Files" key:
		files := form.File["Files"]
		// => []*multipart.FileHeader
		os.MkdirAll(fmt.Sprintf("%s/assets/event/%s/files", os.Getenv("UPLOAD_DIR"), strconv.FormatUint(uint64(event.ID), 10)), 0700)
		// Loop through files:
		for _, file := range files {
			// Save the files to disk:
			if err := c.SaveFile(file, fmt.Sprintf("%s/assets/event/%s/files/%s", os.Getenv("UPLOAD_DIR"), strconv.FormatUint(uint64(event.ID), 10), file.Filename)); err != nil {
				log.Println(err)
			}
		}
	}
}

func getEventFilePaths(event models.Event) map[string]string {
	basePath := fmt.Sprintf("%s/assets/event/%s/files/", os.Getenv("UPLOAD_DIR"), strconv.FormatUint(uint64(event.ID), 10))
	//key filename, value path
	var files = make(map[string]string)
	err := filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {

		if err != nil {
			fmt.Println(err)
			return nil
		}
		if !info.IsDir() {
			files[info.Name()] = strings.Replace(path, os.Getenv("UPLOAD_DIR"), "/upload", 1)
			//files = append(files, strings.Replace(path, paths, "", 1))
		}
		return nil
	})
	if err != nil {
		log.Println(err)
	}
	return files
}

func getCoverPhotoUrl(existingCoverPhoto string, c *fiber.Ctx) string {
	newCoverPhoto, err := c.FormFile("NewCoverPhoto")
	if err != nil {
		fmt.Println("No new cover photo")
		return existingCoverPhoto
	} else {
		os.MkdirAll(fmt.Sprintf("%s/assets/events/", os.Getenv("UPLOAD_DIR")), 0700)
		c.SaveFile(newCoverPhoto, fmt.Sprintf("%s/assets/events/%s", os.Getenv("UPLOAD_DIR"), newCoverPhoto.Filename))
		return fmt.Sprintf("/upload/assets/events/%s", newCoverPhoto.Filename)
	}
}

func getCardPhotoUrl(existingCardPhoto string, c *fiber.Ctx) string {
	newCardPhoto, err := c.FormFile("NewCardPhoto")
	if err != nil {
		fmt.Println("No new card photo")
		return existingCardPhoto
	} else {
		os.MkdirAll(fmt.Sprintf("%s/assets/events/cards", os.Getenv("UPLOAD_DIR")), 0700)
		c.SaveFile(newCardPhoto, fmt.Sprintf("%s/assets/events/cards/%s", os.Getenv("UPLOAD_DIR"), newCardPhoto.Filename))
		return fmt.Sprintf("/upload/assets/events/cards/%s", newCardPhoto.Filename)
	}
}

func Event(c *fiber.Ctx) error {
	id := c.Params("id")
	var event models.Event
	initializers.DB.First(&event, id)
	pageName := "Upcoming Events"
	if event.Date.Before(time.Now()) {
		pageName = "Past Events"
	}
	page := structs.Page{PageName: pageName, Tabs: utils.CurrentTabs(), Classes: utils.Classes}
	files := getEventFilePaths(event)

	return c.Render("event", fiber.Map{
		"Page":        page,
		"Event":       event,
		"Description": event.Description,
		"Files":       files,
		"Location":    queries.GetLocationByName(event.Location),
	})
}

func EditEventGet(c *fiber.Ctx) error {
	id := c.Params("id")
	var event models.Event
	initializers.DB.First(&event, id)
	page := structs.Page{PageName: "Event", Tabs: utils.CurrentTabs(), Classes: utils.Classes}
	files := getEventFilePaths(event)

	eventImagePaths := getExistingEventCoverPhotos()
	eventCardImagePaths := getExistingEventCardPhotos()

	return c.Render("edit_event", fiber.Map{
		"Page":            page,
		"Event":           event,
		"EventPhotos":     eventImagePaths,
		"EventCardPhotos": eventCardImagePaths,
		"Description":     event.Description,
		"Files":           files,
		"Locations":       queries.GetLocations(),
	})
}

func AddEvent(c *fiber.Ctx) error {
	var body struct {
		Name               string
		Date               string
		StartTime          string
		EndTime            string
		Location           string
		Description        template.HTML
		ExistingCoverPhoto string
		ExistingCardPhoto  string
	}

	if err := c.BodyParser(&body); err != nil {
		log.Print(err)
		return err
	}

	date, error := time.ParseInLocation("2006-01-02", body.Date, utils.TZ)
	startTime, error := time.ParseInLocation("2006-01-02 15:04:05", fmt.Sprintf("%s %s", body.Date, body.StartTime), utils.TZ)
	endTime, error := time.ParseInLocation("2006-01-02 15:04:05", fmt.Sprintf("%s %s", body.Date, body.EndTime), utils.TZ)

	if error != nil {
		fmt.Println(error)
		return c.Redirect("/")
	}

	photoUrl := getCoverPhotoUrl(body.ExistingCoverPhoto, c)

	cardPicUrl := getCardPhotoUrl(body.ExistingCardPhoto, c)

	event := models.Event{Title: body.Name, Date: date, StartTime: startTime, EndTime: endTime, Description: body.Description, PictureUrl: photoUrl, CardPicUrl: cardPicUrl, Location: body.Location}
	result := initializers.DB.Create(&event)

	if result.Error != nil {
		log.Print("Error creating Event", result.Error)
		return result.Error
	}

	//Handle Files
	uploadEventFiles(event, c)
	createEventIcs(event, queries.GetLocationByName(event.Location))

	return c.Redirect("/")
}

func EditEventPost(c *fiber.Ctx) error {
	id := c.Params("id")
	var event models.Event
	initializers.DB.First(&event, id)
	//page := structs.Page{PageName: "Event", Tabs: utils.CurrentTabs(), Classes: utils.Classes}

	var body struct {
		Name               string
		Date               string
		StartTime          string
		EndTime            string
		Location           string
		Description        template.HTML
		ExistingCoverPhoto string
		ExistingCardPhoto  string
		DeletedFiles       string
	}

	if err := c.BodyParser(&body); err != nil {
		log.Print(err)
		return err
	}

	if len(body.StartTime) == 5 {
		body.StartTime += ":00"
	}

	if len(body.EndTime) == 5 {
		body.EndTime += ":00"
	}

	date, error := time.ParseInLocation("2006-01-02", body.Date, utils.TZ)
	startTime, error := time.ParseInLocation("2006-01-02 15:04:05", fmt.Sprintf("%s %s", body.Date, body.StartTime), utils.TZ)
	endTime, error := time.ParseInLocation("2006-01-02 15:04:05", fmt.Sprintf("%s %s", body.Date, body.EndTime), utils.TZ)

	filesToDelete := strings.Split(body.DeletedFiles, ",")

	for _, file := range filesToDelete {
		os.Remove(fmt.Sprintf("%s/assets/event/%s/files/%s", os.Getenv("UPLOAD_DIR"), id, file))
	}

	if error != nil {
		fmt.Println(error)
		return c.Redirect("/")
	}

	photoUrl := getCoverPhotoUrl(body.ExistingCoverPhoto, c)
	cardPhotoUrl := getCardPhotoUrl(body.ExistingCardPhoto, c)

	event.Title = body.Name
	event.Date = date
	event.StartTime = startTime
	event.EndTime = endTime
	event.Description = body.Description
	event.PictureUrl = photoUrl
	event.CardPicUrl = cardPhotoUrl
	event.Location = body.Location

	result := initializers.DB.Save(&event)

	if result.Error != nil {
		log.Print("Error creating Event", result.Error)
		return result.Error
	}

	//Handle Files
	uploadEventFiles(event, c)
	createEventIcs(event, queries.GetLocationByName(event.Location))

	//files := getEventFilePaths(event)
	c.Set("HX-Redirect", "/events/"+strconv.FormatUint(uint64(event.ID), 10))
	return c.Next()
}

func DeleteEventPost(c *fiber.Ctx) error {
	id := c.Params("id")
	initializers.DB.Delete(&models.Event{}, id)
	os.RemoveAll(fmt.Sprintf("%s/assets/event/%s/files/", os.Getenv("UPLOAD_DIR"), id))

	c.Set("HX-Redirect", "/")
	return c.Next()
}

func createEventIcs(e models.Event, l models.Location) {
	cal := ics.NewCalendar()
	cal.SetMethod(ics.MethodRequest)
	event := cal.AddEvent(uuid.New().String())
	event.SetCreatedTime(time.Now())
	event.SetDtStampTime(time.Now())
	event.SetModifiedAt(time.Now())
	event.SetStartAt(e.StartTime)
	event.SetEndAt(e.EndTime)
	event.SetSummary(e.Title)
	event.SetLocation(l.Name)
	re := regexp.MustCompile(`\r?\n`)
	htmlDesc := re.ReplaceAllString(string(e.Description), "\n")
	desc := strings.Replace(htmlDesc, "<b>", "*", -1)
	desc = strings.Replace(desc, "</b>", "*", -1)
	for strings.Contains(desc, "<a") {
		start := strings.Index(desc, "<a")
		end := strings.Index(desc, "</a>") + 4
		newDesc := strings.TrimSpace(desc[0:start])
		desc = newDesc + strings.TrimSpace(desc[end:])
	}
	desc = desc + "\n\nFor more information, visit https://shinkyushotokan.us/events/" + strconv.FormatUint(uint64(e.ID), 10)
	event.SetDescription(desc)
	event.SetProperty("X-ALT-DESC;FMTTYPE=text/html", fmt.Sprintf("<!DOCTYPE HTML PUBLIC \"-//W3C//DTD HTML 3.2//EN\"><HTML><BODY>%s</BODY></HTML>", htmlDesc))
	ics := cal.Serialize()
	f, err := os.Create(fmt.Sprintf("%s/assets/event/%s/%s.ics", os.Getenv("UPLOAD_DIR"), strconv.FormatUint(uint64(e.ID), 10), e.Title))
	if err != nil {
		fmt.Println(err)
		return
	}
	f.WriteString(ics)
}

func StartAddEvent(c *fiber.Ctx) error {
	var body struct {
		EventType             string
		Date                  string
		StartTime             string
		EndTime               string
		Location              string
		CheckInTime           string
		PromotionalType       string
		TournamentType        string
		RegistrationCloseDate string
		RegistrationCloseTime string
	}

	if err := c.BodyParser(&body); err != nil {
		log.Print(err)
		return err
	}

	date, err := time.ParseInLocation("2006-01-02", body.Date, utils.TZ)
	if err != nil {
		log.Print(err)
		return err
	}

	startTime, err := time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", body.Date, body.StartTime), utils.TZ)
	if err != nil {
		log.Print(err)
		return err
	}

	endTime, err := time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", body.Date, body.EndTime), utils.TZ)
	if err != nil {
		log.Print(err)
		return err
	}

	checkInTime, err := time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", body.Date, body.CheckInTime), utils.TZ)
	if err != nil {
		log.Print(err)
		return err
	}

	var registrationCloseDate time.Time
	var registrationCloseTime time.Time
	if body.RegistrationCloseDate != "" && body.RegistrationCloseTime != "" {
		registrationCloseDate, err = time.ParseInLocation("2006-01-02", body.RegistrationCloseDate, utils.TZ)
		if err != nil {
			log.Print(err)
			return err
		}
		registrationCloseTime, err = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", body.RegistrationCloseDate, body.RegistrationCloseTime), utils.TZ)
		if err != nil {
			log.Print(err)
			return err
		}
	}

	event := models.Event{
		Title:                 body.EventType,
		Date:                  date,
		StartTime:             startTime,
		EndTime:               endTime,
		Location:              body.Location,
		PromotionalType:       body.PromotionalType,
		TournamentType:        body.TournamentType,
		RegistrationCloseDate: registrationCloseDate,
		RegistrationCloseTime: registrationCloseTime,
	}

	// Set default description for Promotional events
	if body.EventType == "Promotional" {
		description := fmt.Sprintf(`The Promotional is coming up fast!
All students who fulfill their next rank requirements from the kids, teen, and adult classes are welcome!

<b>NOTE</b>
You must know all requirements for the rank you are testing for.
Testing in the promotional is not mandatory. If a student does not want to test, we encourage them to come, watch, and cheer on their classmates.
Turn in your application and fee to your instructor in class.

<b>DEADLINE is %s</b>
If you want to test, but can not turn in your application in person, please email it to <a href="mailto:shinkyu.shotokan.karate@gmail.com" style="text-decoration: underline; color:blue;">shinkyu.shotokan.karate@gmail.com</a>

<b>FEES</b>
Kyu (color belt) Test - $20.00
Belt and certificate are given if passed.

Dan (black belt) Test - $30.00
Menjo is given if passed.

<b>Check In:</b> %s
<b>Test Time:</b> %s - %s`,
			date.AddDate(0, 0, -7).Format("January 2, 2006"),
			checkInTime.Format("3:04PM"),
			startTime.Format("3:04PM"),
			endTime.Format("3:04PM"))

		event.Description = template.HTML(description)
		fmt.Printf("Promotional Type: %s\n", body.PromotionalType)
		if body.PromotionalType == "Pre-Karate" {
			fmt.Println("whoo")
			event.PictureUrl = "/public/events/PreKaratePromotional.jpg"
			event.CardPicUrl = "/public/events/PreKaratePromotional.jpg"
			event.Title = "Pre-Karate Promotional"
		} else {
			event.PictureUrl = "/public/events/PXL_20221114_013356064.MP.jpg"
			event.CardPicUrl = "/public/events/PXL_20221114_013356064.MP.jpg"
			event.Title = "Youth-Adult Promotional"
		}
	} else if body.EventType == "Tournament" {
		// Calculate tournament number (2025 is the 23rd tournament)
		tournamentNumber := date.Year() - 2002
		ordinal := "th"
		if tournamentNumber%10 == 1 && tournamentNumber%100 != 11 {
			ordinal = "st"
		} else if tournamentNumber%10 == 2 && tournamentNumber%100 != 12 {
			ordinal = "nd"
		} else if tournamentNumber%10 == 3 && tournamentNumber%100 != 13 {
			ordinal = "rd"
		}

		event.Title = fmt.Sprintf("%d%s Annual %s Karate Tournament", tournamentNumber, ordinal, body.TournamentType)

		description := fmt.Sprintf(`This event is for participants of %s.

<b>Mandatory kumite safety equipment: Mouthpieces and Hand Pads are required. Head and chest protectors are required for all competitors who are 15 years old or younger and below the rank of brown belt. Groin Cups are required for all Male Competitors.</b>

<b>Order of Events:</b>
-Kata
-Team Kata
-Weapons Kata
-Lunch Break (30 minutes)
-Kumite

<b>Check In:</b> %s
<b>Bow In:</b> %s - %s
<b>Tournament:</b> %s - %s

<b>Fees:</b> 1 event: $25, 2 events: $35, 3 events: $40, 4 events: $45, Late Fee: $10

<b>Registering Online</b>
Click the following button to visit the registration site.
You will need to create an online account to register. If you are registering your child, you will need to add them as an account member before you can continue with the registration.
Select the eligible account member, and click <b>Add To Cart</b>.
Your cart will show with a total of <b>$25.00</b>. This is the base amount for competing in one event. 
Click <b>Checkout</b> and you will be brought to the registration prompt where you will enter information about the participant, how many events they are competing in, and what events they are competing in. 
The total price will change depending on how many events the participant is entering in.`,
			getAgeRangeText(body.TournamentType),
			checkInTime.Format("3:04PM"),
			startTime.Format("3:04PM"),
			startTime.Add(30*time.Minute).Format("3:04PM"),
			startTime.Add(30*time.Minute).Format("3:04PM"),
			endTime.Format("3:04PM"))

		// Add registration close date or "will open later" message
		if !registrationCloseDate.IsZero() {
			description += fmt.Sprintf(`

<b>Online Registration will close %s, at %s</b>

<a class="btn btn-primary shadow-md" href="https://secure.rec1.com/CA/south-san-francisco-ca/catalog?filter=c2VhcmNoPTMxODg2OTU=" target="_blank"><svg height="70%" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg"><g id="SVGRepo_bgCarrier" stroke-width="0"></g><g id="SVGRepo_tracerCarrier" stroke-linecap="round" stroke-linejoin="round"></g><g id="SVGRepo_iconCarrier"> <g id="Interface / External_Link"> <path id="Vector" d="M10.0002 5H8.2002C7.08009 5 6.51962 5 6.0918 5.21799C5.71547 5.40973 5.40973 5.71547 5.21799 6.0918C5 6.51962 5 7.08009 5 8.2002V15.8002C5 16.9203 5 17.4801 5.21799 17.9079C5.40973 18.2842 5.71547 18.5905 6.0918 18.7822C6.5192 19 7.07899 19 8.19691 19H15.8031C16.921 19 17.48 19 17.9074 18.7822C18.2837 18.5905 18.5905 18.2839 18.7822 17.9076C19 17.4802 19 16.921 19 15.8031V14M20 9V4M20 4H15M20 4L13 11" stroke="#000000" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"></path> </g> </g></svg> Register Online</a>`,
				registrationCloseDate.Format("January 2, 2006"),
				registrationCloseTime.Format("3:04PM"))
		} else {
			description += `

<b>Online Registration will open later.</b>`
		}

		event.Description = template.HTML(description)

		// Set specific photos based on tournament type
		if body.TournamentType == "Teen & Adult" {
			event.PictureUrl = "/public/events/TeenAdultBanner.png"
			event.CardPicUrl = "/public/events/cards/TeenAdultCard.png"
		} else if body.TournamentType == "Youth" {
			event.PictureUrl = "/public/events/Youth Karate Tournament Banner.png"
			event.CardPicUrl = "/public/events/cards/YouthTournamentCard.png"
		} else {
			event.PictureUrl = "/public/events/All Ages Tournament Banner.png"
			event.CardPicUrl = "/public/events/cards/AllAgesTournamentCard.png"
		}
	}

	// Add promotional application to event files
	// if event.Title == "Promotional" {
	// 	eventDir := fmt.Sprintf("%s/assets/event/%s/files", os.Getenv("UPLOAD_DIR"), strconv.FormatUint(uint64(event.ID), 10))
	// 	os.MkdirAll(eventDir, 0700)
	// 	// Create a file that points to the promotional application
	// 	f, err := os.Create(eventDir + "/Promotional Application.pdf")
	// 	if err == nil {
	// 		f.WriteString("/public/files/Promotional Application.pdf")
	// 		f.Close()
	// 	}
	// }

	page := structs.Page{PageName: "Event", Tabs: utils.CurrentTabs(), Classes: utils.Classes}

	eventImagePaths := getExistingEventCoverPhotos()
	eventCardImagePaths := getExistingEventCardPhotos()

	return c.Render("edit_event", fiber.Map{
		"Page":            page,
		"Event":           event,
		"EventPhotos":     eventImagePaths,
		"EventCardPhotos": eventCardImagePaths,
		"Description":     event.Description,
		"Locations":       queries.GetLocations(),
	})
}

func getAgeRangeText(tournamentType string) string {
	if tournamentType == "Youth" {
		return "ages 17 and younger."
	} else if tournamentType == "Teen & Adult" {
		return "ages 13 and older."
	}
	return "of all ages."
}
