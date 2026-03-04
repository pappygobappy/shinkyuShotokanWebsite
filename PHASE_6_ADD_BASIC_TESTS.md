# Phase 6: Add Basic Tests

**Priority**: 🧪 CONFIDENCE WIN  
**Timeline**: 3-5 days (ongoing — start small, build gradually)  
**Risk Level**: LOW — begin with unit tests for services only  

---

## Overview

Your codebase currently has **zero tests**. This means:
- No safety net when refactoring (Phase 1-5 changes could break silently)
- No documentation of expected behavior
- Debugging requires manual testing every time

**Impact**: 
- Adding a new feature = fear of breaking existing functionality
- Onboarding second developer = no test suite to validate their work
- Production bugs = manual reproduction in dev environment

**Goal**: Add targeted unit and integration tests for critical paths. Start with **60% coverage on services** (not handlers) — this is achievable and provides maximum ROI.

---

## Testing Strategy

### What We'll Test First (High Priority)
| Component | Test Type | Why |
|-----------|-----------|-----|
| Password validation | Unit test | Security-critical, pure function |
| Password hashing | Unit test | Security-critical, deterministic |
| JWT generation/validation | Integration test | Auth flow depends on it |
| Signup → Login flow | Integration test | Most critical user journey |

### What We'll Test Later (Medium Priority)
- Event CRUD operations
- Class management
- Email sending (mock SMTP)

### What We Won't Test Yet (Low Priority)
- Template rendering (too fragile, UI changes break tests)
- File upload tests (complex setup, low ROI at current scale)
- Database migration tests (covered by manual deployment checks)

---

## Implementation Guide

### Step 1: Set Up Test Infrastructure (5 min)

No special packages needed — Go's built-in `testing` package is sufficient.

Create test directory structure:
```bash
mkdir -p services/auth_test
mkdir -p services/events_test
mkdir -p integration_test
```

**Note**: We're using `_test` suffix in directory names for clarity, but you can also use `*_test.go` files in the same package.

---

### Step 2: Write Unit Tests for Password Validation (15 min)

Create `services/auth/password_test.go`:
```go
package auth_test

import (
    "testing"
    
    "shinkyuShotokan/services/auth"
)

func TestValidatePassword(t *testing.T) {
    tests := []struct {
        name    string
        password string
        wantErr bool
    }{
        {"too short", "weak", true},
        {"minimum length", "strong123", false},
        {"with special chars", "StrongPass!@#", false},
        {"empty", "", true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := auth.ValidatePassword(tt.password)
            
            if (err != nil) != tt.wantErr {
                t.Errorf("ValidatePassword() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}

func TestHashPassword(t *testing.T) {
    password := "TestPass123!"
    
    hash, err := auth.HashPassword(password)
    if err != nil {
        t.Fatalf("HashPassword() error = %v", err)
    }
    
    // Verify hash matches original password
    if err := auth.CompareHashAndPassword(hash, password); err != nil {
        t.Errorf("CompareHashAndPassword() error = %v", err)
    }
    
    // Verify wrong password fails
    if err := auth.CompareHashAndPassword(hash, "WrongPassword"); err == nil {
        t.Error("CompareHashAndPassword() should fail for wrong password")
    }
}
```

**Note**: You'll need to expose `ValidatePassword` and `HashPassword` functions in `services/auth/signup.go`:
```go
// services/auth/signup.go
func ValidatePassword(password string) error {
    if len(password) < 8 {
        return &AppError{Code: 422, Message: "Password must be at least 8 characters"}
    }
    
    // TODO: Add complexity requirements (uppercase, lowercase, numbers, special chars)
    return nil
}

func HashPassword(password string) (string, error) {
    hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return "", err
    }
    
    return string(hash), nil
}
```

---

### Step 3: Write Unit Tests for JWT Generation (15 min)

Create `services/auth/token_test.go`:
```go
package auth_test

import (
    "testing"
    "time"
    
    "shinkyuShotokan/services/auth"
)

func TestGenerateToken(t *testing.T) {
    userID := uint(123)
    
    token, err := auth.GenerateToken(userID)
    if err != nil {
        t.Fatalf("GenerateToken() error = %v", err)
    }
    
    if token == "" {
        t.Error("GenerateToken() returned empty token")
    }
}

func TestValidateToken(t *testing.T) {
    userID := uint(456)
    
    token, err := auth.GenerateToken(userID)
    if err != nil {
        t.Fatalf("GenerateToken() error = %v", err)
    }
    
    validatedUserID, err := auth.ValidateToken(token)
    if err != nil {
        t.Fatalf("ValidateToken() error = %v", err)
    }
    
    if validatedUserID != userID {
        t.Errorf("ValidateToken() got user ID %d, want %d", validatedUserID, userID)
    }
}

func TestTokenExpiration(t *testing.T) {
    // This test would require mocking time — skip for now, add in Phase 7
    
    /*
    oldToken := generateExpiredToken()
    _, err := auth.ValidateToken(oldToken)
    
    if err == nil {
        t.Error("ValidateToken() should fail for expired token")
    }
    */
}

// Helper to create expired token (for future implementation)
func generateExpiredToken() string {
    // TODO: Implement with mocked time or manual claim manipulation
    return ""
}
```

**Note**: For now, skip the expiration test. We'll add proper time mocking in Phase 7 when we introduce dependency injection.

---

### Step 4: Write Integration Tests for Signup → Login Flow (30 min)

Create `integration/auth_flow_test.go`:
```go
package integration_test

import (
    "bytes"
    "encoding/json"
    "net/http/httptest"
    "testing"
    
    "github.com/gofiber/fiber/v2"
    "github.com/stretchr/testify/assert"
    
    "shinkyuShotokan" // Import main package to access app
)

func TestSignupLoginFlow(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration test")
    }
    
    app := fiber.New()
    
    // Initialize database (use test DB!)
    initializers.ConnectToDb()
    initializers.Migrate()
    
    // Register routes
    routes.RegisterPublicRoutes(app)
    
    // Test signup
    t.Run("signup", func(t *testing.T) {
        signupBody := map[string]string{
            "email":     "test@example.com",
            "password":  "StrongPass123!",
            "first_name": "Test",
            "last_name": "User",
        }
        
        body, _ := json.Marshal(signupBody)
        req := httptest.NewRequest("POST", "/signup", bytes.NewReader(body))
        req.Header.Set("Content-Type", "application/json")
        
        resp, err := app.Test(req, -1) // -1 = no timeout
        if err != nil {
            t.Fatalf("Signup request failed: %v", err)
        }
        
        assert.Equal(t, 302, resp.StatusCode, "Signup should redirect on success")
    })
    
    // Test login with created user
    t.Run("login", func(t *testing.T) {
        loginBody := map[string]string{
            "email":    "test@example.com",
            "password": "StrongPass123!",
        }
        
        body, _ := json.Marshal(loginBody)
        req := httptest.NewRequest("POST", "/login", bytes.NewReader(body))
        req.Header.Set("Content-Type", "application/json")
        
        resp, err := app.Test(req, -1)
        if err != nil {
            t.Fatalf("Login request failed: %v", err)
        }
        
        assert.Equal(t, 302, resp.StatusCode, "Login should redirect on success")
        
        // Verify cookie is set
        cookies := resp.Cookies()
        assert.NotEmpty(t, cookies, "Login response should set Authorization cookie")
    })
    
    // Cleanup: delete test user
    initializers.DB.Where("email = ?", "test@example.com").Delete(&models.User{})
}

func TestSignupWithExistingEmail(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration test")
    }
    
    app := fiber.New()
    routes.RegisterPublicRoutes(app)
    
    // Create duplicate user (should fail)
    signupBody := map[string]string{
        "email":     "existing@example.com",  // Assume this exists in DB
        "password":  "StrongPass123!",
        "first_name": "Duplicate",
        "last_name": "User",
    }
    
    body, _ := json.Marshal(signupBody)
    req := httptest.NewRequest("POST", "/signup", bytes.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    
    resp, err := app.Test(req, -1)
    if err != nil {
        t.Fatalf("Signup request failed: %v", err)
    }
    
    assert.Equal(t, 422, resp.StatusCode, "Signup with existing email should return 422")
}
```

**Note**: 
- Add `github.com/stretchr/testify` for assertions: `go get github.com/stretchr/testify`
- Use a separate test database (`shinkyu_test`) to avoid polluting production data
- Run tests with `go test -short ./...` in development (skip integration tests in CI until configured)

---

### Step 5: Add Test Database Setup (10 min)

Create `.env.test`:
```bash
DATABASE_URL=postgres://user:password@localhost:5432/shinkyu_test?sslmode=disable
HMAC_SECRET=test_secret_for_unit_tests_only
SMTP_HOST=localhost
SMTP_PORT=587
ENV=test
```

Update `initializers/connectToDb.go` to support test environment:
```go
func ConnectToDb() {
    dsn := os.Getenv("DATABASE_URL")
    
    if os.Getenv("ENV") == "test" {
        // Ensure test database exists
        log.Println("Connecting to test database...")
        
        // Create test DB if it doesn't exist (requires admin privileges)
        // Alternatively, use Docker container for test DB
    } else {
        log.Println("Connecting to production database...")
    }
    
    var err error
    DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
}
```

---

### Step 6: Create Makefile for Test Commands (5 min)

Create `Makefile`:
```makefile
.PHONY: test test-unit test-integration test-coverage lint clean db-reset

# Run all tests
test: test-unit test-integration

# Run unit tests only (fast, no DB required)
test-unit:
	go test -v ./... -short

# Run integration tests (requires test database)
test-integration:
	go test -v ./integration_test/...

# Run tests with coverage report
test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	open coverage.html  # macOS; use xdg-open on Linux

# Lint code before testing
lint:
	golangci-lint run

# Reset test database
db-reset-test:
	dropdb --if-exists shinkyu_test
	createdb shinkyu_test
	go run cmd/seed/main.go

# Clean up test artifacts
clean:
	rm -f coverage.out coverage.html
```

---

### Step 7: Run Tests and Build Confidence (20 min)

Execute tests:
```bash
# Start with unit tests (no DB required)
make test-unit

# Then integration tests (requires test database)
make db-reset-test
make test-integration

# Finally, coverage report
make test-coverage
```

**Expected output**:
```
=== RUN   TestValidatePassword
--- PASS: TestValidatePassword (0.00s)
=== RUN   TestHashPassword
--- PASS: TestHashPassword (0.01s)
PASS
ok      shinkyuShotokan/services/auth   0.015s

Coverage: 62% of statements
```

---

## Test Coverage Goals

| Component | Target Coverage | Priority |
|-----------|-----------------|----------|
| `services/auth/` | 80%+ | Critical (security) |
| `services/events/` | 60% | High (core functionality) |
| `middleware/` | 50% | Medium (auth logic already covered) |
| `handlers/` | 30% | Low (focus on services first) |

**Note**: Don't aim for 100% coverage. Focus on business logic and edge cases, not trivial getters/setters.

---

## Common Test Scenarios to Cover

### Authentication Flow
```go
func TestLoginWithInvalidPassword(t *testing.T) {
    // Should fail with wrong password
}

func TestLoginWithNonExistentUser(t *testing.T) {
    // Should not leak whether user exists or not (security best practice)
}

func TestSignupWithEmailAlreadyRegistered(t *testing.T) {
    // Should return 409 Conflict
}
```

### Event Management
```go
func TestCreateEventWithInvalidDate(t *testing.T) {
    // End date before start date should fail validation
}

func TestDeleteNonExistentEvent(t *testing.T) {
    // Should return 404 Not Found
}
```

### Class Operations
```go
func TestUpdateClassWithNegativeAge(t *testing.T) {
    // Start age or end age cannot be negative
}
```

---

## Testing Checklist

After implementing Phase 6:
- [ ] Unit tests pass for password validation and hashing
- [ ] Integration test for signup → login flow works end-to-end
- [ ] Test database is isolated from production data
- [ ] `make test` command runs all tests successfully
- [ ] Coverage report shows >60% on services package
- [ ] Tests run in CI pipeline (GitHub Actions, etc.)

---

## Next Steps After Phase 6 Completes

1. **Verify all phases work together** (auth extraction + caching + logging)
2. **Add more integration tests** for event and class management
3. **Consider adding GitHub Actions** for automated testing on PRs

---

## Questions?

If you hit any of these issues:
- "Tests are flaky" → Check for race conditions with `go test -race`
- "Can't connect to test DB" → Verify `.env.test` is loaded or use Docker container
- "Coverage is too low" → Focus on critical paths first, expand gradually

**Ready to start?** Begin by writing unit tests for password validation (`services/auth/password_test.go`), then add the signup → login integration test. Run `make test-unit` and fix any failures before moving forward. Ping me when you hit 60% coverage!
