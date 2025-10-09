package queries

import (
	"shinkyuShotokan/initializers"
	"shinkyuShotokan/models"
)

func GetInstructors() []models.Instructor {
	var instructors []models.Instructor
	initializers.DB.Order("display_order asc, id asc").Find(&instructors)
	return instructors
}

func GetInstructorByID(id uint) models.Instructor {
	var instructor models.Instructor
	initializers.DB.First(&instructor, id)
	return instructor
}

func CreateInstructor(instructor models.Instructor) models.Instructor {
	initializers.DB.Create(&instructor)
	return instructor
}

func UpdateInstructor(instructor models.Instructor) models.Instructor {
	initializers.DB.Save(&instructor)
	return instructor
}

func DeleteInstructor(id uint) {
	initializers.DB.Delete(&models.Instructor{}, id)
}

func swapInstructorOrder(aID uint, bID uint) error {
	var a, b models.Instructor
	if err := initializers.DB.First(&a, aID).Error; err != nil {
		return err
	}
	if err := initializers.DB.First(&b, bID).Error; err != nil {
		return err
	}
	aOrder := a.DisplayOrder
	bOrder := b.DisplayOrder
	if err := initializers.DB.Model(&a).Update("display_order", bOrder).Error; err != nil {
		return err
	}
	if err := initializers.DB.Model(&b).Update("display_order", aOrder).Error; err != nil {
		return err
	}
	return nil
}

func GetVisibleInstructors() []models.Instructor {
	var instructors []models.Instructor
	initializers.DB.Where("hidden = ?", false).Or("hidden IS NULL").Order("display_order asc, id asc").Find(&instructors)
	return instructors
}

func GetHiddenInstructors() []models.Instructor {
	var instructors []models.Instructor
	initializers.DB.Where("hidden = ?", true).Order("id asc").Find(&instructors)
	return instructors
}

func SetInstructorHidden(id uint, hidden bool) error {
	var instructor models.Instructor
	if err := initializers.DB.First(&instructor, id).Error; err != nil {
		return err
	}
	if instructor.Hidden == hidden {
		return nil // no change
	}
	instructor.Hidden = hidden
	if !hidden {
		// If unhidden, move to end of visible list
		instructor.DisplayOrder = GetNextInstructorDisplayOrder()
	}
	return initializers.DB.Save(&instructor).Error
}

func MoveInstructor(id uint, direction string) error {
	var current models.Instructor
	if err := initializers.DB.First(&current, id).Error; err != nil {
		return err
	}
	if current.Hidden {
		return nil // Do not move hidden instructors
	}
	var neighbor models.Instructor
	if direction == "up" {
		if err := initializers.DB.Where("display_order < ? AND (hidden != ? OR hidden is null)", current.DisplayOrder, true).Order("display_order desc").First(&neighbor).Error; err != nil {
			return err
		}
	} else if direction == "down" {
		if err := initializers.DB.Where("display_order > ? AND (hidden != ? OR hidden is null)", current.DisplayOrder, true).Order("display_order asc").First(&neighbor).Error; err != nil {
			return err
		}
	} else {
		return nil
	}
	return swapInstructorOrder(current.ID, neighbor.ID)
}

func GetNextInstructorDisplayOrder() int {
	var max int
	initializers.DB.Model(&models.Instructor{}).Select("COALESCE(MAX(display_order), 0)").Scan(&max)
	return max + 1
}
