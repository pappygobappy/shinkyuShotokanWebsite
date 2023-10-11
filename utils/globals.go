package utils

import (
	"shinkyuShotokan/structs"
)

var Instructors []structs.Instructor

var Tabs []structs.Tab
var Events []structs.Event
var Classes []structs.Class

func Init() {

	Instructors = []structs.Instructor{
		{Name: "Leroy Rodrigues", PictureUrl: "/public/instructors/leroy.jpg", Bio: "Sensei Leroy Rodrigues has studied karate since 1961 and holds the rank of 10th Dan. He founded the Shinkyu Shotokan Dojo in 1983. Sensei Leroy knows approximately 50 karate katas and more than 15 weapons katas. He also has a book published as well as a video, containing older katas of Shorinji-Ryu."},
		{Name: "Sue Miller", PictureUrl: "/public/instructors/940489572.jpg", Bio: "Sensei Sue been training since 1972. In the beginning, she trained with Sensei Leroy Rodrigues as a Okinawan Stylist in Shorinji-Ryu. She is an 8th Dan and currently teaches the Pre-Karate Classes, Youth, Teen and Adults, as well as, Men and Women's Self Defense Classes and is the Head Instructor for our Tournaments and Promotional's."},
		{Name: "Sensei Nobu Kaji", PictureUrl: "/public/instructors/854453422.jpg", Bio: "Sensei Nobu has been training in KobuJutsu and Karate since 1968. He holds the rank of 6th Dan in KobuJutsu and 8th Dan with Shinkyu Shotokan.  His Karate styles include Ryugo-ryu, Magai-ryu, Yamani-ryu, Shito-ryu and Shorin-ryu."},
	}
	Events = []structs.Event{
		{Title: "Promotional", PictureUrl: "/assets/image_carousel/PXL_20221114_013356064.MP.jpg", Alt: "Promotional", Description: "Belt Testing Coming Up fast! Practice Hard!"},
		{Title: "Promotional", PictureUrl: "/assets/image_carousel/PXL_20221114_013356064.MP.jpg", Alt: "Promotional", Description: "Belt Testing Coming Up fast! Practice Hard!"},
		{Title: "Promotional", PictureUrl: "/assets/image_carousel/PXL_20221114_013356064.MP.jpg", Alt: "Promotional", Description: "Belt Testing Coming Up fast! Practice Hard!"},
		{Title: "Promotional", PictureUrl: "/assets/image_carousel/PXL_20221114_013356064.MP.jpg", Alt: "Promotional", Description: "Belt Testing Coming Up fast! Practice Hard!"},
	}
	Classes = []structs.Class{
		{
			Name:          "Pre-Karate",
			GoogleMapsSrc: "https://www.google.com/maps/embed?pb=!1m18!1m12!1m3!1d3158.7109867523395!2d-122.43907052288279!3d37.65599821899114!2m3!1f0!2f0!3f0!3m2!1i1024!2i768!4f13.1!3m3!1m2!1s0x808f7979c0e8e543%3A0xb9afc1672af1c20f!2sMunicipal%20Services%20Building!5e0!3m2!1sen!2sus!4v1696289055039!5m2!1sen!2sus",
			Description:   "An introduction to the discipline of karate in a fun and positive environment.  Focus will be on hand and eye coordination, body awareness, following instructions and social interaction. Gi purchases are non-refundable. Level I students will learn commands in Japanese and become familiar with exercises, kicks and blocks.  Instructor approval is required for promotion to Level II, please note this often takes 3-4 seasons.",
			GetUrl:        "/pre-karate-class",
			StartAge:      4,
			EndAge:        8,
		},
		{
			Name:          "Kids",
			GoogleMapsSrc: "https://www.google.com/maps/embed?pb=!1m18!1m12!1m3!1d3158.7763814395785!2d-122.426385!3d37.654461!2m3!1f0!2f0!3f0!3m2!1i1024!2i768!4f13.1!3m3!1m2!1s0x808f799fc91c36df%3A0xcc99a9bb998b1cae!2sJoseph%20A.%20Fernekes%20Recreation%20Building!5e0!3m2!1sen!2sus!4v1696366282348!5m2!1sen!2sus",
			Description:   "An introduction to the discipline of karate in a fun and positive environment.  Focus will be on hand and eye coordination, body awareness, following instructions and social interaction. Gi purchases are non-refundable. Level I students will learn commands in Japanese and become familiar with exercises, kicks and blocks.  Instructor approval is required for promotion to Level II, please note this often takes 3-4 seasons.",
			GetUrl:        "/kids-class",
			StartAge:      9,
			EndAge:        12,
		},
		{
			Name:          "Teen",
			GoogleMapsSrc: "https://www.google.com/maps/embed?pb=!1m18!1m12!1m3!1d3158.7763814395785!2d-122.426385!3d37.654461!2m3!1f0!2f0!3f0!3m2!1i1024!2i768!4f13.1!3m3!1m2!1s0x808f799fc91c36df%3A0xcc99a9bb998b1cae!2sJoseph%20A.%20Fernekes%20Recreation%20Building!5e0!3m2!1sen!2sus!4v1696366282348!5m2!1sen!2sus",
			Description:   "An introduction to the discipline of karate in a fun and positive environment.  Focus will be on hand and eye coordination, body awareness, following instructions and social interaction. Gi purchases are non-refundable. Level I students will learn commands in Japanese and become familiar with exercises, kicks and blocks.  Instructor approval is required for promotion to Level II, please note this often takes 3-4 seasons.",
			GetUrl:        "/teen-class",
			StartAge:      13,
			EndAge:        17,
		},
		{
			Name:          "Adult",
			GoogleMapsSrc: "https://www.google.com/maps/embed?pb=!1m18!1m12!1m3!1d3158.7763814395785!2d-122.426385!3d37.654461!2m3!1f0!2f0!3f0!3m2!1i1024!2i768!4f13.1!3m3!1m2!1s0x808f799fc91c36df%3A0xcc99a9bb998b1cae!2sJoseph%20A.%20Fernekes%20Recreation%20Building!5e0!3m2!1sen!2sus!4v1696366282348!5m2!1sen!2sus",
			Description:   "An introduction to the discipline of karate in a fun and positive environment.  Focus will be on hand and eye coordination, body awareness, following instructions and social interaction. Gi purchases are non-refundable. Level I students will learn commands in Japanese and become familiar with exercises, kicks and blocks.  Instructor approval is required for promotion to Level II, please note this often takes 3-4 seasons.",
			GetUrl:        "/adult-class",
			StartAge:      18,
		},
	}

	var classesTabs []structs.Tab
	for _, class := range Classes {
		classesTabs = append(classesTabs, structs.Tab{Name: class.Name, GetUrl: class.GetUrl})
	}

	Tabs = []structs.Tab{
		{Name: "Home", GetUrl: "/"},
		{Name: "Instructors", GetUrl: "/instructors"},
		{Name: "History", GetUrl: "/history"},
		{Name: "Classes", SubTabs: classesTabs},
	}
}
