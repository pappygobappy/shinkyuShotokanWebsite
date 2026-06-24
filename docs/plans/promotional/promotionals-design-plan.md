# Promotionals - Design Plan

## Overview

Standalone `Promotional` model for managing belt testing sessions. Each promotional has applicants who can be added/managed by admin, marked checked-in at the door, and exported to a printable roster. Decoupled from the existing `Event` model.

---

## 1. New Models

### `Promotional`

```
Table: promotionals

Columns:
  id              bigint (PK, auto-increment)
  created_at      timestamptz
  updated_at      timestamptz
  deleted_at      timestamptz          (soft delete)
  name            text
  date            date
  start_time      time
  end_time        time
  check_in_time   time
  subtype         text                 (e.g. "Pre-Karate", "Youth & Adult")
  location        text
  description     text
```

**Fields mirror the current promotional event template** but as a first-class model. The existing `Event` model and "Promotional" event type in the template system stay as-is for public-facing event pages.

### `PromotionalApplicant`

```
Table: promotional_applicants

Columns:
  id                 bigint (PK, auto-increment)
  created_at         timestamptz
  updated_at         timestamptz
  deleted_at         timestamptz          (soft delete)
  promotional_id     bigint               (FK → promotionals, cascade delete)
  first_name         text
  last_name          text
  age                integer
  rank_testing_for   text
  belt_size          text                  (nullable)
  checked_in         boolean   DEFAULT false

Junction Table: promotional_applicant_instructors
  promotional_applicant_id  bigint  (FK → promotional_applicants)
  instructor_id             bigint  (FK → instructors)
```

**Design decisions:**

- `RankTestingFor` is a string, not a FK. Ranks aren't a table in the current schema.
- `BeltSize` is nullable — not all applicants need one.
- `CheckedIn` is a boolean — simple true/false for day-of check-in.
- `Instructors` is many-to-many via GORM convention junction table `promotional_applicant_instructors`.
- Soft delete via `gorm.Model.DeletedAt`.

---

## 2. Admin Pages & Routes

### Promotionals CRUD

| Method | Route | Handler | Purpose |
|----------|-------|---------|---------|
| GET | `/admin/promotionals` | `AdminPromotionalsPage` | List all promotionals |
| POST | `/admin/promotionals` | `AddPromotional` | Create a new promotional |
| GET | `/admin/promotionals/:id/edit` | `EditPromotionalGet` | Edit a promotional |
| PUT | `/admin/promotionals/:id` | `EditPromotionalPut` | Update a promotional |
| DELETE | `/admin/promotionals/:id` | `DeletePromotional` | Soft-delete a promotional |

### Applicants (nested under promotional)

| Method | Route | Handler | Purpose |
|----------|-------|---------|---------|
| GET | `/admin/promotionals/:id/applicants` | `PromotionalApplicantsPage` | List applicants |
| POST | `/admin/promotionals/:id/applicants` | `AddPromotionalApplicant` | Add a new applicant |
| POST | `/admin/promotionals/:id/applicants/:applicantID/checked-in` | `ToggleCheckedIn` | Toggle checked-in status |
| POST | `/admin/promotionals/:id/applicants/:applicantID/delete` | `DeletePromotionalApplicant` | Soft-delete an applicant |
| GET | `/admin/promotionals/:id/applicants/print` | `PrintPromotionalApplicantsPage` | Print-optimized roster |

---

## 3. Model Files

### `models/Promotional.go`

```go
package models

import "gorm.io/gorm"

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

### `models/PromotionalApplicant.go`

```go
package models

import "gorm.io/gorm"

type PromotionalApplicant struct {
    gorm.Model
    PromotionalID  uint
    FirstName      string
    LastName       string
    Age            int
    RankTestingFor string
    BeltSize       string
    CheckedIn      bool
    Instructors    []Instructor `gorm:"many2many:promotional_applicant_instructors;"`
}

func (PromotionalApplicant) TableName() string {
    return "promotional_applicants"
}
```

---

## 4. Query Files

### `queries/promotional.go`

```go
package queries

func GetAllPromotionals() []models.Promotional
func GetPromotionalByID(id uint) models.Promotional
func CreatePromotional(p *models.Promotional) error
func UpdatePromotional(p *models.Promotional) error
func SoftDeletePromotional(id uint) error
func RestorePromotional(id uint) error
```

### `queries/promotionalApplicant.go`

```go
package queries

func GetApplicantsByPromotionalID(promoID uint) []models.PromotionalApplicant
func CreateApplicant(applicant *models.PromotionalApplicant) error
func UpdateApplicant(applicant *models.PromotionalApplicant) error
func SoftDeleteApplicant(id uint) error
func RestoreApplicant(id uint) error
func ToggleCheckedIn(id uint) error
func GetCheckedInCountByPromotionalID(promoID uint) (int, error)
```

---

## 5. Handler Files

### `handlers/promotional.go`

```go
func AdminPromotionalsPage(c *fiber.Ctx) error
func AddPromotional(c *fiber.Ctx) error
func EditPromotionalGet(c *fiber.Ctx) error
func EditPromotionalPut(c *fiber.Ctx) error
func DeletePromotional(c *fiber.Ctx) error
```

### `handlers/promotionalApplicant.go`

```go
func PromotionalApplicantsPage(c *fiber.Ctx) error
func AddPromotionalApplicant(c *fiber.Ctx) error
func ToggleCheckedIn(c *fiber.Ctx) error
func DeletePromotionalApplicant(c *fiber.Ctx) error
func PrintPromotionalApplicantsPage(c *fiber.Ctx) error
```

---

## 6. Templates

### `templates/admin_promotionals.html`

List all promotionals:
- Table: Name | Date | SubType | Location | Applicants | Actions
- "Add Promotional" button → opens modal or navigates to edit form
- Date range filter (upcoming, past)

### `templates/admin_promotional_edit.html`

Edit form for a promotional:
- Fields: name, date, start_time, end_time, check_in_time, subtype (dropdown), location, description
- "Back to Promotionals" link
- If promotional has applicants, show count + "Manage Applicants" link

### `templates/admin_promotional_applicants.html`

Main applicant list view:
- Promotional info header (name, date, location)
- Applicant table: Name | Age | Rank | Belt Size | Instructors | Checked In
- Inline checkbox per row to toggle checked-in via HTMX
- Inline form (or modal) to add a new applicant
- "Export Printable Sheet" button → links to `/print` variant
- Sort by checked-in status (unchecked first by default)

### `templates/admin_promotional_applicants_print.html`

Print-optimized variant:
- Strips all nav/sidebar
- Shows: promotional name, date, location
- Full applicant table
- CSS `@media print` for clean page breaks
- "Checked In" column with empty checkboxes for manual marking at the door

---

## 7. Valid Ranks (for UI dropdown)

```
10th kyu
9th kyu
8th kyu
7th kyu
6th kyu
5th kyu
4th kyu
3rd kyu
2nd kyu
1st kyu
1st dan
2nd dan
3rd dan
4th dan
5th dan
6th dan
7th dan
8th dan
9th dan
10th dan
```

---

## 8. Files to Create/Modify

| Action | File |
|----------|------|
| Create | `models/Promotional.go` |
| Create | `models/PromotionalApplicant.go` |
| Create | `queries/promotional.go` |
| Create | `queries/promotionalApplicant.go` |
| Create | `handlers/promotional.go` |
| Create | `handlers/promotionalApplicant.go` |
| Create | `templates/admin_promotionals.html` |
| Create | `templates/admin_promotional_edit.html` |
| Create | `templates/admin_promotional_applicants.html` |
| Create | `templates/admin_promotional_applicants_print.html` |
| Modify | `initializers/syncDb.go` — add both models to AutoMigrate |
| Modify | `routes/admin.go` — register all routes |

---

## 9. Design Decisions Summary

| Question | Decision | Rationale |
|----------|----------|-----------|
| Tie to existing Event model? | No, standalone model | Cleaner separation; promotionals have different data needs than public events |
| Rank stored as FK or string? | String | Ranks aren't a table in the schema; adding one is scope creep |
| Can a student test multiple ranks? | No, one per applicant | Per user requirement |
| Payment tracking? | No | Per user requirement |
| Belt size for all? | Optional field | Only needed for Dan tests where a gi is given |
| Print export PDF or HTML? | HTML with print CSS | Simpler, no external deps, works in all browsers |
| Checked-in column in print? | Empty checkboxes | Lets staff mark off manually at the door |
| Soft delete or hard delete? | Soft delete | Accidental deletions can be restored |
| Instructors many-to-many? | Yes via junction table | Each applicant can have multiple instructors |
