package initializers

import (
	"errors"
	"log"
	"shinkyuShotokan/models"
	"time"

	"gorm.io/gorm"
)

var tz, _ = time.LoadLocation("America/Los_Angeles")

func SyncDb() {
	DB.AutoMigrate(
		&models.CarouselImage{},
		&models.User{}, &models.Event{}, &models.ClassSession{}, &models.ClassPeriod{}, &models.ClassAnnotation{}, &models.Instructor{}, &models.PasswordResetToken{},
	)

	seedLocations()
	seedClasses()
	seedEventSubTypes()
	seedEventTemplates()
	seedInstructors()
	seedCarouselImages()
}

func seedLocations() {
	err := DB.AutoMigrate(&models.Location{})
	if err == nil && DB.Migrator().HasTable(&models.Location{}) {
		if err := DB.First(&models.Location{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			locations := []models.Location{
				{
					Name:             "Municipal Services Building Social Hall",
					Address:          "33 Arroyo Dr\nSouth San Francisco, CA 94080",
					GoogleMapsIframe: "https://www.google.com/maps/embed?pb=!1m14!1m8!1m3!1d370.2848693925495!2d-122.43671174104588!3d37.65611258699785!3m2!1i1024!2i768!4f13.1!3m3!1m2!1s0x808f7979c0e8e543%3A0xb9afc1672af1c20f!2sMunicipal%20Services%20Building!5e0!3m2!1sen!2sus!4v1697495645159!5m2!1sen!2sus",
				},
				{
					Name:             "Joseph A. Fernekes Recreation Building",
					Address:          "781 Tennis Dr\nSouth San Francisco, CA 94080",
					GoogleMapsIframe: "https://www.google.com/maps/embed?pb=!1m18!1m12!1m3!1d754.669097751462!2d-122.4269021943449!3d37.65438985945025!2m3!1f0!2f0!3f0!3m2!1i1024!2i768!4f13.1!3m3!1m2!1s0x808f799fc91c36df%3A0xcc99a9bb998b1cae!2sJoseph%20A.%20Fernekes%20Recreation%20Building!5e0!3m2!1sen!2sus!4v1697495721802!5m2!1sen!2sus",
				},
				{
					Name:             "Westborough Recreation Building",
					Address:          "2380 Galway Dr\nSouth San Francisco, CA 94080",
					GoogleMapsIframe: "https://www.google.com/maps/embed?pb=!1m18!1m12!1m3!1d3159.2318050301537!2d-122.45995728796625!3d37.64375397190207!2m3!1f0!2f0!3f0!3m2!1i1024!2i768!4f13.1!3m3!1m2!1s0x808f7a2d9d450665%3A0x86c39d47310bd2e4!2sWestborough%20Recreation%20Building!5e0!3m2!1sen!2sus!4v1698348935868!5m2!1sen!2sus",
				},
				{
					Name:             "Terrabay Gymnasium & Rec Center",
					Address:          "1121 S San Francisco Dr\nSouth San Francisco, CA 94080",
					GoogleMapsIframe: "https://www.google.com/maps/embed?pb=!1m18!1m12!1m3!1d3158.141657123408!2d-122.41867442351607!3d37.669379072012134!2m3!1f0!2f0!3f0!3m2!1i1024!2i768!4f13.1!3m3!1m2!1s0x808f7911a3751ad7%3A0x1d1ee69af964118e!2sTerrabay%20Gymnasium%20%26%20Rec%20Center!5e0!3m2!1sen!2sus!4v1697499112094!5m2!1sen!2sus",
				},
				{
					Name:             "Library | Parks & Recreation Center, Banquet Hall #130",
					Address:          "901 Civic Campus Wy\nSouth San Francisco, CA 94080",
					GoogleMapsIframe: "https://www.google.com/maps/embed?pb=!1m18!1m12!1m3!1d3158.7044330718954!2d-122.43764732412208!3d37.65615227201559!2m3!1f0!2f0!3f0!3m2!1i1024!2i768!4f13.1!3m3!1m2!1s0x808f798c300cb3a1%3A0x2439615b79f50e47!2sLibrary%20%7C%20Parks%20%26%20Recreation%20Center!5e0!3m2!1sen!2sus!4v1708465837688!5m2!1sen!2sus",
				},
			}

			result := DB.Create(locations)

			if result.Error != nil {
				log.Print("Error creating seed Locations", result.Error)
			}
		}
	}
}

func seedClasses() {
	err := DB.AutoMigrate(&models.Class{})
	if err == nil && DB.Migrator().HasTable(&models.Class{}) {
		if err := DB.First(&models.Class{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			classes := []models.Class{
				{
					Name:        "Pre-Karate",
					Description: "An introduction to the discipline of karate in a fun and positive environment.  Focus will be on hand and eye coordination, body awareness, following instructions and social interaction. Gi purchases are non-refundable. Level I students will learn commands in Japanese and become familiar with exercises, kicks and blocks.  Instructor approval is required for promotion to Level II, please note this often takes 3-4 seasons.",
					GetUrl:      "/pre-karate-class",
					StartAge:    4,
					EndAge:      8,
					LocationID:  "Library | Parks & Recreation Center, Banquet Hall #130",
					Schedule: `Level 1 (Beginners) Session A: Saturday 8:30M - 9:15AM
					Level 1 (Beginners) Session B:  Saturday 8:30AM - 9:15AM
					Level 2 (White, Color Belts): Tuesday 6:00PM - 6:45PM, Saturday 10:30AM - 11:15AM
					Level 2 (Color Belts): Tuesday 6:00PM - 7:15PM, Saturday 10:30AM - 12:00PM`,
					CardPhoto:    "/public/classes/pre-karate/card.png",
					BannerPhoto:  "/public/classes/pre-karate/banner.png",
					BannerAdjust: 65,
				},
				{
					Name:        "Youth",
					Description: "Although self-defense may be the primary reason for taking up karate, this training has much more to offer. A young child can develop self-confidence, concentration, perseverance, good sportsmanship, a respectful attitude, good health along with techniques of self-defense. Parents can stay for first and last classes only. Students will learn commands in Japanese, become familiar with exercises, blocks, kicks and more. Students will be promoted to Level II when ready and promoted by instructor (often it takes 3 - 4 sessions). Karate Gi (uniform) is required and may be purchased from instructors. Sensei Sue Miller has studied Karate since 1972 and holds the rank of 8th Dan degree black belt. ",
					GetUrl:      "/youth-class",
					StartAge:    9,
					EndAge:      12,
					LocationID:  "Joseph A. Fernekes Recreation Building",
					Schedule: `Level 1 (Beginners, Yellow, Blue Belts): Monday/Wednesday 5:30PM - 6:30PM
					Level 2 (Green First Level): Monday/Wednesday 6:00PM - 7:00PM
					Advanced Level (Green, Purple, Brown): Monday/Wednesday 6:00PM - 7:30PM`,
					CardPhoto:    "/public/classes/kids/card.jpg",
					BannerPhoto:  "/public/classes/kids/card.jpg",
					BannerAdjust: 50,
				},
				{
					Name:         "Teen",
					Description:  "Learn Karate as an ancient art form, the traditional way of Shotokan. Attain knowledge of self-defense in a spiritual, mental and physical way. Develop confidence; build up your endurance, focus and self-awareness. This is a great way to keep in shape, in a friendly atmosphere. Students should wear loose clothing. Sensei Leroy Rodriques started this program in 1965. Sensei Sue Miller has studied Karate since 1972 and holds the rank of 8th Dan degree black belt. She will be assisted by Nobu Kaji, 8th degree black belt. ",
					GetUrl:       "/teen-class",
					StartAge:     13,
					EndAge:       17,
					LocationID:   "Joseph A. Fernekes Recreation Building",
					Schedule:     "Tuesday/Thursday 6:30PM - 7:30PM",
					CardPhoto:    "/public/classes/teen/card.jpg",
					BannerPhoto:  "/public/classes/teen/banner.jpg",
					BannerAdjust: 75,
				},
				{
					Name:         "Adult",
					Description:  "Learn Karate as an ancient art form, the traditional way of Shotokan. Attain knowledge of self-defense in a spiritual, mental and physical way. Develop confidence; build up your endurance, focus and self-awareness. This is a great way to keep in shape, in a friendly atmosphere. Students should wear loose clothing. Sensei Leroy Rodriques started this program in 1965. Sensei Sue Miller has studied Karate since 1972 and holds the rank of 8th Dan degree black belt. She will be assisted by Nobu Kaji, 8th degree black belt. ",
					GetUrl:       "/adult-class",
					StartAge:     18,
					LocationID:   "Joseph A. Fernekes Recreation Building",
					Schedule:     "Tuesday/Thursday 7:30PM - 9:00PM",
					CardPhoto:    "/public/classes/adult/card.jpg",
					BannerPhoto:  "/public/classes/adult/banner.jpg",
					BannerAdjust: 65,
				},
			}

			result := DB.Create(classes)

			if result.Error != nil {
				log.Print("Error creating seed Classes", result.Error)
			}
		}
	}
}

func seedEventSubTypes() {
	err := DB.AutoMigrate(&models.EventSubType{})
	if err == nil && DB.Migrator().HasTable(&models.EventSubType{}) {
		if err := DB.First(&models.EventSubType{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			eventSubTypes := []models.EventSubType{
				{
					Model: gorm.Model{ID: 1},
					Name:  "Pre-Karate",
				},
				{
					Model: gorm.Model{ID: 2},
					Name:  "Youth & Adult",
				},
				{
					Model: gorm.Model{ID: 3},
					Name:  "Youth",
				},
				{
					Model: gorm.Model{ID: 4},
					Name:  "Teen & Adult",
				},
				{
					Model: gorm.Model{ID: 5},
					Name:  "All Ages",
				},
			}
			result := DB.Create(eventSubTypes)

			if result.Error != nil {
				log.Print("Error creating seed EventSubType", result.Error)
			}
		}
	}
}

func seedEventTemplates() {
	err := DB.AutoMigrate(&models.EventTemplate{})
	if err == nil && DB.Migrator().HasTable(&models.EventTemplate{}) {
		if err := DB.First(&models.EventTemplate{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			eventTemplates := []models.EventTemplate{
				{
					Name:        "Tournament",
					StartTime:   time.Date(0, 0, 0, 10, 0, 0, 0, tz),
					EndTime:     time.Date(0, 0, 0, 16, 0, 0, 0, tz),
					CheckInTime: time.Date(0, 0, 0, 9, 30, 0, 0, tz),
					Description: `This tournament is for competitors %s.
					
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

<b>Fees:</b> 1 event: $25, 2 events: $35, 3 events: $40, 4 events: $45, Late Fee: $10`,
					LocationID: "Terrabay Gymnasium & Rec Center",
					EventSubTypes: []models.EventSubType{
						{Model: gorm.Model{ID: 3}},
						{Model: gorm.Model{ID: 4}},
						{Model: gorm.Model{ID: 5}},
					},
				},
				{
					Name:        "Promotional",
					StartTime:   time.Date(0, 0, 0, 13, 0, 0, 0, tz),
					EndTime:     time.Date(0, 0, 0, 17, 0, 0, 0, tz),
					CheckInTime: time.Date(0, 0, 0, 12, 0, 0, 0, tz),
					Description: `The Promotional is coming up fast!
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
					LocationID: "Terrabay Gymnasium & Rec Center",
					EventSubTypes: []models.EventSubType{
						{Model: gorm.Model{ID: 2}},
					},
				},
				{
					Name:        "Promotional",
					StartTime:   time.Date(0, 0, 0, 13, 30, 0, 0, tz),
					EndTime:     time.Date(0, 0, 0, 16, 30, 0, 0, tz),
					CheckInTime: time.Date(0, 0, 0, 13, 0, 0, 0, tz),
					Description: `The Pre-Karate Promotional is coming up fast!
All students who fulfill their next rank requirements are welcome!

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
					LocationID: "Terrabay Gymnasium & Rec Center",
					EventSubTypes: []models.EventSubType{
						{Model: gorm.Model{ID: 1}},
					},
				},
			}
			result := DB.Create(eventTemplates)

			if result.Error != nil {
				log.Print("Error creating seed EventTemplates", result.Error)
			}
		}
	}
}

func seedInstructors() {
	err := DB.AutoMigrate(&models.Instructor{})
	if err == nil && DB.Migrator().HasTable(&models.Instructor{}) {
		if err := DB.First(&models.Instructor{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			instructors := []models.Instructor{
				{
					Name:       "Sensei Leroy Rodrigues",
					PictureUrl: "/public/instructors/leroy.jpg",
					Bio: `Sensei Leroy Rodrigues has studied karate since 1961 and holds the rank of 10th Dan. 
					He founded the Shinkyu Shotokan Dojo in 1983. 
					Sensei Leroy knows approximately 50 karate katas and more than 15 weapons katas. 
					He also has a book published as well as a video, containing older katas of Shorinji-Ryu.`,
					DisplayOrder: 1,
				},
				{
					Name:       "Sensei Sue Miller",
					PictureUrl: "/public/instructors/940489572.jpg",
					Bio: `Sensei Sue been training since 1972. 
					In the beginning, she trained with Sensei Leroy Rodrigues as a Okinawan Stylist in Shorinji-Ryu. 
					She is an 8th Dan and currently teaches the Pre-Karate Classes, Youth, Teen and Adults, as well as, Men and Women's Self Defense Classes and is the Head Instructor for our Tournaments and Promotional's.`,
					DisplayOrder: 2,
				},
				{
					Name:       "Sensei Nobu Kaji",
					PictureUrl: "/public/instructors/854453422.jpg",
					Bio: `Sensei Nobu has been training in KobuJutsu and Karate since 1968. 
					He holds the rank of 6th Dan in KobuJutsu and 8th Dan with Shinkyu Shotokan.  
					His Karate styles include Ryugo-ryu, Magai-ryu, Yamani-ryu, Shito-ryu and Shorin-ryu.`,
					DisplayOrder: 3,
				},
				{
					Name:       "Sensei Patrick Dunleavy",
					PictureUrl: "/public/instructors/Patrick.jpg",
					Bio: `Patrick has been continuously studying Shotokan Karate since he was 6 years old. 
					Even from the very beginning, he loved karate and now, almost 30 years later, he is a 5th degree black belt. 
					What has always inspired him was looking up to people who have been taking karate for a long time and seeing how far he could go with his own karate. 
					Today, Patrick is teaching Shotokan Karate in the adult and teen classes while actively continuing his own karate training. He also regularly competes successfully in karate tournaments representing Shinkyu Shotokan.`,
					DisplayOrder: 4,
				},
				{
					Name:       "Senpai Alex Moreno",
					PictureUrl: "/public/instructors/alex.jpeg",
					Bio: `Senpai Alex has been training with Shinkyu Shotokan Karate for over 15 years and earned his Shodan in 2015. 
					Karate has been a constant in his life and he enjoys sharing his knowledge with the next generation of students.`,
					DisplayOrder: 5,
				},
			}

			result := DB.Create(instructors)

			if result.Error != nil {
				log.Print("Error creating seed Instructors", result.Error)
			}
		}
	}
}

func seedCarouselImages() {
	err := DB.AutoMigrate(&models.CarouselImage{})
	if err == nil && DB.Migrator().HasTable(&models.CarouselImage{}) {
		if err := DB.First(&models.CarouselImage{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			images := []models.CarouselImage{
				{
					Path:         "/public/image_carousel/PXL_20210506_005649942.jpg",
					SourceType:   "public",
					DisplayOrder: 1,
				},
				{
					Path:         "/public/image_carousel/PXL_20221114_013356064.MP.jpg",
					SourceType:   "public",
					DisplayOrder: 2,
				},
				{
					Path:         "/public/image_carousel/PXL_20230226_232040512.jpg",
					SourceType:   "public",
					DisplayOrder: 3,
				},
				{
					Path:         "/public/image_carousel/heiansandan2.jpeg",
					SourceType:   "public",
					DisplayOrder: 4,
				},
				{
					Path:         "/public/image_carousel/kizami.jpg",
					SourceType:   "public",
					DisplayOrder: 5,
				},
				{
					Path:         "/public/image_carousel/original_9ed61476-2366-4b52-8b6b-05c698649a55_PXL_20230827_210652281.jpg",
					SourceType:   "public",
					DisplayOrder: 6,
				},
			}

			result := DB.Create(images)

			if result.Error != nil {
				log.Print("Error creating seed Carousel Images", result.Error)
			}
		}
	}
}
