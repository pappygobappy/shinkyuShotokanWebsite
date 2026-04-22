package queries

import (
	"log"
	"shinkyuShotokan/initializers"
	"shinkyuShotokan/models"
)

func GetGearItems() []models.GearItem {
	var items []models.GearItem
	result := initializers.DB.Order("display_order ASC").Find(&items)
	if result.Error != nil {
		log.Print(result.Error)
	}
	return items
}

func CreateGearItem(item *models.GearItem) error {
	var maxOrder int
	initializers.DB.Model(&models.GearItem{}).Select("COALESCE(MAX(display_order), 0)").Scan(&maxOrder)
	item.DisplayOrder = maxOrder + 1
	result := initializers.DB.Create(item)
	return result.Error
}

func GetGearItemById(id uint) models.GearItem {
	var item models.GearItem
	result := initializers.DB.Where("id = ?", id).First(&item)
	if result.Error != nil {
		log.Print(result.Error)
	}
	return item
}

func UpdateGearItem(item *models.GearItem) error {
	result := initializers.DB.Save(item)
	return result.Error
}

func DeleteGearItem(id uint) error {
	result := initializers.DB.Delete(&models.GearItem{}, id)
	return result.Error
}

func swapGearItemOrder(aID uint, bID uint) error {
	var a, b models.GearItem
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
	return initializers.DB.Model(&b).Update("display_order", aOrder).Error
}

func MoveGearItem(id uint, direction string) error {
	var current models.GearItem
	if err := initializers.DB.First(&current, id).Error; err != nil {
		return err
	}
	var neighbor models.GearItem
	if direction == "up" {
		if err := initializers.DB.Where("display_order < ?", current.DisplayOrder).Order("display_order desc").First(&neighbor).Error; err != nil {
			return err
		}
	} else if direction == "down" {
		if err := initializers.DB.Where("display_order > ?", current.DisplayOrder).Order("display_order asc").First(&neighbor).Error; err != nil {
			return err
		}
	} else {
		return nil
	}
	return swapGearItemOrder(current.ID, neighbor.ID)
}