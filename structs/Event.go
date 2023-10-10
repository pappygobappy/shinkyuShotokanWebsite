package structs

import "time"

type Event struct {
	Title       string
	Date        time.Time
	PictureUrl  string
	Alt         string
	Description string
	Href        string
}
