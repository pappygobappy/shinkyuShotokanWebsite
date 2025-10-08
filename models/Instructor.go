package models

import "gorm.io/gorm"

type Instructor struct {
	gorm.Model
	Name         string
	PictureUrl   string
	Bio          string
	DisplayOrder int
	Hidden       bool
}
