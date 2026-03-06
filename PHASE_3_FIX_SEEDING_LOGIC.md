# Phase 3: Externalize Seeding Logic to JSON Files

**Priority**: 🔧 MAINTENABILITY WIN  
**Timeline**: ~1 hour (simple refactor)  
**Risk Level**: LOW — No structural changes, just data externalization  

---

## Overview

Your current seeding logic is crammed into `initializers/syncDb.go` with hardcoded Go structs. The goal is to **externalize seed data to JSON files** so you can update locations, classes, instructors, etc. without recompiling and redeploying.

**Key Design Decision**: Keep the simple "one command" workflow instead of complex CLI tools or release commands. This gives you maintainability without operational complexity.

---

## Architecture Design

### New Structure
```
initializers/
  syncDb.go           # Still runs on startup, but reads JSON files
  
seeds/                # Externalized seed data (JSON files)
  locations.json      # Editable without recompiling!
  classes.json
  instructors.json
  event_templates.json
  event_subtypes.json
  carousel_images.json
```

**How It Works**:
1. `syncDb()` runs AutoMigrate for schema
2. Reads JSON files from `seeds/` directory
3. Inserts data with `ON CONFLICT DO NOTHING` (idempotent)
4. Runs on every deploy but skips if data exists

---

## Implementation Guide

### Step 1: Create JSON Seed Files (10 min)

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
    "schedule": "Level 1 (Beginners) Session A: Saturday 8:30AM...",
    "card_photo": "/public/classes/pre-karate/card.png",
    "banner_photo": "/public/classes/pre-karate/banner.png",
    "banner_adjust": 65
  }
]
```

**Note**: Use snake_case for JSON keys to match GORM's default naming convention.

### Step 2: Update syncDb.go (30 min)

Replace hardcoded structs with JSON loading functions:

```go
func seedLocations() {
    // Check if data already exists (idempotent check)
    var count int64
    DB.Model(&models.Location{}).Count(&count)
    
    if count > 0 {
        log.Println("Locations already seeded, skipping...")
        return
    }
    
    // Load from JSON file
    data, err := os.ReadFile("../seeds/locations.json")
    if err != nil {
        log.Fatalf("Failed to read locations.json: %v", err)
    }
    
    var locationsData []LocationJSON
    if err := json.Unmarshal(data, &locationsData); err != nil {
        log.Fatalf("Failed to parse locations.json: %v", err)
    }
    
    // Convert JSON structs to models
    var locations []models.Location
    for _, loc := range locationsData {
        locations = append(locations, models.Location{
            Name:             loc.Name,
            Address:          loc.Address,
            GoogleMapsIframe: loc.GoogleMapsIframe,
        })
    }
    
    // Insert with ON CONFLICT to prevent duplicates on re-deploy
    if err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&locations).Error; err != nil {
        log.Fatalf("Failed to seed locations: %v", err)
    }
    
    log.Printf("✅ Seeded %d locations", len(locations))
}

// JSON struct for unmarshaling
type LocationJSON struct {
    Name             string `json:"name"`
    Address          string `json:"address"`
    GoogleMapsIframe string `json:"google_maps_iframe"`
}
```

Repeat this pattern for:
- `seedClasses()` → `seeds/classes.json`
- `seedInstructors()` → `seeds/instructors.json`
- `seedEventSubTypes()` → `seeds/event_subtypes.json`
- `seedEventTemplates()` → `seeds/event_templates.json`
- `seedCarouselImages()` → `seeds/carousel_images.json`

### Step 3: Add Required Imports (2 min)

```go
import (
    "encoding/json"
    "os"
    
    "gorm.io/gorm/clause"
)
```

---

## Usage Workflow

### For Development/Testing:
```bash
# Edit a JSON file to add/update seed data
nano seeds/locations.json

# Redeploy - changes will be automatically picked up
fly deploy

# Or run locally to test
go run main.go
```

### For Production Deployment:
```bash
# Just edit the JSON and redeploy
fly deploy

# The app will pick up new seed data on next startup
# (if database is empty) or skip if already seeded
```

---

## Updating Seed Data Without Recompiling

**Before**:
```go
// 1. Edit initializers/syncDb.go
// 2. Add struct to locations slice
// 3. `git commit` + `fly deploy`
// 4. Wait for deployment
# Estimated time: 15-30 minutes
```

**After**:
```bash
# 1. Edit seeds/locations.json
nano seeds/locations.json
# 2. Deploy
fly deploy
# Done! ✅
# Estimated time: 2 minutes
```

---

## Idempotency

The seed functions check if data exists before running:

```go
var count int64
DB.Model(&models.Location{}).Count(&count)

if count > 0 {
    log.Println("Locations already seeded, skipping...")
    return
}
```

Plus `ON CONFLICT DO NOTHING` ensures no duplicate key errors if you're paranoid about running twice.

This means:
- ✅ Existing production data stays intact
- ✅ New test databases can be seeded fresh
- ✅ No risk of overwriting existing data
- ✅ Safe to run on every deploy

---

## Testing Checklist

After implementing Phase 3:
- [ ] `go run main.go` runs migrations and seeds data successfully
- [ ] Running the command twice doesn't duplicate data (idempotent check)
- [ ] Seed data can be updated by editing JSON files (no recompilation needed)
- [ ] Fly.io deployment picks up new seed data automatically

---

## Common Pitfalls to Avoid

❌ **Don't embed sensitive data in seed files** — API keys, passwords should come from env vars  
✅ **Seed only public reference data** — locations, classes, instructors (no user accounts)  
❌ **Don't forget the idempotency check** — you'll get duplicate key errors on re-deploy  
✅ **Validate JSON syntax before committing** — use `python -m json.tool seeds/locations.json`

---

## Next Steps After Phase 3 Completes

1. **Verify seed data loads correctly** — check that all tables have the expected records
2. **Test updating a JSON file** — add a new location and redeploy to confirm it works
3. **Move to Phase 4**: Standardize errors and logging (`utils/errors.go` + structured logging)

---

## Questions?

If you hit any of these issues:
- "JSON parsing fails" → Verify JSON syntax with `python -m json.tool seeds/locations.json`
- "Data already seeded, skipping..." → This is expected behavior if DB has existing data
- "Failed to seed locations" → Check that the JSON file exists and is valid

**Ready to start?** Begin by creating the JSON files from your existing hardcoded data in `syncDb.go`, then refactor each seed function to read from those files. Ping me when ready for Phase 4 guidance.
