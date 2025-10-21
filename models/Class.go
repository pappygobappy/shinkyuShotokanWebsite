package models

import (
	"strings"

	"gorm.io/gorm"
)

type Class struct {
	gorm.Model
	Name         string
	Description  string
	Annotations  []ClassAnnotation
	LocationID   string
	Location     Location `gorm:"references:Name"`
	GetUrl       string
	StartAge     int
	EndAge       int
	Schedule     string
	CardPhoto    string
	BannerPhoto  string
	BannerAdjust int
}

func (class Class) DisplayName() string {
	if strings.Contains(class.Name, "Karate") {
		return class.Name + " Class"
	}
	return class.Name + " Karate Class"
}
