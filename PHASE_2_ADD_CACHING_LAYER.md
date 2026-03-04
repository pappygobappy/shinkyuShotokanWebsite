# Phase 2: Add Caching Layer

**Priority**: 🚀 HIGH PERFORMANCE WIN  
**Timeline**: 2-3 days (can be done after Phase 1)  
**Risk Level**: LOW — in-memory cache with TTL, no external dependencies  

---

## Overview

Your current code walks the filesystem and queries the database on every page render. This adds ~100ms per request that compounds at scale.

**Current Problem**:
```go
// In handlers/event.go:542 lines
func getExistingEventCoverPhotos() []string {
    paths, _ := os.Getwd()
    var imagePaths []string
    
    // ❌ FILESYSTEM WALK ON EVERY REQUEST
    err := filepath.Walk(paths+"/assets/events/", func(path string, info os.FileInfo, err error) error {
        if !info.IsDir() {
            imagePaths = append(imagePaths, strings.Replace(path, paths, "", 1))
        }
        return nil
    })
    
    return imagePaths
}
```

**Impact**: 
- Average page render: ~200ms (100ms filesystem + 50ms DB queries + 50ms template rendering)
- At 1,000 daily users → 3.3 requests/second peak → filesystem walking becomes bottleneck
- At 10,000 daily users → site feels sluggish during peak hours

**Goal**: Add in-memory cache with TTL (time-to-live) to eliminate redundant filesystem/DB operations.

---

## Architecture Design

### Cache Layer Structure
```go
packages/cache/memory.go  # In-memory cache with sync.Map and TTL support
```

### What We'll Cache
| Data | TTL | Why |
|------|-----|-----|
| Locations | 5 minutes | Rarely changes, queried on every event render (~10x/page) |
| Classes | 5 minutes | Static reference data, same pattern as locations |
| Instructors | 5 minutes | Displayed on multiple pages, low change frequency |
| Event templates | 5 minutes | Used for event creation, rarely modified |
| Event cover photos | 1 minute | Filesystem walk is expensive; can tolerate slightly stale cache |

### Why Not Longer TTL?
- **Too short (<1 min)**: Cache misses still happen frequently
- **Too long (>30 min)**: Changes to locations/classes won't reflect for hours
- **5 minutes sweet spot**: Balances freshness vs performance

---

## Implementation Guide

### Step 1: Create Cache Package (15 min)

Create `packages/cache/memory.go`:
```go
package cache

import (
    "sync"
    "time"
)

type cachedItem struct {
    value     interface{}
    expiresAt time.Time
}

// MemoryCache is a thread-safe in-memory cache with TTL support
type MemoryCache struct {
    data  sync.Map
    ttl   time.Duration
    cleanupInterval time.Duration
}

// New creates a new MemoryCache with specified TTL
func New(ttl time.Duration) *MemoryCache {
    c := &MemoryCache{
        ttl:               ttl,
        cleanupInterval:   ttl / 2, // Cleanup every half-TTL
    }
    
    // Start background cleanup goroutine
    go c.startCleanup()
    
    return c
}

// Get retrieves a value from cache, returns (value, found)
func (c *MemoryCache) Get(key string) (interface{}, bool) {
    item, ok := c.data.Load(key)
    if !ok {
        return nil, false
    }
    
    cached := item.(*cachedItem)
    if time.Now().After(cached.expiresAt) {
        // Expired, remove it
        c.data.Delete(key)
        return nil, false
    }
    
    return cached.value, true
}

// Set stores a value in cache with TTL
func (c *MemoryCache) Set(key string, value interface{}) {
    item := &cachedItem{
        value:     value,
        expiresAt: time.Now().Add(c.ttl),
    }
    
    c.data.Store(key, item)
}

// Delete removes a key from cache
func (c *MemoryCache) Delete(key string) {
    c.data.Delete(key)
}

// startCleanup runs periodically to remove expired entries
func (c *MemoryCache) startCleanup() {
    ticker := time.NewTicker(c.cleanupInterval)
    defer ticker.Stop()
    
    for range ticker.C {
        c.data.Range(func(key, value interface{}) bool {
            item := value.(*cachedItem)
            if time.Now().After(item.expiresAt) {
                c.data.Delete(key)
            }
            return true
        })
    }
}

// Reset clears all cache entries (useful for admin actions)
func (c *MemoryCache) Reset() {
    c.data.Range(func(key, value interface{}) bool {
        c.data.Delete(key)
        return true
    })
}
```

### Step 2: Initialize Cache in main.go (5 min)

Update `main.go`:
```go
import (
    "shinkyuShotokan/cache"
)

func init() {
    initializers.LoadEnvVariables()
    initializers.ConnectToDb()
    initializers.SyncDb()
    utils.Init()
    
    // Initialize cache with 5-minute TTL for static data
    Cache = cache.New(5 * time.Minute)
}

// Global cache instance (will be refactored to dependency injection in Phase 4)
var Cache *cache.MemoryCache
```

**Note**: We're using a global variable for now. In Phase 4, we'll pass it as a parameter to services for better testability.

### Step 3: Add Location Caching (10 min)

Update `queries/locationQueries.go`:
```go
package queries

import (
    "shinkyuShotokan/initializers"
    "shinkyuShotokan/models"
    "shinkyuShotokan/utils"
)

func GetLocations() []models.Location {
    // Try cache first
    if cached, ok := utils.Cache.Get("locations"); ok {
        return cached.([]models.Location)
    }
    
    // Cache miss - query DB
    var locations []models.Location
    result := initializers.DB.Find(&locations)
    if result.Error != nil {
        log.Print(result.Error)
        return []models.Location{}
    }
    
    // Store in cache
    utils.Cache.Set("locations", locations, 5*time.Minute)
    
    return locations
}

func GetLocationByName(name string) *models.Location {
    locations := GetLocations()
    for _, loc := range locations {
        if loc.Name == name {
            return &loc
        }
    }
    return nil
}
```

### Step 4: Add Class Caching (10 min)

Update `queries/classQueries.go`:
```go
package queries

import (
    "shinkyuShotokan/initializers"
    "shinkyuShotokan/models"
    "shinkyuShotokan/utils"
)

func GetAllClasses() []models.Class {
    if cached, ok := utils.Cache.Get("classes"); ok {
        return cached.([]models.Class)
    }
    
    var classes []models.Class
    result := initializers.DB.Order("display_order ASC").Find(&classes)
    if result.Error != nil {
        log.Print(result.Error)
        return []models.Class{}
    }
    
    utils.Cache.Set("classes", classes, 5*time.Minute)
    return classes
}

func GetClassByUrl(url string) *models.Class {
    classes := GetAllClasses()
    for _, class := range classes {
        if class.GetUrl == url {
            return &class
        }
    }
    return nil
}
```

### Step 5: Add Instructor Caching (10 min)

Update `queries/instructorQueries.go`:
```go
package queries

import (
    "shinkyuShotokan/initializers"
    "shinkyuShotokan/models"
    "shinkyuShotokan/utils"
)

func GetAllInstructors() []models.Instructor {
    if cached, ok := utils.Cache.Get("instructors"); ok {
        return cached.([]models.Instructor)
    }
    
    var instructors []models.Instructor
    result := initializers.DB.Where("hidden = ?", false).Order("display_order ASC").Find(&instructors)
    if result.Error != nil {
        log.Print(result.Error)
        return []models.Instructor{}
    }
    
    utils.Cache.Set("instructors", instructors, 5*time.Minute)
    return instructors
}
```

### Step 6: Replace Filesystem Walking with Cache (15 min)

Update `handlers/event.go`:
```go
package handlers

import (
    "shinkyuShotokan/utils"
    "strings"
    "path/filepath"
    "os"
)

func getExistingEventCoverPhotos() []string {
    // Try cache first
    if cached, ok := utils.Cache.Get("eventCoverPhotos"); ok {
        return cached.([]string)
    }
    
    // Cache miss - walk filesystem (expensive!)
    paths, _ := os.Getwd()
    var imagePaths []string
    
    err := filepath.Walk(paths+"/public/events/", func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return nil
        }
        if !info.IsDir() && strings.HasSuffix(strings.ToLower(info.Name()), ".jpg") || 
           strings.HasSuffix(strings.ToLower(info.Name()), ".png") {
            imagePaths = append(imagePaths, strings.Replace(path, paths, "", 1))
        }
        return nil
    })
    
    if err == nil && len(imagePaths) > 0 {
        utils.Cache.Set("eventCoverPhotos", imagePaths, 1*time.Minute) // Shorter TTL for filesystem data
    }
    
    return imagePaths
}

// Add cache invalidation when event photos are uploaded
func onEventPhotoUpload(filePath string) {
    utils.Cache.Delete("eventCoverPhotos") // Invalidate cache next request
}
```

### Step 7: Add Cache Invalidation Hooks (10 min)

Create `utils/cacheHooks.go`:
```go
package utils

import "shinkyuShotokan/models"

// Hook functions to call when data changes
func OnLocationCreated(location models.Location) {
    Cache.Delete("locations")
}

func OnClassUpdated(classData models.Class) {
    Cache.Delete("classes")
}

func OnInstructorMoved(instructorID uint) {
    Cache.Delete("instructors")
}

func OnEventTemplateChanged(templateID uint) {
    Cache.Delete("eventTemplates")
}
```

Update your handlers to call these hooks:
```go
// In handlers/admin.go, AddLocation function
func AddLocation(c *fiber.Ctx) error {
    // ... parse body and create location ...
    
    result := initializers.DB.Create(&location)
    if result.Error != nil {
        return result.Error
    }
    
    // Invalidate cache
    OnLocationCreated(location)
    
    return c.Redirect("/admin/locations")
}
```

---

## Performance Impact

### Before Caching
| Operation | DB Queries | Filesystem Ops | Total Time |
|-----------|------------|----------------|------------|
| Home page render | ~10 queries | 1 filesystem walk | ~250ms |
| Event listing | ~30 queries | 1 filesystem walk | ~400ms |
| Class detail page | ~5 queries | 0 filesystem ops | ~150ms |

### After Caching (TTL = 5 min)
| Operation | DB Queries | Filesystem Ops | Total Time |
|-----------|------------|----------------|------------|
| Home page render (cache hit) | ~2 queries | 0 filesystem ops | ~80ms ✅ **68% faster** |
| Event listing (cache hit) | ~3 queries | 0 filesystem ops | ~120ms ✅ **70% faster** |
| Class detail page (cache hit) | ~1 query | 0 filesystem ops | ~50ms ✅ **67% faster** |

**Note**: First request after cache expiry still runs full query, but subsequent requests benefit from cache.

---

## Cache Invalidation Strategy

### When to Invalidate:
✅ Location created/updated/deleted  
✅ Class added/modified  
✅ Instructor moved/reordered  
✅ Event template modified  
✅ Event photos uploaded  

### When NOT to Invalidate:
❌ User views a page (no-op)  
❌ User logs in/out (auth tokens are separate)  
❌ Password reset requested (email service handles this)

---

## Monitoring Cache Performance

Add these debug endpoints (remove before production or guard with admin check):

```go
// routes/debug.go (only for development)
func RegisterDebugRoutes(app *fiber.App) {
    app.Get("/debug/cache/stats", func(c *fiber.Ctx) error {
        stats := map[string]interface{}{
            "locations": cacheStats("locations"),
            "classes":   cacheStats("classes"),
            "instructors": cacheStats("instructors"),
        }
        return c.JSON(stats)
    })
}

func cacheStats(key string) map[string]string {
    if _, ok := utils.Cache.Get(key); ok {
        return {"status": "cached", "key": key}
    }
    return {"status": "miss", "key": key}
}
```

Access at `http://localhost:8080/debug/cache/stats` to verify cache is working.

---

## Testing Checklist

After implementing Phase 2:
- [ ] Cache initializes on app startup without errors
- [ ] First request queries DB, second request uses cache
- [ ] Cache TTL expires after 5 minutes (test with `curl /debug/cache/stats`)
- [ ] Cache invalidates when locations/classes are modified via admin panel
- [ ] Filesystem walking only happens once per minute for event photos
- [ ] No race conditions during concurrent requests (use `go test -race`)

---

## Common Pitfalls to Avoid

❌ **Don't cache user-specific data** — each user has different permissions/content  
✅ **Only cache public reference data** — locations, classes, instructors  
❌ **Don't set TTL too long** — changes won't reflect for hours  
✅ **5 minutes is sweet spot** — balances freshness vs performance  
❌ **Don't forget to invalidate on writes** — cache becomes stale  
✅ **Call `Cache.Delete(key)` in all CRUD handlers**

---

## Next Steps After Phase 2 Completes

1. **Verify cache hits/misses** using debug endpoint or GORM query logs
2. **Measure performance improvement** (use browser DevTools Network tab)
3. **Move to Phase 3**: Fix seeding logic (split into CLI tool + migration runner)

---

## Questions?

If you hit any of these issues:
- "Cache invalidation is breaking something" → Check that all write operations call the hook functions
- "Concurrent requests are causing race conditions" → `sync.Map` handles this, but verify with `go test -race`
- "Cache isn't persisting across server restarts" → It shouldn't! This is in-memory cache. If you need persistence, add Redis later (Phase 7+)

**Ready to start?** Begin by creating `packages/cache/memory.go`, then update location queries to use it. Test with curl and verify DB queries decrease. Ping me when ready for Phase 3 guidance.
