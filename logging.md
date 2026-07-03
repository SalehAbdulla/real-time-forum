# Logging Guide

## Overview

This project uses Go's standard `log/slog` package for structured logging. The logger is initialized in `cmd/main.go` via `pkg/logger.InitLogger()` and stored in `app.Logger`.

## Log Level Configuration

| `LOG_LEVEL` value | Effective `slog` level | Visible messages |
|------------------|----------------------|-----------------|
| (unset/empty)    | `slog.LevelInfo`      | INFO, WARN, ERROR |
| `"debug"`        | `slog.LevelDebug`     | DEBUG, INFO, WARN, ERROR |
| `"warn"`         | `slog.LevelWarn`      | WARN, ERROR |
| `"error"`        | `slog.LevelError`     | ERROR only |

### How to set `LOG_LEVEL`

**For one run:**
```sh
LOG_LEVEL=debug go run ./cmd
```

**For the current shell session:**
```sh
export LOG_LEVEL=debug
go run ./cmd
```

**Persistently (add to your `.zshrc` or `.bash_profile`):**
```sh
echo 'export LOG_LEVEL=debug' >> ~/.zshrc
source ~/.zshrc
```

---

## Log Output Format

### Development (`app.InProduction = false` — default)
Human-readable text output. Source file and line number are included.

```
time=2026-07-03T19:23:01.669+03:00 level=WARN source=/path/to/utils.go:69 msg="request error" error="email already exists" method=POST path=/api/v1/auth/register status=400
```

### Production (`app.InProduction = true`)
JSON output for machine parsing (log aggregators, etc.).

```json
{"time":"2026-07-03T19:23:01.669+03:00","level":"WARN","source":"/path/to/utils.go:69","msg":"request error","error":"email already exists","method":"POST","path":"/api/v1/auth/register","status":400}
```

---

## Log Levels & When to Use Them

| `slog` Level | When to use |
|-------------|-------------|
| `slog.Debug` | Detailed diagnostic information during development. Input format violations, parameter values, internal checks. |
| `slog.Info` | Success events and notable state changes. User registered, logged in, profile updated. |
| `slog.Warn` | Expected/validation failures that are the **client's** fault. Bad request, duplicate email, missing fields, invalid credentials. These are not system problems. |
| `slog.Error` | True internal/server errors where something is **genuinely broken**. Database connection failed, password hash failed, unexpected nil pointer. |

### Stack Trace Behavior

The `stackTraceHandler` in `pkg/logger/logger.go` automatically attaches a goroutine stack trace to **every `ERROR`-level log entry**. This is useful for debugging internal failures but would be noise for validation errors.

- `DEBUG`, `INFO`, `WARN` → **No stack trace**
- `ERROR` → **Stack trace included automatically**

---

## How to Log at Each Layer

### 1. Service Layer (`pkg/service/`)

Use `slog.Debug` for validation violations — these are developer-only details:

```go
slog.Debug("invalid email format", "email", email)
slog.Debug("weak password")
slog.Debug("invalid nickname length", "length", len(nickName))
slog.Debug("passwords do not match")
slog.Debug("invalid age", "age", ageStr)
slog.Debug("invalid gender", "gender", gender)
```

Use `slog.Info` for login audit events (non-critical but notable):

```go
slog.Info("invalid password attempt", "user_id", userID)
slog.Info("login credential lookup failed", "identifier", identifier, "error", err)
```

Use `slog.Error` for true internal failures:

```go
slog.Error("failed to hash password", "error", err)
```

**Do NOT log expected validation errors** (email exists, nickname taken, bad request) in the service layer — just return the error. The handler layer will log it once with full request context.

```go
// GOOD — just return the error
if err := s.db.DoesEmailExists(email); err != nil {
    return 0, err
}
```

```go
// BAD — don't log expected business errors in the service
if err := s.db.DoesEmailExists(email); err != nil {
    slog.Info("email already exists", "email", email)  // NO — handler will log this
    return 0, err
}
```

### 2. Handler Layer (`pkg/handlers/`)

The handler is the **single point of error logging**. Always delegate to `HandleError`:

```go
userID, err := re.AuthService.Register(inputs)
if err != nil {
    re.HandleError(w, r, err)   // logs once with request context, sends response
    return
}

// Success — log at INFO
re.App.Logger.Info("user registered successfully",
    "user_id", userID,
    "email", email,
)
```

**Never** log the error yourself before calling `HandleError`:

```go
// BAD — duplicate logging
re.App.Logger.Info("registration failed", "error", err.Error())  // NO
re.HandleError(w, r, err)

// GOOD — single logging point
re.HandleError(w, r, err)
```

### 3. Repository Layer (`pkg/repositories/`)

Use `slog.Error` for actual database failures:

```go
if err != nil {
    slog.Error("failed to insert user into database",
        "email", email,
        "error", err,
    )
    return 0, realtimeforum.ErrInternal
}
```

Expected cases (no rows found, duplicate entry) should **not** be logged here — they are surfaced via sentinel errors and logged by `HandleError` in the handler.

---

## `HandleError` — The Error Response Utility

`pkg/handlers/utils.go` — `HandleError` maps errors to the correct HTTP status code, log level, and produces a single log line.

### Error-to-Status Mapping

| Error constant | HTTP status | Log level | Stack trace |
|---------------|-------------|-----------|-------------|
| `ErrBadRequest` | 400 Bad Request | WARN | No |
| `ErrInvalidEmail` | 400 Bad Request | WARN | No |
| `ErrEmailExists` | 400 Bad Request | WARN | No |
| `ErrNickName` | 400 Bad Request | WARN | No |
| `ErrNickNameLength` | 400 Bad Request | WARN | No |
| `ErrPasswordLength` | 400 Bad Request | WARN | No |
| `ErrPasswordsDontMatch` | 400 Bad Request | WARN | No |
| `ErrInvalidPassForm` | 400 Bad Request | WARN | No |
| `ErrInvalidAge` | 400 Bad Request | WARN | No |
| `ErrGender` | 400 Bad Request | WARN | No |
| `ErrInvalidCredentials` | 400 Bad Request | WARN | No |
| `ErrUnauthorized` | 401 Unauthorized | WARN | No |
| `ErrForbidden` | 403 Forbidden | WARN | No |
| `ErrNotFound` | 404 Not Found | WARN | No |
| `ErrMethodNotAllowed` | 405 Method Not Allowed | WARN | No |
| `ErrInternal` | 500 Internal Server Error | ERROR | Yes |
| Any unknown error | 500 Internal Server Error | ERROR | Yes |

### Example Log Output

**Validation error (WARN — no stack trace):**
```
time=2026-07-03T19:23:01.669+03:00 level=WARN source=handlers/utils.go:69 msg="request error" error="email already exists" method=POST path=/api/v1/auth/register remote=[::1]:65278 status=400
```

**Internal error (ERROR — includes stack trace):**
```
time=... level=ERROR source=handlers/utils.go:69 msg="request error" error="internal error" method=POST path=/api/v1/auth/register status=500 stack="goroutine 39 [running]:\n..."
```

---

## Request Logging Middleware

The `RequestLogger` middleware in `cmd/middleware.go` automatically logs every incoming HTTP request with:

- Method
- Path
- Remote address
- Response status code
- Request duration
- User agent

You do **not** need to log request information manually.

```
time=... level=INFO source=cmd/middleware.go:51 msg="incoming request" method=POST path=/api/v1/auth/register status=400 duration=1.03ms user_agent=curl/8.7.1
```

---

## Quick Reference — Common Patterns

### Success case
```go
// Handler only
re.App.Logger.Info("user registered successfully",
    "user_id", userID,
    "email", email,
)
```

### Expected/validation error
```go
// Service: just return the error
return 0, realtimeforum.ErrEmailExists

// Handler: single call to HandleError
re.HandleError(w, r, err)
// Produces: WARN "request error" error="email already exists" status=400
```

### Internal/system error
```go
// Where it happens (service or repository):
slog.Error("failed to hash password", "error", err)
return 0, realtimeforum.ErrInternal

// Handler delegates to HandleError
re.HandleError(w, r, err)
// Produces: ERROR "request error" error="internal error" status=500 + stack trace
```

---

## Adding a New Error

1. Define a sentinel error in `errors.go`:

```go
var ErrInvalidUsername = errors.New("invalid username")
```

2. (Optional) If it's a validation error that should map to 400 Bad Request at WARN level, add it to the `default` switch in `HandleError` (`pkg/handlers/utils.go`):

```go
case err == realtimeforum.ErrInvalidUsername,
    err == realtimeforum.ErrInvalidEmail,
    // ...
```

3. Use it in your code:

```go
// Service: return the error
return 0, realtimeforum.ErrInvalidUsername

// Handler: delegate to HandleError
re.HandleError(w, r, err)