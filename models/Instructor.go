package models

import "gorm.io/gorm"

type Instructor struct {
	gorm.Model
	Name         string `json:"name"`
	PictureUrl   string `json:"picture_url"`
	Bio          string `json:"bio"`
	DisplayOrder int    `json:"display_order"`
	Hidden       bool   `json:"hidden"`
	ZoomLevel    int    `json:"zoom_level" gorm:"default:100"`
	OffsetX      int    `json:"offset_x" gorm:"default:0"`
	OffsetY      int    `json:"offset_y" gorm:"default:0"`
}

func (instructor Instructor) InitialZoomLevel() float64 {
	if instructor.ZoomLevel <= 100 {
		return 1
	}
	return float64(instructor.ZoomLevel) / 100.0
}
