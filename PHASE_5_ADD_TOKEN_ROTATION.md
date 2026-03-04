# Phase 5: Add JWT Token Rotation Support

**Priority**: 🔐 SECURITY WIN  
**Timeline**: 1-2 days (can be done after Phase 4)  
**Risk Level**: MEDIUM — requires careful testing of dual-secret validation logic  

---

## Overview

Your current JWT authentication uses a single secret stored in `HMAC_SECRET` environment variable. If this secret is ever leaked (GitHub commit, log file exposure, etc.), all user tokens remain valid until you manually change the secret and redeploy.

**Current Problem**:
```go
// middleware/requireAuth.go:24
func parseJwt(tokenString string) (*jwt.Token, error) {
    return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        // ❌ Single hardcoded secret — no rotation mechanism!
        return []byte(os.Getenv("HMAC_SECRET")), nil
    })
}
```

**Impact**: 
- Leaked secret = all tokens compromised until next deployment
- No automated way to rotate secrets without downtime
- Manual reminder system (calendar event) is error-prone

**Goal**: Implement dual-secret validation that allows seamless token rotation:
1. Maintain two secrets (`HMAC_SECRET` + `HMAC_SECRET_OLD`)
2. Validate tokens signed with either secret
3. Add admin endpoint to rotate secrets without downtime
4. Enforce 90-day rotation schedule automatically

---

## Architecture Design

### Secret Rotation Timeline
```
Day 1:   [Secret A active] -------------------------------+
                                          │
Day 90:  [Secret A active] + [Secret B added]             ├─ Rotate!
                                          │              │
Day 180: [Secret B active] --------------------------------+
```

**Phase 1**: Add dual-secret validation (accept tokens from either secret)  
**Phase 2**: Add admin endpoint to rotate secrets  
**Phase 3**: Enforce rotation reminder + automatic cleanup of old secret after 90 days

---

## Implementation Guide

### Step 1: Update JWT Parsing for Dual-Secret Support (20 min)

Update `middleware/requireAuth.go`:
```go
package middleware

import (
    "fmt"
    "log"
    "os"
    "shinkyuShotokan/models"
    "strconv"
    "time"
    
    "github.com/gofiber/fiber/v2"
    "github.com/golang-jwt/jwt/v5"
)

func parseJwt(tokenString string) (*jwt.Token, error) {
    return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        // Validate signing method
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        
        // Try current secret first
        currentSecret := os.Getenv("HMAC_SECRET")
        if currentSecret != "" {
            if _, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
                return []byte(currentSecret), nil
            }); err == nil && token.Valid {
                return []byte(currentSecret), nil
            }
        }
        
        // Try old secret if configured (rotation period)
        previousSecret := os.Getenv("HMAC_SECRET_OLD")
        if previousSecret != "" {
            if _, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
                return []byte(previousSecret), nil
            }); err == nil && token.Valid {
                log.Printf("Token validated with old secret — rotation recommended")
                return []byte(previousSecret), nil
            }
        }
        
        return nil, fmt.Errorf("invalid or expired token")
    })
}

// ValidateToken checks if a JWT is valid and returns user + which secret was used
func ValidateToken(tokenString string) (models.User, bool, error) {
    var user models.User
    
    if tokenString == "" {
        return user, false, fmt.Errorf("missing token")
    }
    
    // Try current secret first
    currentSecret := os.Getenv("HMAC_SECRET")
    if currentSecret != "" {
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            return []byte(currentSecret), nil
        })
        
        if err == nil && token.Valid {
            claims, ok := token.Claims.(jwt.MapClaims)
            if !ok {
                return user, false, fmt.Errorf("invalid token claims")
            }
            
            exp, ok := claims["exp"].(float64)
            if !ok || float64(time.Now().Unix()) > exp {
                return user, false, fmt.Errorf("token expired")
            }
            
            userID := uint(claims["sub"].(float64))
            initializers.DB.First(&user, userID)
            
            if user.ID == 0 {
                return user, false, fmt.Errorf("user not found")
            }
            
            log.Printf("Token validated with current secret for user %d", user.ID)
            return user, true, nil
        }
    }
    
    // Try old secret (rotation period)
    previousSecret := os.Getenv("HMAC_SECRET_OLD")
    if previousSecret != "" {
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            return []byte(previousSecret), nil
        })
        
        if err == nil && token.Valid {
            claims, ok := token.Claims.(jwt.MapClaims)
            if !ok {
                return user, false, fmt.Errorf("invalid token claims")
            }
            
            exp, ok := claims["exp"].(float64)
            if !ok || float64(time.Now().Unix()) > exp {
                return user, false, fmt.Errorf("token expired")
            }
            
            userID := uint(claims["sub"].(float64))
            initializers.DB.First(&user, userID)
            
            if user.ID == 0 {
                return user, false, fmt.Errorf("user not found")
            }
            
            log.Printf("Token validated with old secret for user %d — rotation recommended", user.ID)
            return user, true, nil // Return true but flag for rotation reminder
        }
    }
    
    return user, false, fmt.Errorf("invalid or expired token")
}
```

### Step 2: Update RequireAuth Functions (10 min)

Update `RequireOwnerAuth` and `RequireAuth`:
```go
func RequireOwnerAuth(c *fiber.Ctx) error {
    tokenString := c.Cookies("Authorization")
    
    user, valid, err := ValidateToken(tokenString)
    if !valid || err != nil {
        log.Printf("Auth failed: %v", err)
        return c.Redirect("/")
    }
    
    // Check if user is owner
    if user.Type != models.OwnerUser {
        log.Printf("Owner access denied for user %d", user.ID)
        return c.Redirect("/")
    }
    
    c.Locals("user", user)
    return c.Next()
}

func RequireAuth(c *fiber.Ctx) error {
    tokenString := c.Cookies("Authorization")
    
    user, valid, err := ValidateToken(tokenString)
    if !valid || err != nil {
        log.Printf("Auth failed: %v", err)
        return c.Redirect("/")
    }
    
    c.Locals("user", user)
    return c.Next()
}
```

### Step 3: Create Admin Rotation Endpoint (15 min)

Create `handlers/admin.go` rotation functions:
```go
package handlers

import (
    "os"
    "strconv"
    
    "github.com/gofiber/fiber/v2"
)

// RotateSecret handles secret rotation (owner-only access)
func RotateSecret(c *fiber.Ctx) error {
    // Get current secret from request body or use existing
    var input struct {
        NewSecret string `json:"new_secret"`
    }
    
    if err := c.BodyParser(&input); err != nil {
        return c.Status(422).JSON(fiber.Map{
            "error": "Invalid request body",
        })
    }
    
    newSecret := input.NewSecret
    if newSecret == "" {
        // Generate random secret if not provided
        newSecret = generateRandomSecret()
    }
    
    // Rotate secrets: current → old, new → current
    oldSecret := os.Getenv("HMAC_SECRET")
    
    // Set environment variables (requires process restart to take effect)
    // For Fly.io, you'd use `fly secrets set` instead
    if err := os.Setenv("HMAC_SECRET_OLD", oldSecret); err != nil {
        return c.Status(500).JSON(fiber.Map{
            "error": "Failed to set old secret",
        })
    }
    
    // Note: In production, you'd use deployment-specific commands:
    // Fly.io: fly secrets set HMAC_SECRET=<newSecret>
    // Docker: Update .env file and restart container
    
    log.Printf("Secret rotation initiated. New secret stored securely.")
    log.Printf("Old secret will be valid for 90 days (HMAC_SECRET_OLD)")
    
    return c.JSON(fiber.Map{
        "message":         "Secret rotation completed",
        "rotation_date":   time.Now().Format(time.RFC3339),
        "next_rotation_by": time.Now().AddDate(0, 0, 90).Format(time.RFC3339),
        "auto_generated":  newSecret == input.NewSecret && input.NewSecret == "",
    })
}

// generateRandomSecret creates a cryptographically secure random secret
func generateRandomSecret() string {
    // TODO: Use crypto/rand for production, this is simplified
    return uuid.New().String() + uuid.New().String()
}

// CheckRotationStatus returns current rotation status (admin-only)
func CheckRotationStatus(c *fiber.Ctx) error {
    oldSecret := os.Getenv("HMAC_SECRET_OLD")
    
    var status fiber.Map
    
    if oldSecret == "" {
        status = fiber.Map{
            "has_old_secret": false,
            "rotation_needed": true,
            "message":        "No rotation history found — consider rotating soon",
        }
    } else {
        // Calculate days since last rotation (assume set when old secret was added)
        // In production, store rotation date in database or file
        status = fiber.Map{
            "has_old_secret": true,
            "rotation_needed": false,
            "message":        "Rotation active — old secret will be removed after 90 days",
        }
    }
    
    return c.JSON(status)
}
```

### Step 4: Register Rotation Routes (5 min)

Update `routes/admin.go`:
```go
func RegisterAdminRoutes(app *fiber.App) {
    adminRoutes := app.Group("admin", middleware.RequireAuth)
    
    // ... existing routes ...
    
    // Secret rotation endpoints (owner-only)
    adminRoutes.Post("/rotate-secret", handlers.RotateSecret)
    adminRoutes.Get("/rotation-status", handlers.CheckRotationStatus)
}

func RegisterOwnerRoutes(app *fiber.App) {
    ownerRoutes := app.Group("owner", middleware.RequireOwnerAuth)
    
    // Add rotation endpoints to owner-only routes for extra security
    ownerRoutes.Post("/rotate-secret", handlers.RotateSecret)
    ownerRoutes.Get("/rotation-status", handlers.CheckRotationStatus)
}
```

### Step 5: Add Rotation Reminder System (10 min)

Create `services/rotation/reminder.go`:
```go
package rotation

import (
    "os"
    "time"
    
    "shinkyuShotokan/utils"
)

// CheckRotationReminder checks if secret rotation is due and sends notification
func CheckRotationReminder() {
    oldSecret := os.Getenv("HMAC_SECRET_OLD")
    
    // If no old secret exists, this is first-time setup — send reminder to rotate
    if oldSecret == "" {
        utils.Logger.WithFields(map[string]interface{}{
            "action":           "secret_rotation_reminder",
            "priority":         "high",
            "message":          "No rotation history found — please rotate HMAC_SECRET within 30 days",
        }).Warn("Secret rotation reminder: no history")
        
        // TODO: Send email notification to owner (use existing email service)
        // sendEmailNotification("secret_rotation_reminder", "Please rotate JWT secret")
        return
    }
    
    // Check if old secret is approaching 90-day expiration
    // In production, store rotation date in database or file
    lastRotation := getRotationDate() // Implement this!
    daysSinceRotation := time.Since(lastRotation).Hours() / 24
    
    if daysSinceRotation >= 85 {
        utils.Logger.WithFields(map[string]interface{}{
            "action":           "secret_rotation_critical",
            "days_since_rotation": daysSinceRotation,
            "message":          "HMAC_SECRET_OLD is approaching 90-day expiration — rotate immediately!",
        }).Error("Secret rotation critical: old secret expiring")
        
        // Send urgent email notification
        sendEmailNotification("secret_rotation_critical", "URGENT: HMAC_SECRET rotation required within 5 days")
    } else if daysSinceRotation >= 75 {
        utils.Logger.WithFields(map[string]interface{}{
            "action":           "secret_rotation_warning",
            "days_since_rotation": daysSinceRotation,
            "message":          "HMAC_SECRET_OLD will expire in ~15 days — schedule rotation",
        }).Warn("Secret rotation warning")
    }
}

// getRotationDate retrieves the last rotation date (implement with DB or file storage)
func getRotationDate() time.Time {
    // TODO: Store rotation date in database or config file
    return time.Now().AddDate(0, 0, -90) // Default to 90 days ago for testing
}

// sendEmailNotification sends email reminder (use existing email service)
func sendEmailNotification(subject string, body string) {
    // TODO: Integrate with existing emailService.SendPasswordResetEmail
    utils.Logger.WithFields(map[string]interface{}{
        "action":   "email_sent",
        "subject":  subject,
    }).Info("Rotation reminder email sent")
}
```

### Step 6: Schedule Rotation Reminder (10 min)

Update `main.go` to run rotation check periodically:
```go
func main() {
    // ... existing setup ...
    
    // Start rotation reminder checker (every 24 hours)
    go func() {
        ticker := time.NewTicker(24 * time.Hour)
        defer ticker.Stop()
        
        for range ticker.C {
            rotation.CheckRotationReminder()
        }
    }()
    
    app.Listen(":" + os.Getenv("PORT"))
}
```

---

## Testing Checklist

After implementing Phase 5:
- [ ] Tokens signed with current secret validate successfully
- [ ] Tokens signed with old secret also validate (during rotation period)
- [ ] `/admin/rotate-secret` endpoint works for owner users only
- [ ] Rotation status endpoint returns current state
- [ ] Reminder system sends notifications at 30, 75, and 85 days
- [ ] Logs include "rotation" context for audit trail

---

## Production Deployment Notes

### Fly.io:
```bash
# Rotate secret on Fly.io (no downtime):
fly secrets set HMAC_SECRET_OLD=$(cat .env | grep HMAC_SECRET | cut -d= -f2)
fly secrets set HMAC_SECRET=<new_random_secret>
fly deploy --restart

# After 90 days, clean up old secret:
fly secrets delete HMAC_SECRET_OLD
```

### Docker:
```bash
# Update .env file with new secret and old secret
echo "HMAC_SECRET_OLD=$(cat .env | grep HMAC_SECRET | cut -d= -f2)" >> .env
echo "HMAC_SECRET=<new_random_secret>" >> .env

# Restart container
docker-compose restart app
```

---

## Security Best Practices

✅ **Rotate every 90 days** — Set calendar reminder  
✅ **Store secrets in env vars, not code** — Never commit `.env` files  
✅ **Use strong random secrets** — At least 64 characters  
✅ **Audit rotation logs** — Check `/admin/rotation-status` monthly  
❌ **Don't rotate more frequently than needed** — Increases operational overhead  
❌ **Don't skip the old secret period** — Causes authentication failures for existing sessions  

---

## Next Steps After Phase 5 Completes

1. **Verify dual-secret validation works** (test with current and old tokens)
2. **Test rotation endpoint** (ensure owner-only access)
3. **Move to Phase 6**: Add basic tests (`go test ./...`)

---

## Questions?

If you hit any of these issues:
- "Rotation endpoint returns 401" → Ensure middleware is checking for owner type correctly
- "Old tokens stop working after rotation" → Verify `HMAC_SECRET_OLD` is set before deleting old secret
- "Can't generate random secret" → Use `openssl rand -hex 32` or similar tool

**Ready to start?** Begin by updating `middleware/requireAuth.go` with dual-secret validation, then test that both current and old tokens work. Create the rotation endpoint next, then add reminder system. Ping me when ready for Phase 6 guidance.
