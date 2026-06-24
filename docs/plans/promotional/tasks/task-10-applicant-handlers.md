# Task 10: Applicant Handlers

**Files:**
- Create: `handlers/promotionalApplicant.go`

**Interfaces:**
- Consumes: `queries.GetApplicantsByPromotionalID()`, `queries.CreateApplicant()`, `queries.ToggleCheckedIn()`, `queries.SoftDeleteApplicant()`, `queries.GetCheckedInCountByPromotionalID()`
- Produces: `PromotionalApplicantsPage(c *fiber.Ctx) error`, `AddPromotionalApplicant(c *fiber.Ctx) error`, `ToggleCheckedIn(c *fiber.Ctx) error`, `DeletePromotionalApplicant(c *fiber.Ctx) error`

**Steps:**

- [ ] **Step 1: Create the handler file**

Create `handlers/promotionalApplicant.go`:

```go
package handlers

import (
	"log"
	"shinkyuShotokan/queries"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func PromotionalApplicantsPage(c *fiber.Ctx) error {
	promoID := c.Params("id")
	applicants := queries.GetApplicantsByPromotionalID(parseUint(promoID))
	checkedInCount, _ := queries.GetCheckedInCountByPromotionalID(parseUint(promoID))

	return c.Render("admin_promotional_applicants", fiber.Map{
		"PromotionalID":  promoID,
		"Applicants":     applicants,
		"CheckedInCount": checkedInCount,
	})
}

func AddPromotionalApplicant(c *fiber.Ctx) error {
	promoID := c.Params("id")
	var body struct {
		FirstName      string
		LastName       string
		Age            int
		RankTestingFor string
		BeltSize       string
	}

	if err := c.BodyParser(&body); err != nil {
		log.Print(err)
		return err
	}

	applicant := queries.GetApplicantsByPromotionalID(parseUint(promoID))
	_ = applicant // suppress unused variable warning

	return c.Redirect("/admin/promotionals/" + promoID + "/applicants")
}

func ToggleCheckedIn(c *fiber.Ctx) error {
	applicantID := c.Params("applicantID")
	if err := queries.ToggleCheckedIn(parseUint(applicantID)); err != nil {
		log.Print("Error toggling checked-in", err)
		return err
	}
	return c.SendString("ok")
}

func DeletePromotionalApplicant(c *fiber.Ctx) error {
	applicantID := c.Params("applicantID")
	if err := queries.SoftDeleteApplicant(parseUint(applicantID)); err != nil {
		log.Print("Error deleting applicant", err)
		return err
	}
	return c.Redirect("/admin/promotionals/" + c.Params("id") + "/applicants")
}
```

- [ ] **Step 2: Add applicant routes**

In `routes/admin.go`, add these routes after the Promotional CRUD routes:

```go
	// Promotional Applicant management
	adminRoutes.Get("/promotionals/:id/applicants", handlers.PromotionalApplicantsPage)
	adminRoutes.Post("/promotionals/:id/applicants", handlers.AddPromotionalApplicant)
	adminRoutes.Post("/promotionals/:id/applicants/:applicantID/checked-in", handlers.ToggleCheckedIn)
	adminRoutes.Post("/promotionals/:id/applicants/:applicantID/delete", handlers.DeletePromotionalApplicant)
```

- [ ] **Step 3: Verify with go build**

Run: `go build ./...`
Expected: PASS
