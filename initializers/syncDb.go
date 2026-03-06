package initializers

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"path/filepath"
	"shinkyuShotokan/models"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var seedDir string

func init() {
	exePath, err := os.Executable()
	if err == nil {
		exeDir := filepath.Dir(exePath)
		seedDir = filepath.Join(exeDir, "..", "seeds")
		log.Printf("Checking executable-based path: %s", seedDir)
	}

	if _, err := os.Stat(seedDir); err != nil {
		cwd, _ := os.Getwd()
		seedDir = filepath.Join(cwd, "seeds")
		log.Printf("Using working directory-based path: %s", seedDir)
	}

	if _, err := os.Stat(seedDir); err != nil {
		log.Fatalf("Seed directory not found at %s", seedDir)
	}
}

var tz, _ = time.LoadLocation("America/Los_Angeles")

type LocationJSON struct {
	Name             string `json:"name"`
	Address          string `json:"address"`
	GoogleMapsIframe string `json:"google_maps_iframe"`
}

type ClassJSON struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	GetUrl       string `json:"get_url"`
	StartAge     int    `json:"start_age"`
	EndAge       int    `json:"end_age"`
	LocationID   string `json:"location_id"`
	Schedule     string `json:"schedule"`
	CardPhoto    string `json:"card_photo"`
	BannerPhoto  string `json:"banner_photo"`
	BannerAdjust int    `json:"banner_adjust"`
}

type InstructorJSON struct {
	Name         string `json:"name"`
	PictureUrl   string `json:"picture_url"`
	Bio          string `json:"bio"`
	DisplayOrder int    `json:"display_order"`
}

type EventSubTypeJSON struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type EventTemplateJSON struct {
	Name          string             `json:"name"`
	StartTime     time.Time          `json:"start_time"`
	EndTime       time.Time          `json:"end_time"`
	CheckInTime   time.Time          `json:"check_in_time"`
	Description   string             `json:"description"`
	LocationID    string             `json:"location_id"`
	EventSubTypes []EventSubTypeJSON `json:"event_sub_types"`
}

type CarouselImageJSON struct {
	Path         string `json:"path"`
	SourceType   string `json:"source_type"`
	DisplayOrder int    `json:"display_order"`
}

func SyncDb() {
	DB.AutoMigrate(
		&models.Location{},
		&models.Class{},
		&models.EventSubType{},
		&models.EventTemplate{},
		&models.CarouselImage{},
		&models.User{},
		&models.Event{},
		&models.ClassSession{},
		&models.ClassPeriod{},
		&models.ClassAnnotation{},
		&models.Instructor{},
		&models.PasswordResetToken{},
		&models.CurrentInstructorsPage{},
	)

	seedLocations()
	seedClasses()
	seedEventSubTypes()
	seedEventTemplates()
	seedInstructors()
	seedCarouselImages()
	setOwner()
	seedCurrentInstructorsPage()
}

func seedCurrentInstructorsPage() {
	if err := DB.First(&models.CurrentInstructorsPage{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		currentInstructorsPage := models.CurrentInstructorsPage{
			PictureUrl: "/public/instructors/original_0dba6fc4-1896-4cc8-926b-56568fb5ea74_PXL_20230827_175124417 (1).jpg",
		}
		DB.Create(&currentInstructorsPage)
	}
}

func setOwner() {
	var count int64
	if err := DB.Model(&models.User{}).Count(&count).Error; err != nil {
		log.Print("Error counting users", err)
		return
	}
	if count != 1 {
		return
	}
	var user models.User
	if err := DB.First(&user).Error; err != nil {
		log.Print("Error retrieving user", err)
		return
	}
	if user.Type != "" {
		return
	}
	if err := DB.Model(&user).Update("type", models.OwnerUser).Error; err != nil {
		log.Print("Error setting user type to owner", err)
	}
}

func seedLocations() {
	var count int64
	DB.Model(&models.Location{}).Count(&count)

	if count > 0 {
		log.Println("Locations already seeded, skipping...")
		return
	}

	jsonData, err := os.ReadFile(filepath.Join(seedDir, "locations.json"))
	if err != nil {
		log.Fatalf("Failed to read locations.json: %v", err)
	}

	var locationsData []LocationJSON
	if err := json.Unmarshal(jsonData, &locationsData); err != nil {
		log.Fatalf("Failed to parse locations.json: %v", err)
	}

	var locations []models.Location
	for _, loc := range locationsData {
		locations = append(locations, models.Location{
			Name:             loc.Name,
			Address:          loc.Address,
			GoogleMapsIframe: loc.GoogleMapsIframe,
		})
	}

	if err := DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&locations).Error; err != nil {
		log.Fatalf("Failed to seed locations: %v", err)
	}

	log.Printf("✅ Seeded %d locations", len(locations))
}

func seedClasses() {
	var count int64
	DB.Model(&models.Class{}).Count(&count)

	if count > 0 {
		log.Println("Classes already seeded, skipping...")
		return
	}

	jsonData, err := os.ReadFile(filepath.Join(seedDir, "classes.json"))
	if err != nil {
		log.Fatalf("Failed to read classes.json: %v", err)
	}

	var classesData []ClassJSON
	if err := json.Unmarshal(jsonData, &classesData); err != nil {
		log.Fatalf("Failed to parse classes.json: %v", err)
	}

	var classes []models.Class
	for _, classData := range classesData {
		classes = append(classes, models.Class{
			Name:         classData.Name,
			Description:  classData.Description,
			GetUrl:       classData.GetUrl,
			StartAge:     classData.StartAge,
			EndAge:       classData.EndAge,
			LocationID:   classData.LocationID,
			Schedule:     classData.Schedule,
			CardPhoto:    classData.CardPhoto,
			BannerPhoto:  classData.BannerPhoto,
			BannerAdjust: classData.BannerAdjust,
		})
	}

	if err := DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&classes).Error; err != nil {
		log.Fatalf("Failed to seed classes: %v", err)
	}

	log.Printf("✅ Seeded %d classes", len(classes))

	var classsWithoutOrder []models.Class
	if err := DB.Where("display_order IS NULL").Find(&classsWithoutOrder).Error; err == nil && len(classsWithoutOrder) > 0 {
		log.Println("Setting order for classes")
		for i := range classesData {
			className := classesData[i].Name
			var displayOrder int
			switch className {
			case "Pre-Karate":
				displayOrder = 1
			case "Youth":
				displayOrder = 2
			case "Teen":
				displayOrder = 3
			case "Adult":
				displayOrder = 4
			default:
				continue
			}

			if err := DB.Model(&models.Class{}).Where("name = ?", className).Update("display_order", displayOrder).Error; err != nil {
				log.Printf("Failed to set display order for class %s: %v", className, err)
			}
		}
	}
}

func seedEventSubTypes() {
	var count int64
	DB.Model(&models.EventSubType{}).Count(&count)

	if count > 0 {
		log.Println("Event subtypes already seeded, skipping...")
		return
	}

	jsonData, err := os.ReadFile(filepath.Join(seedDir, "event_subtypes.json"))
	if err != nil {
		log.Fatalf("Failed to read event_subtypes.json: %v", err)
	}

	var subTypesData []EventSubTypeJSON
	if err := json.Unmarshal(jsonData, &subTypesData); err != nil {
		log.Fatalf("Failed to parse event_subtypes.json: %v", err)
	}

	var subTypes []models.EventSubType
	for _, subTypeData := range subTypesData {
		subTypes = append(subTypes, models.EventSubType{
			Model: gorm.Model{ID: subTypeData.ID},
			Name:  subTypeData.Name,
		})
	}

	if err := DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&subTypes).Error; err != nil {
		log.Fatalf("Failed to seed event subtypes: %v", err)
	}

	log.Printf("✅ Seeded %d event subtypes", len(subTypes))
}

func seedEventTemplates() {
	var count int64
	DB.Model(&models.EventTemplate{}).Count(&count)

	if count > 0 {
		log.Println("Event templates already seeded, skipping...")
		return
	}

	jsonData, err := os.ReadFile(filepath.Join(seedDir, "event_templates.json"))
	if err != nil {
		log.Fatalf("Failed to read event_templates.json: %v", err)
	}

	var templatesData []EventTemplateJSON
	if err := json.Unmarshal(jsonData, &templatesData); err != nil {
		log.Fatalf("Failed to parse event_templates.json: %v", err)
	}

	var templates []models.EventTemplate
	for _, templateData := range templatesData {
		template := models.EventTemplate{
			Name:        templateData.Name,
			StartTime:   templateData.StartTime,
			EndTime:     templateData.EndTime,
			CheckInTime: templateData.CheckInTime,
			Description: templateData.Description,
			LocationID:  templateData.LocationID,
		}

		for _, subTypeID := range templateData.EventSubTypes {
			template.EventSubTypes = append(template.EventSubTypes, models.EventSubType{Model: gorm.Model{ID: subTypeID.ID}})
		}

		templates = append(templates, template)
	}

	if err := DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&templates).Error; err != nil {
		log.Fatalf("Failed to seed event templates: %v", err)
	}

	log.Printf("✅ Seeded %d event templates", len(templates))
}

func seedInstructors() {
	var count int64
	DB.Model(&models.Instructor{}).Count(&count)

	if count > 0 {
		log.Println("Instructors already seeded, skipping...")
		return
	}

	jsonData, err := os.ReadFile(filepath.Join(seedDir, "instructors.json"))
	if err != nil {
		log.Fatalf("Failed to read instructors.json: %v", err)
	}

	var instructorsData []InstructorJSON
	if err := json.Unmarshal(jsonData, &instructorsData); err != nil {
		log.Fatalf("Failed to parse instructors.json: %v", err)
	}

	var instructors []models.Instructor
	for _, instructorData := range instructorsData {
		instructors = append(instructors, models.Instructor{
			Name:         instructorData.Name,
			PictureUrl:   instructorData.PictureUrl,
			Bio:          instructorData.Bio,
			DisplayOrder: instructorData.DisplayOrder,
		})
	}

	if err := DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&instructors).Error; err != nil {
		log.Fatalf("Failed to seed instructors: %v", err)
	}

	log.Printf("✅ Seeded %d instructors", len(instructors))
}

func seedCarouselImages() {
	var count int64
	DB.Model(&models.CarouselImage{}).Count(&count)

	if count > 0 {
		log.Println("Carousel images already seeded, skipping...")
		return
	}

	jsonData, err := os.ReadFile(filepath.Join(seedDir, "carousel_images.json"))
	if err != nil {
		log.Fatalf("Failed to read carousel_images.json: %v", err)
	}

	var imagesData []CarouselImageJSON
	if err := json.Unmarshal(jsonData, &imagesData); err != nil {
		log.Fatalf("Failed to parse carousel_images.json: %v", err)
	}

	var images []models.CarouselImage
	for _, imageData := range imagesData {
		images = append(images, models.CarouselImage{
			Path:         imageData.Path,
			SourceType:   imageData.SourceType,
			DisplayOrder: imageData.DisplayOrder,
		})
	}

	if err := DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&images).Error; err != nil {
		log.Fatalf("Failed to seed carousel images: %v", err)
	}

	log.Printf("✅ Seeded %d carousel images", len(images))
}
