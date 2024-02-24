package structs

import "time"

type CalendarItem struct {
	StartTime time.Time
	Title     string
	Color     string
	Location  string
	Url       string
}
