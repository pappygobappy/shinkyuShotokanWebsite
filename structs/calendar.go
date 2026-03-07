package structs

import (
	"shinkyuShotokan/models"
	"time"
)

type CalendarItem struct {
	StartTime   time.Time
	Title       string
	Color       string
	Location    string
	Url         string
	IsCancelled bool
}

type CalendarDay struct {
	Day               time.Time
	NotInCurrentMonth bool
	Events            []CalendarItem
}

type CalendarWeek map[string]CalendarDay

type CalendarViewResult struct {
	Weeks         []CalendarWeek
	Month         time.Time
	Today         time.Time
	PrevMonth     string
	NextMonth     string
	Locations     []models.Location
	Classes       []ActualClass
	Periods       []models.ClassPeriod
	FilteredClass string
}
