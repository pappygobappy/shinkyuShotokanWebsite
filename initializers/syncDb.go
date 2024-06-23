package initializers

import (
	"errors"
	"log"
	"shinkyuShotokan/models"

	"gorm.io/gorm"
)

func SyncDb() {
	DB.AutoMigrate(&models.User{}, &models.Event{}, &models.ClassSession{}, &models.ClassPeriod{}, &models.ClassAnnotation{})

	seedLocations()
	seedClasses()
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
					Schedule:     "Tuesday/Thursday 6:00PM - 7:00PM",
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
					Schedule:     "Tuesday/Thursday 7:00PM - 8:30PM",
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
