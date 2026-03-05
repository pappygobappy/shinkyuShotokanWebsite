package utils

import (
	"shinkyuShotokan/queries"
)

func InvalidateLocations() {
	if queries.Cache != nil {
		queries.Cache.Delete("locations")
	}
}

func InvalidateClasses() {
	if queries.Cache != nil {
		queries.Cache.Delete("classes")
	}
}

func InvalidateInstructors() {
	if queries.Cache != nil {
		queries.Cache.Delete("instructors")
		queries.Cache.Delete("visibleInstructors")
	}
}

func InvalidateEventCoverPhotos() {
	if queries.Cache != nil {
		queries.Cache.Delete("eventCoverPhotos")
	}
}

func InvalidateEventCardPhotos() {
	if queries.Cache != nil {
		queries.Cache.Delete("eventCardPhotos")
	}
}
