package models

import "gorm.io/gorm"

type ClassAnnotation struct {
	gorm.Model
	Annotation string
	ClassID    uint
}
