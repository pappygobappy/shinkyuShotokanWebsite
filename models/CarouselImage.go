package models

import "gorm.io/gorm"

type CarouselImage struct {
	gorm.Model
	Path         string `gorm:"unique;not null"` // relative path for rendering
	SourceType   string // "upload" or "public" (for reference)
	DisplayOrder int    // display order. Lower is earlier.
}
