package models

import (
	"html/template"
	"regexp"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Event struct {
	gorm.Model
	Title                 string
	Date                  time.Time
	StartTime             time.Time
	EndTime               time.Time
	Location              string
	PictureUrl            string
	CardPicUrl            string
	Alt                   string
	Description           template.HTML
	PromotionalType       string
	TournamentType        string
	RegistrationCloseDate time.Time
	RegistrationCloseTime time.Time
}

func (event Event) SafeDescription() string {
	re := regexp.MustCompile(`\r?\n`)
	htmlDesc := re.ReplaceAllString(string(event.Description), "\n")
	desc := strings.Replace(htmlDesc, "<b>", "*", -1)
	desc = strings.Replace(desc, "</b>", "*", -1)
	for strings.Contains(desc, "<a") {
		start := strings.Index(desc, "<a")
		end := strings.Index(desc, "</a>") + 4
		newDesc := strings.TrimSpace(desc[0:start])
		desc = newDesc + strings.TrimSpace(desc[end:])
	}
	desc = desc + "\n\nFor more information, visit https://shinkyushotokan.us/events/" + strconv.FormatUint(uint64(event.ID), 10)
	return desc
}

func (event Event) GoogleDescription() template.HTML {
	desc := event.Description +
		`

For more information, visit <a href="https://shinkyushotokan.us/events/` + template.HTML(strconv.FormatUint(uint64(event.ID), 10)) + `">https://shinkyushotokan.us/events/` + template.HTML(strconv.FormatUint(uint64(event.ID), 10)) + `</a>`
	return desc
}

func (event Event) OutlookDescription() template.HTML {
	desc := event.Description +
		`

For more information, visit <a href="https://shinkyushotokan.us/events/` + template.HTML(strconv.FormatUint(uint64(event.ID), 10)) + `">https://shinkyushotokan.us/events/` + template.HTML(strconv.FormatUint(uint64(event.ID), 10)) + `</a>`
	desc = template.HTML(strings.Replace(string(desc), "\n", "<br />", -1))
	return desc
}

func (event Event) GetCardPicUrl() string {
	if event.CardPicUrl != "" {
		return event.CardPicUrl
	}
	return event.PictureUrl
}
