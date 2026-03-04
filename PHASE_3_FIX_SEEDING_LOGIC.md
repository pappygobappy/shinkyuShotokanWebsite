# Phase 3: Fix Seeding Logic

**Priority**: 🔧 MAINTENABILITY WIN  
**Timeline**: 1-2 days (can be done after Phase 1 & 2)  
**Risk Level**: LOW — CLI tool runs independently of main app  

---

## Overview

Your current seeding logic is crammed into a 436-line `initializers/syncDb.go` file that runs on **every application startup**. This means:
- Cold starts take 2-5 seconds longer than necessary
- Seed data can't be updated without modifying Go code
- No way to re-seed a test database independently

**Current Problem**:
```go
// initializers/syncDb.go - 436 lines of hell
func SyncDb() {
    DB.AutoMigrate(
        &models.CarouselImage{},
        // ... 10 more models
    )
    
    seedLocations()      // Hardcoded locations in Go code!
    seedClasses()        // More hardcoded data...
    seedEventSubTypes()  // And more...
    seedInstructors()    // A lot more...
}

func seedLocations() {
    // ❌ 10 hardcoded location structs embedded in Go code
    locations := []models.Location{
        {Name: "Municipal Services Building", Address: "..."},
        {Name: "Joseph A. Fernekes Recreation Building", Address: "..."},
        // ... 3 more
    }
    
    result := DB.Create(locations)
}
```

**Goal**: Split seeding into two separate concerns:
1. **Migrations** — Run once at deployment time (AutoMigrate only)
2. **Seed CLI** — Standalone command to populate reference data (JSON-based, editable without recompiling)

---

## Architecture Design

### New Structure
```
initializers/
  migrate.go          # AutoMigrate schema only (runs once)
  
cmd/
  seed/
    main.go           # CLI tool for seeding reference data
    
seeds/                # Externalized seed data (JSON files)
  locations.json      # Editable without recompiling!
  classes.json
  instructors.json
  event_templates.json
```

---

## Implementation Guide

### Step 1: Create Migrate Function (10 min)

Create `initializers/migrate.go`:
```go
package initializers

import (
    "log"
    "shinkyuShotokan/models"
)

// Migrate runs AutoMigrate for all models (schema only, no data seeding)
func Migrate() {
    log.Println("Running database migrations...")
    
    err := DB.AutoMigrate(
        &models.CarouselImage{},
        &models.User{},
        &models.Event{},
        &models.ClassSession{},
        &models.ClassPeriod{},
        &models.ClassAnnotation{},
        &models.Instructor{},
        &models.PasswordResetToken{},
        &models.CurrentInstructorsPage{},
        &models.Location{},
        &models.EventSubType{},
        &models.EventTemplate{},
    )
    
    if err != nil {
        log.Fatalf("Migration failed: %v", err)
    }
    
    log.Println("Database migrations completed successfully")
}
```

### Step 2: Update main.go to Call Migrate (5 min)

Update `main.go`:
```go
func init() {
    LoadEnvVariables()
    ConnectToDb()
    Migrate() // ✅ Only run schema migration, no seeding!
    
    // Don't call SyncDb anymore — it's moved to CLI tool
}
```

### Step 3: Create Seed Data JSON Files (15 min)

Create `seeds/locations.json`:
```json
[
  {
    "name": "Municipal Services Building Social Hall",
    "address": "33 Arroyo Dr\nSouth San Francisco, CA 94080",
    "google_maps_iframe": "https://www.google.com/maps/embed?pb=..."
  },
  {
    "name": "Joseph A. Fernekes Recreation Building",
    "address": "781 Tennis Dr\nSouth San Francisco, CA 94080",
    "google_maps_iframe": "https://www.google.com/maps/embed?pb=..."
  },
  {
    "name": "Westborough Recreation Building",
    "address": "2380 Galway Dr\nSouth San Francisco, CA 94080",
    "google_maps_iframe": "https://www.google.com/maps/embed?pb=..."
  }
]
```

Create `seeds/classes.json`:
```json
[
  {
    "name": "Pre-Karate",
    "description": "An introduction to the discipline of karate...",
    "get_url": "/pre-karate-class",
    "start_age": 4,
    "end_age": 8,
    "location_id": "Library | Parks & Recreation Center, Banquet Hall #130",
    "schedule": "Level 1 (Beginners) Session A: Saturday 8:30M - 9:15AM...",
    "card_photo": "/public/classes/pre-karate/card.png",
    "banner_photo": "/public/classes/pre-karate/banner.png"
  }
]
```

**Note**: Use snake_case for JSON keys to match GORM's default naming convention.

### Step 4: Create Seed CLI Tool (30 min)

Create `cmd/seed/main.go`:
```go
package main

import (
    "encoding/json"
    "log"
    "os"
    "shinkyuShotokan/models"
    "time"
    
    "gorm.io/gorm"
)

func main() {
    // Load environment variables
    if err := godotenv.Load(); err != nil {
        log.Fatal("Error loading .env file")
    }
    
    // Connect to database
    dsn := os.Getenv("DATABASE_URL")
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
    
    log.Println("Connected to database")
    
    // Check if data already exists
    var count int64
    db.Model(&models.Location{}).Count(&count)
    
    if count > 0 {
        log.Printf("Database already seeded (%d locations found), skipping...", count)
        return
    }
    
    log.Println("Seeding database with reference data...")
    
    // Seed locations
    if err := seedLocations(db); err != nil {
        log.Fatalf("Failed to seed locations: %v", err)
    }
    
    // Seed classes
    if err := seedClasses(db); err != nil {
        log.Fatalf("Failed to seed classes: %v", err)
    }
    
    // Seed instructors
    if err := seedInstructors(db); err != nil {
        log.Fatalf("Failed to seed instructors: %v", err)
    }
    
    // Seed event subtypes
    if err := seedEventSubTypes(db); err != nil {
        log.Fatalf("Failed to seed event subtypes: %v", err)
    }
    
    // Seed event templates
    if err := seedEventTemplates(db); err != nil {
        log.Fatalf("Failed to seed event templates: %v", err)
    }
    
    log.Println("✅ Database seeding completed successfully")
}

func seedLocations(db *gorm.DB) error {
    var locationsData []LocationJSON
    
    data, err := os.ReadFile("../seeds/locations.json")
    if err != nil {
        return fmt.Errorf("failed to read locations.json: %w", err)
    }
    
    if err := json.Unmarshal(data, &locationsData); err != nil {
        return fmt.Errorf("failed to parse locations.json: %w", err)
    }
    
    var locations []models.Location
    for _, loc := range locationsData {
        locations = append(locations, models.Location{
            Name:             loc.Name,
            Address:          loc.Address,
            GoogleMapsIframe: loc.GoogleMapsIframe,
        })
    }
    
    if err := db.Create(&locations).Error; err != nil {
        return fmt.Errorf("failed to create locations: %w", err)
    }
    
    log.Printf("✅ Seeded %d locations", len(locations))
    return nil
}

func seedClasses(db *gorm.DB) error {
    var classesData []ClassJSON
    
    data, err := os.ReadFile("../seeds/classes.json")
    if err != nil {
        return fmt.Errorf("failed to read classes.json: %w", err)
    }
    
    if err := json.Unmarshal(data, &classesData); err != nil {
        return fmt.Errorf("failed to parse classes.json: %w", err)
    }
    
    var classes []models.Class
    for _, cls := range classesData {
        class := models.Class{
            Name:         cls.Name,
            Description:  cls.Description,
            GetUrl:       cls.GetUrl,
            StartAge:     cls.StartAge,
            EndAge:       cls.EndAge,
            LocationID:   cls.LocationID,
            Schedule:     cls.Schedule,
        }
        
        // Set display order (1 = Pre-Karate, 2 = Youth, etc.)
        switch cls.Name {
        case "Pre-Karate":
            class.DisplayOrder = 1
        case "Youth":
            class.DisplayOrder = 2
        case "Teen":
            class.DisplayOrder = 3
        case "Adult":
            class.DisplayOrder = 4
        }
        
        classes = append(classes, class)
    }
    
    if err := db.Create(&classes).Error; err != nil {
        return fmt.Errorf("failed to create classes: %w", err)
    }
    
    log.Printf("✅ Seeded %d classes", len(classes))
    return nil
}

// Similar functions for instructors, event subtypes, templates...
```

**Note**: You'll need to define the JSON structs (`LocationJSON`, `ClassJSON`, etc.) at the bottom of the file.

### Step 5: Add JSON Struct Definitions (10 min)

Append to `cmd/seed/main.go`:
```go
type LocationJSON struct {
    Name             string `json:"name"`
    Address          string `json:"address"`
    GoogleMapsIframe string `json:"google_maps_iframe"`
}

type ClassJSON struct {
    Name         string `json:"name"`
    Description  string `json:"description"`
    GetUrl       string `json:"get_url"`
    StartAge     int    `json:"start_age"`
    EndAge       int    `json:"end_age,omitempty"`
    LocationID   string `json:"location_id"`
    Schedule     string `json:"schedule"`
}

type InstructorJSON struct {
    Name         string `json:"name"`
    PictureUrl   string `json:"picture_url"`
    Bio          string `json:"bio"`
    DisplayOrder int    `json:"display_order"`
}

// Add similar structs for event subtypes and templates...
```

### Step 6: Update go.mod Dependencies (5 min)

Add required dependencies:
```bash
go get github.com/joho/godotenv
go get gorm.io/driver/postgres
```

### Step 7: Test Seed CLI (10 min)

Run the seed tool independently:
```bash
# From project root
cd cmd/seed
go run main.go

# Expected output:
# Connected to database
# Seeding database with reference data...
# ✅ Seeded 5 locations
# ✅ Seeded 4 classes
# ✅ Seeded 5 instructors
# ...
# ✅ Database seeding completed successfully
```

---

## Usage Workflow

### For Development/Testing:
```bash
# Reset test database
dropdb shinkyu_test
createdb shinkyu_test

# Run migrations (schema only)
go run main.go  # This calls Migrate() in init()

# Seed reference data
cd cmd/seed && go run main.go
```

### For Production Deployment:
```bash
# Fly.io deployment example (fly.toml)
[env]
  SKIP_SEED = "true"  # Set this to skip seeding on deploy

# If you need to re-seed production:
fly ssh console
cd /app/cmd/seed
go run main.go
```

---

## Updating Seed Data Without Recompiling

**Before (hardcoded in Go)**:
```go
// To add a new location, you must:
1. Edit initializers/syncDb.go
2. Add struct to locations slice
3. `git commit` + `fly deploy`
4. Wait for deployment to complete
5. Hope nothing broke

# Estimated time: 15-30 minutes (including review process)
```

**After (edit JSON, run CLI)**:
```bash
# To add a new location:
1. Edit seeds/locations.json
2. `cd cmd/seed && go run main.go`
3. Done! ✅

# Estimated time: 2 minutes
```

---

## Migration Strategy for Existing Database

If you already have seed data in production, you don't need to migrate it. The CLI tool only runs if the database is empty:

```go
var count int64
db.Model(&models.Location{}).Count(&count)

if count > 0 {
    log.Printf("Database already seeded (%d locations found), skipping...")
    return
}
```

This means:
- ✅ Existing production data stays intact
- ✅ New test databases can be seeded fresh
- ✅ No risk of overwriting existing data

---

## Testing Checklist

After implementing Phase 3:
- [ ] `go run main.go` runs migrations without seeding (cold start ~2 seconds faster)
- [ ] `cd cmd/seed && go run main.go` successfully seeds reference data
- [ ] Running seed CLI twice doesn't duplicate data (idempotent check)
- [ ] Seed data can be updated by editing JSON files (no recompilation needed)
- [ ] Migrations and seeding work in both development and production environments

---

## Common Pitfalls to Avoid

❌ **Don't embed sensitive data in seed files** — API keys, passwords should come from env vars  
✅ **Seed only public reference data** — locations, classes, instructors (no user accounts)  
❌ **Don't run seed on every deployment** — use `SKIP_SEED` flag or check if DB is empty  
✅ **Test seed CLI in isolation** — ensure it works without starting the full app  

---

## Next Steps After Phase 3 Completes

1. **Verify cold start time improvement** (should be ~2-3 seconds faster)
2. **Update deployment scripts** to include seed command if needed
3. **Move to Phase 4**: Standardize errors and logging (`utils/errors.go` + structured logging)

---

## Questions?

If you hit any of these issues:
- "Seed CLI can't connect to database" → Check `DATABASE_URL` in `.env` file
- "JSON parsing fails" → Verify JSON syntax with `python -m json.tool seeds/locations.json`
- "Migrations are running every time anyway" → Make sure you removed the call to `SyncDb()` from `main.go`

**Ready to start?** Begin by creating `initializers/migrate.go`, then update `main.go` to call it. Test that migrations run without seeding data, then create the seed CLI tool. Ping me when ready for Phase 4 guidance.
