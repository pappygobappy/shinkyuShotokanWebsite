# Phase 4: Standardize Errors & Logging

**Priority**: 📊 OBSERVABILITY WIN  
**Timeline**: 2-3 days (can be done after Phase 1)  
**Risk Level**: LOW — drop-in replacement for `log.Print()` calls  

---

## Overview

Your current logging approach is scattered across handlers with bare `log.Print()` calls that provide no context:
```go
// handlers/userHandler.go:247
if err != nil {
    log.Print(err)  // ❌ What error? Which user? What action failed?
}

// handlers/event.go:89
log.Print("Error creating Location", result.Error)  // ❌ Inconsistent format!
```

**Impact**: When something breaks at 3 AM, you're staring at stack traces with no context about which user triggered the error or what action was being performed.

**Goal**: Define standardized error types and add structured logging (JSON output) for better debugging and observability.

---

## Architecture Design

### New Structure
```go
utils/
  errors.go           # AppError type + standard error variables
  logger.go           # Structured logger setup
  
middleware/
  errors.go           # Centralized error handler middleware
```

### Error Standardization
Instead of returning bare `error` or panicking:
```go
type AppError struct {
    Code    int        // HTTP status code (401, 422, 500)
    Message string     // User-friendly message
    Err     error      // Underlying error (optional, for debugging)
}

func (e *AppError) Error() string { return e.Message }

// Standard errors
var (
    ErrNotFound      = &AppError{Code: 404, Message: "Resource not found"}
    ErrUnauthorized  = &AppError{Code: 401, Message: "Unauthorized"}
    ErrForbidden     = &AppError{Code: 403, Message: "Access denied"}
    ErrValidation    = &AppError{Code: 422, Message: "Validation failed"}
    ErrInternal      = &AppError{Code: 500, Message: "Internal server error"}
)
```

### Structured Logging
Instead of `log.Print()`, use logrus with context:
```go
import "github.com/sirupsen/logrus"

var logger = logrus.New()

logger.WithFields(logrus.Fields{
    "user_id":   user.ID,
    "action":    "login_failed",
    "email":     email,
}).WithError(err).Error("Authentication failed")

// Output (JSON format for production):
// {"level":"error","msg":"Authentication failed","user_id":123,"action":"login_failed","email":"test@example.com"}
```

---

## Implementation Guide

### Step 1: Create Error Types in utils/errors.go (15 min)

Create `utils/errors.go`:
```go
package utils

import "fmt"

// AppError represents a standardized application error
type AppError struct {
    Code    int         // HTTP status code
    Message string      // User-friendly message
    Err     interface{} // Underlying cause (optional)
}

func (e *AppError) Error() string {
    if e.Err != nil {
        return fmt.Sprintf("%s: %v", e.Message, e.Err)
    }
    return e.Message
}

// Helper to create AppError with underlying error
func NewAppError(code int, message string, err interface{}) *AppError {
    return &AppError{Code: code, Message: message, Err: err}
}

// Standard error variables
var (
    ErrNotFound      = &AppError{Code: 404, Message: "Resource not found"}
    ErrUnauthorized  = &AppError{Code: 401, Message: "Unauthorized"}
    ErrForbidden     = &AppError{Code: 403, Message: "Access denied"}
    ErrValidation    = &AppError{Code: 422, Message: "Validation failed"}
    ErrInternal      = &AppError{Code: 500, Message: "Internal server error"}
    ErrConflict      = &AppError{Code: 409, Message: "Resource already exists"}
)

// IsAppError checks if an error is an AppError
func IsAppError(err error) bool {
    _, ok := err.(*AppError)
    return ok
}
```

### Step 2: Create Logger Setup in utils/logger.go (10 min)

Create `utils/logger.go`:
```go
package utils

import (
    "os"
    
    "github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

func InitLogger() {
    Logger = logrus.New()
    
    // Set output to stdout (Kubernetes-friendly)
    Logger.SetOutput(os.Stdout)
    
    // Use JSON format for structured logging
    Logger.SetFormatter(&logrus.JSONFormatter{
        TimestampFormat: time.RFC3339,
        PrettyPrint:     false, // Set true for local development
    })
    
    // Set log level based on environment
    if os.Getenv("ENV") == "production" {
        Logger.SetLevel(logrus.ErrorLevel)
    } else {
        Logger.SetLevel(logrus.DebugLevel)
    }
}

// WithContext adds common fields to all logs
func WithContext(fields map[string]interface{}) *logrus.Entry {
    return Logger.WithFields(logrus.Fields{
        "app":  "shinkyu-shotokan",
        "env":  os.Getenv("ENV"),
    }.Merge(fields))
}
```

Update `main.go` to call logger initialization:
```go
func init() {
    LoadEnvVariables()
    ConnectToDb()
    Migrate()
    
    InitLogger() // ✅ Initialize structured logger
    
    Cache = cache.New(5 * time.Minute)
}
```

### Step 3: Create Error Handler Middleware (10 min)

Create `middleware/errors.go`:
```go
package middleware

import (
    "shinkyuShotokan/utils"
    
    "github.com/gofiber/fiber/v2"
)

// ErrorHandler handles errors with standardized responses
func ErrorHandler(c *fiber.Ctx) error {
    err := c.Next()
    
    if err != nil {
        // Handle AppError types
        if appErr, ok := err.(*utils.AppError); ok {
            logFields := map[string]interface{}{
                "status":  appErr.Code,
                "message": appErr.Message,
            }
            
            // Add user context if available
            if user, ok := c.Locals("user").(models.User); ok {
                logFields["user_id"] = user.ID
            }
            
            utils.Logger.WithFields(logFields).Error(appErr.Error())
            
            // Return appropriate response based on request type
            hxRequest := c.Locals("hxRequest") == true
            
            if hxRequest {
                // For HTMX requests, render error template
                return c.Status(appErr.Code).Render("partials/error", fiber.Map{
                    "Error": appErr.Message,
                })
            }
            
            // For regular/API requests, return JSON
            return c.Status(appErr.Code).JSON(fiber.Map{
                "error": appErr.Message,
            })
        }
        
        // Handle other errors (panic recovery, DB errors, etc.)
        utils.Logger.WithFields(logrus.Fields{
            "path":    c.Path(),
            "method":  c.Method(),
            "stack":   getStackTrace(),
        }).Error(err.Error())
        
        return c.Status(500).JSON(fiber.Map{
            "error": "Internal server error",
        })
    }
    
    return nil
}

// getStackTrace is a simple stack trace helper (replace with better solution in production)
func getStackTrace() string {
    // TODO: Add proper stack trace collection using runtime/debug
    return "stack trace not available"
}
```

### Step 4: Register Error Middleware in main.go (5 min)

Update `main.go`:
```go
func main() {
    engine := html.New("./templates", ".html")
    addEngineFuncs(engine)
    
    app := fiber.New(fiber.Config{
        Views:             engine,
        PassLocalsToViews: true,
        BodyLimit:         16 * 1024 * 1024,
    })
    
    // Register error handler middleware (MUST be first!)
    app.Use(middleware.ErrorHandler)
    
    // ... rest of your setup ...
}
```

### Step 5: Replace log.Print() Calls (30-60 min)

Start replacing scattered `log.Print()` calls with structured logging. Here's a pattern to follow:

**Before**:
```go
// handlers/userHandler.go:247
if err != nil {
    log.Print(err)
    return err
}
```

**After**:
```go
// handlers/auth/login.go
func LoginPost(c *fiber.Ctx) error {
    var input auth.LoginInput
    
    if err := c.BodyParser(&input); err != nil {
        utils.Logger.WithFields(logrus.Fields{
            "action": "login_parse_failed",
            "path":   c.Path(),
        }).WithError(err).Warn("Failed to parse login request")
        
        return &utils.AppError{Code: 422, Message: "Invalid request body"}
    }
    
    result, appErr := auth.Login(input)
    if appErr != nil {
        utils.Logger.WithFields(logrus.Fields{
            "action": "login_failed",
            "email":  input.Email,
        }).WithError(appErr).Warning(appErr.Message)
        
        return handleAppError(c, appErr)
    }
    
    // Success log (info level for normal operations)
    utils.Logger.WithFields(logrus.Fields{
        "action":   "login_success",
        "user_id":  result.User.ID,
        "email":    result.User.Email,
    }).Info("User logged in successfully")
    
    return c.Redirect("/admin")
}
```

**Priority replacements**:
1. ✅ Authentication failures (login, signup, password reset)
2. ✅ Database errors (query failures, constraint violations)
3. ✅ File upload/download operations
4. ✅ Email sending failures
5. ✅ Permission/access denied events

### Step 6: Add Request ID Middleware (Optional but Recommended) (10 min)

Create `middleware/requestID.go`:
```go
package middleware

import (
    "github.com/gofiber/fiber/v2"
    "github.com/google/uuid"
)

// RequestID adds a unique ID to each request for tracing
func RequestID() fiber.Handler {
    return func(c *fiber.Ctx) error {
        // Check if X-Request-ID header is set
        requestID := c.Get("X-Request-ID")
        
        if requestID == "" {
            requestID = uuid.New().String()
        }
        
        c.Set("X-Request-ID", requestID)
        c.Locals("requestID", requestID)
        
        return c.Next()
    }
}
```

Register in `main.go`:
```go
app.Use(middleware.RequestID())
```

Now all logs include the request ID:
```go
utils.Logger.WithFields(logrus.Fields{
    "request_id": c.Locals("requestID").(string),
}).Info("Processing request")
```

---

## Logging Levels Reference

| Level | When to Use | Example |
|-------|-------------|---------|
| `Debug` | Detailed debugging info (dev only) | Query parameters, raw request body |
| `Info` | Normal operations | User logged in, event created |
| `Warn` | Unexpected but handled | Rate limit approaching, deprecated API used |
| `Error` | Errors that need attention | Login failed, DB query error |
| `Fatal` | Application cannot continue | Database connection lost, config missing |

---

## Example: Complete Handler with Structured Logging

**Before** (`handlers/userHandler.go`):
```go
func SignupPost(c *fiber.Ctx) error {
    var body struct { ... }
    
    if err := c.BodyParser(&body); err != nil {
        log.Print(err)  // ❌ No context
        return err
    }
    
    if len(body.Password) < 8 {
        return c.Status(422).Render("signup", fiber.Map{"Error": "Password too short"})
    }
    
    var existingUser models.User
    initializers.DB.Where("email = ?", body.Email).First(&existingUser)
    if existingUser.ID > 0 {
        log.Print("user already exists")  // ❌ Inconsistent format!
        return c.Status(422).Render("signup", fiber.Map{"Error": "Email already registered"})
    }
    
    hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
    if err != nil {
        log.Print(err)  // ❌ No user context!
        return c.Status(500).Render("signup", fiber.Map{"Error": "Failed to create account"})
    }
    
    user := models.User{...}
    result := initializers.DB.Create(&user)
    if result.Error != nil {
        log.Print(result.Error)  // ❌ No action context!
        return result.Error
    }
    
    return c.Redirect("/admin")
}
```

**After** (extracted to `services/auth/signup.go` with logging in handler):
```go
func SignupPost(c *fiber.Ctx) error {
    var input auth.SignupInput
    
    if err := c.BodyParser(&input); err != nil {
        utils.Logger.WithFields(logrus.Fields{
            "action":  "signup_parse_failed",
            "path":    c.Path(),
            "request": c.Locals("requestID"),
        }).WithError(err).Warn("Failed to parse signup request")
        
        return &utils.AppError{Code: 422, Message: "Invalid request body"}
    }
    
    result, appErr := auth.Signup(input)
    if appErr != nil {
        utils.Logger.WithFields(logrus.Fields{
            "action":   "signup_failed",
            "email":    input.Email,
            "request":  c.Locals("requestID"),
            "status":   appErr.Code,
        }).WithError(appErr).Warning(appErr.Message)
        
        return handleAppError(c, appErr)
    }
    
    utils.Logger.WithFields(logrus.Fields{
        "action":   "signup_success",
        "user_id":  result.User.ID,
        "email":    input.Email,
        "request":  c.Locals("requestID"),
    }).Info("New user signed up")
    
    return c.Redirect("/admin")
}
```

---

## Testing Checklist

After implementing Phase 4:
- [ ] `go run main.go` starts without logger errors
- [ ] Logs output in JSON format (check with `go run main.go 2>&1 | head`)
- [ ] AppError types are returned consistently across handlers
- [ ] Error middleware handles both HTMX and regular requests
- [ ] Structured logs include user context where available
- [ ] Request ID is added to all log entries

---

## Production Logging Setup

For production deployment (Fly.io), configure log aggregation:

```toml
# fly.toml
[logs]
  format = "json"
  
[env]
  ENV = "production"
```

This enables Fly.io's built-in log viewer with JSON parsing. For more advanced setups, integrate with:
- **Datadog** — Full-stack observability
- **Sentry** — Error tracking + alerts
- **Prometheus + Grafana** — Metrics dashboard

---

## Next Steps After Phase 4 Completes

1. **Verify log output format** (should be JSON in production)
2. **Check for missing context** in error logs (user IDs, request IDs)
3. **Move to Phase 5**: Add JWT token rotation support (`middleware/requireAuth.go`)

---

## Questions?

If you hit any of these issues:
- "JSON logs are too verbose" → Lower log level to `Error` or filter debug output
- "AppError isn't being caught" → Make sure error middleware is registered first in `main.go`
- "Can't import logrus" → Run `go get github.com/sirupsen/logrus`

**Ready to start?** Begin by creating `utils/errors.go` and `utils/logger.go`, then update one handler (like login) to use structured logging. Test the output format, then proceed to replace other `log.Print()` calls. Ping me when ready for Phase 5 guidance.
