package utils

import (
	"html/template"
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
		classesTabs = append(classesTabs, structs.Tab{Name: template.HTML(class.Name), GetUrl: class.GetUrl})
	}

	requirementsTabs := []structs.Tab{
		{Name: `<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 120 120" width="30">
		<rect width="100" height="20" x="10" y="37" fill="#FFD700" rx="5" class="color2A2A2A svgShape" stroke="black" stroke-width="3"></rect>
		<rect width="53.898" height="20" x="18.574" y="75.284" fill="#FFD700" rx="5"
			transform="rotate(-45 18.574 75.284)" class="color393939 svgShape" stroke="black" stroke-width="3"></rect>
		<rect width="53.898" height="20" fill="#FFD700" rx="5"
			transform="scale(-1 1) rotate(-45 40.463 159.351)" class="color393939 svgShape" stroke="black" stroke-width="3"></rect>
		<rect width="34" height="30" x="43" y="32" fill="#FFD700" rx="15" class="color393939 svgShape" stroke="black" stroke-width="3"></rect>
	</svg>
		10th kyu`, GetUrl: "/requirements/10thkyu"},
		{Name: `<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 120 120" width="30">
		<rect width="100" height="20" x="10" y="37" fill="#4169E1" rx="5" class="color2A2A2A svgShape" stroke="black" stroke-width="3"></rect>
		<rect width="53.898" height="20" x="18.574" y="75.284" fill="#4169E1" rx="5"
			transform="rotate(-45 18.574 75.284)" class="color393939 svgShape" stroke="black" stroke-width="3"></rect>
		<rect width="53.898" height="20" fill="#4169E1" rx="5"
			transform="scale(-1 1) rotate(-45 40.463 159.351)" class="color393939 svgShape" stroke="black" stroke-width="3"></rect>
		<rect width="34" height="30" x="43" y="32" fill="#4169E1" rx="15" class="color393939 svgShape" stroke="black" stroke-width="3"></rect>
	</svg>
		9th & 8th kyu`, GetUrl: "/requirements/9th8thkyu"},
		{Name: `<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 120 120" width="30">
		<rect width="100" height="20" x="10" y="37" fill="#2E8B57" rx="5" class="color2A2A2A svgShape" stroke="black" stroke-width="3"></rect>
		<rect width="53.898" height="20" x="18.574" y="75.284" fill="#2E8B57" rx="5"
			transform="rotate(-45 18.574 75.284)" class="color393939 svgShape" stroke="black" stroke-width="3"></rect>
		<rect width="53.898" height="20" fill="#2E8B57" rx="5"
			transform="scale(-1 1) rotate(-45 40.463 159.351)" class="color393939 svgShape" stroke="black" stroke-width="3"></rect>
		<rect width="34" height="30" x="43" y="32" fill="#2E8B57" rx="15" class="color393939 svgShape" stroke="black" stroke-width="3"></rect>
	</svg>
		6th & 7th kyu`, GetUrl: "/requirements/7th6thkyu"},
		{Name: `<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 120 120" width="30">
		<rect width="100" height="20" x="10" y="37" fill="#990099" rx="5" class="color2A2A2A svgShape" stroke="black" stroke-width="3"></rect>
		<rect width="53.898" height="20" x="18.574" y="75.284" fill="#990099" rx="5"
			transform="rotate(-45 18.574 75.284)" class="color393939 svgShape" stroke="black" stroke-width="3"></rect>
		<rect width="53.898" height="20" fill="#990099" rx="5"
			transform="scale(-1 1) rotate(-45 40.463 159.351)" class="color393939 svgShape" stroke="black" stroke-width="3"></rect>
		<rect width="34" height="30" x="43" y="32" fill="#990099" rx="15" class="color393939 svgShape" stroke="black" stroke-width="3"></rect>
	</svg>
		5th & 4th kyu`, GetUrl: "/requirements/5th4thkyu"},
		{Name: `<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 120 120" width="30">
		<rect width="100" height="20" x="10" y="37" fill="#65401a" rx="5" class="color2A2A2A svgShape" stroke="black" stroke-width="3"></rect>
		<rect width="53.898" height="20" x="18.574" y="75.284" fill="#65401a" rx="5"
			transform="rotate(-45 18.574 75.284)" class="color393939 svgShape" stroke="black" stroke-width="3"></rect>
		<rect width="53.898" height="20" fill="#65401a" rx="5"
			transform="scale(-1 1) rotate(-45 40.463 159.351)" class="color393939 svgShape" stroke="black" stroke-width="3"></rect>
		<rect width="34" height="30" x="43" y="32" fill="#65401a" rx="15" class="color393939 svgShape" stroke="black" stroke-width="3"></rect>
	</svg>
		3rd & 2nd kyu`, GetUrl: "/requirements/3rd2ndkyu"},
		{Name: `<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 120 120" width="30">
		<rect width="100" height="20" x="10" y="37" fill="#65401a" rx="5" class="color2A2A2A svgShape" stroke="black" stroke-width="3"></rect>
		<rect width="53.898" height="20" x="18.574" y="75.284" fill="#65401a" rx="5"
			transform="rotate(-45 18.574 75.284)" class="color393939 svgShape" stroke="black" stroke-width="3"></rect>
		<rect width="53.898" height="20" fill="#65401a" rx="5"
			transform="scale(-1 1) rotate(-45 40.463 159.351)" class="color393939 svgShape" stroke="black" stroke-width="3"></rect>
		<rect width="34" height="30" x="43" y="32" fill="#65401a" rx="15" class="color393939 svgShape" stroke="black" stroke-width="3"></rect>
	</svg>
		1st kyu`, GetUrl: "/requirements/1stkyu"},
		{Name: `<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 120 120" width="30" height="30"><path fill="#000000" d="M45 54c0 7.078 6.117 13 13.235 13h3.53C68.856 67 75 61.107 75 54c0-7.18-5.926-13-13.235-13h-3.53C51.1 41 45 46.764 45 54zm-.706-10a16.89 16.89 0 0 0-3.292 10.23L31.218 64H15a5 5 0 0 1-5-5V49a5 5 0 0 1 5-5h29.294zm0 20a16.811 16.811 0 0 1-2.506-4.898L22.109 78.75a5 5 0 0 0 0 7.07l7.071 7.072a5 5 0 0 0 7.071 0l20.621-20.62c-.857-.618-1.65-1.174-2.382-1.675-4.152-.91-7.76-3.3-10.196-6.596zm43.944 0-9.253-9.282A16.911 16.911 0 0 0 75.706 44H105a5 5 0 0 1 5 5v10a5 5 0 0 1-5 5H88.238zm-10.15-4.525L82.598 64h-.054l14.749 14.749a5 5 0 0 1 0 7.07l-7.071 7.072a5 5 0 0 1-7.071 0L62.529 72.27 64 70.857A17.3 17.3 0 0 0 75.706 64a16.828 16.828 0 0 0 2.383-4.525z" class="color000 svgShape"></path></svg>
		1st dan`, GetUrl: "/requirements/1stdan"},
	}

	Tabs = []structs.Tab{
		{Name: "Home", GetUrl: "/"},
		{Name: "Instructors", GetUrl: "/instructors"},
		{Name: "History", GetUrl: "/history"},
		{Name: "Classes", SubTabs: classesTabs},
		{Name: "Requirements", SubTabs: requirementsTabs},
	}
}
