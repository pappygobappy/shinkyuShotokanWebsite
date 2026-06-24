# Task 3: Promotional Query Layer

**Files:**
- Create: `queries/promotional.go`

**Interfaces:**
- Consumes: `initializers.DB` (GORM instance)
- Produces: `GetAllPromotionals() []models.Promotional`, `GetPromotionalByID(id uint) models.Promotional`, `CreatePromotional(p *models.Promotional) error`, `UpdatePromotional(p *models.Promotional) error`, `SoftDeletePromotional(id uint) error`, `RestorePromotional(id uint) error`

**Steps:**

- [ ] **Step 1: Create the query file**

Create `queries/promotional.go`:

```go
package queries

import (
	"log"
	"shinkyuShotokan/initializers"
	"shinkyuShotokan/models"
)

func GetAllPromotionals() []models.Promotional {
	var promotionals []models.Promotional
	result := initializers.DB.Order("date DESC").Find(&promotionals)
	if result.Error != nil {
		log.Print(result.Error)
	}
	return promotionals
}

func GetPromotionalByID(id uint) models.Promotional {
	var promotional models.Promotional
	result := initializers.DB.Where("id = ?", id).First(&promotional)
	if result.Error != nil {
		log.Print(result.Error)
	}
	return promotional
}

func CreatePromotional(p *models.Promotional) error {
	result := initializers.DB.Create(p)
	return result.Error
}

func UpdatePromotional(p *models.Promotional) error {
	result := initializers.DB.Save(p)
	return result.Error
}

func SoftDeletePromotional(id uint) error {
	var promotional models.Promotional
	if err := initializers.DB.First(&promotional, id).Error; err != nil {
		return err
	}
	result := initializers.DB.Delete(&promotional)
	return result.Error
}

func RestorePromotional(id uint) error {
	return initializers.DB.Unscoped().Model(&models.Promotional{}).Where("id = ?", id).Update("deleted_at", nil).Error
}
```

- [ ] **Step 2: Verify with go build**

Run: `go build ./...`
Expected: PASS
