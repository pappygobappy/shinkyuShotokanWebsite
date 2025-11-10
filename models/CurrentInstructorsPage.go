package models

import "gorm.io/gorm"

type CurrentInstructorsPage struct {
	gorm.Model
	PictureUrl string
}
