# Task 2: Promotional Models + DB Migration

**Files:**
- Create: `models/Promotional.go`
- Create: `models/PromotionalApplicant.go`
- Modify: `initializers/syncDb.go:87-102`

**Interfaces:**
- Consumes: `gorm.Model` from `gorm.io/gorm`
- Produces: `models.Promotional` struct with `TableName() string` returning `"promotionals"`; `models.PromotionalApplicant` struct with `TableName() string` returning `"promotional_applicants"`

**Steps:**

- [ ] **Step 1: Create Promotional model**

Create `models/Promotional.go`:

```go
package models

import (
	"time"

	"gorm.io/gorm"
)

type Promotional struct {
	gorm.Model
	Name        string
	Date        time.Time
	StartTime   time.Time
	EndTime     time.Time
	CheckInTime time.Time
	SubType     string
	Location    string
	Description string
}

func (Promotional) TableName() string {
	return "promotionals"
}
```

- [ ] **Step 2: Create PromotionalApplicant model**

Create `models/PromotionalApplicant.go`:

```go
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
```

Note: `BeltSize` is `*string` (pointer) to handle nullable — nil means no value.

- [ ] **Step 3: Add models to AutoMigrate**

In `initializers/syncDb.go`, find the `AutoMigrate` call (around line 87). After `&models.GearItem{},` add two lines:

```go
		&models.Promotional{},
		&models.PromotionalApplicant{},
```

The full AutoMigrate call should end with:

```go
	DB.AutoMigrate(
		&models.Location{},
		&models.Class{},
		&models.EventSubType{},
		&models.EventTemplate{},
		&models.CarouselImage{},
		&models.User{},
		&models.Event{},
		&models.ClassSession{},
		&models.ClassPeriod{},
		&models.ClassAnnotation{},
		&models.Instructor{},
		&models.PasswordResetToken{},
		&models.CurrentInstructorsPage{},
		&models.GearItem{},
		&models.Promotional{},
		&models.PromotionalApplicant{},
	)
```

- [ ] **Step 4: Verify with go build**

Run: `go build ./...`
Expected: PASS

- [ ] **Step 5: Verify DB migration**

Run: `go run main.go`
Check logs for "AutoMigrate" output — should include `promotionals` and `promotional_applicants` tables
Expected: Server starts with no migration errors
