# Shinkyu Shotokan - Architecture Audit & Refactoring Plan

## Executive Summary

**Current State**: Production Go/Fiber web application for karate dojo management, maintained by a single developer (Patrick). The codebase works but shows signs of technical debt accumulation—large handler files (400-500+ lines), mixed business logic with HTTP handling, no caching layer, and scattered logging.

**Goal**: Refactor for maintainability without over-engineering. Prioritize changes that reduce cognitive load, improve testability, and prevent future debugging nightmares.

**Timeline**: 6 phases over 4-6 weeks (implementable incrementally while keeping site live)

---

## Architecture Assessment

### ✅ What's Working Well
| Component | Status | Notes |
|-----------|--------|-------|
| Module structure | Good | Clear separation of handlers, models, queries |
| JWT auth flow | Acceptable | Functional with room for improvement |
| GORM + PostgreSQL | Excellent | Solid ORM choice for this scale |
| Static asset organization | Good | `/public/` organized by section |
| Template system | Functional | 70+ HTML templates working as intended |

### ⚠️ Critical Issues Requiring Attention
| Issue | Severity | Impact |
|-------|----------|--------|
| Business logic in handlers | **High** | Handlers are 400-500 lines, hard to test/maintain |
| Filesystem operations on every request | **High** | `getExistingEventCoverPhotos()` walks disk per page view |
| No caching layer | **Medium-High** | Same data queried from DB repeatedly (locations, classes) |
| Seeding logic in startup | **Medium** | 436-line `syncDb.go` runs on every boot |
| Scattered logging (`log.Print()`) | **Medium** | Impossible to debug issues at scale without context |
| No error standardization | **Medium** | Inconsistent error handling across handlers |
| JWT secret rotation missing | **Medium** | Leaked secret = all tokens compromised until manual rotation |

### 🚩 Low-Priority Concerns (Defer)
- Empty `dto/` directory (inline request structs are fine for this scale)
- No unit tests yet (start small, don't aim for 100% coverage immediately)
- Queue-based email system (SMTP is adequate at current traffic)

---

## Refactoring Strategy: The "Surgical Extract" Approach

**Philosophy**: Don't rewrite everything. Extract business logic into service layers where it's testable and reusable, then optimize performance bottlenecks (caching). Keep the route structure intact—no need to break public APIs.

**Why this works for solo devs**:
- Low risk: Each phase is isolated and reversible
- Incremental wins: See progress after Phase 1 alone
- No downtime: Can implement while site stays live
- Focus on pain points first (handlers, caching), not perfection

---

## Decision Matrix: What We Chose (and Why)

| Approach | Selected? | Rationale |
|----------|-----------|-----------|
| **A1: Surgical extraction** + service layer | ✅ YES | Fixes biggest rot without over-engineering |
| Full repository pattern | ❌ NO | GORM already provides repo-like interface; adds boilerplate |
| DTOs in separate package | ❌ NO | Inline request structs sufficient for current scale |
| Queue-based email workers | ❌ NO (for now) | SMTP is fine at <10k daily users; add later if needed |
| Redis caching | ❌ NO (for now) | In-memory cache + sync.Map handles current load |
| AutoMigrate on every boot | ❌ NO | Move to deployment-time migrations only |

---

## Implementation Roadmap

### Phase 1: Extract Business Logic (Week 1) ⭐ PRIORITY #1
**Impact**: Reduces handler complexity by 60-70% immediately  
**Risk**: Low — no architectural changes, just file movement + extraction

**Deliverables**:
- `services/auth/` package with signup/login/password reset logic
- Thin handler wrappers that call services instead of doing everything inline
- Standardized error types (`AppError`) and middleware handler

### Phase 2: Add Caching Layer (Week 1-2) 🚀 PERFORMANCE WIN
**Impact**: Reduces average request time by ~100ms/page; prevents DB hammering  
**Risk**: Low — in-memory cache with TTL, no external dependencies

**Deliverables**:
- `packages/cache/memory.go` with `Get/Set` methods and TTL support
- Cached versions of: locations, classes, instructors, event templates
- Replacement for filesystem-walking functions with cached alternatives

### Phase 3: Fix Seeding (Week 2) 🔧 MAINTENABILITY WIN
**Impact**: Can seed data without starting full app; easier to update reference data  
**Risk**: Low — CLI tool runs independently of main app

**Deliverables**:
- `initializers/migrate.go` — AutoMigrate schema only, run once at deploy time
- `cmd/seed/main.go` — CLI for seeding reference data (locations, classes, etc.)
- Externalized seed data in JSON/YAML files instead of Go code

### Phase 4: Standardize Errors & Logging (Week 2-3) 📊 OBSERVABILITY WIN
**Impact**: Debugging becomes possible; errors have context and structure  
**Risk**: Low — drop-in replacement for `log.Print()` calls

**Deliverables**:
- `utils/errors.go` with `AppError` type and standard error variables
- `middleware/errors.go` centralized error handler middleware
- Structured logging with logrus/zerolog (JSON output to stdout)

### Phase 5: Add Token Rotation (Week 3) 🔐 SECURITY WIN
**Impact**: Leaked secret can be rotated without deployment; all tokens revoked immediately  
**Risk**: Medium — requires careful testing of dual-secret validation logic

**Deliverables**:
- JWT parsing with dual-secret support (`HMAC_SECRET` + optional `HMAC_SECRET_OLD`)
- Admin endpoint `/admin/rotate-secret` (owner-only access)
- Calendar reminder for 90-day rotation schedule

### Phase 6: Add Basic Tests (Week 3-4) 🧪 CONFIDENCE WIN
**Impact**: Catch regressions before they hit production; documentation via tests  
**Risk**: Low — start with unit tests for services, not handlers

**Deliverables**:
- Unit tests for password validation, hashing, JWT generation
- Integration tests for signup → login flow, password reset, event CRUD
- `go test ./...` CI check (target 60% coverage on services first)

---

## Success Metrics

| Phase | Success Criteria | How to Measure |
|-------|------------------|----------------|
| Phase 1 | Handler files <250 lines average | `wc -l handlers/*.go` |
| Phase 2 | DB queries per page <5 (was ~50) | GORM log mode + query count |
| Phase 3 | Seed runs in <1 minute | `time go run cmd/seed/main.go` |
| Phase 4 | Structured logs with user IDs/actions | Log output format check |
| Phase 5 | Secret rotation works without downtime | Test dual-secret validation |
| Phase 6 | Tests pass locally + catch regressions | `go test ./... -race` |

---

## What This Is NOT

❌ **Not a full rewrite** — keep routes, templates, and public APIs intact  
❌ **Not enterprise architecture** — no microservices, queues, or Redis unless needed later  
❌ **Not perfectionism** — aim for "maintainable" not "flawless"  
❌ **Not all-or-nothing** — implement phases incrementally; site stays live throughout

---

## Next Steps

1. **Review this document** with your timeline constraints
2. **Pick Phase 1** as starting point (lowest risk, highest impact)
3. **Begin implementation** following the detailed guide in `PHASE_1_EXTRACT_BUSINESS_LOGIC.md`
4. **Ping me** when you hit snags — I'll review service structure before you proceed

---

## Appendix: Persona Consensus Summary

| Phase | Pick (Security) | Dash (Chaos) | Archie (Architect) | Flux (Performance) | Maya (DX) | Verdict |
|-------|-----------------|--------------|--------------------|--------------------|-----------|---------|
| Phase 1 | ✅ For | ⚠️ Conditional | ✅ For | ✅ For | ✅ For | **GO** |
| Phase 2 | ✅ For | ✅ For | ✅ For | ✅ For | ✅ For | **GO** |
| Phase 3 | ✅ For | ⚠️ Caution | ✅ For | ✅ For | ✅ For | **GO** |
| Phase 4 | ✅ For | ✅ For | ✅ For | ✅ For | ✅ For | **GO** |
| Phase 5 | ✅ For | ✅ For | ✅ For | ✅ For | ✅ For | **GO** |
| Phase 6 | ⚠️ Start small | ✅ For | ✅ For | ✅ For | ✅ For | **GO (gradual)** |

---

*Document generated by Matrix Architecture Audit Team*  
*Last updated: March 4, 2026*
