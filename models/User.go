package models

import "gorm.io/gorm"

type UserType string

const (
	AdminUser UserType = "admin"
	Owner     UserType = "owner"
)

type User struct {
	gorm.Model
	Email        string `gorm:"unique"`
	PasswordHash string
	FirstName    string
	LastName     string
	Type         UserType
}
