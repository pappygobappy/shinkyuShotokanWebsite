# Task 5: Promotional CRUD Routes

**Files:**
- Modify: `routes/admin.go`

**Interfaces:**
- Consumes: `handlers.AdminPromotionalsPage`, `handlers.AddPromotional`, `handlers.EditPromotionalGet`, `handlers.EditPromotionalPut`, `handlers.DeletePromotional`
- Produces: 4 new routes under `adminRoutes` group

**Steps:**

- [ ] **Step 1: Add CRUD routes**

In `routes/admin.go`, add this block after the existing Gear management routes (after line 92):

```go
	// Promotional CRUD
	adminRoutes.Post("/promotionals", handlers.AddPromotional)
	adminRoutes.Get("/promotionals/:id/edit", handlers.EditPromotionalGet)
	adminRoutes.Put("/promotionals/:id", handlers.EditPromotionalPut)
	adminRoutes.Delete("/promotionals/:id", handlers.DeletePromotional)
```

Combined with Task 1's route, the Promotional section in `routes/admin.go` should look like:

```go
	// Promotional management
	adminRoutes.Get("/promotionals", handlers.AdminPromotionalsPage)
	adminRoutes.Post("/promotionals", handlers.AddPromotional)
	adminRoutes.Get("/promotionals/:id/edit", handlers.EditPromotionalGet)
	adminRoutes.Put("/promotionals/:id", handlers.EditPromotionalPut)
	adminRoutes.Delete("/promotionals/:id", handlers.DeletePromotional)
```

- [ ] **Step 2: Verify with go build**

Run: `go build ./...`
Expected: PASS
