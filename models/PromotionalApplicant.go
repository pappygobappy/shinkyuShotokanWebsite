package models

import (
	"gorm.io/gorm"
)

type PromotionalApplicant struct {
	gorm.Model
	PromotionalID  uint
	FirstName      string
	LastName       string
	Age            int
	RankTestingFor string
	BeltSize       *string
	CheckedIn      bool
	Instructors    []Instructor `gorm:"many2many:promotional_applicant_instructors;"`
}

func (PromotionalApplicant) TableName() string {
	return "promotional_applicants"
}