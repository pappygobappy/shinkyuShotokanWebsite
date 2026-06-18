# Shinkyu Shotokan Website - AI Agent Guidelines

## Project Context

**What this is**: A Go web application for Shinkyu Shotokan, a Shotokan karate dojo in South San Francisco operating since 1965. The site manages class schedules (Pre-Karate through Adult programs), tournament/event calendar, instructor profiles, student accounts, and admin tools for staff to manage content.

**Why it matters**: This isn't just another website—it's the digital home of a 60-year-old martial arts school. The tone should respect tradition while delivering modern functionality. Performance matters less than correctness: a parent checking class times or a student looking up tournament dates shouldn't hit dead ends.

---

## Go & Fiber Best Practices for This Codebase

### Architecture Patterns (What We're Doing Now)

```
HTTP Handler → Service Layer → Query Layer → Database
      ↑            ↑                ↑
    Routes    Business Logic      DB Queries (GORM)
```

**Key principles**:
- **Handlers should be thin**: They parse requests, call services, return responses. No business logic here.
- **Services own the rules**: Password validation, token generation, email sending—these live in `services/`.
- **Queries are dumb**: Just database operations. No validation, no transformations.
- **Models match the schema**: GORM structs that map directly to PostgreSQL tables.

### File Organization Rules

| Directory | Purpose | What Goes Here | What NOT to Put Here |
|-----------|---------|----------------|---------------------|
| `handlers/` | HTTP layer only | Request parsing, response formatting, status codes | Business logic, DB queries, file system ops |
| `services/` | Business rules | Signup flow, password reset, token generation, email sending | HTTP-specific code (fiber.Ctx), template rendering |
| `queries/` | Database access | GORM CRUD operations, raw SQL | Validation, business rules, external API calls |
| `models/` | Data structures | User, Event, Class, Instructor structs | Logic-heavy methods (keep these in services) |
| `middleware/` | Cross-cutting concerns | Auth checks, CORS, error handling | Request-specific logic that belongs in handlers |
| `utils/` | Shared helpers | Error types, time formatting, constants | Anything with side effects or external dependencies |

### Code Style Requirements

**Imports**: Use groupings exactly as they appear in `main.go`:
```go
import (
	"encoding/json"  // standard library first
	"log"
	"os"
	
	"github.com/gofiber/fiber/v2"  // third-party, alphabetized
	
	"shinkyuShotokan/models"  // local package (module name: shinkyuShotokan)
)
```

**Error handling**: Never ignore errors. Always check and either:
- Return them up the call stack (services → handlers)
- Log them with context (`log.Printf("login failed for %s: %v", email, err)`)
- Convert to user-friendly messages in handlers only

**Naming conventions**:
- Functions: `camelCase` starting with lowercase (`handleLogin`, `generateToken`)
- Types/Structs: `PascalCase` (`User`, `AppError`, `MemoryCache`)
- Variables: Match their purpose (`user` not `u`, `emailAddress` not `ea`)
- Files: Lowercase, descriptive (`auth.go`, not `AuthenticationHandler_v2.go`)

### Commit Message Conventions

Use [Conventional Commits](https://www.conventionalcommits.org/) format:

```
<type>(<scope>): <description>

[optional body]

[optional footer breaking/change notes]
```

**Types**:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `refactor`: Code refactoring without behavior change
- `chore`: Maintenance tasks, dependencies, config
- `perf`: Performance improvements
- `test`: Adding or updating tests

**Examples**:
```
feat(auth): add password reset email flow
fix(event): correct timezone conversion for testing events
refactor(handler): extract validation logic to service layer
chore(deps): bump fiber to v2.50.0
test(services): add login credential validation tests
```

**Best practices**:
- Use imperative mood ("add" not "added", "fix" not "fixed")
- Keep subject line under 72 characters
- Include scope when relevant (e.g., `feat(auth)`, `fix(event)`)
- Reference issues where appropriate (`Closes #123`)

### JWT & Authentication Patterns

**Current flow**:
1. User submits email/password → handler calls `services/auth/login.go`
2. Service validates credentials, generates JWT with claims `{sub: userID, exp: timestamp}`
3. Handler sets HTTP-only cookie with token
4. Middleware extracts cookie on subsequent requests, attaches user to context

**Security rules**:
- Never log full tokens or secrets (even in development)
- Use `os.Getenv("HMAC_SECRET")` only—never hardcode
- Token expiration: 30 days max (see `.env.example`)
- Password hashing must use bcrypt with cost ≥10

### Database Access Rules

**GORM usage**:
```go
// ✅ Good: prepared statements via query layer
var events []models.Event
if err := queries.DB.Find(&events).Where("start_time > ?", time.Now()).Error; err != nil {
    return utils.AppError{Code: 500, Msg: "Failed to fetch events"}
}

// ❌ Bad: raw DB queries in handlers
db.Raw("SELECT * FROM events WHERE...") // avoid unless absolutely necessary
```

**Query layer expectations**:
- Each query function returns `(result, error)`—never panic
- Use `queries.DB` (pre-configured GORM instance from initializers)
- Add indexes on frequently queried fields (`start_time`, `event_type`, `user_id`)

### Caching Strategy (Phase 2)

**What gets cached**: Static/reference data that changes infrequently:
- Locations, classes, instructors, event templates
- TTL: 5 minutes (see `main.go` cache initialization)
- Cache invalidation: Manual clear on admin updates (no complex pub/sub yet)

```go
// ✅ Good pattern from packages/cache/
key := fmt.Sprintf("location_%d", id)
if cached, exists := cache.Get(key); exists {
    return cached.(*models.Location), nil
}
// else query DB and store
cache.Set(key, location, 5*time.Minute)
```

### Logging Standards

**Before Phase 4 (current state)**: Use `log.Print()` sparingly with context
```go
log.Printf("event updated id=%d user_id=%d", eventID, userID)
```

**After Phase 4**: Migrate to structured logging (JSON output preferred for production)

---

## Existing Architecture Patterns in This Codebase

### Known Refactoring Targets

**Large handler files that need extraction to services layer:**
| File | Lines | What to Extract | Priority |
|------|-------|-----------------|----------|
| `handlers/calendar.go` | ~12,000 bytes | Calendar view logic, date range queries, color-coding rules | High |
| `handlers/event.go` | ~17,000 bytes | Event validation, template merging, timezone handling | Medium |

**Pattern**: Extract business logic to `services/`, keep handlers thin for parsing + response formatting.

### Template System Patterns

**Location**: `templates/` directory (70+ HTML files)  
**Engine**: `github.com/gofiber/template/html/v2`

**Helper functions registered in `main.go`**:
| Function | Purpose | Example Usage |
|----------|---------|---------------|
| `makeMap` | Creates maps from key/value pairs | `{{ makeMap "key1" "value1" "key2" "value2" }}` |
| `htmlRender` | Escapes HTML safely | `{{ htmlRender $event.Description }}` |
| `gmtRfc5545` | Time formatting for iCal exports | Used in calendar (.ics) file generation |
| `yahooDateFormat` | Yahoo/Outlook date format | Calendar invite URLs |
| `outlookCalInvite` | Generates Outlook calendar link URLs | `{{ outlookCalInvite $event }}` |
| `startTimePSTString` | Convert to PST timezone string | Display times in local time zone |
| `formatTimePST` | Format PST time for display | Show class/event start times |
| `isToday` | Template conditionals for event dates | `{{ if isToday $event.StartTime }}<span class="today">Today</span>{{ end }}` |
| `minus` | Arithmetic in templates | `{{ minus $event.EndTime $event.StartTime }}` (avoid complex math) |

**Template best practices**:
- Keep logic minimal: no loops deeper than 2 levels, no function calls with >2 args
- Use partials for repeated sections (`_header.html`, `_footer.html`)
- Pass only what's needed from handlers (don't dump entire structs)

### File Upload Patterns

**Upload directories** (configured via `UPLOAD_DIR` env var):
| Purpose | Path Pattern | Example |
|---------|--------------|---------|
| Carousel images | `/upload/assets/image_carousel/` | `/upload/assets/image_carousel/homepage1.jpg` |
| Event covers | `/public/events/` | `/public/events/tournament-2024.png` |
| Class photos | `/public/classes/` | `/public/classes/pre-karate-banner.jpg` |
| Instructor headshots | `/public/instructors/` | `/public/instructors/sensei-sue.jpg` |

**Key patterns**:
- Uploaded images go to `upload/assets/` directory (persistent volume on Fly.io)
- Static assets served from `public/` directory (embedded in binary)
- Always store relative paths: `/public/events/tournament.png`, not absolute filesystem paths
- Use `filesService.go` helpers for file operations

### Event Template System

**Many-to-Many relationship**: `EventTemplate` ↔ `EventSubType` (tournament divisions)

```go
// EventTemplate struct has this association
type EventTemplate struct {
    ID            uint
    Name          string // "Tournament", "Promotional"
    EventSubTypes []EventSubType // Many-to-many via junction table
}

// EventSubType represents tournament divisions
type EventSubType struct {
    ID   uint
    Name string // "Pre-Karate", "Youth & Adult", "All Ages"
}
```

**Usage**: When creating an event, admins select a template which pre-fills:
- Default start/end times (1970-01-01 placeholders for recurring events)
- Check-in time settings
- Description with `%s` placeholders for dynamic content
- Location reference
- Associated divisions/subtypes

### Middleware Patterns

**Current middleware in `middleware/requireAuth.go`**:
```go
// Three types:
RequireOwnerAuth()  // Only users with Type=owner can access
RequireAuth()       // Any authenticated user
AttachUser()        // Optional: attach user if logged in, otherwise continue anonymously
```

---

## Business Domain Knowledge

### User Types & Permissions

| Type | Access Level | Use Case |
|------|--------------|----------|
| `owner` | Full admin access | Dojo owners manage all content, users, events |
| `admin` | (future) Limited admin | Staff who can update classes/events but not delete accounts |
| `student` | Read-only + profile | View schedules, register for events, manage own account |

### Event Types & Workflow

**Event categories**:
- Tournament: External competition with registration
- Testing: Belt promotion testing (requires prerequisites)
- Seminar: Guest instructor workshop
- Social: Dojo gathering, potluck, etc.

**Key fields**: `title`, `startTime`, `endTime`, `location`, `event_type`, `description`, `prerequisites` (for testing events)

### Class Structure

**Programs by age group**:
- Pre-Karate (ages 4-6)
- Youth Karate (ages 7-12)
- Teen Karate (ages 13-17)
- Adult Karate (ages 18+)

**Schedule attributes**: `day_of_week`, `start_time`, `end_time`, `location_id`, `instructor_id`

### Email Flows

**Password reset flow** (already implemented):
1. User submits email on `/forgot-password` page
2. Service generates token, stores in DB with TTL
3. Email sent with link: `/reset-password?token=xxx`
4. User submits new password with valid token
5. Password updated, token invalidated, login redirect

---

## Routes Overview

**Public routes** (no auth required, user attached if logged in):
| Path | Handler | Purpose |
|------|---------|---------|
| `/` | `Home()` | Homepage with carousel and upcoming events |
| `/instructors`, `/instructors/:id` | `Instructors()`, `SenseiSue()` | Instructor listing and bios |
| `/history` | `History()` | Dojo history page |
| `/events/:id` | `Event()` | Individual event detail page |
| `/requirements/:rank` | `Requirements()` | Belt requirement pages (10th kyu through 1st dan) |
| `/contact-us` | `ContactUs()` | Contact form page |
| `/calendar`, `/calendar/:id` | `Calendar()`, `CalendarItemView()` | Full calendar view and single event context |
| `/pre-karate-class`, `/youth-class`, etc. | `Classes()` | Individual class pages (dynamic from classes.json) |

**Authentication routes**:
| Path | Handler | Purpose |
|------|---------|---------|
| `/login` | `LoginGet()`, `LoginPost()` | Login page and form submission |
| `/signup` | `SignupGet()`, `SignupPost()` | Registration page and account creation |
| `/forgot-password`, `/reset-password/:token` | Password reset flow | Token-based password recovery |

**Admin routes** (require authentication):
- Location management: `/admin/locations/*`
- Class management: `/admin/classes/:id`, `/admin/calendar/:id`, `/admin/classPeriod`
- Event CRUD: `/admin/events/*`, `/admin/event-templates/*`
- Instructor management: `/admin/instructors/*`, `/admin/instructors/upload-page-image`
- Carousel images: `/admin/carousel-images/*`, `/admin/upload-carousel-image`
- User profile: `/admin/userProfile*`

**Owner-only routes**:
| Path | Handler | Purpose |
|------|---------|---------|
| `/owner/users` | `AdminUsersPage()` | List all users in system |

### Common Tasks & How to Approach Them

### Adding a New Public Page (e.g., "About Us")

1. **Create template**: `templates/about.html` with existing layout structure
2. **Add handler**: `handlers/about.go`
3. **Register route**: `routes/public.go` → `app.Get("/about", handlers.About)`
4. **Keep it simple**: No service layer needed unless fetching data from DB

### Adding a New Admin Feature (e.g., "Edit Instructor")

1. **Model check**: Does `models.Instructor` have all required fields? Add if needed.
2. **Query layer**: Create `queries/instructor.go` functions for CRUD operations
3. **Service layer** (if business logic exists): e.g., validation rules, notifications
4. **Handler**: Parse form, call service/query, return success/error response
5. **Route**: Register in `routes/admin.go` with `RequireOwnerAuth()` middleware
6. **Template**: Add edit form to existing template or create new one

**Carousel Image Soft-Delete Pattern**:
```go
// Three operations for carousel images:
POST /admin/carousel-images/:id/remove   // Soft delete (sets DeletedAt)
POST /admin/carousel-images/:id/restore  // Restore from soft delete
POST /admin/carousel-images/:id/hard-delete // Permanently remove from DB
```

### Fixing a Bug in Authentication Flow

**Debugging checklist**:
1. Check logs: `log.Print()` statements in handlers/services
2. Verify token: Use JWT debugger to inspect decoded token claims
3. Database check: Confirm user exists, password hash matches bcrypt cost 10
4. Middleware trace: Add temporary log in `requireAuth.go` to see where flow breaks

**Common authentication bugs**:
- Token expiration too short (check `exp` claim vs. actual time)
- Wrong secret key mismatch (`HMAC_SECRET` env var not set correctly)
- Cookie domain/path issues (affects cross-origin requests)

### Performance Issue: Slow Page Load

**Diagnosis steps**:
1. Check query count in logs (Phase 4 logging will show this)
2. Profile cache hits: Is static data being queried repeatedly?
3. Filesystem ops: Does page walk directories for photos/covers?
4. Template complexity: Are there nested loops rendering thousands of items?

**Optimization priorities**:
1. Add caching for repeated DB queries (Phase 2)
2. Pre-render templates that don't change often
3. Lazy-load images and defer non-critical JS

---

## Testing Guidelines

**Priority order**:
1. **Services layer** (pure functions, no HTTP dependencies): password validation, token generation, email formatting
2. **Query layer** (database operations): can be mocked with test DB
3. **Handlers** (last priority): use fiber's `test.Ctx` for integration tests

**Test file naming**:
```
services/auth/login_test.go      # Tests for login service
queries/event_test.go            # Tests for event queries
middleware/requireAuth_test.go   # Tests for auth middleware
handlers/admin_test.go           # Integration tests for admin endpoints
```

---

## Deployment & Environment

### Fly.io Configuration

**Key file**: `fly.toml`  
**Region**: `sjc` (San Jose for low latency on West Coast)  
**Scaling**: Single instance adequate for current traffic; add more if needed

**Required secrets** (set in Fly.io dashboard):
- `HMAC_SECRET`: 64-character hex string (generate with `openssl rand -hex 32`)
- `SMTP_USERNAME`, `SMTP_PASSWORD`: Gmail app password or other SMTP credentials
- `UPLOAD_DIR`: Absolute path for file uploads (e.g., `/data/uploads`)

### Local Development Setup

1. **PostgreSQL**: Run locally or use Docker (`docker run --name shinkyu-db -e POSTGRES_PASSWORD=devpass -p 5432:5432 postgres`)
2. **Environment variables**: Copy `.env.example` with all required keys
3. **Run migrations & seed**: `go run main.go` auto-runs on startup via `initializers/syncDb.go`

---

## What to Avoid (Anti-Patterns)

### ❌ Don't Add Logic to Handlers

**Bad**:
```go
func CreateEvent(c *fiber.Ctx) error {
    // Business logic here = bad
    if event.StartTime.Before(time.Now()) {
        return c.Status(400).SendString("Cannot create past events")
    }
    if len(event.Title) < 3 {
        return c.Status(400).SendString("Title too short")
    }
    // ... DB query ...
}
```

**Good**:
```go
func CreateEvent(c *fiber.Ctx) error {
    var req CreateEventRequest
    if err := c.Bind(&req); err != nil {
        return utils.SendError(c, 400, "Invalid request")
    }

    event, err := services.CreateEvent(req) // validation happens here
    if err != nil {
        return utils.SendError(c, err.Code, err.Msg)
    }

    return c.JSON(event)
}
```

### ❌ Don't Mix Concerns Across Layers

- **Handlers** don't query DB directly (use `queries/`)
- **Services** don't render templates or set cookies
- **Queries** don't validate business rules (e.g., "can this student register?")
- **Models** are dumb data containers (no methods that call external APIs)

### ❌ Don't Over-Engineer Early

This is a dojo management site, not AWS. Avoid:
- Microservices architecture (single monolith is fine)
- Redis/queue systems until you have scaling problems
- Complex caching strategies (in-memory + TTL works for now)
- Full test coverage on day one (start with critical paths)

### ❌ Don't Ignore the Business Context

Remember: karate instructors and dojo staff are your real users. Features should:
- Work offline-first where possible (students might check schedules at tournaments without WiFi)
- Be mobile-friendly (many parents use phones to look up class times)
- Respect tradition (don't make dramatic UI changes that confuse long-time students)

---

## Quick Reference

### Useful Commands

```bash
# Start development server
go run main.go

# Run tests with race detector
go test ./... -race

# Format code (Go standard)
gofmt -w .

# Check for issues
go vet ./...
```

### Key File Locations

| Concern | Files |
|---------|-------|
| App entry point | `main.go` |
| Route registration | `routes/*.go` |
| HTTP handlers | `handlers/*.go` |
| Business logic | `services/auth/*.go`, `services/emailService.go` |
| Database queries | `queries/*.go` |
| Data models | `models/*.go` |
| Middleware | `middleware/requireAuth.go` |
| Templates | `templates/*.html` |
| Static assets | `public/`, `/upload/` (from env) |

### Common Error Codes & Messages

| Code | Meaning | When to Use |
|------|---------|-------------|
| 400 | Bad Request | Invalid form data, missing required fields |
| 401 | Unauthorized | Wrong password, expired token, no auth cookie |
| 403 | Forbidden | User lacks permission (e.g., student accessing admin) |
| 404 | Not Found | Event/class/instructor doesn't exist |
| 500 | Internal Server Error | DB connection failed, unexpected panic in service |

---

*Built with Go, Fiber, and respect for Shotokan tradition.*
