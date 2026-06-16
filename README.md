#User Management REST API

A production-quality Go REST API for managing users, built with clean architecture principles. Designed as an internship assignment demonstrating strong engineering practices, idiomatic Go, and interview-ready code quality.

---

##  Tech Stack

| Technology | Purpose |
|---|---|
| **Go 1.23+** | Core language |
| **Fiber v2** | High-performance HTTP framework |
| **PostgreSQL 16** | Primary database |
| **SQLC** | Type-safe SQL code generation |
| **Uber Zap** | Structured logging |
| **go-playground/validator** | Input validation |
| **Docker & Docker Compose** | Containerized deployment |

---

##  Architecture

The project follows **clean architecture** with strict layer separation:

```
HTTP Request
    │
    ▼
┌─────────────────────┐
│    Middleware        │  ← Request ID + Duration Logging
├─────────────────────┤
│    Handler Layer     │  ← HTTP concerns only (parse, respond)
├─────────────────────┤
│    Validator         │  ← Input validation (go-playground/validator)
├─────────────────────┤
│    Service Layer     │  ← Business logic + age calculation
├─────────────────────┤
│    Repository Layer  │  ← Database access (interface)
├─────────────────────┤
│    SQLC Queries      │  ← Type-safe generated SQL
├─────────────────────┤
│    PostgreSQL        │  ← Data persistence
└─────────────────────┘
```

### Why This Architecture?

- **Handler** never contains business logic — only HTTP concerns (parsing body, returning status codes).
- **Service** owns all business rules — calculating age, validating existence, enriching responses.
- **Repository** wraps SQLC behind an interface — enables dependency injection and testability.
- **Middleware** is composable — request ID and logging are orthogonal concerns.

---

##  Folder Structure

```
ainyx/
├── cmd/
│   └── server/
│       └── main.go              # Entry point, DI wiring, graceful shutdown
├── internal/
│   ├── handler/
│   │   └── user.go              # HTTP route handlers
│   ├── service/
│   │   └── user.go              # Business logic layer
│   ├── repository/
│   │   └── user.go              # DB access via interface
│   ├── middleware/
│   │   ├── request_id.go        # X-Request-ID middleware
│   │   └── logger.go            # Request duration logging
│   ├── validator/
│   │   └── validator.go         # Input validation
│   ├── utils/
│   │   └── age.go               # CalculateAge utility
│   └── models/
│       ├── user.go              # Request/Response DTOs
│       └── response.go          # Error & pagination wrappers
├── db/
│   ├── migrations/
│   │   ├── 000001_create_users.up.sql
│   │   └── 000001_create_users.down.sql
│   ├── query/
│   │   └── user.sql             # SQLC query definitions
│   └── sqlc/                    # Generated type-safe Go code
│       ├── db.go
│       ├── models.go
│       ├── querier.go
│       └── user.sql.go
├── configs/
│   └── config.go                # Environment-based configuration
├── tests/
│   └── age_test.go              # Unit tests for CalculateAge
├── .env.example
├── .gitignore
├── Dockerfile
├── docker-compose.yml
├── sqlc.yaml
├── go.mod
└── go.sum
```

---

##  Environment Variables

| Variable | Description | Default |
|---|---|---|
| `PORT` | Server port | `3000` |
| `DATABASE_URL` | PostgreSQL connection string | `postgres://postgres:postgres@localhost:5432/ainyx?sslmode=disable` |
| `APP_ENV` | Environment (`development` / `production`) | `development` |

Copy `.env.example` to `.env` and adjust values:

```bash
cp .env.example .env
```

---

##  Running with Docker (Recommended)

The fastest way to get everything running — no manual setup required.

```bash
docker compose up --build
```

This starts:
- **PostgreSQL 16** on port `5432`
- **Go application** on port `3000`

Migrations run automatically on startup.

---

##  Running Locally

### Prerequisites
- Go 1.23+
- PostgreSQL 16+

### Steps

1. **Clone the repository**
   ```bash
   git clone https://github.com/ganesh/ainyx.git
   cd ainyx
   ```

2. **Set up environment**
   ```bash
   cp .env.example .env
   # Edit .env with your PostgreSQL connection string
   ```

3. **Start PostgreSQL** (if not using Docker)
   ```bash
   # Create the database
   createdb ainyx
   ```

4. **Install dependencies**
   ```bash
   go mod download
   ```

5. **Run the application**
   ```bash
   go run ./cmd/server
   ```

The server starts on `http://localhost:3000`.

---

##  API Endpoints

### Health Check

```bash
GET /health
```

### Create User

```bash
curl -X POST http://localhost:3000/users \
  -H "Content-Type: application/json" \
  -d '{"name": "Ganesh", "dob": "2004-01-15"}'
```

**Response (201 Created):**
```json
{
  "id": 1,
  "name": "Ganesh",
  "dob": "2004-01-15",
  "age": 22
}
```

### Get User

```bash
curl http://localhost:3000/users/1
```

**Response (200 OK):**
```json
{
  "id": 1,
  "name": "Ganesh",
  "dob": "2004-01-15",
  "age": 22
}
```

### Update User

```bash
curl -X PUT http://localhost:3000/users/1 \
  -H "Content-Type: application/json" \
  -d '{"name": "Ganesh Updated", "dob": "2004-01-15"}'
```

**Response (200 OK):**
```json
{
  "id": 1,
  "name": "Ganesh Updated",
  "dob": "2004-01-15",
  "age": 22
}
```

### Delete User

```bash
curl -X DELETE http://localhost:3000/users/1
```

**Response: 204 No Content**

### List Users (with Pagination)

```bash
curl "http://localhost:3000/users?page=1&limit=10"
```

**Response (200 OK):**
```json
{
  "data": [
    {
      "id": 1,
      "name": "Ganesh",
      "dob": "2004-01-15",
      "age": 22
    }
  ],
  "page": 1,
  "limit": 10,
  "total": 1
}
```

---

##  SQLC Usage

SQLC generates type-safe Go code from raw SQL queries. No ORM, no runtime reflection.

### Configuration (`sqlc.yaml`)

```yaml
version: "2"
sql:
  - engine: "postgresql"
    queries: "db/query/"
    schema: "db/migrations/"
    gen:
      go:
        package: "sqlc"
        out: "db/sqlc"
        sql_package: "pgx/v5"
```

### Regenerating Code

```bash
# Install sqlc
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# Generate
sqlc generate
```

### Query Definitions (`db/query/user.sql`)

Six queries are defined:
- `CreateUser` — INSERT with RETURNING
- `GetUser` — SELECT by ID
- `ListUsers` — SELECT with LIMIT/OFFSET
- `UpdateUser` — UPDATE with RETURNING
- `DeleteUser` — DELETE by ID
- `CountUsers` — COUNT for pagination

---

##  Testing

Run unit tests:

```bash
go test ./tests/ -v
```

### Test Coverage

```bash
go test ./tests/ -v -cover
```

### What's Tested

| Test Case | Description |
|---|---|
| Birthday already passed | Age = year difference |
| Birthday is today | Age = year difference (inclusive) |
| Birthday not yet reached | Age = year difference - 1 |
| Born Jan 1st | Edge case: start of year |
| Born Dec 31st | Edge case: end of year |
| Leap year (Feb 29) | Handles non-existent day in non-leap years |
| Newborn (born today) | Returns 0 |
| One year old | Exact anniversary |

---

##  Design Decisions

### 1. Age Is Never Stored

Age changes daily — storing it would require a cron job and risk stale data. Instead, `CalculateAge(dob)` computes it on every response. This is a O(1) operation with zero maintenance burden.

### 2. Repository Interface Pattern

The `UserRepository` interface decouples the service layer from the database implementation. This enables:
- Swapping PostgreSQL for another store without touching business logic
- Easy mocking in unit tests
- Clear dependency direction (service → interface ← implementation)

### 3. SQLC Over GORM

SQLC generates code at build time from raw SQL. Benefits:
- No runtime reflection or magic
- SQL is visible, reviewable, and optimizable
- Type-safe parameters and return values
- Compile-time errors for schema mismatches

### 4. Fiber v2

Chosen for performance (built on fasthttp) and Express-like ergonomics. The handler layer stays thin because all logic lives in the service.

### 5. Structured Logging with Zap

Every log entry includes structured fields (`request_id`, `method`, `path`, `duration`, `status_code`). This makes logs searchable and parseable by tools like ELK, Datadog, or Grafana Loki.

### 6. Middleware Composition

Request ID and duration logging are orthogonal concerns that compose cleanly. The Request ID middleware runs first, stores the ID in Fiber locals, and the logger middleware reads it downstream.

### 7. Validation at the Boundary

Validation happens once at the handler layer using `go-playground/validator`. The service layer trusts that inputs are already validated. This keeps the service focused on business logic.

### 8. Graceful Shutdown

The server listens for SIGINT/SIGTERM and calls `app.Shutdown()` before exiting. This ensures in-flight requests complete and database connections are properly closed.

---

##  Package Guide

| Package | Purpose |
|---|---|
| `cmd/server` | Application entry point — wires DI, starts server |
| `internal/handler` | HTTP route handlers — parse requests, return responses |
| `internal/service` | Business logic — age calculation, CRUD orchestration |
| `internal/repository` | Database access — wraps SQLC behind an interface |
| `internal/middleware` | Cross-cutting concerns — request ID, logging |
| `internal/validator` | Input validation — go-playground/validator rules |
| `internal/utils` | Pure utility functions — `CalculateAge` |
| `internal/models` | DTOs — request/response structs |
| `db/sqlc` | Generated type-safe database code |
| `db/migrations` | SQL migration files |
| `db/query` | SQLC query definitions |
| `configs` | Environment-based configuration |

---

##  Interview Preparation

### Q: Why not store age in the database?

Age changes daily. Storing it creates a stale data problem that requires a daily cron job to fix. Computing it on the fly from `dob` is O(1), always correct, and eliminates an entire class of bugs.

### Q: Why use an interface for the repository?

It inverts the dependency — the service depends on an abstraction, not a concrete implementation. This makes the service layer unit-testable without a database and allows swapping storage backends.

### Q: How does the Request ID middleware work?

It checks for an existing `X-Request-ID` header (for distributed tracing). If absent, it generates a UUID. The ID is stored in Fiber locals and set on the response header, making it available for downstream logging.

### Q: Why SQLC instead of GORM?

GORM uses runtime reflection and hides SQL behind method chains. SQLC generates code from actual SQL at build time. This gives compile-time safety, visible SQL, and zero runtime overhead.

### Q: How do you handle errors without exposing internals?

The service layer wraps database errors with user-friendly messages (`"failed to create user"`) and logs the real error with Zap. The handler maps sentinel errors like `ErrUserNotFound` to proper HTTP status codes.

### Q: Explain the pagination approach.

Two queries run: `ListUsers` (with LIMIT/OFFSET) and `CountUsers` (total rows). The response includes `page`, `limit`, and `total` so clients can compute total pages. OFFSET-based pagination is simple and sufficient at this scale.

### Q: Why Fiber over standard library?

Fiber provides routing, middleware composition, and JSON handling with significantly less boilerplate than `net/http`. It's built on fasthttp for performance. The handler layer stays thin regardless of framework choice.

---

##  License

MIT
