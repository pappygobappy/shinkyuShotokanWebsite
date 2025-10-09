package queries

import (
	"log"
	"shinkyuShotokan/initializers"
	"shinkyuShotokan/models"

	"github.com/gofiber/fiber/v2"
)

func GetClasses() []models.Class {
	var classes []models.Class
	result := initializers.DB.Find(&classes)
	if result.Error != nil {
		log.Print(result.Error)
	}
	return classes
}

func FindClassByID(id string) models.Class {
	var class models.Class
	result := initializers.DB.Preload("Location").Preload("Annotations").First(&class, id)
	if result.Error != nil {
		log.Print(result.Error)
	}
	return class
}

func FindClassByPath(path string) models.Class {
	var class models.Class
	result := initializers.DB.Preload("Location").Preload("Annotations").Where("get_url = ?", path).First(&class)
	if result.Error != nil {
		log.Print(result.Error)
	}
	return class
}

func FindClassByName(name string) models.Class {
	var class models.Class
	result := initializers.DB.Preload("Location").Preload("Annotations").Where("name = ?", name).First(&class)
	if result.Error != nil {
		log.Print(result.Error)
	}
	return class
}

// Carousel image DB access
func GetCarouselImages() []models.CarouselImage {
	var images []models.CarouselImage
	initializers.DB.Order("display_order asc, id asc").Find(&images)
	return images
}

func GetCarouselImagePaths() []string {
	var paths []string
	initializers.DB.Model(&models.CarouselImage{}).
		Order("display_order asc, id asc").
		Pluck("path", &paths)
	return paths
}

func AddCarouselImage(path string, sourceType string) models.CarouselImage {
	img := models.CarouselImage{Path: path, SourceType: sourceType, DisplayOrder: GetNextCarouselImageOrder()}
	initializers.DB.Create(&img)
	return img
}

func GetNextCarouselImageOrder() int {
	var max int
	initializers.DB.Model(&models.CarouselImage{}).Select("COALESCE(MAX(display_order), 0)").Scan(&max)
	return max + 1
}

func MoveCarouselImage(id uint, direction string) error {
	var current models.CarouselImage
	if err := initializers.DB.First(&current, id).Error; err != nil {
		return err
	}
	var neighbor models.CarouselImage
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
	aOrder, bOrder := current.DisplayOrder, neighbor.DisplayOrder
	if err := initializers.DB.Model(&current).Update("display_order", bOrder).Error; err != nil {
		return err
	}
	if err := initializers.DB.Model(&neighbor).Update("display_order", aOrder).Error; err != nil {
		return err
	}
	return nil
}

func SoftDeleteCarouselImage(id string) error {
	var image models.CarouselImage
	if err := initializers.DB.First(&image, id).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Carousel image not found")
	}
	if err := initializers.DB.Delete(&image).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return nil
}

func RestoreCarouselImage(id string) error {
	if err := initializers.DB.Unscoped().Model(&models.CarouselImage{}).
		Where("id = ?", id).
		Update("deleted_at", nil).Update("display_order", GetNextCarouselImageOrder()).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return nil
}

func GetDeletedCarouselImages() []models.CarouselImage {
	var images []models.CarouselImage
	initializers.DB.Unscoped().Where("deleted_at IS NOT NULL").Order("display_order asc, id asc").Find(&images)
	return images
}
