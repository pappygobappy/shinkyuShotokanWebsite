# Phase 1: Extract Business Logic from Handlers

**Priority**: ⭐ HIGHEST  
**Timeline**: 3-5 days (can be done in a single focused sprint)  
**Risk Level**: LOW — no architectural changes, just extraction and file reorganization  

---

## Overview

Handler files are currently doing too much: validating input → querying DB → rendering templates. This makes them hard to test, hard to debug, and easy to break when adding new features.

**Goal**: Extract business logic into dedicated service packages while keeping HTTP handling in handlers thin and focused.

---

## Current State Example

### Before (userHandler.go - 439 lines)
```go
func SignupPost(c *fiber.Ctx) error {
    // 1. Parse form
    var body struct {
        Email    string `form:"email"`
        Password string `form:"password"`
        FirstName string `form:"first_name"`
        LastName  string `form:"last_name"`
    }
    
    if err := c.BodyParser(&body); err != nil {
        log.Print(err)
        return err
    }
    
    // 2. Validate email format
    if !isValidEmail(body.Email) {
        return c.Status(422).Render("signup", fiber.Map{"Error": "Invalid email"})
    }
    
    // 3. Validate password strength
    if len(body.Password) < 8 {
        return c.Status(422).Render("signup", fiber.Map{"Error": "Password too short"})
    }
    
    // 4. Check if user exists (DB query)
    var existingUser models.User
    initializers.DB.Where("email = ?", body.Email).First(&existingUser)
    if existingUser.ID > 0 {
        return c.Status(422).Render("signup", fiber.Map{"Error": "Email already registered"})
    }
    
    // 5. Hash password (bcrypt)
    hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
    if err != nil {
        log.Print(err)
        return c.Status(500).Render("signup", fiber.Map{"Error": "Failed to create account"})
    }
    
    // 6. Create user (DB query)
    user := models.User{
        Email:        body.Email,
        PasswordHash: string(hash),
        FirstName:    body.FirstName,
        LastName:     body.LastName,
        Type:         "admin", // hardcoded!
    }
    result := initializers.DB.Create(&user)
    
    // 7. Generate JWT token
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "sub": user.ID,
        "exp": time.Now().Add(72 * time.Hour).Unix(),
    })
    tokenString, _ := token.SignedString([]byte(os.Getenv("HMAC_SECRET")))
    
    // 8. Set cookie
    c.Cookie(&fiber.Cookie{
        Name:     "Authorization",
        Value:    tokenString,
        Path:     "/",
        MaxAge:   259200, // 72 hours
        HTTPOnly: true,
    })
    
    return c.Redirect("/admin")
}
```

**Problems**:
- 8 different concerns mixed together (validation, DB, auth, cookies)
- Can't test password hashing without spinning up HTTP server
- Hardcoded `Type: "admin"` — where's the logic for owner vs admin?
- Error handling scattered across function

---

## After (Extracted Structure)

### File Organization
```
services/
  auth/
    signup.go          # Signup business logic
    login.go           # Login business logic  
    passwordReset.go   # Password reset flow
    token.go           # JWT generation/validation helpers
  events/              # Will be Phase 1b
  classes/             # Will be Phase 1c

handlers/
  auth.go              # Thin wrappers calling services
```

### After (auth/signup.go - ~60 lines)
```go
package auth

import (
    "shinkyuShotokan/models"
    "shinkyuShotokan/utils"
    
    "github.com/golang-jwt/jwt/v5"
    "golang.org/x/crypto/bcrypt"
)

type SignupInput struct {
    Email     string
    Password  string
    FirstName string
    LastName  string
}

type SignupResult struct {
    User      models.User
    Token     string
    ExpiresAt int // seconds until expiry
}

type AppError struct {
    Code    int
    Message string
}

func (e *AppError) Error() string { return e.Message }

var (
    ErrInvalidEmail     = &AppError{Code: 422, Message: "Invalid email format"}
    ErrPasswordTooShort = &AppError{Code: 422, Message: "Password must be at least 8 characters"}
    ErrUserExists       = &AppError{Code: 409, Message: "Email already registered"}
    ErrInternal         = &AppError{Code: 500, Message: "Failed to create account"}
)

func Signup(input SignupInput) (*SignupResult, *AppError) {
    // 1. Validate email format
    if !utils.IsValidEmail(input.Email) {
        return nil, ErrInvalidEmail
    }
    
    // 2. Validate password strength
    if len(input.Password) < 8 {
        return nil, ErrPasswordTooShort
    }
    
    // 3. Check if user exists
    var existingUser models.User
    result := initializers.DB.Where("email = ?", input.Email).First(&existingUser)
    if result.Error == nil && existingUser.ID > 0 {
        return nil, ErrUserExists
    }
    
    // 4. Hash password
    hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
    if err != nil {
        utils.Logger.WithError(err).Error("Failed to hash password")
        return nil, ErrInternal
    }
    
    // 5. Create user
    user := models.User{
        Email:        input.Email,
        PasswordHash: string(hash),
        FirstName:    input.FirstName,
        LastName:     input.LastName,
        Type:         "admin", // TODO: make configurable via input or default
    }
    
    if err := initializers.DB.Create(&user).Error; err != nil {
        utils.Logger.WithError(err).Error("Failed to create user")
        return nil, ErrInternal
    }
    
    // 6. Generate JWT token
    tokenString, err := generateToken(user.ID)
    if err != nil {
        utils.Logger.WithError(err).Error("Failed to generate token")
        return nil, ErrInternal
    }
    
    return &SignupResult{
        User:      user,
        Token:     tokenString,
        ExpiresAt: 259200, // 72 hours in seconds
    }, nil
}

func generateToken(userID uint) (string, error) {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "sub": userID,
        "exp": time.Now().Add(72 * time.Hour).Unix(),
    })
    
    secret := []byte(os.Getenv("HMAC_SECRET"))
    return token.SignedString(secret)
}
```

### After (handlers/auth.go - ~30 lines)
```go
package handlers

import (
    "shinkyuShotokan/services/auth"
    "strconv"
    
    "github.com/gofiber/fiber/v2"
)

func SignupPost(c *fiber.Ctx) error {
    var input auth.SignupInput
    
    if err := c.BodyParser(&input); err != nil {
        return handleAppError(c, &auth.AppError{Code: 422, Message: "Invalid request body"})
    }
    
    result, appErr := auth.Signup(input)
    if appErr != nil {
        return handleAppError(c, appErr)
    }
    
    // Set cookie
    c.Cookie(&fiber.Cookie{
        Name:     "Authorization",
        Value:    result.Token,
        Path:     "/",
        MaxAge:   result.ExpiresAt,
        HTTPOnly: true,
    })
    
    // Return appropriate response based on HX-request header
    hxRequest, _ := strconv.ParseBool(c.Get("hx-request"))
    if hxRequest {
        return c.Redirect("/admin")
    }
    
    return c.Redirect("/admin")
}

func handleAppError(c *fiber.Ctx, err *auth.AppError) error {
    if c.Locals("hxRequest") == true {
        // For HTMX requests, render error template
        return c.Status(err.Code).Render("partials/error", fiber.Map{"Error": err.Message})
    }
    
    // For regular requests, redirect with flash message
    return c.Status(err.Code).JSON(fiber.Map{
        "error": err.Message,
    })
}
```

---

## Implementation Steps

### Step 1: Create Service Directory Structure (5 min)
```bash
mkdir -p services/auth
mkdir -p services/events
mkdir -p services/classes
mkdir -p services/users
```

### Step 2: Define Error Types in utils/errors.go (10 min)
Create new file `utils/errors.go`:
```go
package utils

import "github.com/gofiber/fiber/v2"

type AppError struct {
    Code    int
    Message string
}

func (e *AppError) Error() string { return e.Message }

// Standard error variables
var (
    ErrNotFound      = &AppError{Code: 404, Message: "Resource not found"}
    ErrUnauthorized  = &AppError{Code: 401, Message: "Unauthorized"}
    ErrForbidden     = &AppError{Code: 403, Message: "Access denied"}
    ErrValidation    = &AppError{Code: 422, Message: "Validation failed"}
    ErrInternal      = &AppError{Code: 500, Message: "Internal server error"}
)

// Helper for structured logging (will replace log.Print())
func Logger() *logrus.Entry {
    return logrus.WithFields(logrus.Fields{
        "app": "shinkyu-shotokan",
    })
}
```

**Note**: Add `github.com/sirupsen/logrus` to go.mod later in Phase 4. For now, use `fmt.Printf()` or keep `log.Print()`.

### Step 3: Extract Signup Logic (30 min)
Create `services/auth/signup.go` with the code example above.

**Key changes to make**:
- Move validation logic into service functions (`IsValidEmail`, password length check)
- Replace direct `initializers.DB` calls with consistent error handling
- Return typed errors instead of panicking or returning bare `error`

### Step 4: Create Thin Handler Wrapper (15 min)
Create `handlers/auth.go`:
```go
package handlers

import (
    "shinkyuShotokan/services/auth"
    "strconv"
    
    "github.com/gofiber/fiber/v2"
)

func SignupGet(c *fiber.Ctx) error {
    return c.Render("signup", fiber.Map{
        "Page": structs.Page{PageName: "Sign Up"},
    })
}

func SignupPost(c *fiber.Ctx) error {
    var input auth.SignupInput
    
    if err := c.BodyParser(&input); err != nil {
        return handleAppError(c, &auth.AppError{Code: 422, Message: "Invalid request body"})
    }
    
    result, appErr := auth.Signup(input)
    if appErr != nil {
        return handleAppError(c, appErr)
    }
    
    // Set cookie
    c.Cookie(&fiber.Cookie{
        Name:     "Authorization",
        Value:    result.Token,
        Path:     "/",
        MaxAge:   result.ExpiresAt,
        HTTPOnly: true,
    })
    
    return c.Redirect("/admin")
}

func handleAppError(c *fiber.Ctx, err *auth.AppError) error {
    hxRequest := c.Locals("hxRequest") == true
    
    if hxRequest {
        return c.Status(err.Code).Render("partials/error", fiber.Map{"Error": err.Message})
    }
    
    return c.Status(err.Code).JSON(fiber.Map{
        "error": err.Message,
    })
}
```

### Step 5: Update main.go to Import New Packages (2 min)
No changes needed — Go auto-discovers packages. Just ensure you're not importing old handler functions that no longer exist.

### Step 6: Test Signup Flow (10 min)
```bash
# Start server
go run main.go

# Test with curl
curl -X POST http://localhost:8080/signup \
    -H "Content-Type: application/x-www-form-urlencoded" \
    -d "email=test@example.com&password=strongpass123&first_name=Test&last_name=User"

# Should redirect to /admin with cookie set
```

### Step 7: Delete Old Signup Code from userHandler.go (5 min)
Once verified, remove the old `SignupPost` function from `handlers/userHandler.go`. Keep it commented out for a few days as safety net.

---

## What to Extract Next (After Signup)

### Priority Order:
1. ✅ **Login** → `services/auth/login.go` + handler wrapper
2. ✅ **Password reset flow** → `services/auth/passwordReset.go`
3. ⏭️ **Event CRUD** → `services/events/eventService.go`
4. ⏭️ **Class management** → `services/classes/classService.go`

---

## Common Pitfalls to Avoid

❌ **Don't extract everything at once** — start with auth, verify it works, then move on  
❌ **Don't create circular imports** — services shouldn't import handlers; handlers import services  
❌ **Don't over-engineer interfaces** — keep service functions simple, add interfaces only if you need mocking for tests (Phase 6)  
❌ **Don't forget to handle HTMX requests** — your app uses hx-request headers for partial updates

---

## Testing Checklist

After implementing Phase 1:
- [ ] Signup flow works end-to-end
- [ ] Login flow works end-to-end
- [ ] Password reset token generation and validation works
- [ ] Error messages are user-friendly (not stack traces)
- [ ] HTMX partial renders still work after redirects
- [ ] JWT cookies are set correctly with proper flags (HTTPOnly, Secure in production)

---

## Metrics for Success

**Before Phase 1**:
```bash
wc -l handlers/*.go
# Expected: userHandler.go = 439 lines, event.go = 542 lines, etc.
```

**After Phase 1**:
```bash
wc -l handlers/auth.go          # ~30-50 lines (thin wrapper)
wc -l services/auth/signup.go   # ~60-80 lines (business logic)
wc -l services/auth/login.go    # ~40-60 lines
wc -l services/auth/passwordReset.go  # ~50-70 lines

# Total: handlers should be 2-3x smaller, services ~150-200 lines each
```

---

## Next Steps After Phase 1 Completes

1. **Verify all auth flows still work** (signup, login, logout, password reset)
2. **Run `go test ./...`** to ensure no regressions in other packages
3. **Move to Phase 2**: Add caching layer for static data (locations, classes, instructors)

---

## Questions?

If you hit any of these issues:
- "My service functions are still too big" → Split into smaller sub-functions with clear responsibilities
- "I can't test because it depends on DB" → We'll add mock interfaces in Phase 6; for now, use integration tests
- "Handler wrapper is getting fat again" → Extract more logic into the service; handler should only handle HTTP concerns

**Ready to start?** Begin by creating `utils/errors.go` and `services/auth/signup.go`, then test the signup flow. Ping me when you're ready for Phase 2 guidance.
