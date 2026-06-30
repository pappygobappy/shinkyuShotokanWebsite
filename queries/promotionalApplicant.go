package queries

import (
	"log"
	"shinkyuShotokan/initializers"
	"shinkyuShotokan/models"
)

func GetApplicantsByPromotionalID(promoID uint) []models.PromotionalApplicant {
	var applicants []models.PromotionalApplicant
	result := initializers.DB.Where("promotional_id = ? AND deleted_at IS NULL", promoID).Order("first_name ASC, last_name ASC").Find(&applicants)
	if result.Error != nil {
		log.Print(result.Error)
	}
	return applicants
}

func CreateApplicant(applicant *models.PromotionalApplicant) error {
	result := initializers.DB.Create(applicant)
	return result.Error
}

func UpdateApplicant(applicant *models.PromotionalApplicant) error {
	result := initializers.DB.Save(applicant)
	return result.Error
}

func SoftDeleteApplicant(id uint) error {
	var applicant models.PromotionalApplicant
	if err := initializers.DB.First(&applicant, id).Error; err != nil {
		return err
	}
	result := initializers.DB.Delete(&applicant)
	return result.Error
}

func ToggleCheckedIn(id uint) error {
	var applicant models.PromotionalApplicant
	if err := initializers.DB.First(&applicant, id).Error; err != nil {
		return err
	}
	applicant.CheckedIn = !applicant.CheckedIn
	result := initializers.DB.Save(&applicant)
	return result.Error
}

func GetCheckedInCountByPromotionalID(promoID uint) (int, error) {
	var count int64
	result := initializers.DB.Model(&models.PromotionalApplicant{}).
		Where("promotional_id = ? AND checked_in = ? AND deleted_at IS NULL", promoID, true).
		Count(&count)
	return int(count), result.Error
}