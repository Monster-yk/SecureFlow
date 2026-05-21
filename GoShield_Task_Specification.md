# GoShield — Complete Function Task List

> **207 tasks** across **86 files** spanning **8 weeks** of learning.

## How to use this file

Each task = one function, struct, constant, or component in one file. Read the task description first. Understand **WHY** the function exists. Then go learn the concepts listed. Then implement it yourself. Tasks are ordered so dependencies always come before the task that needs them. Do one task per session. Commit after every file is complete.

**Legend:**
- **Signature** — The exact function/type/component signature you'll write
- **Purpose** — What it does in one sentence
- **Why it exists** — What breaks without it (2–3 sentences)
- **Inputs / Returns** — Parameter and return value explanations
- **Concepts to learn** — What to study before coding (with search terms)
- **Depends on** — Prerequisite TASK numbers

---
---

## ⚙️ WEEK 1 — Authentication & Server Foundation

### FILE: `backend/internal/config/config.go`
**Purpose of this file:** Centralizes all environment-variable-based configuration into a single typed struct so the rest of the application never reads `os.Getenv` directly.
**Week:** Week 1
**Layer:** Backend

---

#### TASK-001: Load
**Signature:** `func Load() (*Config, error)`
**File Location:** `backend/internal/config/config.go`
**Purpose:** Reads environment variables and returns a populated Config struct.
**Why it exists:** Every service needs configuration (DB URL, JWT secret, ports). Centralizing this into one function means if an env var is missing, the app fails fast at startup with a clear error instead of crashing randomly at runtime. It also makes testing easier — you can set env vars once and get a config object.
**Inputs:** None (reads from `os.Getenv` internally). Expected env vars: `DATABASE_URL`, `JWT_SECRET`, `PORT`, `REDIS_URL`, `KAFKA_BROKERS`, `CLAUDE_API_KEY`.
**Returns:** `*Config` — a pointer to a Config struct with all fields populated. `error` — returned if any required env var is missing.
**Concepts to learn before implementing:**
- Go structs and field tags → search "Go struct tutorial"
- `os.Getenv` and environment variables in Go → search "golang os.Getenv"
- Error handling patterns in Go (returning error as second value) → search "golang error handling best practices"
- The 12-Factor App config methodology → search "12 factor app config"
**Depends on:** None (this is the first task)

---

### FILE: `backend/internal/model/user.go`
**Purpose of this file:** Defines the User domain model and role constants that every layer of the application (handlers, repositories, auth) shares.
**Week:** Week 1
**Layer:** Backend

---

#### TASK-002: Role Constants
**Signature:** `const RoleAnalyst = "analyst"` and `const RoleAdmin = "admin"`
**File Location:** `backend/internal/model/user.go`
**Purpose:** Define the two authorization roles used throughout GoShield's RBAC system.
**Why it exists:** Hardcoding role strings like `"analyst"` in multiple files leads to typos and inconsistency. Defining them as constants means the compiler catches misspellings and every file references the same value. These constants are used in JWT claims, middleware authorization checks, and database seeding.
**Inputs:** N/A (constants have no inputs)
**Returns:** N/A (constants are values, not functions)
**Concepts to learn before implementing:**
- Go constants and `const` blocks → search "golang constants iota"
- Role-Based Access Control (RBAC) concept → search "RBAC explained"
- Why string constants beat magic strings → search "magic strings anti-pattern"
**Depends on:** None

---

#### TASK-003: User Struct
**Signature:** `type User struct { ... }`
**File Location:** `backend/internal/model/user.go`
**Purpose:** Represents a user record as stored in the PostgreSQL `users` table.
**Why it exists:** Go is statically typed — you need a struct to scan database rows into and to serialize as JSON in API responses. This struct is the single source of truth for what a "user" looks like across the entire backend.
**Fields to define:**
- `ID` (uuid.UUID) — primary key, universally unique identifier for each user
- `Email` (string) — unique login identifier, used in FindByEmail queries
- `PasswordHash` (string) — bcrypt-hashed password, NEVER the plaintext password. Tagged `json:"-"` so it's never sent in API responses
- `Role` (string) — either RoleAnalyst or RoleAdmin, controls what endpoints the user can access
- `CreatedAt` (time.Time) — timestamp of account creation, auto-set by PostgreSQL
- `UpdatedAt` (time.Time) — timestamp of last modification
**Inputs:** N/A (struct definition)
**Returns:** N/A (struct definition)
**Concepts to learn before implementing:**
- Go struct field tags (`json:"..."`, `db:"..."`) → search "golang struct tags json"
- UUID type in Go (github.com/google/uuid) → search "golang uuid package"
- Why password hashes, not passwords → search "bcrypt password hashing explained"
- time.Time in Go → search "golang time package"
**Depends on:** TASK-002

---

### FILE: `backend/internal/auth/jwt.go`
**Purpose of this file:** Implements JWT token generation and validation. This is the core authentication mechanism — it creates tokens on login and verifies them on every protected request.
**Week:** Week 1
**Layer:** Backend

---

#### TASK-004: NewService
**Signature:** `func NewService(secret string, expiry time.Duration) *Service`
**File Location:** `backend/internal/auth/jwt.go`
**Purpose:** Constructor that creates a new JWT Service with the signing secret and token expiry duration.
**Why it exists:** Go doesn't have traditional constructors, so the convention is a `New*` function. This ensures every JWT Service is created with a valid secret and expiry. Without it, you'd have uninitialized Service structs with empty secrets, which would generate invalid tokens.
**Inputs:**
- `secret` (string) — the HMAC-SHA256 signing key from config. Must be long and random in production.
- `expiry` (time.Duration) — how long tokens remain valid (e.g., `24 * time.Hour`).
**Returns:** `*Service` — a pointer to a fully initialized Service struct.
**Concepts to learn before implementing:**
- Go constructor pattern (New* functions) → search "golang constructor pattern"
- Pointer receivers vs value receivers → search "golang pointer vs value receiver"
- time.Duration in Go → search "golang time.Duration"
**Depends on:** TASK-001

---

#### TASK-005: Generate
**Signature:** `func (s *Service) Generate(userID uuid.UUID, role string) (string, error)`
**File Location:** `backend/internal/auth/jwt.go`
**Purpose:** Creates a signed JWT token containing the user's ID and role as claims.
**Why it exists:** After a user logs in with correct credentials, you need to give them a token they can send with future requests to prove their identity. This function creates that token. The token embeds the user's ID (so you know WHO they are) and role (so you know WHAT they can do) without hitting the database on every request.
**Inputs:**
- `userID` (uuid.UUID) — the authenticated user's database ID, stored as the `sub` (subject) claim.
- `role` (string) — the user's role ("analyst" or "admin"), stored as a custom `role` claim.
**Returns:**
- `string` — the signed JWT token string (header.payload.signature format)
- `error` — returned if signing fails (rare, but possible with invalid keys)
**Concepts to learn before implementing:**
- What is a JWT and how it works (header, payload, signature) → search "JWT explained simply"
- HMAC-SHA256 signing → search "HMAC SHA256 JWT"
- Go JWT library (github.com/golang-jwt/jwt/v5) → search "golang-jwt v5 tutorial"
- JWT claims (sub, exp, iat, custom claims) → search "JWT claims explained"
**Depends on:** TASK-004, TASK-002

---

#### TASK-006: Validate
**Signature:** `func (s *Service) Validate(tokenStr string) (*Claims, error)`
**File Location:** `backend/internal/auth/jwt.go`
**Purpose:** Parses a JWT token string, verifies its signature and expiration, and extracts the claims.
**Why it exists:** Every protected API request includes a JWT in the Authorization header. This function is called by the middleware to verify the token is genuine (not tampered with), not expired, and to extract the user's ID and role so the handler knows who's making the request.
**Inputs:**
- `tokenStr` (string) — the raw JWT string from the Authorization header (after stripping "Bearer ").
**Returns:**
- `*Claims` — a pointer to a Claims struct containing UserID and Role extracted from the token.
- `error` — returned if the token is expired, has an invalid signature, or is malformed.
**Concepts to learn before implementing:**
- JWT validation and signature verification → search "JWT token validation golang"
- Token expiration checking → search "JWT exp claim validation"
- Custom claims struct in Go → search "golang jwt custom claims"
- Error handling for invalid tokens → search "golang-jwt parse token errors"
**Depends on:** TASK-004

---

### FILE: `backend/internal/auth/middleware.go`
**Purpose of this file:** Provides HTTP middleware that intercepts every request to protected routes, validates the JWT token, and injects user information into the request context.
**Week:** Week 1
**Layer:** Backend

---

#### TASK-007: JWTMiddleware
**Signature:** `func JWTMiddleware(authService *Service) func(http.Handler) http.Handler`
**File Location:** `backend/internal/auth/middleware.go`
**Purpose:** Returns a middleware function that validates JWT tokens from the Authorization header and stores user claims in the request context.
**Why it exists:** Without this middleware, every single handler would need to copy-paste the same token extraction and validation code. The middleware pattern lets you write this logic once and apply it to all protected routes. It follows the "chain of responsibility" pattern — if the token is invalid, the request is rejected before it ever reaches the handler.
**Inputs:**
- `authService` (*Service) — the JWT service used to validate tokens (injected via closure).
- Inner function receives `next` (http.Handler) — the next handler in the chain to call if auth succeeds.
- Inner handler receives `w` (http.ResponseWriter) and `r` (*http.Request) — the standard HTTP pair.
**Returns:**
- Outer: `func(http.Handler) http.Handler` — a middleware-compatible function signature that Chi router expects.
- Inner: Writes 401 Unauthorized JSON response if token is missing/invalid, or calls `next.ServeHTTP(w, r.WithContext(ctx))` with claims in context if valid.
**Concepts to learn before implementing:**
- HTTP middleware pattern in Go → search "golang http middleware pattern"
- Chi router middleware → search "go-chi middleware tutorial"
- Go context.WithValue for request-scoped data → search "golang context.WithValue"
- Authorization header format ("Bearer \<token\>") → search "HTTP Authorization Bearer token"
- Closure pattern for dependency injection → search "golang closure dependency injection"
**Depends on:** TASK-006

---

### FILE: `backend/internal/repository/user_repo.go`
**Purpose of this file:** Data access layer for the users table. All SQL queries related to users live here, keeping database logic separated from business logic.
**Week:** Week 1
**Layer:** Backend

---

#### TASK-008: NewUserRepo
**Signature:** `func NewUserRepo(db *sql.DB) *UserRepo`
**File Location:** `backend/internal/repository/user_repo.go`
**Purpose:** Constructor that creates a UserRepo with a database connection.
**Why it exists:** The repository needs a database connection to execute queries. This constructor takes the connection as a parameter (dependency injection) rather than creating it internally, which makes the repo testable — in tests you can pass a mock or test database.
**Inputs:**
- `db` (*sql.DB) — an open PostgreSQL connection pool from the database/sql package.
**Returns:** `*UserRepo` — a pointer to the initialized repository.
**Concepts to learn before implementing:**
- Repository pattern in Go → search "golang repository pattern"
- database/sql package basics → search "golang database sql tutorial"
- Dependency injection in Go → search "golang dependency injection"
- `*sql.DB` is a connection pool, not a single connection → search "golang sql.DB connection pool"
**Depends on:** TASK-001

---

#### TASK-009: Create
**Signature:** `func (r *UserRepo) Create(ctx context.Context, user *model.User) error`
**File Location:** `backend/internal/repository/user_repo.go`
**Purpose:** Inserts a new user record into the users table.
**Why it exists:** When someone registers, their email, hashed password, and role need to be persisted. This function executes an INSERT query with RETURNING to populate the user's auto-generated ID and timestamps. Without it, registered users would vanish on server restart.
**Inputs:**
- `ctx` (context.Context) — request-scoped context for cancellation and timeouts.
- `user` (*model.User) — a User struct with Email, PasswordHash, and Role populated. ID and timestamps are set by PostgreSQL.
**Returns:** `error` — returned if the INSERT fails (e.g., duplicate email violates unique constraint).
**Concepts to learn before implementing:**
- Go context.Context and why every DB call needs it → search "golang context database"
- SQL INSERT with RETURNING clause (PostgreSQL) → search "postgresql INSERT RETURNING"
- `db.QueryRowContext` for single-row results → search "golang QueryRowContext"
- Scanning results into struct fields with `row.Scan` → search "golang sql row scan"
**Depends on:** TASK-008, TASK-003

---

#### TASK-010: FindByEmail
**Signature:** `func (r *UserRepo) FindByEmail(ctx context.Context, email string) (*model.User, error)`
**File Location:** `backend/internal/repository/user_repo.go`
**Purpose:** Retrieves a user from the database by their email address.
**Why it exists:** During login, the user provides an email. You need to look up their record to compare the provided password against the stored hash. This is the first step of authentication — find the user, then verify their credentials.
**Inputs:**
- `ctx` (context.Context) — request-scoped context.
- `email` (string) — the email address to search for.
**Returns:**
- `*model.User` — the found user (with PasswordHash for credential verification), or nil if not found.
- `error` — returned on database errors. If no user found, return a specific sentinel error or nil user.
**Concepts to learn before implementing:**
- SQL SELECT WHERE clause → search "postgresql SELECT WHERE"
- `sql.ErrNoRows` sentinel error → search "golang sql.ErrNoRows"
- Returning pointer vs error for "not found" scenarios → search "golang not found pattern"
**Depends on:** TASK-008, TASK-003

---

#### TASK-011: FindByID
**Signature:** `func (r *UserRepo) FindByID(ctx context.Context, id uuid.UUID) (*model.User, error)`
**File Location:** `backend/internal/repository/user_repo.go`
**Purpose:** Retrieves a user from the database by their UUID primary key.
**Why it exists:** The `/me` endpoint needs to fetch the current user's full profile using the user ID extracted from their JWT token. The JWT only contains the ID — this function hydrates it into a full User object with email, role, and timestamps.
**Inputs:**
- `ctx` (context.Context) — request-scoped context.
- `id` (uuid.UUID) — the user's unique identifier from the JWT claims.
**Returns:**
- `*model.User` — the found user (with PasswordHash excluded via `json:"-"` tag in responses).
- `error` — returned on database errors or if no user exists with that ID.
**Concepts to learn before implementing:**
- UUID as primary key in PostgreSQL → search "postgresql uuid primary key"
- Parameterized queries to prevent SQL injection → search "golang sql parameterized query"
**Depends on:** TASK-008, TASK-003

---

### FILE: `backend/internal/handler/health.go`
**Purpose of this file:** Provides a simple health check endpoint that load balancers and Kubernetes probes use to determine if the server is alive.
**Week:** Week 1
**Layer:** Backend

---

#### TASK-012: Health
**Signature:** `func Health(w http.ResponseWriter, r *http.Request)`
**File Location:** `backend/internal/handler/health.go`
**Purpose:** Returns a 200 OK JSON response indicating the server is healthy.
**Why it exists:** In production, Kubernetes liveness probes and load balancers periodically hit a health endpoint to check if the service is running. If this endpoint stops responding, the orchestrator restarts the container. It's also the simplest possible handler — a good first handler to write to learn the pattern.
**Inputs:**
- `w` (http.ResponseWriter) — the response writer to send the JSON response.
- `r` (*http.Request) — the incoming request (unused but required by the HandlerFunc signature).
**Returns:** Writes `{"status": "ok"}` as JSON with Content-Type header and 200 status code.
**Concepts to learn before implementing:**
- Go `http.HandlerFunc` signature → search "golang http.HandlerFunc"
- `json.NewEncoder(w).Encode()` for JSON responses → search "golang json encode http response"
- Setting response headers (`Content-Type`) → search "golang http response headers"
**Depends on:** None

---

### FILE: `backend/internal/handler/auth_handler.go`
**Purpose of this file:** HTTP handlers for authentication endpoints — registration, login, and fetching the current user's profile. Each function maps to one REST API route.
**Week:** Week 1
**Layer:** Backend

---

#### TASK-013: NewAuthHandler
**Signature:** `func NewAuthHandler(userRepo *repository.UserRepo, authService *auth.Service) *AuthHandler`
**File Location:** `backend/internal/handler/auth_handler.go`
**Purpose:** Constructor that creates an AuthHandler with its dependencies (user repository and JWT service).
**Why it exists:** The auth handler needs two dependencies: a way to read/write users (UserRepo) and a way to create/validate tokens (auth.Service). This constructor wires those dependencies in, following the dependency injection pattern. It's called once during server startup.
**Inputs:**
- `userRepo` (*repository.UserRepo) — for database operations on users.
- `authService` (*auth.Service) — for generating JWT tokens after successful login.
**Returns:** `*AuthHandler` — fully initialized handler ready to be registered on routes.
**Concepts to learn before implementing:**
- Handler struct pattern in Go (grouping related handlers) → search "golang handler struct pattern"
- Dependency injection without frameworks → search "golang manual dependency injection"
**Depends on:** TASK-008, TASK-004

---

#### TASK-014: Register
**Signature:** `func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request)`
**File Location:** `backend/internal/handler/auth_handler.go`
**Purpose:** Handles POST /api/register — creates a new user account.
**Why it exists:** New users need a way to create accounts. This handler receives email, password, and role from the request body, hashes the password with bcrypt, stores the user in the database, and returns a success response. Without password hashing, a database breach would expose all passwords.
**Inputs:**
- `w` (http.ResponseWriter) — response writer.
- `r` (*http.Request) — contains JSON body with `email`, `password`, and `role` fields.
**Returns:** Writes 201 Created with the new user's ID and email (no password), or 400/409 on validation failure / duplicate email.
**Concepts to learn before implementing:**
- Parsing JSON request bodies with `json.NewDecoder(r.Body).Decode()` → search "golang parse json request body"
- bcrypt password hashing → search "golang bcrypt hash password"
- HTTP status codes (201 Created, 400 Bad Request, 409 Conflict) → search "HTTP status codes explained"
- Input validation patterns → search "golang request validation"
**Depends on:** TASK-013, TASK-009

---

#### TASK-015: Login
**Signature:** `func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request)`
**File Location:** `backend/internal/handler/auth_handler.go`
**Purpose:** Handles POST /api/login — authenticates a user and returns a JWT token.
**Why it exists:** This is the entry point of the authentication flow. The user provides email and password. The handler looks up the user by email, compares the password against the stored bcrypt hash, and if valid, generates a JWT token. This token is then used for all subsequent authenticated requests.
**Inputs:**
- `w` (http.ResponseWriter) — response writer.
- `r` (*http.Request) — contains JSON body with `email` and `password` fields.
**Returns:** Writes 200 OK with `{"token": "...", "user": {...}}` on success, or 401 Unauthorized if credentials are wrong.
**Concepts to learn before implementing:**
- bcrypt.CompareHashAndPassword → search "golang bcrypt compare"
- Authentication flow (lookup → compare → issue token) → search "JWT login flow"
- 401 Unauthorized vs 403 Forbidden → search "401 vs 403 HTTP"
**Depends on:** TASK-013, TASK-010, TASK-005

---

#### TASK-016: Me
**Signature:** `func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request)`
**File Location:** `backend/internal/handler/auth_handler.go`
**Purpose:** Handles GET /api/me — returns the currently authenticated user's profile.
**Why it exists:** The frontend needs to fetch the logged-in user's details (email, role) to display in the UI and make authorization decisions. This handler extracts the user ID from the JWT claims stored in the request context (put there by the JWT middleware), fetches the full user from the database, and returns it.
**Inputs:**
- `w` (http.ResponseWriter) — response writer.
- `r` (*http.Request) — the JWT middleware has already placed claims in `r.Context()`.
**Returns:** Writes 200 OK with the user's profile as JSON (ID, email, role, timestamps — no password hash).
**Concepts to learn before implementing:**
- Extracting values from context → search "golang context.Value"
- Type assertion for context values → search "golang type assertion"
- Protected endpoint pattern (middleware → handler) → search "golang protected routes middleware"
**Depends on:** TASK-013, TASK-007, TASK-011

---

### FILE: `backend/cmd/server/main.go`
**Purpose of this file:** The application entry point. Wires all dependencies together (config, database, repositories, handlers, routes) and starts the HTTP server. This is the "glue" that connects every other component.
**Week:** Week 1
**Layer:** Backend

---

#### TASK-017: main — Load Config
**Signature:** Part of `func main()`
**File Location:** `backend/cmd/server/main.go`
**Purpose:** Calls `config.Load()` to read all environment variables into a typed Config struct.
**Why it exists:** The very first thing the server does is load configuration. If any required env var is missing, the server should crash immediately with a helpful error rather than starting and failing later. This is the "fail fast" principle.
**Inputs:** N/A (reads environment variables)
**Returns:** A `*config.Config` used by all subsequent setup steps.
**Concepts to learn before implementing:**
- Go main package and main function → search "golang main package entry point"
- log.Fatal for startup errors → search "golang log.Fatal"
- Fail-fast startup pattern → search "fail fast application startup"
**Depends on:** TASK-001

---

#### TASK-018: main — Connect to PostgreSQL
**Signature:** Part of `func main()`
**File Location:** `backend/cmd/server/main.go`
**Purpose:** Opens a connection pool to PostgreSQL using the DATABASE_URL from config.
**Why it exists:** The database is the backbone of the application — user data, logs, and alerts all live there. This step establishes the connection pool and verifies connectivity with `db.Ping()`. If the database is unreachable, the server fails at startup.
**Inputs:** `cfg.DatabaseURL` — PostgreSQL connection string.
**Returns:** `*sql.DB` — an open connection pool.
**Concepts to learn before implementing:**
- `sql.Open` and connection strings → search "golang sql.Open postgres"
- `db.Ping()` to verify connectivity → search "golang db.Ping"
- PostgreSQL connection strings → search "postgresql connection string format"
- `defer db.Close()` for cleanup → search "golang defer"
- lib/pq driver import with blank identifier → search "golang blank import database driver"
**Depends on:** TASK-017

---

#### TASK-019: main — Wire Repositories
**Signature:** Part of `func main()`
**File Location:** `backend/cmd/server/main.go`
**Purpose:** Creates repository instances by passing the database connection to their constructors.
**Why it exists:** Repositories are the data access layer. They need the database connection pool to execute queries. This step creates a UserRepo (and later LogRepo, AlertRepo) that handlers will use. It's the dependency injection wiring step.
**Inputs:** `db` (*sql.DB) from the previous step.
**Returns:** `*repository.UserRepo` (and later other repos).
**Concepts to learn before implementing:**
- Manual dependency injection wiring → search "golang wire dependencies main"
- Repository pattern purpose → search "repository pattern explained"
**Depends on:** TASK-018, TASK-008

---

#### TASK-020: main — Wire Handlers
**Signature:** Part of `func main()`
**File Location:** `backend/cmd/server/main.go`
**Purpose:** Creates handler instances by passing repositories and services to their constructors.
**Why it exists:** Handlers need repositories to access data and services (like JWT) to perform authentication. This step creates an AuthHandler with the UserRepo and auth.Service. Each handler gets exactly the dependencies it needs — nothing more.
**Inputs:** Repository instances and the auth.Service.
**Returns:** `*handler.AuthHandler` (and later other handlers).
**Concepts to learn before implementing:**
- Layered architecture (handler → service → repository) → search "layered architecture golang"
**Depends on:** TASK-019, TASK-004, TASK-013

---

#### TASK-021: main — Register Routes
**Signature:** Part of `func main()`
**File Location:** `backend/cmd/server/main.go`
**Purpose:** Creates a Chi router, registers middleware, and maps URL paths to handler functions.
**Why it exists:** The router is the traffic director — it decides which handler function processes each incoming HTTP request based on the URL path and HTTP method. Without route registration, the server would accept connections but not know what to do with any of them.
**Inputs:** Chi router `r := chi.NewRouter()`, middleware functions, handler instances.
**Returns:** A fully configured `*chi.Mux` router.
**Concepts to learn before implementing:**
- Chi router basics (NewRouter, Route, Get, Post) → search "go-chi router tutorial"
- Route grouping with `r.Route("/api", func(r chi.Router) {...})` → search "go-chi route groups"
- Applying middleware to route groups → search "go-chi middleware groups"
- RESTful URL design → search "REST API URL design best practices"
**Depends on:** TASK-020, TASK-007

---

#### TASK-022: main — Start HTTP Server
**Signature:** Part of `func main()`
**File Location:** `backend/cmd/server/main.go`
**Purpose:** Starts the HTTP server on the configured port, listening for incoming requests.
**Why it exists:** This is the final step — actually starting the server. It uses `http.ListenAndServe` with the Chi router and the port from config. A log message indicates the server is ready. In production, you'd also handle graceful shutdown.
**Inputs:** `cfg.Port` and the configured router.
**Returns:** Blocks forever (or until error). If `ListenAndServe` returns, it means the server crashed.
**Concepts to learn before implementing:**
- `http.ListenAndServe` → search "golang http.ListenAndServe"
- Graceful shutdown with signal handling (optional, advanced) → search "golang graceful shutdown http server"
- Logging server startup information → search "golang log Printf"
**Depends on:** TASK-021

---

### FILE: `migrations/001_create_users.sql`
**Purpose of this file:** SQL migration that creates the database schema for user authentication — the custom enum type for roles, the users table, and performance indexes.
**Week:** Week 1
**Layer:** Backend

---

#### TASK-023: CREATE TYPE user_role
**Signature:** `CREATE TYPE user_role AS ENUM ('analyst', 'admin');`
**File Location:** `migrations/001_create_users.sql`
**Purpose:** Defines a PostgreSQL custom enum type for user roles.
**Why it exists:** Using an enum instead of a plain VARCHAR enforces valid values at the database level. If someone tries to insert `role = "hacker"`, PostgreSQL will reject it. This is defense in depth — the application validates roles AND the database enforces them.
**Inputs:** N/A (DDL statement)
**Returns:** N/A (creates a type in the database)
**Concepts to learn before implementing:**
- PostgreSQL custom enum types → search "postgresql CREATE TYPE ENUM"
- Database-level validation vs application-level → search "database constraints vs application validation"
- SQL migrations and why they're versioned → search "database migration versioning"
**Depends on:** None

---

#### TASK-024: CREATE TABLE users
**Signature:** `CREATE TABLE users (...);`
**File Location:** `migrations/001_create_users.sql`
**Purpose:** Creates the users table with columns for ID, email, password hash, role, and timestamps.
**Why it exists:** The users table is the foundation of authentication. Every registered user has a row here. The table enforces constraints: UUID primary key with default generation, unique email (no duplicate accounts), NOT NULL on required fields, and the user_role enum type for the role column.
**Fields to define:**
- `id` UUID PRIMARY KEY DEFAULT gen_random_uuid() — auto-generated unique identifier
- `email` VARCHAR(255) UNIQUE NOT NULL — login identifier, uniqueness prevents duplicate accounts
- `password_hash` VARCHAR(255) NOT NULL — bcrypt hash, never plaintext
- `role` user_role NOT NULL DEFAULT 'analyst' — uses the custom enum, new users default to analyst
- `created_at` TIMESTAMPTZ NOT NULL DEFAULT NOW() — auto-set on insert
- `updated_at` TIMESTAMPTZ NOT NULL DEFAULT NOW() — updated via trigger or application code
**Inputs:** N/A (DDL statement)
**Returns:** N/A (creates a table)
**Concepts to learn before implementing:**
- PostgreSQL CREATE TABLE syntax → search "postgresql CREATE TABLE"
- UUID generation with gen_random_uuid() → search "postgresql gen_random_uuid"
- UNIQUE constraints → search "postgresql UNIQUE constraint"
- TIMESTAMPTZ vs TIMESTAMP → search "postgresql timestamptz vs timestamp"
- DEFAULT values in PostgreSQL → search "postgresql DEFAULT column value"
**Depends on:** TASK-023

---

#### TASK-025: CREATE INDEX on users
**Signature:** `CREATE INDEX idx_users_email ON users(email);`
**File Location:** `migrations/001_create_users.sql`
**Purpose:** Creates a B-tree index on the email column for fast lookups during login.
**Why it exists:** Without an index, every login attempt would require PostgreSQL to scan every row in the users table to find the matching email (a "sequential scan"). With an index, it jumps directly to the matching row. This is the difference between O(n) and O(log n) lookup time as the user count grows.
**Inputs:** N/A (DDL statement)
**Returns:** N/A (creates an index)
**Concepts to learn before implementing:**
- Database indexes and B-trees → search "postgresql index explained"
- When to create indexes → search "database index best practices"
- `CREATE INDEX` syntax → search "postgresql CREATE INDEX"
- EXPLAIN ANALYZE to verify index usage → search "postgresql EXPLAIN ANALYZE"
**Depends on:** TASK-024

---
---

## 📊 WEEK 2 — Security Logs, Alerts & Detection

### FILE: `migrations/002_create_logs.sql`
**Purpose of this file:** SQL migration that creates the security_logs table — the core data store for all ingested security events — along with performance indexes for common query patterns.
**Week:** Week 2
**Layer:** Backend

---

#### TASK-026: CREATE TABLE security_logs
**Signature:** `CREATE TABLE security_logs (...);`
**File Location:** `migrations/002_create_logs.sql`
**Purpose:** Creates the table that stores every security event ingested into GoShield.
**Why it exists:** Security logs are the primary data that GoShield analyzes. Every firewall event, IDS alert, or authentication attempt is stored here. The table is designed for high-volume inserts (from Kafka) and flexible querying (by threat level, source, time range).
**Fields to define:**
- `id` UUID PRIMARY KEY DEFAULT gen_random_uuid() — unique log entry identifier
- `source` VARCHAR(50) NOT NULL — where the log came from (e.g., "firewall", "ids", "auth_service")
- `message` TEXT NOT NULL — the raw log message or description
- `threat_level` VARCHAR(20) NOT NULL DEFAULT 'low' — severity classification: low, medium, high, critical
- `ip_address` VARCHAR(45) — source IP (supports IPv6 length)
- `user_agent` VARCHAR(500) — browser/client user agent string
- `ai_summary` TEXT — Claude AI's analysis of this log (nullable, populated asynchronously)
- `created_at` TIMESTAMPTZ NOT NULL DEFAULT NOW() — when the log was ingested
**Inputs:** N/A (DDL statement)
**Returns:** N/A (creates a table)
**Concepts to learn before implementing:**
- TEXT vs VARCHAR in PostgreSQL → search "postgresql TEXT vs VARCHAR"
- Nullable columns and when to use them → search "postgresql nullable columns"
- Designing tables for high write throughput → search "postgresql high write performance"
**Depends on:** TASK-024 (must run after users table exists for migration ordering)

---

#### TASK-027: CREATE INDEX on security_logs
**Signature:** `CREATE INDEX idx_logs_level ON security_logs(threat_level);` and `CREATE INDEX idx_logs_created ON security_logs(created_at DESC);`
**File Location:** `migrations/002_create_logs.sql`
**Purpose:** Creates indexes for the two most common query patterns: filtering by threat level and sorting by time.
**Why it exists:** The dashboard constantly queries logs filtered by threat level ("show me all critical logs") and sorted by creation time ("most recent first"). Without these indexes, every dashboard load would trigger slow full-table scans. The DESC index on created_at optimizes the default "newest first" sort order.
**Inputs:** N/A (DDL statement)
**Returns:** N/A (creates indexes)
**Concepts to learn before implementing:**
- Composite vs single-column indexes → search "postgresql composite index"
- DESC index for sort optimization → search "postgresql descending index"
- Index trade-offs (faster reads, slower writes) → search "database index write penalty"
**Depends on:** TASK-026

---

### FILE: `migrations/003_create_alerts.sql`
**Purpose of this file:** SQL migration that creates the alerts table — stores threat alerts generated when the rule engine detects a suspicious log entry.
**Week:** Week 2
**Layer:** Backend

---

#### TASK-028: CREATE TABLE alerts
**Signature:** `CREATE TABLE alerts (...);`
**File Location:** `migrations/003_create_alerts.sql`
**Purpose:** Creates the table that stores alerts triggered by the threat detection system.
**Why it exists:** When a log is classified as high or critical threat, an alert is generated and stored here. Alerts are the actionable items that analysts review and resolve. Each alert references the security log that triggered it, creating a foreign key relationship.
**Fields to define:**
- `id` UUID PRIMARY KEY DEFAULT gen_random_uuid() — unique alert identifier
- `log_id` UUID NOT NULL REFERENCES security_logs(id) — the log that triggered this alert
- `rule_name` VARCHAR(255) NOT NULL — which detection rule fired (e.g., "SQL Injection Detected")
- `severity` VARCHAR(20) NOT NULL — threat severity (mirrors log's threat_level)
- `ai_summary` TEXT — Claude AI's analysis and recommendation (nullable, async)
- `resolved` BOOLEAN NOT NULL DEFAULT FALSE — whether an analyst has addressed this alert
- `resolved_at` TIMESTAMPTZ — when the alert was resolved (nullable)
- `created_at` TIMESTAMPTZ NOT NULL DEFAULT NOW() — when the alert was generated
**Inputs:** N/A (DDL statement)
**Returns:** N/A (creates a table)
**Concepts to learn before implementing:**
- Foreign keys in PostgreSQL → search "postgresql FOREIGN KEY REFERENCES"
- BOOLEAN columns with defaults → search "postgresql BOOLEAN column"
- Nullable TIMESTAMPTZ for optional timestamps → search "postgresql nullable timestamp"
**Depends on:** TASK-026

---

#### TASK-029: CREATE INDEX on alerts
**Signature:** `CREATE INDEX idx_alerts_resolved ON alerts(resolved);` and `CREATE INDEX idx_alerts_log_id ON alerts(log_id);`
**File Location:** `migrations/003_create_alerts.sql`
**Purpose:** Creates indexes for filtering unresolved alerts and looking up alerts by their source log.
**Why it exists:** The alerts page defaults to showing only unresolved alerts (WHERE resolved = FALSE). The log_id index supports the reverse lookup: "does this log already have an alert?" Both queries run frequently and benefit from indexes.
**Inputs:** N/A (DDL statement)
**Returns:** N/A (creates indexes)
**Concepts to learn before implementing:**
- Boolean column indexing → search "postgresql index boolean column"
- Foreign key indexes → search "postgresql index foreign key"
**Depends on:** TASK-028

---

### FILE: `backend/internal/model/log.go`
**Purpose of this file:** Defines the SecurityLog domain model, threat level constants, log source constants, and the ingest request DTO. Used by handlers, repositories, and the detection engine.
**Week:** Week 2
**Layer:** Backend

---

#### TASK-030: ThreatLevel Constants
**Signature:** `const ThreatLow = "low"`, `const ThreatMedium = "medium"`, `const ThreatHigh = "high"`, `const ThreatCritical = "critical"`
**File Location:** `backend/internal/model/log.go`
**Purpose:** Define the four severity levels for security log classification.
**Why it exists:** Threat levels are used everywhere — the rule engine classifies logs, the dashboard filters by level, the frontend colors badges by level. Having constants prevents typos ("critcal" vs "critical") and provides a single reference point. The levels follow standard security severity naming.
**Inputs:** N/A (constants)
**Returns:** N/A (constants)
**Concepts to learn before implementing:**
- Go const blocks → search "golang const block"
- Security severity levels (INFO, LOW, MEDIUM, HIGH, CRITICAL) → search "security severity levels CVSS"
**Depends on:** None

---

#### TASK-031: LogSource Constants
**Signature:** `const SourceFirewall = "firewall"`, `const SourceIDS = "ids"`, `const SourceAuthService = "auth_service"`, `const SourceEndpoint = "endpoint"`
**File Location:** `backend/internal/model/log.go`
**Purpose:** Define the valid sources from which security logs can originate.
**Why it exists:** In a real SOC, logs come from different systems — firewalls, intrusion detection systems, authentication services, and endpoint agents. These constants ensure consistent source labeling across ingestion, storage, and filtering. The frontend's SourcePie chart groups logs by these values.
**Inputs:** N/A (constants)
**Returns:** N/A (constants)
**Concepts to learn before implementing:**
- Security Operations Center (SOC) log sources → search "SOC log sources explained"
- Network security device types → search "firewall IDS IPS explained"
**Depends on:** None

---

#### TASK-032: SecurityLog Struct
**Signature:** `type SecurityLog struct { ... }`
**File Location:** `backend/internal/model/log.go`
**Purpose:** Represents a security log entry as stored in the security_logs table.
**Why it exists:** This struct mirrors the database table and serves as the data transfer object between the repository layer (database) and the handler layer (HTTP). JSON tags control how it serializes in API responses. DB tags control how it maps to database columns.
**Fields to define:**
- `ID` (uuid.UUID) — `json:"id" db:"id"`
- `Source` (string) — `json:"source" db:"source"` — one of the LogSource constants
- `Message` (string) — `json:"message" db:"message"` — raw log message
- `ThreatLevel` (string) — `json:"threat_level" db:"threat_level"` — one of the ThreatLevel constants
- `IPAddress` (string) — `json:"ip_address" db:"ip_address"`
- `UserAgent` (string) — `json:"user_agent" db:"user_agent"`
- `AISummary` (*string) — `json:"ai_summary" db:"ai_summary"` — pointer because nullable
- `CreatedAt` (time.Time) — `json:"created_at" db:"created_at"`
**Inputs:** N/A (struct definition)
**Returns:** N/A (struct definition)
**Concepts to learn before implementing:**
- Nullable fields as pointer types in Go → search "golang nullable field pointer"
- Multiple struct tags (json + db) → search "golang multiple struct tags"
- sql.NullString alternative → search "golang sql.NullString vs pointer"
**Depends on:** TASK-030, TASK-031

---

#### TASK-033: IngestRequest Struct
**Signature:** `type IngestRequest struct { ... }`
**File Location:** `backend/internal/model/log.go`
**Purpose:** Defines the expected JSON shape for incoming log ingestion API requests.
**Why it exists:** Separating the ingest request DTO from the SecurityLog model is important because they have different fields. The API client sends source, message, ip_address, and user_agent. The server adds the ID, threat_level (from the rule engine), ai_summary (from Claude), and created_at (from PostgreSQL). Using a separate struct means you validate only the fields the client should send.
**Fields to define:**
- `Source` (string) — `json:"source"` — required
- `Message` (string) — `json:"message"` — required
- `IPAddress` (string) — `json:"ip_address"` — optional
- `UserAgent` (string) — `json:"user_agent"` — optional
**Inputs:** N/A (struct definition)
**Returns:** N/A (struct definition)
**Concepts to learn before implementing:**
- DTO (Data Transfer Object) pattern → search "DTO pattern explained"
- Why separate request/response structs from database models → search "golang request struct vs model"
- JSON binding and struct validation → search "golang validate struct fields"
**Depends on:** None

---

### FILE: `backend/internal/model/alert.go`
**Purpose of this file:** Defines the Alert domain model representing a threat alert generated by the detection engine.
**Week:** Week 2
**Layer:** Backend

---

#### TASK-034: Alert Struct
**Signature:** `type Alert struct { ... }`
**File Location:** `backend/internal/model/alert.go`
**Purpose:** Represents an alert record in the alerts table.
**Why it exists:** Alerts are the analyst-facing output of the threat detection pipeline. When a log is classified as dangerous, an alert is created linking back to that log. This struct carries all the information an analyst needs to triage the threat: what rule fired, how severe it is, what the AI thinks, and whether it's been resolved.
**Fields to define:**
- `ID` (uuid.UUID) — `json:"id" db:"id"`
- `LogID` (uuid.UUID) — `json:"log_id" db:"log_id"` — references the triggering security log
- `RuleName` (string) — `json:"rule_name" db:"rule_name"` — which detection rule matched
- `Severity` (string) — `json:"severity" db:"severity"` — threat level
- `AISummary` (*string) — `json:"ai_summary" db:"ai_summary"` — nullable Claude analysis
- `Resolved` (bool) — `json:"resolved" db:"resolved"` — has an analyst addressed this?
- `ResolvedAt` (*time.Time) — `json:"resolved_at" db:"resolved_at"` — nullable timestamp
- `CreatedAt` (time.Time) — `json:"created_at" db:"created_at"`
**Inputs:** N/A (struct definition)
**Returns:** N/A (struct definition)
**Concepts to learn before implementing:**
- Nullable time.Time as *time.Time → search "golang nullable time"
- Foreign key relationships in domain models → search "golang model relationships"
**Depends on:** TASK-032

---

### FILE: `backend/internal/repository/log_repo.go`
**Purpose of this file:** Data access layer for the security_logs table. Contains all SQL queries for inserting, querying, filtering, and aggregating security logs.
**Week:** Week 2
**Layer:** Backend

---

#### TASK-035: NewLogRepo
**Signature:** `func NewLogRepo(db *sql.DB) *LogRepo`
**File Location:** `backend/internal/repository/log_repo.go`
**Purpose:** Constructor that creates a LogRepo with a database connection pool.
**Why it exists:** Same pattern as NewUserRepo — injects the database dependency for testability and clean architecture.
**Inputs:**
- `db` (*sql.DB) — open PostgreSQL connection pool.
**Returns:** `*LogRepo` — initialized repository.
**Concepts to learn before implementing:**
- (Same constructor pattern as TASK-008)
**Depends on:** TASK-018

---

#### TASK-036: Insert
**Signature:** `func (r *LogRepo) Insert(ctx context.Context, log *model.SecurityLog) error`
**File Location:** `backend/internal/repository/log_repo.go`
**Purpose:** Inserts a new security log into the security_logs table.
**Why it exists:** Every ingested security event needs to be persisted. This function executes an INSERT with RETURNING to populate the log's auto-generated ID and timestamp. It's called by the ingest handler (synchronous path) and later by the Kafka consumer (async path).
**Inputs:**
- `ctx` (context.Context) — request context.
- `log` (*model.SecurityLog) — log entry with Source, Message, ThreatLevel, IPAddress, UserAgent populated.
**Returns:** `error` — returned on database failures.
**Concepts to learn before implementing:**
- INSERT with multiple columns → search "postgresql INSERT multiple columns"
- Scanning RETURNING values → search "golang sql scan returning"
**Depends on:** TASK-035, TASK-032

---

#### TASK-037: InsertWithSummary
**Signature:** `func (r *LogRepo) InsertWithSummary(ctx context.Context, log *model.SecurityLog, summary string) error`
**File Location:** `backend/internal/repository/log_repo.go`
**Purpose:** Inserts a security log along with its AI-generated summary in a single database operation.
**Why it exists:** When the AI analysis is available at insert time (during the Kafka consumer pipeline), we can write both the log and its AI summary in one INSERT instead of doing an INSERT followed by an UPDATE. This reduces database round-trips and prevents a race condition where the log exists briefly without its summary.
**Inputs:**
- `ctx` (context.Context) — request context.
- `log` (*model.SecurityLog) — the security log entry.
- `summary` (string) — the Claude AI analysis text.
**Returns:** `error` — returned on database failures.
**Concepts to learn before implementing:**
- Reducing database round-trips → search "database round trip optimization"
- When to use INSERT vs INSERT+UPDATE → search "upsert pattern"
**Depends on:** TASK-035, TASK-032

---

#### TASK-038: List
**Signature:** `func (r *LogRepo) List(ctx context.Context, level string, source string, limit int, offset int) ([]model.SecurityLog, int, error)`
**File Location:** `backend/internal/repository/log_repo.go`
**Purpose:** Retrieves a paginated, filtered list of security logs for the Logs page.
**Why it exists:** The logs page needs to display logs in a table with filters (by threat level, by source) and pagination (page 1, page 2, etc.). This function builds a dynamic SQL query with optional WHERE clauses based on which filters are active, plus LIMIT/OFFSET for pagination. It also returns the total count for the pagination UI.
**Inputs:**
- `ctx` (context.Context) — request context.
- `level` (string) — threat level filter (empty string = no filter).
- `source` (string) — source filter (empty string = no filter).
- `limit` (int) — number of rows per page (e.g., 20).
- `offset` (int) — number of rows to skip (e.g., page 2 with limit 20 → offset 20).
**Returns:**
- `[]model.SecurityLog` — the page of log entries.
- `int` — total count of matching logs (for "Showing 1-20 of 150").
- `error` — returned on database failures.
**Concepts to learn before implementing:**
- Dynamic SQL query building in Go → search "golang dynamic sql where clause"
- LIMIT and OFFSET pagination → search "postgresql LIMIT OFFSET pagination"
- COUNT(*) for total results → search "postgresql COUNT with pagination"
- Two-query pagination pattern (count + data) → search "pagination total count query"
**Depends on:** TASK-035, TASK-032

---

#### TASK-039: CountByLevel
**Signature:** `func (r *LogRepo) CountByLevel(ctx context.Context) (map[string]int, error)`
**File Location:** `backend/internal/repository/log_repo.go`
**Purpose:** Returns a count of logs grouped by threat level for the dashboard statistics.
**Why it exists:** The dashboard displays stat cards showing "15 Critical", "42 High", etc. and a bar chart of threat level distribution. This function runs a GROUP BY query to get these counts in a single database call instead of four separate queries.
**Inputs:**
- `ctx` (context.Context) — request context.
**Returns:**
- `map[string]int` — maps threat level names to their counts, e.g., `{"critical": 15, "high": 42, ...}`
- `error` — returned on database failures.
**Concepts to learn before implementing:**
- SQL GROUP BY and COUNT → search "postgresql GROUP BY COUNT"
- Go maps → search "golang map tutorial"
- Scanning multiple rows with `rows.Next()` → search "golang sql rows iteration"
**Depends on:** TASK-035

---

#### TASK-040: FindByID
**Signature:** `func (r *LogRepo) FindByID(ctx context.Context, id uuid.UUID) (*model.SecurityLog, error)`
**File Location:** `backend/internal/repository/log_repo.go`
**Purpose:** Retrieves a single security log by its UUID.
**Why it exists:** When viewing log details or when the alert system needs to reference the original log, this function fetches the complete log entry. It's also used by the alert creation flow to verify the log exists before creating an alert for it.
**Inputs:**
- `ctx` (context.Context) — request context.
- `id` (uuid.UUID) — the log's unique identifier.
**Returns:**
- `*model.SecurityLog` — the log entry, or nil if not found.
- `error` — returned on database failures.
**Concepts to learn before implementing:**
- Single-row queries with QueryRowContext → search "golang QueryRowContext"
- Handling sql.ErrNoRows → search "golang sql.ErrNoRows handling"
**Depends on:** TASK-035, TASK-032

---

### FILE: `backend/internal/repository/alert_repo.go`
**Purpose of this file:** Data access layer for the alerts table. Handles creating, listing, resolving, and querying alerts.
**Week:** Week 2
**Layer:** Backend

---

#### TASK-041: NewAlertRepo
**Signature:** `func NewAlertRepo(db *sql.DB) *AlertRepo`
**File Location:** `backend/internal/repository/alert_repo.go`
**Purpose:** Constructor that creates an AlertRepo with a database connection pool.
**Why it exists:** Same dependency injection pattern as other repos.
**Inputs:**
- `db` (*sql.DB) — open PostgreSQL connection pool.
**Returns:** `*AlertRepo` — initialized repository.
**Concepts to learn before implementing:**
- (Same constructor pattern as TASK-008)
**Depends on:** TASK-018

---

#### TASK-042: Create
**Signature:** `func (r *AlertRepo) Create(ctx context.Context, alert *model.Alert) error`
**File Location:** `backend/internal/repository/alert_repo.go`
**Purpose:** Inserts a new alert into the alerts table.
**Why it exists:** When the threat detector classifies a log as high/critical, an alert is generated and persisted via this function. The alert links to the triggering log via log_id and stores the rule name and severity. This is the "write" side of the detection pipeline.
**Inputs:**
- `ctx` (context.Context) — request context.
- `alert` (*model.Alert) — alert with LogID, RuleName, Severity populated.
**Returns:** `error` — returned on insert failure.
**Concepts to learn before implementing:**
- Foreign key inserts → search "postgresql insert with foreign key"
- Transactional integrity when creating related records → search "golang sql transaction"
**Depends on:** TASK-041, TASK-034

---

#### TASK-043: List
**Signature:** `func (r *AlertRepo) List(ctx context.Context, resolvedFilter *bool) ([]model.Alert, error)`
**File Location:** `backend/internal/repository/alert_repo.go`
**Purpose:** Retrieves alerts, optionally filtered by resolved status.
**Why it exists:** The alerts page shows active (unresolved) alerts by default but allows viewing resolved ones too. The optional boolean filter (nil = all, true = resolved only, false = unresolved only) provides this flexibility. Results are ordered by creation time descending (newest first).
**Inputs:**
- `ctx` (context.Context) — request context.
- `resolvedFilter` (*bool) — nil returns all; pointer-to-true returns resolved; pointer-to-false returns unresolved.
**Returns:**
- `[]model.Alert` — list of matching alerts.
- `error` — returned on database failures.
**Concepts to learn before implementing:**
- Optional filter parameters with pointer types → search "golang optional parameter pointer"
- Dynamic WHERE clause building → search "golang conditional SQL where"
**Depends on:** TASK-041, TASK-034

---

#### TASK-044: Resolve
**Signature:** `func (r *AlertRepo) Resolve(ctx context.Context, id uuid.UUID) error`
**File Location:** `backend/internal/repository/alert_repo.go`
**Purpose:** Marks an alert as resolved by setting resolved=TRUE and resolved_at=NOW().
**Why it exists:** When an analyst investigates a threat and determines it's handled (or a false positive), they click "Resolve" in the UI. This function updates the alert record to reflect that action. It's an UPDATE, not a DELETE — resolved alerts are kept for audit trail purposes.
**Inputs:**
- `ctx` (context.Context) — request context.
- `id` (uuid.UUID) — the alert's unique identifier.
**Returns:** `error` — returned on update failure or if no alert found with that ID.
**Concepts to learn before implementing:**
- SQL UPDATE SET WHERE → search "postgresql UPDATE SET WHERE"
- NOW() for timestamps → search "postgresql NOW function"
- Audit trail pattern (update vs delete) → search "soft delete audit trail pattern"
**Depends on:** TASK-041

---

#### TASK-045: FindByLogID
**Signature:** `func (r *AlertRepo) FindByLogID(ctx context.Context, logID uuid.UUID) (*model.Alert, error)`
**File Location:** `backend/internal/repository/alert_repo.go`
**Purpose:** Checks if an alert already exists for a given log entry.
**Why it exists:** When the detection pipeline processes a log, it needs to check if an alert has already been created for that log (idempotency). Without this check, reprocessing the same log (e.g., after a Kafka consumer restart) would create duplicate alerts.
**Inputs:**
- `ctx` (context.Context) — request context.
- `logID` (uuid.UUID) — the security log's ID to check.
**Returns:**
- `*model.Alert` — the existing alert, or nil if none exists for this log.
- `error` — returned on database failures.
**Concepts to learn before implementing:**
- Idempotency in event processing → search "idempotent consumer pattern"
- Foreign key lookups → search "postgresql query by foreign key"
**Depends on:** TASK-041, TASK-034

---

### FILE: `backend/internal/detector/rules.go`
**Purpose of this file:** Implements the rule-based threat detection engine. Contains threat detection rules and a classifier that scans log messages for suspicious patterns.
**Week:** Week 2
**Layer:** Backend

---

#### TASK-046: Rule Struct
**Signature:** `type Rule struct { Name string; Pattern string; Level string }`
**File Location:** `backend/internal/detector/rules.go`
**Purpose:** Defines the structure of a single threat detection rule.
**Why it exists:** Each rule has a human-readable name (shown in alerts), a keyword/pattern to search for in log messages, and a threat level to assign when matched. The struct makes rules data-driven — you can add new detection rules by adding to the rules slice without changing the Classify logic.
**Fields:**
- `Name` (string) — human-readable rule name, e.g., "SQL Injection Detected"
- `Pattern` (string) — keyword or pattern to look for in log messages, e.g., "SELECT * FROM", "DROP TABLE"
- `Level` (string) — threat level to assign when this rule matches (uses ThreatLevel constants)
**Inputs:** N/A (struct definition)
**Returns:** N/A (struct definition)
**Concepts to learn before implementing:**
- Rule engine pattern → search "rule engine design pattern"
- Data-driven design → search "data driven programming"
**Depends on:** TASK-030

---

#### TASK-047: Default Rules Slice
**Signature:** `var rules = []Rule{ ... }`
**File Location:** `backend/internal/detector/rules.go`
**Purpose:** Defines the default set of threat detection rules that ship with GoShield.
**Why it exists:** These are the built-in detection rules that classify common attack patterns. Each rule targets a specific attack vector. Example rules:
- "SQL Injection" — patterns: "SELECT", "DROP TABLE", "UNION SELECT", "OR 1=1" → critical
- "XSS Attack" — patterns: "\<script\>", "javascript:" → high
- "Path Traversal" — patterns: "../", "/etc/passwd" → high
- "Brute Force" — patterns: "failed login", "invalid password" → medium
- "Port Scan" — patterns: "port scan", "SYN flood" → medium
**Inputs:** N/A (package-level variable)
**Returns:** N/A (package-level variable)
**Concepts to learn before implementing:**
- Package-level variables in Go → search "golang package level variable"
- Common cyberattack patterns (SQL injection, XSS, path traversal) → search "OWASP top 10 attacks"
- Keyword-based threat detection → search "keyword pattern matching IDS"
**Depends on:** TASK-046

---

#### TASK-048: Classify
**Signature:** `func Classify(message string) (string, string)`
**File Location:** `backend/internal/detector/rules.go`
**Purpose:** Scans a log message against all rules and returns the highest threat level match.
**Why it exists:** This is the core detection function. When a new log arrives, Classify is called with the log message. It iterates through all rules, checks if any pattern matches (case-insensitive substring match), and returns the matching rule name and threat level. If no rule matches, it returns an empty rule name and "low" threat level.
**Inputs:**
- `message` (string) — the raw log message to analyze.
**Returns:**
- `string` — the name of the matched rule (empty string if no match).
- `string` — the threat level ("low" if no match, otherwise the matched rule's level).
**Concepts to learn before implementing:**
- strings.Contains for substring matching → search "golang strings.Contains"
- strings.ToLower for case-insensitive comparison → search "golang case insensitive string compare"
- Multiple return values in Go → search "golang multiple return values"
- Iterating slices with range → search "golang range loop"
**Depends on:** TASK-047

---

### FILE: `backend/internal/ratelimit/redis.go`
**Purpose of this file:** Implements a Redis-based rate limiter using the sliding window pattern. Prevents API abuse by limiting requests per client per time window.
**Week:** Week 2
**Layer:** Backend

---

#### TASK-049: New
**Signature:** `func New(redisURL string, limit int, window time.Duration) (*Limiter, error)`
**File Location:** `backend/internal/ratelimit/redis.go`
**Purpose:** Creates a new rate Limiter connected to Redis.
**Why it exists:** Rate limiting is essential for API security — without it, an attacker could flood the ingest endpoint with millions of fake logs, overwhelming the database and Kafka. Redis is used as the backing store because it's fast (in-memory), atomic (INCR is thread-safe), and shared across multiple API server instances.
**Inputs:**
- `redisURL` (string) — Redis connection URL from config.
- `limit` (int) — maximum number of requests allowed per window (e.g., 100).
- `window` (time.Duration) — the time window for rate limiting (e.g., 1 minute).
**Returns:**
- `*Limiter` — initialized rate limiter.
- `error` — returned if Redis connection fails.
**Concepts to learn before implementing:**
- Rate limiting concepts (sliding window, token bucket) → search "rate limiting algorithms explained"
- Redis basics (key-value store, INCR, EXPIRE) → search "redis tutorial basics"
- go-redis client library → search "golang go-redis v9 tutorial"
**Depends on:** TASK-001

---

#### TASK-050: Allow
**Signature:** `func (l *Limiter) Allow(ctx context.Context, key string) (bool, error)`
**File Location:** `backend/internal/ratelimit/redis.go`
**Purpose:** Checks if a request from the given key (usually client IP) should be allowed or rate-limited.
**Why it exists:** Called on every incoming request (as middleware or within handlers). It increments a counter in Redis for the given key, sets an expiry on the key if it's new, and checks if the count exceeds the limit. Returns true if the request is allowed, false if rate-limited. Redis INCR is atomic, so this works correctly even with concurrent requests.
**Inputs:**
- `ctx` (context.Context) — request context.
- `key` (string) — rate limit key, typically the client's IP address or API key.
**Returns:**
- `bool` — true if request is allowed, false if rate limit exceeded.
- `error` — returned on Redis communication failure.
**Concepts to learn before implementing:**
- Redis INCR command (atomic increment) → search "redis INCR command"
- Redis EXPIRE for TTL → search "redis EXPIRE TTL"
- Sliding window rate limiting with Redis → search "redis sliding window rate limiter"
- Extracting client IP from HTTP request → search "golang get client IP address"
**Depends on:** TASK-049

---

### FILE: `backend/internal/handler/log_handler.go`
**Purpose of this file:** HTTP handlers for security log operations — ingesting new logs, listing logs with filters, and fetching individual log details.
**Week:** Week 2
**Layer:** Backend

---

#### TASK-051: NewLogHandler
**Signature:** `func NewLogHandler(logRepo *repository.LogRepo, alertRepo *repository.AlertRepo, detector func(string) (string, string), limiter *ratelimit.Limiter) *LogHandler`
**File Location:** `backend/internal/handler/log_handler.go`
**Purpose:** Constructor that creates a LogHandler with all its dependencies.
**Why it exists:** The log handler is the most complex handler — it needs the log repository (to store logs), the alert repository (to create alerts when threats are detected), the detector function (to classify logs), and the rate limiter (to prevent abuse). This constructor wires all four dependencies.
**Inputs:**
- `logRepo` (*repository.LogRepo) — for log CRUD operations.
- `alertRepo` (*repository.AlertRepo) — for creating alerts on threat detection.
- `detector` (func(string) (string, string)) — the Classify function from the detector package.
- `limiter` (*ratelimit.Limiter) — for rate limiting the ingest endpoint.
**Returns:** `*LogHandler` — fully initialized handler.
**Concepts to learn before implementing:**
- Function types as parameters → search "golang function as parameter"
- Handler with multiple dependencies → search "golang handler multiple dependencies"
**Depends on:** TASK-035, TASK-041, TASK-048, TASK-049

---

#### TASK-052: Ingest
**Signature:** `func (h *LogHandler) Ingest(w http.ResponseWriter, r *http.Request)`
**File Location:** `backend/internal/handler/log_handler.go`
**Purpose:** Handles POST /api/logs — ingests a new security log, classifies it, and optionally creates an alert.
**Why it exists:** This is the primary data entry point for GoShield. External systems (firewalls, IDS, etc.) send logs here. The handler: (1) rate-limits the request, (2) parses the JSON body into an IngestRequest, (3) runs the message through the threat detector, (4) inserts the log into the database, (5) if threat is high/critical, creates an alert, (6) returns the created log. This is the synchronous ingestion path — Week 3 adds the Kafka-based async path.
**Inputs:**
- `w` (http.ResponseWriter) — response writer.
- `r` (*http.Request) — contains JSON body matching IngestRequest struct.
**Returns:** Writes 201 Created with the created SecurityLog as JSON, or 429 Too Many Requests if rate-limited, or 400 Bad Request on invalid input.
**Concepts to learn before implementing:**
- HTTP 429 Too Many Requests → search "HTTP 429 status code"
- Orchestrating multiple operations in a handler → search "golang handler business logic"
- Error response patterns → search "golang http error response JSON"
**Depends on:** TASK-051, TASK-036, TASK-042, TASK-048, TASK-050

---

#### TASK-053: List
**Signature:** `func (h *LogHandler) List(w http.ResponseWriter, r *http.Request)`
**File Location:** `backend/internal/handler/log_handler.go`
**Purpose:** Handles GET /api/logs — returns a paginated, filtered list of security logs.
**Why it exists:** The logs page needs to fetch logs with optional filters (threat level, source) and pagination (page number, page size). This handler extracts query parameters from the URL, passes them to the repository, and returns the results along with pagination metadata (total count, current page, total pages).
**Inputs:**
- `w` (http.ResponseWriter) — response writer.
- `r` (*http.Request) — URL query parameters: `level`, `source`, `page`, `limit`.
**Returns:** Writes 200 OK with `{"data": [...], "total": 150, "page": 1, "limit": 20}` JSON response.
**Concepts to learn before implementing:**
- Extracting URL query parameters → search "golang http request query parameters"
- `strconv.Atoi` for string-to-int conversion → search "golang strconv.Atoi"
- Pagination response format → search "REST API pagination response format"
- Calculating total pages: `math.Ceil(total / limit)` → search "golang math.Ceil"
**Depends on:** TASK-051, TASK-038

---

#### TASK-054: GetByID
**Signature:** `func (h *LogHandler) GetByID(w http.ResponseWriter, r *http.Request)`
**File Location:** `backend/internal/handler/log_handler.go`
**Purpose:** Handles GET /api/logs/{id} — returns a single security log by its UUID.
**Why it exists:** When a user clicks on a log row to see its details (including the full message and AI summary), the frontend fetches the complete log entry by ID. This handler extracts the UUID from the URL path, queries the database, and returns the full log object.
**Inputs:**
- `w` (http.ResponseWriter) — response writer.
- `r` (*http.Request) — URL path parameter `id` extracted via Chi's `chi.URLParam(r, "id")`.
**Returns:** Writes 200 OK with the SecurityLog as JSON, or 404 Not Found if no log with that ID exists.
**Concepts to learn before implementing:**
- Chi URL parameters → search "go-chi URL parameters chi.URLParam"
- UUID parsing from string → search "golang uuid.Parse"
- HTTP 404 Not Found handling → search "golang http 404 response"
**Depends on:** TASK-051, TASK-040

---

### FILE: `backend/internal/handler/alert_handler.go`
**Purpose of this file:** HTTP handlers for alert operations — listing active alerts and resolving them.
**Week:** Week 2
**Layer:** Backend

---

#### TASK-055: NewAlertHandler
**Signature:** `func NewAlertHandler(alertRepo *repository.AlertRepo) *AlertHandler`
**File Location:** `backend/internal/handler/alert_handler.go`
**Purpose:** Constructor that creates an AlertHandler with the alert repository.
**Why it exists:** The alert handler only needs the alert repository — it doesn't need the log repo or detector because alerts are queried and resolved, not created through this handler (alerts are created in the log ingestion flow).
**Inputs:**
- `alertRepo` (*repository.AlertRepo) — for alert query and update operations.
**Returns:** `*AlertHandler` — initialized handler.
**Concepts to learn before implementing:**
- (Same constructor pattern as previous handlers)
**Depends on:** TASK-041

---

#### TASK-056: List
**Signature:** `func (h *AlertHandler) List(w http.ResponseWriter, r *http.Request)`
**File Location:** `backend/internal/handler/alert_handler.go`
**Purpose:** Handles GET /api/alerts — returns alerts, optionally filtered by resolved status.
**Why it exists:** The alerts page displays active threats that need analyst attention. By default it shows unresolved alerts, but analysts can toggle to see resolved ones for reference. The optional `resolved` query parameter controls the filter.
**Inputs:**
- `w` (http.ResponseWriter) — response writer.
- `r` (*http.Request) — optional query parameter `resolved` ("true" or "false").
**Returns:** Writes 200 OK with an array of Alert objects as JSON.
**Concepts to learn before implementing:**
- Optional query parameters → search "golang optional query parameter"
- `strconv.ParseBool` for boolean query params → search "golang strconv.ParseBool"
**Depends on:** TASK-055, TASK-043

---

#### TASK-057: Resolve
**Signature:** `func (h *AlertHandler) Resolve(w http.ResponseWriter, r *http.Request)`
**File Location:** `backend/internal/handler/alert_handler.go`
**Purpose:** Handles PATCH /api/alerts/{id}/resolve — marks an alert as resolved.
**Why it exists:** When an analyst investigates an alert and determines it's been handled (real threat mitigated or false positive confirmed), they resolve it. This handler takes the alert UUID from the URL, calls the repository's Resolve function to update the database, and returns a success response. Using PATCH (not PUT) because we're partially updating the resource.
**Inputs:**
- `w` (http.ResponseWriter) — response writer.
- `r` (*http.Request) — URL path parameter `id` (the alert's UUID).
**Returns:** Writes 200 OK with `{"message": "alert resolved"}` on success, or 404 if alert not found.
**Concepts to learn before implementing:**
- HTTP PATCH vs PUT → search "HTTP PATCH vs PUT difference"
- RESTful resource actions → search "REST API action endpoints design"
- Chi URL parameter extraction → search "go-chi URLParam"
**Depends on:** TASK-055, TASK-044


---
---

## 🔄 WEEK 3 — Kafka Log Ingestion Pipeline

### FILE: `backend/internal/kafka/producer.go`
**Purpose of this file:** Implements a Kafka producer that sends security log messages to a Kafka topic for asynchronous processing. Decouples the HTTP ingestion endpoint from the heavy processing pipeline.
**Week:** Week 3
**Layer:** Backend

---

#### TASK-058: NewProducer
**Signature:** `func NewProducer(brokers []string, topic string) (*Producer, error)`
**File Location:** `backend/internal/kafka/producer.go`
**Purpose:** Creates a new Kafka producer connected to the specified brokers for a given topic.
**Why it exists:** Decoupling log ingestion from processing is critical at scale. The API server pushes logs to Kafka instead of processing them synchronously, which improves throughput (the API returns immediately) and reliability (if the consumer crashes, messages are retained in Kafka until processed).
**Inputs:**
- `brokers` ([]string) — Kafka broker addresses from config, e.g., `["localhost:9092"]`.
- `topic` (string) — Kafka topic name, e.g., `"security-logs"`.
**Returns:**
- `*Producer` — initialized Kafka producer.
- `error` — returned if broker connection fails.
**Concepts to learn before implementing:**
- Apache Kafka basics (topics, producers, consumers, brokers) → search "kafka basics tutorial"
- franz-go or confluent-kafka-go library → search "golang kafka producer franz-go"
- Message queues and async processing → search "message queue pattern"
- Kafka topic partitioning → search "kafka partitions explained"
**Depends on:** TASK-001

---

#### TASK-059: Send
**Signature:** `func (p *Producer) Send(ctx context.Context, log *model.SecurityLog) error`
**File Location:** `backend/internal/kafka/producer.go`
**Purpose:** Serializes a SecurityLog to JSON and publishes it to the Kafka topic.
**Why it exists:** This is how logs enter the async pipeline. The API handler calls Send after receiving a log, putting it on the Kafka topic for the consumer to process with threat detection and AI analysis. JSON serialization is used because it's human-readable and easy to debug (you can inspect Kafka messages directly).
**Inputs:**
- `ctx` (context.Context) — request context for cancellation.
- `log` (*model.SecurityLog) — the log to publish.
**Returns:** `error` — returned if JSON serialization or Kafka publish fails.
**Concepts to learn before implementing:**
- JSON serialization with json.Marshal → search "golang json.Marshal"
- Kafka message publishing → search "kafka produce message golang"
- Async vs sync processing trade-offs → search "synchronous vs asynchronous processing"
- Message keys for partitioning → search "kafka message key partitioning"
**Depends on:** TASK-058, TASK-032

---

#### TASK-060: Close
**Signature:** `func (p *Producer) Close()`
**File Location:** `backend/internal/kafka/producer.go`
**Purpose:** Flushes any pending messages and closes the Kafka producer connection.
**Why it exists:** Kafka producers buffer messages for efficiency (batching multiple messages into one network call). Close ensures all buffered messages are sent before the application shuts down. Without it, you'd lose the last batch of messages on shutdown — potentially dropping security events.
**Inputs:** None
**Returns:** None (may log errors internally)
**Concepts to learn before implementing:**
- Resource cleanup in Go (defer pattern) → search "golang defer resource cleanup"
- Kafka producer flush → search "kafka producer flush close"
- Graceful shutdown → search "golang graceful shutdown pattern"
**Depends on:** TASK-058

---

### FILE: `backend/internal/kafka/consumer.go`
**Purpose of this file:** Implements a Kafka consumer that reads security logs from the topic, classifies them using the rule engine, and persists them to the database with generated alerts.
**Week:** Week 3
**Layer:** Backend

---

#### TASK-061: NewConsumer
**Signature:** `func NewConsumer(brokers []string, topic string, groupID string, logRepo *repository.LogRepo, alertRepo *repository.AlertRepo) (*Consumer, error)`
**File Location:** `backend/internal/kafka/consumer.go`
**Purpose:** Creates a new Kafka consumer subscribed to the security logs topic.
**Why it exists:** The consumer is the processing engine of the async pipeline. It reads logs from Kafka, runs threat detection, and stores results. The `groupID` enables consumer groups — multiple consumer instances with the same group ID share the topic's partitions, enabling horizontal scaling.
**Inputs:**
- `brokers` ([]string) — Kafka broker addresses.
- `topic` (string) — topic to consume from.
- `groupID` (string) — Kafka consumer group ID for load balancing.
- `logRepo` (*repository.LogRepo) — for inserting classified logs.
- `alertRepo` (*repository.AlertRepo) — for creating alerts on threat detection.
**Returns:**
- `*Consumer` — initialized consumer.
- `error` — returned if subscription fails.
**Concepts to learn before implementing:**
- Kafka consumer groups → search "kafka consumer group explained"
- Consumer group rebalancing → search "kafka consumer rebalancing"
- At-least-once delivery semantics → search "kafka delivery semantics"
- Partition assignment → search "kafka consumer partition assignment"
**Depends on:** TASK-001, TASK-035, TASK-041

---

#### TASK-062: Start
**Signature:** `func (c *Consumer) Start(ctx context.Context) error`
**File Location:** `backend/internal/kafka/consumer.go`
**Purpose:** Begins consuming messages from Kafka in a continuous loop. For each message: (1) deserializes the JSON log, (2) classifies it with the rule engine, (3) inserts it into the database, (4) creates an alert if threat level is high/critical.
**Why it exists:** This is the heart of the async processing pipeline. The Start method runs indefinitely, processing every log that enters the Kafka topic. The inner EachRecord callback handles the actual logic for each individual message. Using a context allows graceful shutdown — canceling the context stops the loop.
**Inner callback flow:** Deserialize JSON → call `detector.Classify(message)` → call `logRepo.Insert()` → if high/critical, call `alertRepo.Create()`
**Inputs:**
- `ctx` (context.Context) — used for cancellation to stop the consumer gracefully.
**Returns:** `error` — returned if the consumer loop encounters a fatal error.
**Concepts to learn before implementing:**
- Long-running goroutines → search "golang long running goroutine"
- JSON deserialization with json.Unmarshal → search "golang json.Unmarshal"
- Context cancellation for graceful stop → search "golang context.Done cancel"
- Kafka offset commit (message acknowledgment) → search "kafka offset commit"
**Depends on:** TASK-061, TASK-048, TASK-036, TASK-042

---

#### TASK-063: Close
**Signature:** `func (c *Consumer) Close()`
**File Location:** `backend/internal/kafka/consumer.go`
**Purpose:** Stops the consumer loop and closes the Kafka connection.
**Why it exists:** Clean shutdown prevents duplicate message processing. When the consumer closes cleanly, it commits its current offset so it won't reprocess messages on restart. Without clean shutdown, the next consumer start would replay all uncommitted messages.
**Inputs:** None
**Returns:** None
**Concepts to learn before implementing:**
- Kafka offset management → search "kafka consumer offset commit"
- Clean consumer shutdown → search "kafka consumer graceful shutdown golang"
**Depends on:** TASK-061

---

### FILE: `backend/cmd/consumer/main.go`
**Purpose of this file:** Entry point for the standalone Kafka consumer process. This runs as a separate binary/container from the API server, following the microservice separation pattern.
**Week:** Week 3
**Layer:** Backend

---

#### TASK-064: main
**Signature:** `func main()`
**File Location:** `backend/cmd/consumer/main.go`
**Purpose:** Loads config, connects to PostgreSQL, creates repos, creates a Kafka consumer, and starts consuming in a long-running loop.
**Why it exists:** The consumer runs as a separate process from the API server. This separation of concerns means the API server handles HTTP requests while the consumer handles async log processing. They scale independently — you can run 1 API server and 3 consumers if processing is the bottleneck. The separate binary pattern (`cmd/consumer/`) is standard in Go monorepos.
**Inputs:** Environment variables (via config.Load).
**Returns:** Blocks forever (consumer.Start runs continuously). Uses signal handling (SIGTERM, SIGINT) for graceful shutdown.
**Concepts to learn before implementing:**
- Separate binaries in Go (cmd/ directory pattern) → search "golang cmd directory multiple binaries"
- Signal handling (SIGTERM, SIGINT) → search "golang os signal notify"
- Microservice separation → search "microservice process separation"
- `go build ./cmd/consumer` to build specific binary → search "golang build specific package"
**Depends on:** TASK-001, TASK-018, TASK-035, TASK-041, TASK-061, TASK-062

---
---

## 🧠 WEEK 4 — gRPC Threat Classification Service

### FILE: `backend/proto/classifier.proto`
**Purpose of this file:** Defines the gRPC service contract (protocol buffer definition) for the threat classification microservice. This is the interface definition that both the server and client use — like a shared API schema.
**Week:** Week 4
**Layer:** Backend

---

#### TASK-065: ClassifyRequest Message
**Signature:** `message ClassifyRequest { string message = 1; string source = 2; string ip_address = 3; }`
**File Location:** `backend/proto/classifier.proto`
**Purpose:** Defines the input message for a classification RPC call.
**Why it exists:** Protocol Buffers define strongly-typed, language-agnostic message formats. The ClassifyRequest describes exactly what data the client sends to the classification service. Field numbers (`= 1`, `= 2`) are used for binary encoding — they must never change once deployed, as they're how protobuf identifies fields on the wire.
**Fields to define:**
- `message` (string, field 1) — the log message to classify
- `source` (string, field 2) — log source for context
- `ip_address` (string, field 3) — source IP for context
**Inputs:** N/A (message definition)
**Returns:** N/A (message definition)
**Concepts to learn before implementing:**
- Protocol Buffers (protobuf) basics → search "protobuf tutorial basics"
- Proto3 syntax → search "proto3 syntax guide"
- Field numbering rules → search "protobuf field number rules"
- .proto file structure (syntax, package, option, import) → search "protobuf file structure"
**Depends on:** None

---

#### TASK-066: ClassifyResponse Message
**Signature:** `message ClassifyResponse { string threat_level = 1; string rule_name = 2; float confidence = 3; }`
**File Location:** `backend/proto/classifier.proto`
**Purpose:** Defines the output message returned after classification.
**Why it exists:** The response tells the caller what threat level was assigned, which rule matched, and a confidence score. The confidence field (0.0–1.0) could be used by the AI layer to decide if human review is needed — low confidence classifications get sent to Claude for a second opinion.
**Fields to define:**
- `threat_level` (string, field 1) — the assigned severity level
- `rule_name` (string, field 2) — which detection rule fired
- `confidence` (float, field 3) — certainty of classification (0.0–1.0)
**Inputs:** N/A (message definition)
**Returns:** N/A (message definition)
**Concepts to learn before implementing:**
- Protobuf scalar types (string, float, int32, bool) → search "protobuf scalar types"
- Response message design → search "protobuf response message best practices"
**Depends on:** None

---

#### TASK-067: ClassifierService Service Definition
**Signature:** `service ClassifierService { rpc Classify(ClassifyRequest) returns (ClassifyResponse); }`
**File Location:** `backend/proto/classifier.proto`
**Purpose:** Defines the gRPC service interface with its RPC methods.
**Why it exists:** The service definition is the contract between client and server. The `protoc` compiler generates Go server interfaces and client stubs from this definition. Any server implementing ClassifierService must provide a Classify method. Any client using the generated stub can call Classify — this is how gRPC achieves language-agnostic communication.
**Inputs:** N/A (service definition)
**Returns:** N/A (service definition)
**Concepts to learn before implementing:**
- gRPC service definition → search "grpc service definition protobuf"
- RPC (Remote Procedure Call) concept → search "RPC explained simply"
- Service-oriented architecture → search "service oriented architecture microservices"
**Depends on:** TASK-065, TASK-066

---

#### TASK-068: Classify RPC Method
**Signature:** `rpc Classify(ClassifyRequest) returns (ClassifyResponse);`
**File Location:** `backend/proto/classifier.proto`
**Purpose:** Declares the Classify RPC method within the ClassifierService.
**Why it exists:** This declares a unary RPC — the simplest gRPC pattern where the client sends one request and gets one response (like a regular function call over the network). For each log message, the client calls Classify and gets back a threat level and rule name. Future extensions could add streaming RPCs for batch classification.
**Inputs:** ClassifyRequest message
**Returns:** ClassifyResponse message
**Concepts to learn before implementing:**
- Unary RPC vs streaming RPC → search "grpc unary vs streaming"
- Code generation with protoc → search "protoc generate go code grpc"
- protoc-gen-go and protoc-gen-go-grpc plugins → search "golang protoc gen go grpc"
- Running protoc command → search "protobuf compile proto file golang"
**Depends on:** TASK-067

---

### FILE: `backend/classifier/server.go`
**Purpose of this file:** Implements the gRPC ClassifierService server — receives classification requests over gRPC and returns threat analysis results using the same rule engine as the HTTP path.
**Week:** Week 4
**Layer:** Backend

---

#### TASK-069: Server Struct
**Signature:** `type Server struct { pb.UnimplementedClassifierServiceServer }`
**File Location:** `backend/classifier/server.go`
**Purpose:** Defines the gRPC server struct that implements the ClassifierService interface.
**Why it exists:** In gRPC-Go, you implement a service by embedding the `Unimplemented*` server struct (generated by protoc). This satisfies the interface and provides default "unimplemented" responses for any methods you haven't written yet — enabling forward compatibility when new RPCs are added to the proto file without breaking existing servers.
**Fields:**
- `pb.UnimplementedClassifierServiceServer` — embedded struct for forward compatibility
**Inputs:** N/A (struct definition)
**Returns:** N/A (struct definition)
**Concepts to learn before implementing:**
- gRPC server implementation in Go → search "golang grpc server implementation"
- Embedded structs in Go → search "golang embedded struct"
- Forward compatibility with UnimplementedServer → search "grpc unimplemented server golang"
**Depends on:** TASK-067

---

#### TASK-070: Classify Method
**Signature:** `func (s *Server) Classify(ctx context.Context, req *pb.ClassifyRequest) (*pb.ClassifyResponse, error)`
**File Location:** `backend/classifier/server.go`
**Purpose:** Implements the Classify RPC — receives a log message and returns its threat classification.
**Why it exists:** This is the actual business logic of the gRPC service. It takes the message from the request, runs it through the same `detector.Classify` function used in the sync HTTP path, and returns the result as a ClassifyResponse. This means the classification logic is reusable — both the HTTP handler and the gRPC service use the same rules engine.
**Inputs:**
- `ctx` (context.Context) — gRPC request context.
- `req` (*pb.ClassifyRequest) — contains message, source, ip_address fields.
**Returns:**
- `*pb.ClassifyResponse` — populated with threat_level, rule_name, and confidence.
- `error` — returned as a gRPC status error if something fails.
**Concepts to learn before implementing:**
- gRPC method implementation → search "golang implement grpc method"
- gRPC status codes (OK, InvalidArgument, Internal) → search "grpc status codes golang"
- Protobuf generated code usage → search "golang use protobuf generated code"
- Reusing business logic across transports (HTTP + gRPC) → search "hexagonal architecture ports adapters"
**Depends on:** TASK-069, TASK-048

---

### FILE: `backend/cmd/classifier/main.go`
**Purpose of this file:** Entry point for the gRPC classifier microservice. Runs as a separate process from the API server and consumer, following the single-responsibility principle.
**Week:** Week 4
**Layer:** Backend

---

#### TASK-071: main
**Signature:** `func main()`
**File Location:** `backend/cmd/classifier/main.go`
**Purpose:** Creates a gRPC server, registers the ClassifierService, and starts listening on a TCP port.
**Why it exists:** Like the consumer (TASK-064), the classifier runs as its own process. This follows the microservice pattern — each service has a single responsibility. The classifier only classifies threats. It's independently deployable and scalable — if classification becomes a bottleneck, you can run more classifier instances behind a gRPC load balancer.
**Steps:** (1) Load config, (2) create a `net.Listener` on a TCP port (e.g., `:50051`), (3) create `grpc.NewServer()`, (4) register the ClassifierService server with `pb.RegisterClassifierServiceServer()`, (5) call `grpcServer.Serve(listener)` to start accepting gRPC connections.
**Inputs:** Environment variables (via config.Load), TCP port for gRPC.
**Returns:** Blocks forever (Serve runs continuously). Exits on fatal errors.
**Concepts to learn before implementing:**
- gRPC server setup in Go → search "golang grpc server setup tutorial"
- net.Listen for TCP → search "golang net.Listen tcp"
- gRPC default port conventions → search "grpc default port 50051"
- Registering gRPC services → search "golang register grpc service"
**Depends on:** TASK-001, TASK-069, TASK-070

---
---

## 🤖 WEEK 5 — AI-Powered Anomaly Analysis

### FILE: `backend/internal/ai/claude.go`
**Purpose of this file:** HTTP client for the Claude AI API. Sends security log data to Claude for intelligent anomaly analysis and receives natural-language summaries that help analysts understand and respond to threats.
**Week:** Week 5
**Layer:** Backend

---

#### TASK-072: Client Struct
**Signature:** `type Client struct { apiKey string; httpClient *http.Client; model string }`
**File Location:** `backend/internal/ai/claude.go`
**Purpose:** Defines the Claude API client with authentication and HTTP configuration.
**Why it exists:** Encapsulates all Claude API interaction details — the API key for authentication, an HTTP client (with timeout configuration to prevent hanging requests), and the model name (e.g., `"claude-sonnet-4-20250514"`). Having a dedicated struct means the rest of the application calls simple methods like `AnalyzeLog` without knowing HTTP details.
**Fields:**
- `apiKey` (string) — Claude API key from config, sent as `x-api-key` header
- `httpClient` (*http.Client) — configured with timeout (e.g., 30 seconds)
- `model` (string) — Claude model identifier (e.g., `"claude-sonnet-4-20250514"`)
**Inputs:** N/A (struct definition)
**Returns:** N/A (struct definition)
**Concepts to learn before implementing:**
- HTTP client configuration in Go → search "golang http.Client timeout"
- API client struct pattern → search "golang API client pattern"
- Claude AI API overview → search "anthropic claude API documentation"
**Depends on:** TASK-001

---

#### TASK-073: NewClient
**Signature:** `func NewClient(apiKey string) *Client`
**File Location:** `backend/internal/ai/claude.go`
**Purpose:** Constructor that creates a Claude client with the API key and sensible defaults.
**Why it exists:** Sets up the HTTP client with a reasonable timeout (e.g., 30 seconds — Claude can take a few seconds to respond) and the default Claude model. The API key is injected from config rather than hardcoded, supporting different keys for dev/staging/prod.
**Inputs:**
- `apiKey` (string) — from config.ClaudeAPIKey.
**Returns:** `*Client` — initialized Claude client ready to make API calls.
**Concepts to learn before implementing:**
- Constructor pattern with sensible defaults → search "golang constructor default values"
- HTTP client timeout best practices → search "golang http client timeout best practices"
**Depends on:** TASK-072

---

#### TASK-074: request Struct (unexported)
**Signature:** `type request struct { Model string \`json:"model"\`; MaxTokens int \`json:"max_tokens"\`; Messages []message \`json:"messages"\` }`
**File Location:** `backend/internal/ai/claude.go`
**Purpose:** Defines the JSON structure for Claude API requests.
**Why it exists:** Claude's API expects a specific JSON format with model, max_tokens, and a messages array (conversation format). Defining this as a struct lets you use `json.Marshal` to serialize it correctly. The struct is unexported (lowercase `request`) because only the claude package uses it — it's an internal implementation detail.
**Fields:**
- `Model` (string) — `json:"model"` — which Claude model to use
- `MaxTokens` (int) — `json:"max_tokens"` — maximum response length
- `Messages` ([]message) — `json:"messages"` — array of `{role, content}` pairs
**Inputs:** N/A (struct definition)
**Returns:** N/A (struct definition)
**Concepts to learn before implementing:**
- Unexported types in Go (lowercase) → search "golang exported vs unexported"
- Claude API request format → search "claude API messages format"
- Nested struct for JSON serialization → search "golang nested struct json"
**Depends on:** None

---

#### TASK-075: response Struct (unexported)
**Signature:** `type response struct { Content []contentBlock \`json:"content"\` }`
**File Location:** `backend/internal/ai/claude.go`
**Purpose:** Defines the JSON structure for Claude API responses.
**Why it exists:** Claude returns a response with a `content` array containing text blocks. This struct captures the relevant parts of the response for deserialization. Only the text content is needed — other metadata (usage, stop_reason) can be ignored for this use case. Go's JSON decoder silently ignores unrecognized fields.
**Fields:**
- `Content` ([]contentBlock) — each block has `Type` and `Text` fields
**Inputs:** N/A (struct definition)
**Returns:** N/A (struct definition)
**Concepts to learn before implementing:**
- Selective JSON deserialization (ignoring extra fields) → search "golang json unmarshal ignore fields"
- Claude API response format → search "claude API response structure"
**Depends on:** None

---

#### TASK-076: AnalyzeLog
**Signature:** `func (c *Client) AnalyzeLog(ctx context.Context, log *model.SecurityLog) (string, error)`
**File Location:** `backend/internal/ai/claude.go`
**Purpose:** Sends a security log to Claude for AI analysis and returns the analysis as a text string.
**Why it exists:** This is the AI integration point. For high/critical threat logs, Claude analyzes the log message, IP address, source, and threat level to provide: (1) a plain-English explanation of the threat, (2) potential impact assessment, (3) recommended response actions. This analysis is stored as `ai_summary` in both the log and alert records, giving analysts expert-level context instantly.
**Implementation flow:** (1) Build a prompt describing the security log with context, (2) create a request struct with the prompt as a user message, (3) serialize to JSON with json.Marshal, (4) POST to `https://api.anthropic.com/v1/messages` with headers (`x-api-key`, `anthropic-version`, `content-type`), (5) read and deserialize the response body, (6) extract the text content from the first content block.
**Inputs:**
- `ctx` (context.Context) — request context for timeout/cancellation.
- `log` (*model.SecurityLog) — the log to analyze.
**Returns:**
- `string` — Claude's analysis text.
- `error` — returned on HTTP failures, API errors, or timeout.
**Concepts to learn before implementing:**
- Making HTTP POST requests in Go → search "golang http.Post JSON"
- Setting custom HTTP headers → search "golang http request headers"
- Claude API authentication (x-api-key header) → search "claude API authentication"
- Prompt engineering for security analysis → search "LLM prompt engineering security"
- Reading and closing response bodies → search "golang http response body read close"
**Depends on:** TASK-073, TASK-032

---

### FILE: `backend/internal/ai/analyzer.go`
**Purpose of this file:** Helper functions that decide when to invoke AI analysis and how to parse the response. Keeps decision logic separate from HTTP client logic for cleaner testing.
**Week:** Week 5
**Layer:** Backend

---

#### TASK-077: ShouldAnalyze
**Signature:** `func ShouldAnalyze(threatLevel string) bool`
**File Location:** `backend/internal/ai/analyzer.go`
**Purpose:** Determines whether a log's threat level warrants AI analysis.
**Why it exists:** Calling Claude for every single log would be expensive ($$$) and slow. This function implements the decision logic: only analyze logs classified as "high" or "critical" by the rule engine. Low and medium threats don't get AI analysis — they're too routine. This is both a cost optimization (fewer API calls) and a performance optimization (most logs skip the 2-3 second AI delay).
**Inputs:**
- `threatLevel` (string) — one of the ThreatLevel constants.
**Returns:** `bool` — true if the log should be sent to Claude, false otherwise.
**Concepts to learn before implementing:**
- Feature gating logic → search "feature flag conditional logic"
- API cost optimization → search "LLM API cost optimization"
- Separation of concerns (decision vs execution) → search "separation of concerns programming"
**Depends on:** TASK-030

---

#### TASK-078: ParseResponse
**Signature:** `func ParseResponse(raw string) string`
**File Location:** `backend/internal/ai/analyzer.go`
**Purpose:** Cleans and formats the raw Claude response text for storage and display.
**Why it exists:** Claude's responses may include markdown formatting, excessive whitespace, or preamble text ("Here's my analysis..."). This function normalizes the response — trimming whitespace, removing unnecessary formatting, and potentially extracting structured sections. Clean summaries display better in the frontend AlertCard component and are more useful in reports.
**Inputs:**
- `raw` (string) — the raw text from Claude's response.
**Returns:** `string` — cleaned and formatted analysis text.
**Concepts to learn before implementing:**
- String manipulation in Go (strings.TrimSpace, strings.Replace) → search "golang strings package"
- Text cleaning and normalization → search "text normalization programming"
**Depends on:** None

---

### FILE: `backend/internal/kafka/consumer.go` (UPDATE)
**Purpose of this file:** Updated consumer that integrates AI analysis into the processing pipeline, enriching high-threat logs with Claude's analysis before storage.
**Week:** Week 5
**Layer:** Backend

---

#### TASK-079: Updated Start with AI Integration
**Signature:** `func (c *Consumer) Start(ctx context.Context) error` (modified)
**File Location:** `backend/internal/kafka/consumer.go`
**Purpose:** Updates the existing Start method to add AI analysis between classification and database insertion.
**Why it exists:** The original Start (TASK-062) classified logs and stored them. Now the pipeline gains a new enrichment step: after classification, if `ShouldAnalyze` returns true, call `claude.AnalyzeLog` to get an AI summary, then use `InsertWithSummary` instead of `Insert`. This enriches high-threat logs with AI-generated context that helps analysts understand and respond to threats faster.
**Updated flow:** Deserialize → Classify → **ShouldAnalyze?** → **If yes: AnalyzeLog → ParseResponse** → InsertWithSummary (or Insert if no analysis) → Create alert if high/critical → **Broadcast alert via WebSocket**
**Inputs:** Same as TASK-062
**Returns:** Same as TASK-062
**Concepts to learn before implementing:**
- Pipeline pattern (sequential processing stages) → search "pipeline pattern golang"
- Enrichment pattern in event processing → search "event enrichment pattern"
- Error handling for optional steps (AI failure shouldn't block log storage) → search "golang graceful error handling non-critical"
- Fallback logic (if AI fails, still insert without summary) → search "graceful degradation pattern"
**Depends on:** TASK-062, TASK-077, TASK-076, TASK-078, TASK-037


---
---

## 🌐 WEEK 6 — WebSocket & Full Frontend

### FILE: `backend/internal/ws/hub.go`
**Purpose of this file:** Manages WebSocket connections and broadcasts real-time alerts to all connected frontend clients. This is the server-side engine for the live alert feed.
**Week:** Week 6
**Layer:** Backend

---

#### TASK-080: Hub Struct
**Signature:** `type Hub struct { clients map[*websocket.Conn]bool; broadcast chan []byte; register chan *websocket.Conn; unregister chan *websocket.Conn }`
**File Location:** `backend/internal/ws/hub.go`
**Purpose:** Defines the central WebSocket connection manager with channels for concurrent-safe communication.
**Why it exists:** The Hub is the central coordinator for all WebSocket communication. It uses Go channels to safely manage concurrent connection/disconnection events and message broadcasting without race conditions. The `clients` map tracks which connections are alive. Using channels instead of mutexes is idiomatic Go — "don't communicate by sharing memory; share memory by communicating."
**Fields:**
- `clients` (map[*websocket.Conn]bool) — set of active WebSocket connections
- `broadcast` (chan []byte) — channel for outgoing messages to all clients
- `register` (chan *websocket.Conn) — channel for new connections
- `unregister` (chan *websocket.Conn) — channel for disconnections
**Inputs:** N/A (struct definition)
**Returns:** N/A (struct definition)
**Concepts to learn before implementing:**
- Go channels → search "golang channels tutorial"
- Go maps as sets → search "golang map set pattern"
- WebSocket protocol → search "websocket protocol explained"
- Concurrency-safe data structures → search "golang concurrent map access"
**Depends on:** None

---

#### TASK-081: NewHub
**Signature:** `func NewHub() *Hub`
**File Location:** `backend/internal/ws/hub.go`
**Purpose:** Constructor that initializes all channels and the clients map, then starts the run goroutine.
**Why it exists:** Channels must be initialized with `make()` before use, otherwise sends/receives will panic (nil channel deadlock). The constructor ensures all channels are ready and kicks off the `run` goroutine that processes events in the background.
**Inputs:** None
**Returns:** `*Hub` — fully initialized hub with the run goroutine already started.
**Concepts to learn before implementing:**
- make() for channels and maps → search "golang make channel map"
- Starting goroutines with `go` keyword → search "golang go keyword goroutine"
**Depends on:** TASK-080

---

#### TASK-082: run
**Signature:** `func (h *Hub) run()`
**File Location:** `backend/internal/ws/hub.go`
**Purpose:** Internal goroutine that continuously listens on register, unregister, and broadcast channels using a select loop.
**Why it exists:** The `select` statement is Go's multiplexer for channels. The `run` loop continuously waits for one of three events: a new client connecting (add to map), a client disconnecting (remove from map and close connection), or a message to broadcast (send to every connected client). This single goroutine is the **only** one that reads/writes the clients map, eliminating the need for locks or mutexes.
**Inputs:** None (reads from Hub's channels internally)
**Returns:** None (runs forever as a goroutine)
**Concepts to learn before implementing:**
- Go select statement → search "golang select channel"
- Goroutine lifecycle → search "golang goroutine tutorial"
- Fan-out pattern (one message to many clients) → search "golang fan out pattern"
- Closing WebSocket connections → search "gorilla websocket close"
**Depends on:** TASK-080

---

#### TASK-083: Register
**Signature:** `func (h *Hub) Register(conn *websocket.Conn)`
**File Location:** `backend/internal/ws/hub.go`
**Purpose:** Sends a connection to the register channel so the run goroutine adds it to the clients map.
**Why it exists:** Public method called by the WebSocket upgrade handler (TASK-086) when a new client connects. Sends the connection through the channel so the `run` goroutine can safely add it to the clients map without race conditions.
**Inputs:**
- `conn` (*websocket.Conn) — the newly upgraded WebSocket connection.
**Returns:** None
**Concepts to learn before implementing:**
- Channel send operation (`ch <- value`) → search "golang channel send"
**Depends on:** TASK-080

---

#### TASK-084: Unregister
**Signature:** `func (h *Hub) Unregister(conn *websocket.Conn)`
**File Location:** `backend/internal/ws/hub.go`
**Purpose:** Sends a connection to the unregister channel so the run goroutine removes it from the clients map.
**Why it exists:** Called when a client disconnects or a write error occurs. The `run` goroutine removes the connection from the map and closes it. Without proper unregistration, the server would attempt to write to dead connections, causing errors and resource leaks.
**Inputs:**
- `conn` (*websocket.Conn) — the disconnected WebSocket connection.
**Returns:** None
**Depends on:** TASK-080

---

#### TASK-085: Broadcast
**Signature:** `func (h *Hub) Broadcast(message []byte)`
**File Location:** `backend/internal/ws/hub.go`
**Purpose:** Sends a message to the broadcast channel, which the run goroutine delivers to every connected client.
**Why it exists:** Called whenever a new alert is created in the processing pipeline. The message (serialized alert JSON) is sent to every connected frontend client in real-time. This is how the dashboard's live alert feed gets updated without polling — true push-based updates.
**Inputs:**
- `message` ([]byte) — JSON-encoded alert data to send to all clients.
**Returns:** None
**Depends on:** TASK-080

---

### FILE: `backend/cmd/server/main.go` (UPDATE)
**Purpose of this file:** Adds a WebSocket upgrade endpoint to the existing API server routes.
**Week:** Week 6
**Layer:** Backend

---

#### TASK-086: WebSocket Upgrade Handler
**Signature:** Part of `func main()` — registers a `GET /ws` route
**File Location:** `backend/cmd/server/main.go`
**Purpose:** Upgrades HTTP connections to WebSocket protocol and registers them with the Hub.
**Why it exists:** WebSocket connections start as regular HTTP GET requests that are "upgraded" to a persistent, bidirectional connection via the HTTP 101 Switching Protocols response. The `gorilla/websocket.Upgrader` handles the protocol switch. After upgrading, the connection is registered with the Hub, and a read pump goroutine listens for client disconnection (or incoming messages from the client).
**Inputs:** Standard HTTP handler inputs (`w`, `r`). Uses the Upgrader's `CheckOrigin` to allow cross-origin requests in development.
**Returns:** Upgrades the connection (no HTTP response body — the connection becomes a WebSocket).
**Concepts to learn before implementing:**
- WebSocket upgrade process → search "websocket upgrade http"
- gorilla/websocket library → search "golang gorilla websocket tutorial"
- HTTP 101 Switching Protocols → search "HTTP 101 switching protocols"
- WebSocket read pump pattern → search "gorilla websocket read pump"
**Depends on:** TASK-022, TASK-081, TASK-083

---

### FILE: `frontend/src/types/user.ts`
**Purpose of this file:** TypeScript type definitions for user-related data structures, mirroring the backend User model for type-safe frontend development.
**Week:** Week 6
**Layer:** Frontend

---

#### TASK-087: Role Type
**Signature:** `export type Role = 'analyst' | 'admin';`
**File Location:** `frontend/src/types/user.ts`
**Purpose:** Defines a TypeScript union type restricting role values to exactly two valid options.
**Why it exists:** TypeScript union types provide compile-time safety — if you try to compare `role === 'superadmin'`, TypeScript flags it as an error immediately. This mirrors the backend's RoleAnalyst/RoleAdmin constants, ensuring frontend and backend agree on valid values.
**Inputs:** N/A (type definition)
**Returns:** N/A (type definition)
**Concepts to learn before implementing:**
- TypeScript union types → search "typescript union type"
- Literal types in TypeScript → search "typescript literal types"
- Shared types between frontend and backend → search "typescript api types"
**Depends on:** None

---

#### TASK-088: User Interface
**Signature:** `export interface User { id: string; email: string; role: Role; createdAt: string; updatedAt: string; }`
**File Location:** `frontend/src/types/user.ts`
**Purpose:** Defines the shape of a User object as returned by the API.
**Why it exists:** TypeScript interfaces define the expected shape of objects. Every API response containing user data is typed with this interface, giving you autocompletion in your editor and compile-time checking that you're accessing valid properties. Timestamps are `string` (not Date) because JSON serialization converts dates to ISO strings.
**Inputs:** N/A (interface definition)
**Returns:** N/A (interface definition)
**Concepts to learn before implementing:**
- TypeScript interfaces → search "typescript interface tutorial"
- Mapping backend models to frontend types → search "typescript api response types"
**Depends on:** TASK-087

---

### FILE: `frontend/src/types/log.ts`
**Purpose of this file:** TypeScript types for security logs, mirroring the backend SecurityLog model.
**Week:** Week 6
**Layer:** Frontend

---

#### TASK-089: ThreatLevel Type
**Signature:** `export type ThreatLevel = 'low' | 'medium' | 'high' | 'critical';`
**File Location:** `frontend/src/types/log.ts`
**Purpose:** Union type for the four threat severity levels.
**Why it exists:** Mirrors the backend ThreatLevel constants. Used for type-safe filtering, badge coloring (each level maps to a specific color), and chart data categorization. If a new level is added to the backend but not the frontend, TypeScript will catch the mismatch.
**Inputs:** N/A (type definition)
**Returns:** N/A (type definition)
**Depends on:** None

---

#### TASK-090: LogSource Type
**Signature:** `export type LogSource = 'firewall' | 'ids' | 'auth_service' | 'endpoint';`
**File Location:** `frontend/src/types/log.ts`
**Purpose:** Union type for valid log source identifiers.
**Why it exists:** Mirrors the backend LogSource constants. Used in filter dropdowns (so users can only select valid sources) and pie chart categorization (each source gets a distinct slice).
**Inputs:** N/A (type definition)
**Returns:** N/A (type definition)
**Depends on:** None

---

#### TASK-091: SecurityLog Interface
**Signature:** `export interface SecurityLog { id: string; source: LogSource; message: string; threatLevel: ThreatLevel; ipAddress: string; userAgent: string; aiSummary: string | null; createdAt: string; }`
**File Location:** `frontend/src/types/log.ts`
**Purpose:** Complete type definition for a security log entry.
**Why it exists:** This is the frontend's representation of a security log. The `aiSummary` is nullable (`string | null`) matching the backend's `*string` pointer — not all logs have AI analysis. Field names use camelCase (TypeScript convention) which should match the JSON response keys.
**Inputs:** N/A (interface definition)
**Returns:** N/A (interface definition)
**Concepts to learn before implementing:**
- Nullable types in TypeScript → search "typescript null type"
- Interface vs type alias → search "typescript interface vs type"
**Depends on:** TASK-089, TASK-090

---

### FILE: `frontend/src/types/alert.ts`
**Purpose of this file:** TypeScript type for the Alert model.
**Week:** Week 6
**Layer:** Frontend

---

#### TASK-092: Alert Interface
**Signature:** `export interface Alert { id: string; logId: string; ruleName: string; severity: ThreatLevel; aiSummary: string | null; resolved: boolean; resolvedAt: string | null; createdAt: string; }`
**File Location:** `frontend/src/types/alert.ts`
**Purpose:** Defines the shape of an Alert object as returned by the API and received via WebSocket.
**Why it exists:** Mirrors the backend Alert struct. Used in AlertCard rendering, AlertFeed list mapping, Redux alertSlice state, and WebSocket message parsing. Both `aiSummary` and `resolvedAt` are nullable because they're only set under certain conditions.
**Inputs:** N/A (interface definition)
**Returns:** N/A (interface definition)
**Depends on:** TASK-089

---

### FILE: `frontend/src/api/client.ts`
**Purpose of this file:** Configures the Axios HTTP client instance with base URL, authentication interceptors, and global error handling. Every API call in the application goes through this configured instance.
**Week:** Week 6
**Layer:** Frontend

---

#### TASK-093: axiosInstance Configuration
**Signature:** `export const axiosInstance = axios.create({ baseURL: '/api', headers: { 'Content-Type': 'application/json' } });`
**File Location:** `frontend/src/api/client.ts`
**Purpose:** Creates a pre-configured Axios instance with the API base URL and default headers.
**Why it exists:** Creating a configured Axios instance means every API call automatically uses the right base URL and headers without repetition. The `/api` base URL works with the Vite proxy (TASK-175) in development and with Nginx reverse proxy in production.
**Inputs:** N/A (configuration object)
**Returns:** An Axios instance used by all API modules.
**Concepts to learn before implementing:**
- Axios instance creation → search "axios create instance"
- Base URL configuration → search "axios baseURL"
- Default headers → search "axios default headers"
**Depends on:** None

---

#### TASK-094: Request Interceptor
**Signature:** `axiosInstance.interceptors.request.use((config) => { ... return config; });`
**File Location:** `frontend/src/api/client.ts`
**Purpose:** Attaches the JWT token from localStorage to every outgoing request's Authorization header.
**Why it exists:** After login, the JWT is stored in localStorage. Instead of manually adding `Authorization: Bearer <token>` to every API call, the interceptor automatically attaches it. This is the DRY principle — write the auth header logic once, apply it everywhere.
**Inputs:** Axios request config object (modified in-place to add the Authorization header).
**Returns:** Modified config with the Bearer token attached.
**Concepts to learn before implementing:**
- Axios interceptors → search "axios interceptors tutorial"
- localStorage for JWT tokens → search "jwt localstorage"
- Authorization Bearer header format → search "authorization bearer token header"
**Depends on:** TASK-093

---

#### TASK-095: Response Interceptor
**Signature:** `axiosInstance.interceptors.response.use((response) => response, (error) => { ... });`
**File Location:** `frontend/src/api/client.ts`
**Purpose:** Handles 401 Unauthorized errors globally by clearing auth state and redirecting to the login page.
**Why it exists:** If a token expires while the user is on the dashboard, subsequent API calls return 401. The response interceptor catches this globally, clears the stored token from localStorage, and redirects to `/login`. Without it, the user would see confusing error messages on every component that fetches data.
**Inputs:** Axios response (pass-through on success) or error (handled on failure).
**Returns:** The response on success; a rejected promise on error (after handling 401).
**Concepts to learn before implementing:**
- Axios error interceptor → search "axios response interceptor error handling"
- Token expiration handling in SPAs → search "jwt token expiration frontend"
- Global error handling → search "axios global error handler"
**Depends on:** TASK-093

---

### FILE: `frontend/src/api/auth.api.ts`
**Purpose of this file:** API functions for authentication endpoints — login, register, and fetching the current user's profile.
**Week:** Week 6
**Layer:** Frontend

---

#### TASK-096: login
**Signature:** `export const login = (email: string, password: string) => axiosInstance.post<{ token: string; user: User }>('/login', { email, password });`
**File Location:** `frontend/src/api/auth.api.ts`
**Purpose:** Sends a login request and returns the JWT token and user profile.
**Why it exists:** Wraps the login POST request with TypeScript generics for type-safe responses. The caller gets back a typed object with `token` and `user`, not an untyped `any`. This enables autocompletion and compile-time checks.
**Inputs:**
- `email` (string) — user's email address.
- `password` (string) — user's plaintext password (sent over HTTPS).
**Returns:** Promise resolving to `{ data: { token: string; user: User } }`.
**Concepts to learn before implementing:**
- Axios POST with TypeScript generics → search "axios typescript generic post"
- API service layer pattern → search "frontend api service layer"
**Depends on:** TASK-093, TASK-088

---

#### TASK-097: register
**Signature:** `export const register = (email: string, password: string, role: Role) => axiosInstance.post<{ user: User }>('/register', { email, password, role });`
**File Location:** `frontend/src/api/auth.api.ts`
**Purpose:** Sends a registration request to create a new user account.
**Why it exists:** Same pattern as login — typed API function that the useAuth hook calls.
**Inputs:** `email` (string), `password` (string), `role` (Role — restricted to 'analyst' | 'admin').
**Returns:** Promise resolving to `{ data: { user: User } }`.
**Depends on:** TASK-093, TASK-088

---

#### TASK-098: me
**Signature:** `export const me = () => axiosInstance.get<User>('/me');`
**File Location:** `frontend/src/api/auth.api.ts`
**Purpose:** Fetches the currently authenticated user's profile.
**Why it exists:** After page refresh, the app needs to verify the stored token is still valid and fetch the user's profile. The request interceptor (TASK-094) automatically attaches the JWT token, so this function just makes a simple GET request.
**Inputs:** None (token is auto-attached by interceptor).
**Returns:** Promise resolving to `{ data: User }`.
**Depends on:** TASK-093, TASK-094, TASK-088

---

### FILE: `frontend/src/api/logs.api.ts`
**Purpose of this file:** API functions for security log endpoints — fetching, filtering, paginating, and ingesting logs.
**Week:** Week 6
**Layer:** Frontend

---

#### TASK-099: getLogs
**Signature:** `export const getLogs = (params: { level?: string; source?: string; page?: number; limit?: number }) => axiosInstance.get<{ data: SecurityLog[]; total: number; page: number; limit: number }>('/logs', { params });`
**File Location:** `frontend/src/api/logs.api.ts`
**Purpose:** Fetches a paginated, filtered list of security logs.
**Why it exists:** The `params` object maps directly to URL query parameters. Axios automatically serializes `{ level: 'high', page: 2 }` into `?level=high&page=2`. The generic type ensures the response shape is properly typed — you get autocompletion on `response.data.total`.
**Inputs:**
- `params` — optional filters: `level`, `source`, `page`, `limit`.
**Returns:** Promise with typed response containing `data` (log array), `total`, `page`, `limit`.
**Concepts to learn before implementing:**
- Axios query parameters → search "axios get request params"
- Optional properties in TypeScript → search "typescript optional properties"
**Depends on:** TASK-093, TASK-091

---

#### TASK-100: ingestLog
**Signature:** `export const ingestLog = (log: { source: string; message: string; ipAddress?: string; userAgent?: string }) => axiosInstance.post<SecurityLog>('/logs', log);`
**File Location:** `frontend/src/api/logs.api.ts`
**Purpose:** Sends a new security log to the ingest endpoint.
**Why it exists:** Used by the demo seeder or admin interface to manually inject test logs.
**Inputs:** Log object matching the IngestRequest shape.
**Returns:** Promise with the created SecurityLog.
**Depends on:** TASK-093, TASK-091

---

#### TASK-101: getLogById
**Signature:** `export const getLogById = (id: string) => axiosInstance.get<SecurityLog>(\`/logs/${id}\`);`
**File Location:** `frontend/src/api/logs.api.ts`
**Purpose:** Fetches a single security log by its UUID.
**Why it exists:** Used when clicking on a log row to see its full details, including the complete message and AI summary.
**Inputs:** `id` (string) — the log's UUID.
**Returns:** Promise with the SecurityLog.
**Depends on:** TASK-093, TASK-091

---

### FILE: `frontend/src/api/alerts.api.ts`
**Purpose of this file:** API functions for alert endpoints — listing alerts and resolving them.
**Week:** Week 6
**Layer:** Frontend

---

#### TASK-102: getAlerts
**Signature:** `export const getAlerts = (resolved?: boolean) => axiosInstance.get<Alert[]>('/alerts', { params: resolved !== undefined ? { resolved } : {} });`
**File Location:** `frontend/src/api/alerts.api.ts`
**Purpose:** Fetches alerts, optionally filtered by resolved status.
**Why it exists:** The conditional params logic ensures the `resolved` parameter is only included in the URL when explicitly specified. Calling `getAlerts()` fetches all alerts; `getAlerts(false)` fetches only unresolved ones.
**Inputs:** `resolved` (boolean | undefined) — optional filter.
**Returns:** Promise with an array of Alert objects.
**Depends on:** TASK-093, TASK-092

---

#### TASK-103: resolveAlert
**Signature:** `export const resolveAlert = (id: string) => axiosInstance.patch<{ message: string }>(\`/alerts/${id}/resolve\`);`
**File Location:** `frontend/src/api/alerts.api.ts`
**Purpose:** Marks an alert as resolved by calling the backend PATCH endpoint.
**Why it exists:** Uses PATCH (not PUT) to match the backend endpoint. The response confirms resolution with a message.
**Inputs:** `id` (string) — the alert's UUID.
**Returns:** Promise with `{ message: "alert resolved" }`.
**Depends on:** TASK-093

---

### FILE: `frontend/src/store/index.ts`
**Purpose of this file:** Configures the Redux Toolkit store, combining all slice reducers, and exports typed hooks for use throughout the application.
**Week:** Week 6
**Layer:** Frontend

---

#### TASK-104: configureStore
**Signature:** `export const store = configureStore({ reducer: { auth: authReducer, alerts: alertReducer, logs: logReducer } });`
**File Location:** `frontend/src/store/index.ts`
**Purpose:** Creates the Redux store combining all feature slices into a single state tree.
**Why it exists:** Redux Toolkit's `configureStore` combines all slice reducers, automatically sets up Redux DevTools for debugging, and includes default middleware (thunk for async actions). This is the single source of truth for all application state.
**Inputs:** Reducer map object.
**Returns:** The configured Redux store.
**Concepts to learn before implementing:**
- Redux Toolkit configureStore → search "redux toolkit configureStore tutorial"
- State management concepts → search "redux state management basics"
- Combining reducers → search "redux toolkit combine reducers"
**Depends on:** None (but will import from TASK-107, TASK-112, TASK-116 slices)

---

#### TASK-105: RootState Type
**Signature:** `export type RootState = ReturnType<typeof store.getState>;`
**File Location:** `frontend/src/store/index.ts`
**Purpose:** TypeScript utility type that automatically infers the shape of the entire Redux state.
**Why it exists:** Instead of manually defining the state shape, `ReturnType<typeof store.getState>` automatically derives it from the store configuration. Every `useSelector` call uses this type, so TypeScript knows exactly what properties exist in `state.auth.user` or `state.logs.filters`.
**Inputs:** N/A (type definition)
**Returns:** N/A (type definition)
**Concepts to learn before implementing:**
- ReturnType utility type → search "typescript ReturnType"
- Typed Redux hooks → search "redux toolkit typescript typed hooks"
**Depends on:** TASK-104

---

#### TASK-106: AppDispatch Type
**Signature:** `export type AppDispatch = typeof store.dispatch;`
**File Location:** `frontend/src/store/index.ts`
**Purpose:** The typed dispatch function for dispatching actions with full type safety.
**Why it exists:** Using the raw `useDispatch` hook returns a generic dispatch function. By creating `AppDispatch`, you ensure that only valid actions can be dispatched. Combined with a custom `useAppDispatch = () => useDispatch<AppDispatch>()` hook, it provides full type safety.
**Inputs:** N/A (type definition)
**Returns:** N/A (type definition)
**Depends on:** TASK-104

---

### FILE: `frontend/src/store/authSlice.ts`
**Purpose of this file:** Redux slice managing authentication state — the current user profile, JWT token, and login/logout actions.
**Week:** Week 6
**Layer:** Frontend

---

#### TASK-107: initialState
**Signature:** `const initialState: AuthState = { user: null, token: localStorage.getItem('token') };`
**File Location:** `frontend/src/store/authSlice.ts`
**Purpose:** Defines the initial auth state, restoring the JWT token from localStorage.
**Why it exists:** On app load, the token is restored from localStorage so the user stays logged in across page refreshes (session persistence). The user starts as `null` and is populated by calling the `/me` endpoint with the stored token. If no token exists in localStorage, the user is unauthenticated.
**Inputs:** N/A (constant definition)
**Returns:** N/A (constant definition)
**Concepts to learn before implementing:**
- Redux initial state → search "redux toolkit initial state"
- Persisting auth state → search "redux persist token localstorage"
- Session persistence patterns → search "spa session persistence"
**Depends on:** TASK-088

---

#### TASK-108: setCredentials Action
**Signature:** `setCredentials: (state, action: PayloadAction<{ user: User; token: string }>) => { state.user = action.payload.user; state.token = action.payload.token; localStorage.setItem('token', action.payload.token); }`
**File Location:** `frontend/src/store/authSlice.ts`
**Purpose:** Stores user profile and JWT token after successful login.
**Why it exists:** Called after successful login. Stores both the user profile in Redux state (for UI rendering — showing user email, role) and the token in localStorage (for persistence across page refreshes). Immer (built into Redux Toolkit) allows the direct mutation syntax (`state.user = ...`) — it produces an immutable update under the hood.
**Inputs:** PayloadAction with `{ user: User; token: string }` payload.
**Returns:** N/A (reducer — mutates state via Immer)
**Concepts to learn before implementing:**
- Redux Toolkit createSlice → search "redux toolkit createSlice tutorial"
- PayloadAction type → search "redux toolkit PayloadAction"
- Immer immutable updates → search "immer redux toolkit"
**Depends on:** TASK-107

---

#### TASK-109: logout Action
**Signature:** `logout: (state) => { state.user = null; state.token = null; localStorage.removeItem('token'); }`
**File Location:** `frontend/src/store/authSlice.ts`
**Purpose:** Clears all auth state and removes the persisted token.
**Why it exists:** When the user logs out (or when a 401 response is intercepted), all traces of authentication must be removed — from Redux state (so the UI updates immediately) and from localStorage (so refreshing the page doesn't restore the old session).
**Inputs:** No payload.
**Returns:** N/A (reducer)
**Depends on:** TASK-107

---

#### TASK-110: selectUser Selector
**Signature:** `export const selectUser = (state: RootState) => state.auth.user;`
**File Location:** `frontend/src/store/authSlice.ts`
**Purpose:** Type-safe selector for accessing the current user from any component.
**Why it exists:** Selectors provide a clean, reusable way to extract specific pieces of state. Instead of writing `useSelector((state: RootState) => state.auth.user)` in every component, you write `useSelector(selectUser)`. If the state shape changes, you only update the selector.
**Inputs:** `state` (RootState) — the entire Redux state tree.
**Returns:** `User | null` — the current user or null if not authenticated.
**Depends on:** TASK-105

---

#### TASK-111: selectToken Selector
**Signature:** `export const selectToken = (state: RootState) => state.auth.token;`
**File Location:** `frontend/src/store/authSlice.ts`
**Purpose:** Type-safe selector for accessing the JWT token.
**Why it exists:** Used by the request interceptor (to attach the token to requests) and by ProtectedRoute (to check if the user is authenticated). Having a selector means the token access logic is centralized.
**Inputs:** `state` (RootState).
**Returns:** `string | null`.
**Depends on:** TASK-105

---

### FILE: `frontend/src/store/alertSlice.ts`
**Purpose of this file:** Redux slice managing real-time alert state — receives new alerts from WebSocket and stores them for display.
**Week:** Week 6
**Layer:** Frontend

---

#### TASK-112: initialState
**Signature:** `const initialState: AlertState = { alerts: [] };`
**File Location:** `frontend/src/store/alertSlice.ts`
**Purpose:** Starts with an empty alerts array.
**Why it exists:** Alerts are populated from two sources: initial API fetch (on page load) and real-time WebSocket messages (pushed from server). The array holds Alert objects in reverse chronological order (newest first) for the live feed.
**Inputs:** N/A (constant)
**Returns:** N/A (constant)
**Depends on:** TASK-092

---

#### TASK-113: addAlert Action
**Signature:** `addAlert: (state, action: PayloadAction<Alert>) => { state.alerts.unshift(action.payload); }`
**File Location:** `frontend/src/store/alertSlice.ts`
**Purpose:** Adds a new alert to the beginning of the alerts array.
**Why it exists:** `unshift` adds the new alert at index 0 (newest first). Called by the `useAlerts` WebSocket hook when a real-time alert arrives. The UI updates instantly without a page refresh — this is the core of the real-time experience.
**Inputs:** PayloadAction with a single Alert object.
**Returns:** N/A (reducer)
**Concepts to learn before implementing:**
- Array unshift → search "javascript array unshift"
- Real-time state updates → search "redux real time updates websocket"
**Depends on:** TASK-112

---

#### TASK-114: clearAlerts Action
**Signature:** `clearAlerts: (state) => { state.alerts = []; }`
**File Location:** `frontend/src/store/alertSlice.ts`
**Purpose:** Resets the alerts array to empty.
**Why it exists:** Used when the user logs out (clear stale data) or when switching between alert views (resolved vs unresolved) to prevent mixing data from different queries.
**Inputs:** No payload.
**Returns:** N/A (reducer)
**Depends on:** TASK-112

---

#### TASK-115: selectAlerts Selector
**Signature:** `export const selectAlerts = (state: RootState) => state.alerts.alerts;`
**File Location:** `frontend/src/store/alertSlice.ts`
**Purpose:** Type-safe selector for the alerts array.
**Why it exists:** Used by AlertFeed and DashboardPage to access the current list of alerts.
**Inputs:** `state` (RootState).
**Returns:** `Alert[]`.
**Depends on:** TASK-105

---

### FILE: `frontend/src/store/logSlice.ts`
**Purpose of this file:** Redux slice managing security log state — the current page of logs, filters, pagination metadata.
**Week:** Week 6
**Layer:** Frontend

---

#### TASK-116: initialState
**Signature:** `const initialState: LogState = { logs: [], total: 0, page: 1, limit: 20, filters: { level: '', source: '' } };`
**File Location:** `frontend/src/store/logSlice.ts`
**Purpose:** Defines the initial state for the logs feature.
**Why it exists:** Stores the current page of logs, total count (for pagination controls), current page number, page size, and active filters. This centralized state lets multiple components (LogTable, LogFilters, pagination controls) share the same data source and stay in sync.
**Inputs:** N/A (constant)
**Returns:** N/A (constant)
**Depends on:** TASK-091

---

#### TASK-117: setLogs Action
**Signature:** `setLogs: (state, action: PayloadAction<{ data: SecurityLog[]; total: number }>) => { state.logs = action.payload.data; state.total = action.payload.total; }`
**File Location:** `frontend/src/store/logSlice.ts`
**Purpose:** Replaces the current log page data when a new page or filter result is fetched.
**Why it exists:** When the user changes pages or applies filters, the API returns a new set of logs. This action replaces the current logs array entirely (not append) and updates the total count for accurate pagination.
**Inputs:** PayloadAction with `{ data: SecurityLog[]; total: number }`.
**Returns:** N/A (reducer)
**Depends on:** TASK-116

---

#### TASK-118: setFilters Action
**Signature:** `setFilters: (state, action: PayloadAction<{ level?: string; source?: string }>) => { state.filters = { ...state.filters, ...action.payload }; state.page = 1; }`
**File Location:** `frontend/src/store/logSlice.ts`
**Purpose:** Updates filter values and resets pagination to page 1.
**Why it exists:** When filters change, the page must reset to 1. You don't want to stay on page 5 when the total results might now be only 3 pages. The spread operator merges new filter values with existing ones, so you can update just one filter at a time.
**Inputs:** PayloadAction with partial filter object.
**Returns:** N/A (reducer)
**Depends on:** TASK-116

---

#### TASK-119: setPage Action
**Signature:** `setPage: (state, action: PayloadAction<number>) => { state.page = action.payload; }`
**File Location:** `frontend/src/store/logSlice.ts`
**Purpose:** Updates the current page number for pagination.
**Why it exists:** When the user clicks "Next" or "Previous" or a specific page number, this action updates the page state, which triggers a re-fetch in the useLogs hook via useEffect dependency tracking.
**Inputs:** PayloadAction with the new page number.
**Returns:** N/A (reducer)
**Depends on:** TASK-116

---

#### TASK-120: selectLogs Selector
**Signature:** `export const selectLogs = (state: RootState) => state.logs.logs;`
**File Location:** `frontend/src/store/logSlice.ts`
**Purpose:** Type-safe selector for the current page of logs.
**Depends on:** TASK-105

---

#### TASK-121: selectFilters Selector
**Signature:** `export const selectFilters = (state: RootState) => state.logs.filters;`
**File Location:** `frontend/src/store/logSlice.ts`
**Purpose:** Type-safe selector for the active filter values.
**Why it exists:** Used by LogFilters component to display current filter selections and by useLogs hook to build API query parameters.
**Depends on:** TASK-105

---

### FILE: `frontend/src/hooks/useAuth.ts`
**Purpose of this file:** Custom React hook encapsulating all authentication logic — login, logout, and authentication status checks. Components call this hook instead of interacting with Redux and API modules directly.
**Week:** Week 6
**Layer:** Frontend

---

#### TASK-122: useAuth — login function
**Signature:** `const login = async (email: string, password: string) => { const { data } = await authApi.login(email, password); dispatch(setCredentials(data)); };`
**File Location:** `frontend/src/hooks/useAuth.ts`
**Purpose:** Combines the API call and Redux dispatch into one reusable function.
**Why it exists:** Any component that needs login functionality calls `useAuth().login()` instead of duplicating the API call + dispatch logic. The hook handles the full flow: make the API request, extract the token and user from the response, and store both in Redux.
**Inputs:** `email` (string), `password` (string).
**Returns:** Promise (void) — resolves on success, throws on failure.
**Concepts to learn before implementing:**
- Custom React hooks → search "react custom hooks tutorial"
- Async/await in hooks → search "react hook async function"
- Composing multiple operations → search "react hook composition pattern"
**Depends on:** TASK-096, TASK-108

---

#### TASK-123: useAuth — logout function
**Signature:** `const logout = () => { dispatch(logoutAction()); navigate('/login'); };`
**File Location:** `frontend/src/hooks/useAuth.ts`
**Purpose:** Clears all auth state and redirects to the login page.
**Why it exists:** Logout needs to happen in two places simultaneously: (1) clear Redux state and localStorage (via the logout action), and (2) navigate to the login page. This function combines both steps.
**Inputs:** None.
**Returns:** void.
**Concepts to learn before implementing:**
- React Router useNavigate → search "react router useNavigate"
**Depends on:** TASK-109

---

#### TASK-124: useAuth — isAuthenticated
**Signature:** `const isAuthenticated = !!token;`
**File Location:** `frontend/src/hooks/useAuth.ts`
**Purpose:** A derived boolean indicating whether the user is currently authenticated.
**Why it exists:** Used by ProtectedRoute (TASK-170) to gate access to protected pages. The double negation (`!!`) converts a possibly-null token string to a boolean — null/empty becomes `false`, any string becomes `true`.
**Inputs:** Token value from `useSelector(selectToken)`.
**Returns:** `boolean`.
**Concepts to learn before implementing:**
- Double negation for boolean coercion → search "javascript double negation"
- Derived state from selectors → search "redux derived state"
**Depends on:** TASK-111

---

### FILE: `frontend/src/hooks/useAlerts.ts`
**Purpose of this file:** Custom hook that establishes a WebSocket connection to the backend and dispatches incoming alerts to the Redux store in real-time.
**Week:** Week 6
**Layer:** Frontend

---

#### TASK-125: useAlerts — WebSocket setup
**Signature:** Inside `useEffect`: `const ws = new WebSocket(\`ws://${window.location.host}/ws\`);`
**File Location:** `frontend/src/hooks/useAlerts.ts`
**Purpose:** Creates a WebSocket connection to the backend's `/ws` endpoint when the component mounts.
**Why it exists:** The WebSocket URL uses `window.location.host` so it automatically works in both development (proxied through Vite to localhost:8080) and production (same origin). The connection is established inside a useEffect to tie its lifecycle to the component.
**Inputs:** None (uses browser's WebSocket API).
**Returns:** A WebSocket instance.
**Concepts to learn before implementing:**
- WebSocket API in JavaScript → search "javascript websocket tutorial"
- React useEffect for subscriptions → search "react useEffect websocket"
- WebSocket URL construction → search "websocket url scheme ws wss"
**Depends on:** TASK-086

---

#### TASK-126: useAlerts — onmessage handler
**Signature:** `ws.onmessage = (event) => { const alert: Alert = JSON.parse(event.data); dispatch(addAlert(alert)); };`
**File Location:** `frontend/src/hooks/useAlerts.ts`
**Purpose:** Handles incoming WebSocket messages by parsing them as Alert objects and dispatching to Redux.
**Why it exists:** When the backend broadcasts a new alert via WebSocket, this handler deserializes the JSON and dispatches it to the Redux store. Any component using `selectAlerts` instantly sees the new alert appear — no polling, no manual refresh. This is true real-time push-based updates.
**Inputs:** WebSocket MessageEvent with JSON string in `event.data`.
**Returns:** N/A (dispatches Redux action as side effect).
**Depends on:** TASK-125, TASK-113

---

#### TASK-127: useAlerts — cleanup on unmount
**Signature:** `return () => { ws.close(); };`
**File Location:** `frontend/src/hooks/useAlerts.ts`
**Purpose:** Closes the WebSocket connection when the component unmounts.
**Why it exists:** React's useEffect cleanup function runs when the component unmounts (or before re-running the effect). Closing the WebSocket prevents memory leaks and ghost connections. Without cleanup, every navigation away and back would create a new connection without closing the old one, eventually exhausting server resources.
**Inputs:** None.
**Returns:** N/A (cleanup function).
**Concepts to learn before implementing:**
- useEffect cleanup function → search "react useEffect cleanup"
- WebSocket close → search "javascript websocket close"
- Memory leak prevention in React → search "react memory leak useEffect"
**Depends on:** TASK-125

---

### FILE: `frontend/src/hooks/useLogs.ts`
**Purpose of this file:** Custom hook for fetching and managing paginated, filtered security logs. Encapsulates the API call, Redux dispatch, and reactive refetching.
**Week:** Week 6
**Layer:** Frontend

---

#### TASK-128: useLogs — fetch function
**Signature:** `const fetchLogs = async () => { const { data } = await logsApi.getLogs({ level, source, page, limit }); dispatch(setLogs(data)); };`
**File Location:** `frontend/src/hooks/useLogs.ts`
**Purpose:** Fetches logs from the API with current filters and pagination, then updates Redux state.
**Why it exists:** Encapsulates the API call, query parameter assembly, and Redux dispatch into one reusable function. Called on mount and whenever filters or page number change (via useEffect dependency tracking).
**Inputs:** Uses filter and pagination values from Redux state.
**Returns:** Promise (void).
**Depends on:** TASK-099, TASK-117

---

#### TASK-129: useLogs — pagination
**Signature:** `const goToPage = (p: number) => { dispatch(setPage(p)); };`
**File Location:** `frontend/src/hooks/useLogs.ts`
**Purpose:** Updates the current page in Redux, which triggers a re-fetch via useEffect.
**Why it exists:** Changing the page in Redux triggers a useEffect dependency, which calls fetchLogs with the new page number. This reactive pattern means the displayed data always matches the current page — change the state, data follows automatically.
**Inputs:** `p` (number) — the target page number.
**Returns:** void.
**Concepts to learn before implementing:**
- useEffect dependency array → search "react useEffect dependency array"
- Reactive data fetching → search "react refetch on state change"
**Depends on:** TASK-119, TASK-128

---

#### TASK-130: useLogs — filter params
**Signature:** `const updateFilters = (f: Partial<Filters>) => { dispatch(setFilters(f)); };`
**File Location:** `frontend/src/hooks/useLogs.ts`
**Purpose:** Updates filter values in Redux and resets to page 1 (handled inside setFilters reducer).
**Why it exists:** When a user selects "Critical" in the threat level dropdown, this function dispatches the filter change. The setFilters reducer resets the page to 1 (since the result set changed), and the useEffect re-fetches with the new filters.
**Inputs:** `f` (Partial\<Filters\>) — partial filter object, e.g., `{ level: 'critical' }`.
**Returns:** void.
**Depends on:** TASK-118, TASK-128

---

### FILE: `frontend/src/hooks/useStats.ts`
**Purpose of this file:** Custom hook for fetching dashboard statistics — log counts grouped by threat level.
**Week:** Week 6
**Layer:** Frontend

---

#### TASK-131: useStats — fetch counts by threat level
**Signature:** `export const useStats = () => { ... return { stats, loading, error }; };`
**File Location:** `frontend/src/hooks/useStats.ts`
**Purpose:** Fetches threat level distribution data for the dashboard stat cards and bar chart.
**Why it exists:** The dashboard displays stat cards ("Critical: 15", "High: 42") and a bar chart showing threat level distribution. This hook fetches the aggregated counts once on mount and provides loading/error states for the UI to render appropriate feedback (spinner while loading, error message on failure).
**Inputs:** None.
**Returns:** `{ stats: { low: number; medium: number; high: number; critical: number }; loading: boolean; error: string | null }`.
**Concepts to learn before implementing:**
- React data fetching patterns → search "react data fetching useEffect"
- Loading and error states → search "react loading error state pattern"
- Custom hook return objects → search "react custom hook return object"
**Depends on:** TASK-099

---

### FILE: `frontend/src/utils/threatColors.ts`
**Purpose of this file:** Utility functions mapping threat levels to consistent colors across the entire UI — badges, charts, stat cards, and alert cards.
**Week:** Week 6
**Layer:** Frontend

---

#### TASK-132: getThreatColor
**Signature:** `export const getThreatColor = (level: ThreatLevel): string => { ... };`
**File Location:** `frontend/src/utils/threatColors.ts`
**Purpose:** Returns a text/border color hex code for each threat level.
**Why it exists:** Threat level colors must be consistent across the entire application — badges in log tables, bar chart fills, stat card accents, alert card borders. A single utility function ensures one source of truth for the color palette. Example mapping: critical → `'#EF4444'` (red), high → `'#F97316'` (orange), medium → `'#EAB308'` (yellow), low → `'#22C55E'` (green).
**Inputs:** `level` (ThreatLevel) — one of the four threat levels.
**Returns:** `string` — hex color code.
**Depends on:** TASK-089

---

#### TASK-133: getThreatBgColor
**Signature:** `export const getThreatBgColor = (level: ThreatLevel): string => { ... };`
**File Location:** `frontend/src/utils/threatColors.ts`
**Purpose:** Returns a lighter/translucent background color for each threat level.
**Why it exists:** Background colors are softer versions of the text colors, used for badge backgrounds and card highlights. Having separate text and background color functions lets you create accessible contrast ratios — dark text on a light background of the same hue.
**Inputs:** `level` (ThreatLevel).
**Returns:** `string` — hex or rgba color for backgrounds.
**Depends on:** TASK-089

---

### FILE: `frontend/src/utils/formatDate.ts`
**Purpose of this file:** Date formatting utilities for consistent timestamp display throughout the application.
**Week:** Week 6
**Layer:** Frontend

---

#### TASK-134: formatRelative
**Signature:** `export const formatRelative = (date: string): string => { ... };`
**File Location:** `frontend/src/utils/formatDate.ts`
**Purpose:** Converts an ISO timestamp string to a human-friendly relative time like "2 minutes ago", "3 hours ago", "yesterday".
**Why it exists:** Relative times are more intuitive in real-time feeds ("2 min ago" is immediately understandable; "2024-01-15T14:23:00Z" requires mental math). Used in AlertFeed and LogTable for the createdAt column.
**Inputs:** `date` (string) — ISO 8601 timestamp from the API.
**Returns:** `string` — relative time string.
**Concepts to learn before implementing:**
- Date manipulation in JavaScript → search "javascript date formatting"
- Relative time calculation → search "javascript relative time ago function"
- Intl.RelativeTimeFormat API → search "intl relativetimeformat"
**Depends on:** None

---

#### TASK-135: formatAbsolute
**Signature:** `export const formatAbsolute = (date: string): string => { ... };`
**File Location:** `frontend/src/utils/formatDate.ts`
**Purpose:** Formats an ISO timestamp to a readable absolute format like "Jan 15, 2024 2:23 PM".
**Why it exists:** Used in log detail views and settings page where exact timestamps matter more than relative times. Uses `Intl.DateTimeFormat` for locale-aware formatting.
**Inputs:** `date` (string) — ISO 8601 timestamp.
**Returns:** `string` — formatted date string.
**Depends on:** None

---

### FILE: `frontend/src/components/ui/Badge.tsx`
**Purpose of this file:** Reusable badge component that displays a colored pill/tag based on threat level — part of the design system.
**Week:** Week 6
**Layer:** Frontend

---

#### TASK-136: Badge Component
**Signature:** `export const Badge: React.FC<{ level: ThreatLevel; label?: string }>`
**File Location:** `frontend/src/components/ui/Badge.tsx`
**Purpose:** Renders a small colored pill showing the threat level.
**Why it exists:** Badges appear in log tables (every row), alert cards (severity indicator), and stat cards. A reusable component ensures consistent styling everywhere. If the color scheme changes, you update one component. Uses `getThreatColor` and `getThreatBgColor` for dynamic styling based on the threat level.
**Inputs (Props):**
- `level` (ThreatLevel) — determines the color.
- `label` (string, optional) — custom text to display instead of the level name.
**Returns:** JSX — a styled `<span>` element with colored background and text.
**Concepts to learn before implementing:**
- React functional components with props → search "react functional component props"
- Conditional/dynamic CSS styling → search "react dynamic inline styles"
**Depends on:** TASK-089, TASK-132, TASK-133

---

### FILE: `frontend/src/components/ui/Button.tsx`
**Purpose of this file:** Reusable button component with variants (primary, secondary, danger) and loading state — part of the design system.
**Week:** Week 6
**Layer:** Frontend

---

#### TASK-137: Button Component
**Signature:** `export const Button: React.FC<ButtonProps>` where ButtonProps extends `ButtonHTMLAttributes<HTMLButtonElement>` with `{ variant?: 'primary' | 'secondary' | 'danger'; isLoading?: boolean }`
**File Location:** `frontend/src/components/ui/Button.tsx`
**Purpose:** A design-system button with consistent sizing, colors, hover/active states, and a loading spinner.
**Why it exists:** Every form and action in the app needs buttons. A shared component ensures consistent appearance and behavior. Extending `ButtonHTMLAttributes` means all native button props (onClick, disabled, type, etc.) are automatically available without manually forwarding each one.
**Inputs (Props):**
- `variant` ('primary' | 'secondary' | 'danger', default 'primary') — visual style.
- `isLoading` (boolean, default false) — shows spinner and disables click.
- All native HTML button attributes.
**Returns:** JSX — a styled `<button>` element.
**Concepts to learn before implementing:**
- React component composition → search "react component composition"
- Extending HTML attributes in TypeScript → search "react extend html attributes typescript"
- Design system components → search "design system react button"
**Depends on:** None

---

### FILE: `frontend/src/components/ui/Input.tsx`
**Purpose of this file:** Reusable form input component with label and error display.
**Week:** Week 6
**Layer:** Frontend

---

#### TASK-138: Input Component
**Signature:** `export const Input: React.FC<InputProps>` where InputProps extends `InputHTMLAttributes<HTMLInputElement>` with `{ label?: string; error?: string }`
**File Location:** `frontend/src/components/ui/Input.tsx`
**Purpose:** A styled text input with optional label above and error message below.
**Why it exists:** Login, register, and settings forms all need inputs with labels and validation error display. A shared component ensures consistent border styling, focus ring, padding, and error state appearance across all forms.
**Inputs (Props):**
- `label` (string, optional) — renders a `<label>` above the input.
- `error` (string, optional) — renders a red error message below the input.
- All native HTML input attributes (type, placeholder, value, onChange, etc.).
**Returns:** JSX — label + input + error wrapper.
**Depends on:** None

---

### FILE: `frontend/src/components/ui/Modal.tsx`
**Purpose of this file:** Reusable modal/dialog component with backdrop overlay.
**Week:** Week 6
**Layer:** Frontend

---

#### TASK-139: Modal Component
**Signature:** `export const Modal: React.FC<{ isOpen: boolean; onClose: () => void; title?: string; children: React.ReactNode }>`
**File Location:** `frontend/src/components/ui/Modal.tsx`
**Purpose:** Renders a centered overlay dialog with backdrop click-to-close behavior.
**Why it exists:** Modals are used for confirmation dialogs ("Are you sure you want to resolve this alert?") and detail views. Using React Portal (`ReactDOM.createPortal`) renders the modal outside the component tree, preventing z-index stacking issues. The `isOpen` prop controls visibility — when false, nothing renders.
**Inputs (Props):**
- `isOpen` (boolean) — whether the modal is visible.
- `onClose` (() => void) — callback when user clicks backdrop or close button.
- `title` (string, optional) — modal header text.
- `children` (React.ReactNode) — modal body content.
**Returns:** JSX — portal-rendered overlay with centered content box.
**Concepts to learn before implementing:**
- React Portal → search "react portal modal"
- Overlay/backdrop pattern → search "modal overlay css"
- React.ReactNode type → search "react children typescript ReactNode"
**Depends on:** None

---

### FILE: `frontend/src/components/ui/Spinner.tsx`
**Purpose of this file:** Loading spinner component for indicating async operations.
**Week:** Week 6
**Layer:** Frontend

---

#### TASK-140: Spinner Component
**Signature:** `export const Spinner: React.FC<{ size?: 'sm' | 'md' | 'lg' }>`
**File Location:** `frontend/src/components/ui/Spinner.tsx`
**Purpose:** Renders an animated spinning circle as a loading indicator.
**Why it exists:** Used inside Button (isLoading state), on pages during data fetch, and for WebSocket connection status. CSS keyframe animation provides smooth rotation without JavaScript. Size variants allow using the same spinner in different contexts (small inside a button, large on a full-page loading screen).
**Inputs (Props):**
- `size` ('sm' | 'md' | 'lg', default 'md') — controls width/height.
**Returns:** JSX — a `<div>` with CSS spin animation.
**Concepts to learn before implementing:**
- CSS keyframe animations → search "css spinner animation keyframes"
- CSS animation property → search "css animation rotate"
**Depends on:** None

---

### FILE: `frontend/src/components/ui/StatCard.tsx`
**Purpose of this file:** Dashboard stat card component showing a single metric with label and colored accent.
**Week:** Week 6
**Layer:** Frontend

---

#### TASK-141: StatCard Component
**Signature:** `export const StatCard: React.FC<{ label: string; value: number | string; color: string; icon?: React.ReactNode }>`
**File Location:** `frontend/src/components/ui/StatCard.tsx`
**Purpose:** Renders a card with a large value, descriptive label, and a colored accent indicator.
**Why it exists:** The dashboard displays four stat cards in a row — one per threat level ("Critical: 15", "High: 42", "Medium: 88", "Low: 210"). Each card uses the threat level's color for its accent bar/border. The component is generic enough to display any metric.
**Inputs (Props):**
- `label` (string) — metric name, e.g., "Critical Threats".
- `value` (number | string) — the metric value.
- `color` (string) — accent color (from getThreatColor).
- `icon` (React.ReactNode, optional) — optional icon element.
**Returns:** JSX — a styled card with accent and content.
**Depends on:** None

---

### FILE: `frontend/src/components/layout/AppShell.tsx`
**Purpose of this file:** The main application shell providing the persistent sidebar + topbar + content area layout for all authenticated pages.
**Week:** Week 6
**Layer:** Frontend

---

#### TASK-142: AppShell Component
**Signature:** `export const AppShell: React.FC<{ children: React.ReactNode }>`
**File Location:** `frontend/src/components/layout/AppShell.tsx`
**Purpose:** Renders a flex layout with Sidebar on the left, TopBar at the top, and a scrollable content area for page children.
**Why it exists:** Every authenticated page shares the same layout — sidebar for navigation, topbar for context, content area for the page. The AppShell wraps all protected routes (via ProtectedRoute in TASK-170) so you write the layout once. The content area scrolls independently while the sidebar and topbar stay fixed.
**Inputs (Props):**
- `children` (React.ReactNode) — the page content to render in the main area.
**Returns:** JSX — flex container with Sidebar, TopBar, and content area.
**Concepts to learn before implementing:**
- React layout components → search "react layout component pattern"
- CSS flexbox layout (sidebar + content) → search "css flexbox layout sidebar"
- React Router Outlet → search "react router outlet"
**Depends on:** TASK-143, TASK-145

---

### FILE: `frontend/src/components/layout/Sidebar.tsx`
**Purpose of this file:** Navigation sidebar with links to all main pages and a logout button.
**Week:** Week 6
**Layer:** Frontend

---

#### TASK-143: Sidebar — navLinks Map
**Signature:** `const navLinks = [{ path: '/dashboard', label: 'Dashboard', icon: ... }, { path: '/logs', label: 'Logs', icon: ... }, { path: '/alerts', label: 'Alerts', icon: ... }, { path: '/settings', label: 'Settings', icon: ... }];`
**File Location:** `frontend/src/components/layout/Sidebar.tsx`
**Purpose:** Defines the navigation links as a data array.
**Why it exists:** Data-driven navigation — adding a new page means adding one object to this array instead of writing new JSX. The active link is highlighted by comparing each link's `path` with the current URL from `useLocation()`. This pattern scales better than hardcoding each `<NavLink>`.
**Inputs:** N/A (constant array)
**Returns:** N/A (constant array)
**Concepts to learn before implementing:**
- Data-driven UI rendering → search "react data driven rendering"
- React Router useLocation → search "react router useLocation active link"
- NavLink active styling → search "react router NavLink active"
**Depends on:** None

---

#### TASK-144: Sidebar — logout Handler
**Signature:** `const handleLogout = () => { logout(); };`
**File Location:** `frontend/src/components/layout/Sidebar.tsx`
**Purpose:** Calls the useAuth logout function when the logout button/link is clicked.
**Why it exists:** The sidebar always shows a logout option at the bottom. Clicking it triggers the auth hook's logout, which clears Redux state, removes the token from localStorage, and navigates to /login.
**Inputs:** None (click event handler).
**Returns:** void.
**Depends on:** TASK-123

---

### FILE: `frontend/src/components/layout/TopBar.tsx`
**Purpose of this file:** Top navigation bar showing the current page title and WebSocket connection status.
**Week:** Week 6
**Layer:** Frontend

---

#### TASK-145: TopBar Component
**Signature:** `export const TopBar: React.FC<{ title?: string }>`
**File Location:** `frontend/src/components/layout/TopBar.tsx`
**Purpose:** Renders a horizontal bar at the top of the content area with the page title, user info, and a live WebSocket status badge.
**Why it exists:** Gives users context about which page they're on and whether the real-time alert feed is active. The WebSocket status badge (green dot = connected, red = disconnected) provides immediate feedback about the live feed status.
**Inputs (Props):**
- `title` (string, optional) — the current page name (e.g., "Dashboard", "Security Logs").
**Returns:** JSX — horizontal bar with title, user info, and status indicators.
**Depends on:** TASK-136

---

### FILE: `frontend/src/components/charts/ThreatLevelBar.tsx`
**Purpose of this file:** Bar chart component visualizing the distribution of logs across threat levels.
**Week:** Week 6
**Layer:** Frontend

---

#### TASK-146: ThreatLevelBar — data transform function
**Signature:** `const chartData = transformStats(stats);`
**File Location:** `frontend/src/components/charts/ThreatLevelBar.tsx`
**Purpose:** Converts the stats object from useStats into the array format Recharts expects.
**Why it exists:** Recharts expects data as an array of objects like `[{ name: 'Low', count: 100, fill: '#22C55E' }, ...]`. The API returns an object like `{ low: 100, medium: 50, ... }`. This transform function bridges the gap and adds colors from `getThreatColor`.
**Inputs:** `stats` — object with `{ low: number, medium: number, high: number, critical: number }`.
**Returns:** Array of `{ name: string; count: number; fill: string }` objects.
**Depends on:** TASK-131, TASK-132

---

#### TASK-147: ThreatLevelBar — render function
**Signature:** Component return JSX with `<ResponsiveContainer>`, `<BarChart>`, `<Bar>`, `<XAxis>`, `<YAxis>`, `<Tooltip>`
**File Location:** `frontend/src/components/charts/ThreatLevelBar.tsx`
**Purpose:** Renders the bar chart with styled axes, tooltips, and responsive sizing.
**Why it exists:** The bar chart is a key dashboard visualization — analysts can instantly see which threat levels are most common. ResponsiveContainer ensures the chart resizes with its container.
**Inputs:** Chart data from the transform function.
**Returns:** JSX — Recharts BarChart inside a ResponsiveContainer.
**Concepts to learn before implementing:**
- Recharts BarChart → search "recharts barchart tutorial"
- ResponsiveContainer → search "recharts responsive container"
- Recharts customization (colors, labels) → search "recharts customize bar chart"
**Depends on:** TASK-146

---

### FILE: `frontend/src/components/charts/LogTimelineLine.tsx`
**Purpose of this file:** Line chart showing log ingestion volume over time — helps analysts spot activity spikes.
**Week:** Week 6
**Layer:** Frontend

---

#### TASK-148: LogTimelineLine Component
**Signature:** `export const LogTimelineLine: React.FC<{ logs: SecurityLog[] }>`
**File Location:** `frontend/src/components/charts/LogTimelineLine.tsx`
**Purpose:** Aggregates logs by time bucket (hour or day) and renders a line chart showing activity trends over time.
**Why it exists:** Activity spikes often indicate ongoing attacks. A line chart showing "logs per hour" lets analysts visually detect anomalies — a sudden spike from 10 logs/hour to 500 logs/hour suggests something is happening. The chart aggregates raw log timestamps into time buckets for visualization.
**Inputs (Props):**
- `logs` (SecurityLog[]) — the raw log data to aggregate and chart.
**Returns:** JSX — Recharts LineChart with time on X-axis and count on Y-axis.
**Concepts to learn before implementing:**
- Recharts LineChart → search "recharts linechart"
- Time-series data aggregation → search "javascript group by date"
- Date bucketing (group by hour/day) → search "javascript date floor hour"
**Depends on:** TASK-091

---

### FILE: `frontend/src/components/charts/SourcePie.tsx`
**Purpose of this file:** Pie chart showing the distribution of logs across different sources.
**Week:** Week 6
**Layer:** Frontend

---

#### TASK-149: SourcePie Component
**Signature:** `export const SourcePie: React.FC<{ logs: SecurityLog[] }>`
**File Location:** `frontend/src/components/charts/SourcePie.tsx`
**Purpose:** Groups logs by source (firewall, IDS, auth_service, endpoint) and renders a pie chart showing the proportional distribution.
**Why it exists:** Understanding which sources generate the most logs helps with capacity planning and identifies which systems are most active or under most pressure. The pie chart provides an instant visual breakdown.
**Inputs (Props):**
- `logs` (SecurityLog[]) — the raw log data to aggregate.
**Returns:** JSX — Recharts PieChart with labeled slices.
**Concepts to learn before implementing:**
- Recharts PieChart → search "recharts piechart tutorial"
- Data aggregation with reduce → search "javascript reduce group by"
**Depends on:** TASK-091

---

### FILE: `frontend/src/components/alerts/AlertCard.tsx`
**Purpose of this file:** Card component displaying a single alert with severity badge, rule name, timestamps, resolve action, and an expandable AI summary section.
**Week:** Week 6
**Layer:** Frontend

---

#### TASK-150: AlertCard — expandAISummary toggle function
**Signature:** `const [expanded, setExpanded] = useState(false); const toggleExpand = () => setExpanded(!expanded);`
**File Location:** `frontend/src/components/alerts/AlertCard.tsx`
**Purpose:** Manages the expand/collapse state for the AI summary section.
**Why it exists:** AI summaries can be several paragraphs long. Showing them collapsed by default keeps the alert feed compact and scannable. Users click "Show AI Analysis" to expand and read the full analysis. The toggle pattern using `useState` is a fundamental React interaction pattern.
**Inputs:** None (internal state).
**Returns:** N/A (state management).
**Concepts to learn before implementing:**
- React useState → search "react useState tutorial"
- Toggle pattern → search "react toggle state"
**Depends on:** None

---

#### TASK-151: AlertCard — render function
**Signature:** `export const AlertCard: React.FC<{ alert: Alert; onResolve?: (id: string) => void }>`
**File Location:** `frontend/src/components/alerts/AlertCard.tsx`
**Purpose:** Renders the complete alert card with all information and actions.
**Why it exists:** Each alert card displays: severity badge (colored by level), rule name (what was detected), creation time (how recent), a "Resolve" button (if not yet resolved), and a collapsible AI summary section. The `onResolve` callback is optional — the card can be used in read-only mode (dashboard feed) or interactive mode (alerts page).
**Inputs (Props):**
- `alert` (Alert) — the alert data to display.
- `onResolve` ((id: string) => void, optional) — callback when "Resolve" is clicked.
**Returns:** JSX — styled card with badge, text, actions, and collapsible section.
**Depends on:** TASK-150, TASK-136, TASK-134, TASK-092

---

### FILE: `frontend/src/components/alerts/AlertFeed.tsx`
**Purpose of this file:** Real-time scrolling feed of alert cards — the live feed on the dashboard.
**Week:** Week 6
**Layer:** Frontend

---

#### TASK-152: AlertFeed Component
**Signature:** `export const AlertFeed: React.FC`
**File Location:** `frontend/src/components/alerts/AlertFeed.tsx`
**Purpose:** Maps over the alerts array from Redux and renders an AlertCard for each.
**Why it exists:** This is the live feed on the dashboard — new alerts appear at the top as they arrive via WebSocket. It's a simple mapping component: `alerts.map(a => <AlertCard alert={a} />)`. The key prop uses `alert.id` for efficient React reconciliation. Shows "No active alerts" when the array is empty.
**Inputs:** None (reads from Redux via `useSelector(selectAlerts)`).
**Returns:** JSX — a scrollable container of AlertCard components.
**Concepts to learn before implementing:**
- React list rendering with key prop → search "react list rendering key"
- Array.map for component lists → search "react array map components"
- Empty state handling → search "react empty state component"
**Depends on:** TASK-151, TASK-115

---

### FILE: `frontend/src/components/logs/LogFilters.tsx`
**Purpose of this file:** Filter controls for the logs table — dropdowns for selecting threat level and source filters.
**Week:** Week 6
**Layer:** Frontend

---

#### TASK-153: LogFilters — onFilterChange handler
**Signature:** `const handleFilterChange = (field: string, value: string) => { updateFilters({ [field]: value }); };`
**File Location:** `frontend/src/components/logs/LogFilters.tsx`
**Purpose:** A generic handler for all filter field changes.
**Why it exists:** Using computed property names (`[field]: value`) means one handler function works for both the threat level and source dropdowns. When the user selects "Critical" from the level dropdown, it calls `handleFilterChange('level', 'critical')`, which dispatches `setFilters({ level: 'critical' })`.
**Inputs:** `field` (string — 'level' or 'source'), `value` (string — the selected filter value).
**Returns:** void.
**Concepts to learn before implementing:**
- Computed property names in JavaScript → search "javascript computed property names"
- Controlled form components in React → search "react controlled components"
**Depends on:** TASK-130

---

#### TASK-154: LogFilters — filter state
**Signature:** Reads current filter values from Redux via `useSelector(selectFilters)`.
**File Location:** `frontend/src/components/logs/LogFilters.tsx`
**Purpose:** Keeps the filter dropdowns in sync with the Redux state.
**Why it exists:** The filter dropdowns are controlled components — their displayed values are driven by Redux state, not internal component state. This ensures the dropdowns always show the active filters, even if filters are changed programmatically (e.g., a "Clear Filters" button).
**Inputs:** Redux state via `selectFilters`.
**Returns:** Current filter values for rendering dropdown selections.
**Depends on:** TASK-121

---

### FILE: `frontend/src/components/logs/LogRow.tsx`
**Purpose of this file:** A single row in the log table with an expandable detail panel.
**Week:** Week 6
**Layer:** Frontend

---

#### TASK-155: LogRow — expand toggle
**Signature:** `const [expanded, setExpanded] = useState(false);`
**File Location:** `frontend/src/components/logs/LogRow.tsx`
**Purpose:** Manages the expand/collapse state for the log detail panel.
**Why it exists:** Each log row shows a summary (source, truncated message, threat level, time). Clicking the row toggles an expandable detail panel showing the full message, IP address, user agent, and AI summary. This keeps the table compact while providing detail on demand.
**Inputs:** None (internal state).
**Returns:** N/A (state management).
**Depends on:** None

---

#### TASK-156: LogRow — render function
**Signature:** `export const LogRow: React.FC<{ log: SecurityLog }>`
**File Location:** `frontend/src/components/logs/LogRow.tsx`
**Purpose:** Renders a table row with summary data and a conditionally rendered detail panel.
**Why it exists:** Each row shows: source icon/label, truncated message (first 80 chars), threat level badge, relative creation time. When expanded, a detail panel slides open below with: full message text, IP address, user agent string, and AI summary (if available). The AI summary section shows "No AI analysis" for low/medium threats.
**Inputs (Props):**
- `log` (SecurityLog) — the log entry to display.
**Returns:** JSX — table row + optional expanded detail panel.
**Depends on:** TASK-155, TASK-136, TASK-134, TASK-091

---

### FILE: `frontend/src/components/logs/LogTable.tsx`
**Purpose of this file:** The main log table component with sortable column headers.
**Week:** Week 6
**Layer:** Frontend

---

#### TASK-157: LogTable — sort handler
**Signature:** `const [sortField, setSortField] = useState('createdAt'); const [sortDir, setSortDir] = useState<'asc'|'desc'>('desc');`
**File Location:** `frontend/src/components/logs/LogTable.tsx`
**Purpose:** Manages the current sort field and direction, toggling when column headers are clicked.
**Why it exists:** Clicking a column header sorts the table by that column. Clicking the same header again reverses the sort direction. This provides client-side sorting for the current page of results.
**Inputs:** Column header click events.
**Returns:** N/A (state management).
**Concepts to learn before implementing:**
- Client-side table sorting → search "react table sorting"
- Sort direction toggle logic → search "toggle sort order asc desc"
**Depends on:** None

---

#### TASK-158: LogTable — render function
**Signature:** `export const LogTable: React.FC`
**File Location:** `frontend/src/components/logs/LogTable.tsx`
**Purpose:** Renders the complete log table with column headers (Source, Message, Threat Level, Time) and maps logs to LogRow components.
**Why it exists:** This component assembles the table structure: thead with clickable column headers (showing sort direction indicators), tbody mapping logs to LogRow components, and a "No logs found" empty state message when the array is empty.
**Inputs:** Reads logs from Redux via `useSelector(selectLogs)`.
**Returns:** JSX — complete HTML table with headers and rows.
**Depends on:** TASK-157, TASK-156, TASK-120

---

### FILE: `frontend/src/pages/LoginPage.tsx`
**Purpose of this file:** The login form page — the first page unauthenticated users see.
**Week:** Week 6
**Layer:** Frontend

---

#### TASK-159: LoginPage — form state
**Signature:** `const [email, setEmail] = useState(''); const [password, setPassword] = useState(''); const [error, setError] = useState('');`
**File Location:** `frontend/src/pages/LoginPage.tsx`
**Purpose:** Manages the login form's controlled input state and error message.
**Why it exists:** React controlled components require state for each input field. Each keystroke updates state via onChange, and the input's value is always driven by state. The error state holds any login failure message to display to the user.
**Inputs:** N/A (state initialization).
**Returns:** N/A (state management).
**Depends on:** None

---

#### TASK-160: LoginPage — handleSubmit
**Signature:** `const handleSubmit = async (e: FormEvent) => { e.preventDefault(); try { await login(email, password); navigate('/dashboard'); } catch { setError('Invalid credentials'); } };`
**File Location:** `frontend/src/pages/LoginPage.tsx`
**Purpose:** Handles form submission — authenticates the user and redirects on success.
**Why it exists:** Calls the `useAuth` hook's login function, which handles the API call and Redux dispatch. On success, navigates to the dashboard. On failure, sets the error state to display a message. `e.preventDefault()` stops the browser's default form submission (which would cause a page reload).
**Inputs:** FormEvent (from the form's onSubmit).
**Returns:** void.
**Concepts to learn before implementing:**
- React form submission → search "react form handleSubmit"
- e.preventDefault() → search "javascript preventDefault form"
- React Router programmatic navigation → search "react router navigate programmatic"
- try/catch for async errors → search "javascript try catch async await"
**Depends on:** TASK-159, TASK-122, TASK-138, TASK-137

---

### FILE: `frontend/src/pages/RegisterPage.tsx`
**Purpose of this file:** The registration form page for new user account creation.
**Week:** Week 6
**Layer:** Frontend

---

#### TASK-161: RegisterPage — form state
**Signature:** `const [email, setEmail] = useState(''); const [password, setPassword] = useState(''); const [role, setRole] = useState<Role>('analyst');`
**File Location:** `frontend/src/pages/RegisterPage.tsx`
**Purpose:** Manages the registration form state including role selection.
**Why it exists:** Same pattern as LoginPage, with the addition of a role selector. The role defaults to 'analyst' (the most common role) and can be changed to 'admin' via a dropdown.
**Inputs:** N/A (state initialization).
**Returns:** N/A (state management).
**Depends on:** TASK-087

---

#### TASK-162: RegisterPage — handleSubmit
**Signature:** `const handleSubmit = async (e: FormEvent) => { e.preventDefault(); await authApi.register(email, password, role); navigate('/login'); };`
**File Location:** `frontend/src/pages/RegisterPage.tsx`
**Purpose:** Handles registration form submission — creates the account and redirects to login.
**Why it exists:** After successful registration, the user is redirected to the login page to authenticate with their new credentials (rather than auto-logging in, which is a design choice for security).
**Inputs:** FormEvent.
**Returns:** void.
**Depends on:** TASK-161, TASK-097, TASK-138, TASK-137

---

### FILE: `frontend/src/pages/DashboardPage.tsx`
**Purpose of this file:** The main dashboard page — the command center showing live alert feed, stat cards, and charts.
**Week:** Week 6
**Layer:** Frontend

---

#### TASK-163: DashboardPage — useEffect for stats
**Signature:** `useEffect(() => { fetchStats(); }, []);`
**File Location:** `frontend/src/pages/DashboardPage.tsx`
**Purpose:** Fetches threat level statistics when the dashboard page mounts.
**Why it exists:** Dashboard stats need to load when the page first renders. The empty dependency array `[]` ensures fetchStats runs exactly once on mount, not on every re-render. The stats data populates the StatCards and ThreatLevelBar chart.
**Inputs:** None (uses useStats hook internally).
**Returns:** N/A (side effect).
**Depends on:** TASK-131

---

#### TASK-164: DashboardPage — render sections
**Signature:** Component return JSX with a grid layout containing StatCards, AlertFeed, and charts.
**File Location:** `frontend/src/pages/DashboardPage.tsx`
**Purpose:** Assembles all dashboard components into a cohesive layout.
**Why it exists:** The dashboard is the command center — analysts see everything at a glance. The layout uses CSS Grid for responsive positioning: a row of 4 StatCards across the top, AlertFeed (live) on one side, and the three charts (ThreatLevelBar, LogTimelineLine, SourcePie) arranged in a grid below.
**Inputs:** Stats data from useStats, alerts from Redux.
**Returns:** JSX — responsive grid with all dashboard components.
**Concepts to learn before implementing:**
- CSS Grid layout → search "css grid responsive dashboard"
- Dashboard UI design patterns → search "dashboard layout best practices"
**Depends on:** TASK-163, TASK-141, TASK-152, TASK-147, TASK-148, TASK-149

---

### FILE: `frontend/src/pages/LogsPage.tsx`
**Purpose of this file:** Page displaying the filterable, paginated security log table.
**Week:** Week 6
**Layer:** Frontend

---

#### TASK-165: LogsPage — pagination handler
**Signature:** Pagination controls (Previous/Next buttons, page number display) that dispatch setPage actions.
**File Location:** `frontend/src/pages/LogsPage.tsx`
**Purpose:** Assembles the pagination UI and wires it to the useLogs hook.
**Why it exists:** The LogsPage is the container that composes LogFilters + LogTable + pagination controls. The pagination handler calculates total pages from `total` and `limit`, renders page numbers, and calls `goToPage()` when clicked.
**Inputs:** Uses useLogs hook values (page, total, limit, goToPage).
**Returns:** JSX — pagination controls.
**Depends on:** TASK-129

---

#### TASK-166: LogsPage — filter handler
**Signature:** Passes `updateFilters` from useLogs to the LogFilters component as a prop.
**File Location:** `frontend/src/pages/LogsPage.tsx`
**Purpose:** Wires the LogFilters component to the useLogs data fetching hook.
**Why it exists:** The LogsPage orchestrates the connection between the filter UI (LogFilters), the data display (LogTable), and the data management hook (useLogs). It passes `updateFilters` down as a prop so LogFilters can trigger refetches.
**Inputs:** Uses useLogs hook values.
**Returns:** JSX — page with filters + table + pagination.
**Depends on:** TASK-130, TASK-153, TASK-158

---

### FILE: `frontend/src/pages/AlertsPage.tsx`
**Purpose of this file:** Page displaying all alerts with the ability to resolve them.
**Week:** Week 6
**Layer:** Frontend

---

#### TASK-167: AlertsPage — handleResolve function
**Signature:** `const handleResolve = async (id: string) => { await alertsApi.resolveAlert(id); refetchAlerts(); };`
**File Location:** `frontend/src/pages/AlertsPage.tsx`
**Purpose:** Handles the "Resolve" button click on an alert card.
**Why it exists:** When an analyst clicks "Resolve", this function calls the API endpoint, then refetches the alerts list to update the UI (removing the alert from the unresolved list). The refetch ensures the UI accurately reflects the current server state.
**Inputs:** `id` (string) — the UUID of the alert to resolve.
**Returns:** void.
**Depends on:** TASK-103, TASK-102, TASK-151

---

### FILE: `frontend/src/pages/SettingsPage.tsx`
**Purpose of this file:** User settings/profile page.
**Week:** Week 6
**Layer:** Frontend

---

#### TASK-168: SettingsPage — handleUpdateProfile function
**Signature:** `const handleUpdateProfile = async () => { ... };`
**File Location:** `frontend/src/pages/SettingsPage.tsx`
**Purpose:** Placeholder for future profile update functionality.
**Why it exists:** For the initial version, this page displays the current user's information (email, role, account creation date) from the `useAuth` hook. The handleUpdateProfile function is a placeholder for future enhancements like password change or notification preferences.
**Inputs:** Form state (future implementation).
**Returns:** void.
**Depends on:** TASK-110, TASK-138, TASK-137

---

### FILE: `frontend/src/App.tsx`
**Purpose of this file:** Root application component — sets up React Router with route definitions and the ProtectedRoute guard for authenticated pages.
**Week:** Week 6
**Layer:** Frontend

---

#### TASK-169: App — router setup
**Signature:** `const router = createBrowserRouter([...]);` or `<BrowserRouter><Routes>...</Routes></BrowserRouter>`
**File Location:** `frontend/src/App.tsx`
**Purpose:** Initializes the React Router for client-side navigation.
**Why it exists:** React Router v6 maps URL paths to page components, enabling client-side navigation without full page reloads. The router handles forward/back browser buttons, direct URL access, and programmatic navigation.
**Inputs:** Route configuration array.
**Returns:** Router instance or JSX router component.
**Concepts to learn before implementing:**
- React Router v6 setup → search "react router v6 tutorial"
- createBrowserRouter → search "react router createBrowserRouter"
- Client-side routing → search "single page app routing"
**Depends on:** None

---

#### TASK-170: App — ProtectedRoute component
**Signature:** `const ProtectedRoute = ({ children }: { children: React.ReactNode }) => { if (!isAuthenticated) return <Navigate to="/login" />; return <AppShell>{children}</AppShell>; };`
**File Location:** `frontend/src/App.tsx`
**Purpose:** Route guard that redirects unauthenticated users and wraps authenticated pages in the AppShell layout.
**Why it exists:** Serves two purposes: (1) authentication check — if the user isn't logged in, redirect to /login, (2) layout wrapper — authenticated pages render inside the AppShell (with sidebar and topbar) automatically. This dual responsibility keeps route definitions clean.
**Inputs:**
- `children` (React.ReactNode) — the protected page component.
**Returns:** JSX — either a `<Navigate>` redirect or `<AppShell>` with children.
**Concepts to learn before implementing:**
- Route guards in React → search "react protected route"
- React Router Navigate component → search "react router navigate redirect"
- Conditional rendering → search "react conditional rendering"
**Depends on:** TASK-124, TASK-142

---

#### TASK-171: App — route definitions
**Signature:** Route configuration mapping paths to page components.
**File Location:** `frontend/src/App.tsx`
**Purpose:** Defines all application routes and their associated page components.
**Why it exists:** Maps every URL to its page component. Public routes (login, register) render without the AppShell. Protected routes render inside ProtectedRoute (which adds AppShell and auth checks). The root path `/` redirects to `/dashboard`.
**Routes:**
- `/login` → LoginPage (public)
- `/register` → RegisterPage (public)
- `/dashboard` → ProtectedRoute → DashboardPage
- `/logs` → ProtectedRoute → LogsPage
- `/alerts` → ProtectedRoute → AlertsPage
- `/settings` → ProtectedRoute → SettingsPage
- `/` → redirect to `/dashboard`
**Depends on:** TASK-169, TASK-170, TASK-160, TASK-162, TASK-164, TASK-166, TASK-167, TASK-168

---

### FILE: `frontend/src/main.tsx`
**Purpose of this file:** Application entry point — bootstraps React with Redux Provider and renders the root component into the DOM.
**Week:** Week 6
**Layer:** Frontend

---

#### TASK-172: main — Provider setup
**Signature:** `<Provider store={store}>`
**File Location:** `frontend/src/main.tsx`
**Purpose:** Wraps the entire application in the Redux Provider component.
**Why it exists:** The Redux Provider makes the store available to all React components via React's context API. Without it, `useSelector` and `useDispatch` hooks would throw errors because they wouldn't be able to find the store.
**Inputs:** `store` — the configured Redux store from TASK-104.
**Returns:** JSX wrapper.
**Concepts to learn before implementing:**
- Redux Provider component → search "redux provider react"
- React Context API → search "react context api"
**Depends on:** TASK-104

---

#### TASK-173: main — StrictMode
**Signature:** `<React.StrictMode>`
**File Location:** `frontend/src/main.tsx`
**Purpose:** Enables additional development-mode checks for finding bugs early.
**Why it exists:** StrictMode activates extra development checks: double-rendering components to detect impure renders, warning about deprecated APIs, and detecting unexpected side effects. These checks only run in development and are stripped in production builds. Double rendering helps catch bugs where components rely on side effects in the render phase.
**Inputs:** None (wrapper component).
**Returns:** JSX wrapper.
**Concepts to learn before implementing:**
- React StrictMode → search "react strict mode explained"
- Why React renders twice in dev → search "react strict mode double render"
**Depends on:** None

---

#### TASK-174: main — render call
**Signature:** `ReactDOM.createRoot(document.getElementById('root')!).render(<StrictMode><Provider store={store}><App /></Provider></StrictMode>);`
**File Location:** `frontend/src/main.tsx`
**Purpose:** Mounts the entire React application tree into the DOM.
**Why it exists:** This is the bootstrap code that connects the React virtual DOM to the browser's real DOM. `createRoot` is React 18's rendering API (replacing the older `ReactDOM.render`). The `!` non-null assertion tells TypeScript that `document.getElementById('root')` will never be null (the element exists in index.html).
**Inputs:** The root DOM element from index.html, the component tree.
**Returns:** void (initiates React rendering).
**Concepts to learn before implementing:**
- React 18 createRoot → search "react 18 createRoot"
- DOM mounting → search "reactdom createRoot tutorial"
- Non-null assertion in TypeScript → search "typescript non-null assertion"
**Depends on:** TASK-172, TASK-173, TASK-169

---

### FILE: `frontend/vite.config.ts`
**Purpose of this file:** Vite bundler configuration — sets up the development server with API proxy rules.
**Week:** Week 6
**Layer:** Frontend

---

#### TASK-175: defineConfig — proxy rule
**Signature:** `export default defineConfig({ server: { proxy: { '/api': { target: 'http://localhost:8080', changeOrigin: true }, '/ws': { target: 'ws://localhost:8080', ws: true } } } });`
**File Location:** `frontend/vite.config.ts`
**Purpose:** Configures Vite's dev server to proxy `/api` and `/ws` requests to the Go backend.
**Why it exists:** In development, the React dev server (port 5173) and Go API server (port 8080) run on different ports — different origins. Without a proxy, API calls would fail due to CORS restrictions. The proxy forwards `/api/*` and `/ws` requests from the Vite dev server to the Go backend, making it look like everything is on the same origin. In production, both are served from the same origin via Nginx, so the proxy isn't needed.
**Inputs:** Proxy configuration object.
**Returns:** Vite configuration object.
**Concepts to learn before implementing:**
- Vite proxy configuration → search "vite proxy config"
- CORS in development → search "cors development proxy"
- WebSocket proxy → search "vite websocket proxy"
- Vite config file → search "vite config tutorial"
**Depends on:** None


---
---

## 🐳 WEEK 7 — Containerization, Orchestration & CI/CD

### FILE: `Dockerfile.backend`
**Purpose of this file:** Multi-stage Docker build for the Go backend — compiles the source code in a full SDK image and copies only the tiny binary to a minimal production image.
**Week:** Week 7
**Layer:** DevOps

---

#### TASK-176: Stage 1 — Builder
**Signature:** `FROM golang:1.22-alpine AS builder` → `WORKDIR /app` → `COPY go.mod go.sum ./` → `RUN go mod download` → `COPY . .` → `RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server`
**File Location:** `Dockerfile.backend`
**Purpose:** Compiles the Go source code into a static binary inside a full Go SDK image.
**Why it exists:** Multi-stage builds separate the build environment from the runtime environment. The builder stage contains the Go compiler (~1GB), all source code, and all dependencies — but none of that is needed at runtime. By building here and copying only the binary to the next stage, the final image shrinks from ~1GB to ~10MB. `CGO_ENABLED=0` produces a statically linked binary with no C library dependencies, and `GOOS=linux` ensures the binary runs on Linux containers regardless of your development OS.
**Inputs:** Go source code, go.mod, go.sum.
**Returns:** A compiled binary at `/app/server`.
**Concepts to learn before implementing:**
- Docker multi-stage builds → search "docker multi-stage build tutorial"
- Go static binary compilation (CGO_ENABLED=0) → search "golang static binary cgo"
- Alpine Linux images → search "docker alpine image"
- Docker COPY and layer caching → search "docker layer caching optimization"
- Copying go.mod first for cache optimization → search "docker go mod download cache"
**Depends on:** TASK-022

---

#### TASK-177: Stage 2 — Runner
**Signature:** `FROM alpine:3.19` → `RUN apk add --no-cache ca-certificates` → `COPY --from=builder /app/server /server` → `EXPOSE 8080` → `CMD ["/server"]`
**File Location:** `Dockerfile.backend`
**Purpose:** Creates the final minimal production image with just the compiled binary and CA certificates.
**Why it exists:** The runner stage uses a bare Alpine image (~5MB) with only the binary and CA certificates (needed for HTTPS calls to the Claude AI API and any other TLS connections). No Go compiler, no source code, no build tools. This minimal image reduces attack surface (fewer packages = fewer potential vulnerabilities) and speeds up deployment (smaller image = faster pulls).
**Inputs:** The compiled binary from the builder stage.
**Returns:** A runnable Docker image.
**Concepts to learn before implementing:**
- Minimal Docker images → search "docker minimal production image"
- ca-certificates for HTTPS → search "docker alpine ca-certificates"
- EXPOSE vs port publishing → search "docker EXPOSE vs -p flag"
- CMD vs ENTRYPOINT → search "dockerfile CMD vs ENTRYPOINT"
- `COPY --from=builder` for cross-stage copying → search "docker multi-stage copy from"
**Depends on:** TASK-176

---

### FILE: `Dockerfile.frontend`
**Purpose of this file:** Multi-stage Docker build for the React frontend — builds optimized static assets and serves them with Nginx.
**Week:** Week 7
**Layer:** DevOps

---

#### TASK-178: Stage 1 — Builder
**Signature:** `FROM node:20-alpine AS builder` → `WORKDIR /app` → `COPY package*.json ./` → `RUN npm ci` → `COPY . .` → `RUN npm run build`
**File Location:** `Dockerfile.frontend`
**Purpose:** Installs dependencies and builds the React app into optimized static files (HTML, CSS, JS).
**Why it exists:** `npm ci` (clean install) is used instead of `npm install` because it installs exact versions from package-lock.json, ensuring reproducible builds (the same dependencies every time). The `npm run build` command runs Vite's production build, which minifies JavaScript, tree-shakes unused code, and generates optimized bundles in a `dist/` directory.
**Inputs:** package.json, package-lock.json, source code.
**Returns:** Static build artifacts in `/app/dist/`.
**Concepts to learn before implementing:**
- npm ci vs npm install → search "npm ci vs npm install"
- React/Vite production build → search "vite build production"
- Node Alpine image → search "docker node alpine"
- Package-lock.json and reproducible builds → search "package-lock.json purpose"
**Depends on:** TASK-174

---

#### TASK-179: Stage 2 — Nginx Runner
**Signature:** `FROM nginx:alpine` → `COPY --from=builder /app/dist /usr/share/nginx/html` → `COPY nginx.conf /etc/nginx/conf.d/default.conf` → `EXPOSE 80`
**File Location:** `Dockerfile.frontend`
**Purpose:** Serves the built static files using the Nginx web server with SPA-friendly configuration.
**Why it exists:** React apps compile down to static HTML/CSS/JS files that need a web server. Nginx is the industry standard for static file serving — it's fast, lightweight, and highly configurable. The custom `nginx.conf` configures: (1) `try_files $uri /index.html` — serves index.html for all routes so React Router handles client-side routing, (2) caching headers for static assets, (3) optionally, reverse proxy for `/api` requests.
**Inputs:** Built static files from builder stage, nginx.conf.
**Returns:** A runnable Docker image serving the frontend on port 80.
**Concepts to learn before implementing:**
- Nginx for SPAs → search "nginx react single page app"
- SPA history fallback (try_files) → search "nginx try_files spa"
- Nginx reverse proxy → search "nginx reverse proxy configuration"
- Nginx configuration basics → search "nginx config tutorial"
**Depends on:** TASK-178

---

### FILE: `docker-compose.yml`
**Purpose of this file:** Orchestrates all services for local development — defines the entire stack (database, cache, message queue, backend services, frontend) with one `docker compose up` command.
**Week:** Week 7
**Layer:** DevOps

---

#### TASK-180: postgres Service
**Signature:** Service block: `image: postgres:16-alpine`, `ports: 5432:5432`, `volumes: pgdata:/var/lib/postgresql/data`, `environment: POSTGRES_USER, POSTGRES_PASSWORD, POSTGRES_DB`, `healthcheck: pg_isready`
**File Location:** `docker-compose.yml`
**Purpose:** Runs PostgreSQL database for storing users, logs, and alerts.
**Why it exists:** PostgreSQL is the persistent data store. The named volume `pgdata` ensures data survives container restarts (without it, every `docker compose down` would wipe all data). The healthcheck (`pg_isready`) ensures dependent services (backend-api) don't start until Postgres is ready to accept connections — preventing startup race conditions.
**Inputs:** Environment variables for database credentials.
**Returns:** PostgreSQL available on port 5432.
**Concepts to learn before implementing:**
- docker-compose service definition → search "docker compose service"
- Named volumes for persistence → search "docker compose named volumes"
- Healthchecks → search "docker compose healthcheck"
- PostgreSQL Docker image → search "postgres docker image"
**Depends on:** TASK-024, TASK-026, TASK-028

---

#### TASK-181: redis Service
**Signature:** Service block: `image: redis:7-alpine`, `ports: 6379:6379`
**File Location:** `docker-compose.yml`
**Purpose:** Runs Redis for the rate limiter's sliding window counters.
**Why it exists:** Redis backs the rate limiter. The Alpine image keeps it lightweight (~30MB). No persistence volume is needed — rate limit counters are ephemeral by design (they reset on restart, which is acceptable).
**Inputs:** None (uses default configuration).
**Returns:** Redis available on port 6379.
**Depends on:** TASK-049

---

#### TASK-182: zookeeper Service
**Signature:** Service block: `image: confluentinc/cp-zookeeper:7.5.0`, `ports: 2181:2181`, `environment: ZOOKEEPER_CLIENT_PORT=2181`
**File Location:** `docker-compose.yml`
**Purpose:** Runs ZooKeeper for Kafka cluster coordination.
**Why it exists:** Kafka (in traditional mode) requires ZooKeeper for broker registration, topic metadata management, and leader election. ZooKeeper must start before Kafka. In newer Kafka versions, KRaft mode eliminates this dependency, but the traditional ZooKeeper-based setup is more widely documented and supported.
**Inputs:** Environment variable for client port.
**Returns:** ZooKeeper available on port 2181.
**Concepts to learn before implementing:**
- ZooKeeper's role in Kafka → search "kafka zookeeper explained"
- Confluent Docker images → search "confluent kafka docker images"
**Depends on:** None

---

#### TASK-183: kafka Service
**Signature:** Service block: `image: confluentinc/cp-kafka:7.5.0`, `ports: 9092:9092`, `depends_on: zookeeper`, `environment: KAFKA_BROKER_ID, KAFKA_ZOOKEEPER_CONNECT, KAFKA_ADVERTISED_LISTENERS, KAFKA_AUTO_CREATE_TOPICS_ENABLE`
**File Location:** `docker-compose.yml`
**Purpose:** Runs the Kafka message broker for the async log processing pipeline.
**Why it exists:** Kafka is the message broker between the API server (producer) and consumer service. `KAFKA_ADVERTISED_LISTENERS` must be configured correctly — it tells clients (producers and consumers) how to connect to the broker. `AUTO_CREATE_TOPICS_ENABLE=true` simplifies development by automatically creating topics when they're first used.
**Inputs:** Environment variables for broker configuration.
**Returns:** Kafka broker available on port 9092.
**Concepts to learn before implementing:**
- Kafka broker configuration → search "kafka docker compose configuration"
- Advertised listeners → search "kafka advertised listeners docker"
- Topic auto-creation → search "kafka auto create topics"
**Depends on:** TASK-182

---

#### TASK-184: backend-api Service
**Signature:** Service block: `build: ./backend`, `ports: 8080:8080`, `depends_on: [postgres, redis, kafka]` (with conditions), `environment: DATABASE_URL, JWT_SECRET, REDIS_URL, KAFKA_BROKERS, CLAUDE_API_KEY`
**File Location:** `docker-compose.yml`
**Purpose:** Builds and runs the main Go API server.
**Why it exists:** The API server is the primary backend service handling HTTP requests. `depends_on` with healthcheck conditions ensures it starts only after PostgreSQL, Redis, and Kafka are healthy — preventing connection errors at startup. Environment variables are defined here, centralizing configuration in one file.
**Inputs:** Environment variables for all service connections.
**Returns:** API server available on port 8080.
**Concepts to learn before implementing:**
- Docker Compose build context → search "docker compose build context"
- depends_on with condition → search "docker compose depends_on condition"
- Environment variables in Compose → search "docker compose environment variables"
**Depends on:** TASK-177, TASK-180, TASK-181, TASK-183

---

#### TASK-185: backend-consumer Service
**Signature:** Service block: same build as API but with `command: /consumer`, `depends_on: [postgres, kafka]`
**File Location:** `docker-compose.yml`
**Purpose:** Runs the Kafka consumer as a separate container.
**Why it exists:** Same Docker image as the API server but with a different `command` to run the consumer binary. This is the microservice pattern — one image, multiple entry points. The consumer doesn't need Redis (it doesn't rate-limit) or its own port (it doesn't serve HTTP requests — it only reads from Kafka and writes to PostgreSQL).
**Inputs:** Environment variables for Kafka and PostgreSQL.
**Returns:** Consumer process running in background.
**Concepts to learn before implementing:**
- Multiple binaries from one Docker image → search "docker multi binary golang"
- CMD override in Compose → search "docker compose command override"
**Depends on:** TASK-177, TASK-180, TASK-183

---

#### TASK-186: backend-classifier Service
**Signature:** Service block: same build with `command: /classifier`, `ports: 50051:50051`
**File Location:** `docker-compose.yml`
**Purpose:** Runs the gRPC classifier as a separate container.
**Why it exists:** The gRPC classifier runs independently. It doesn't need database or Kafka access — it only needs the detection rules which are compiled into the binary. It's stateless and CPU-bound, making it easy to scale. Port 50051 is the gRPC convention.
**Inputs:** Minimal environment variables.
**Returns:** gRPC classifier available on port 50051.
**Depends on:** TASK-177

---

#### TASK-187: frontend Service
**Signature:** Service block: `build: ./frontend`, `ports: 3000:80`, `depends_on: [backend-api]`
**File Location:** `docker-compose.yml`
**Purpose:** Builds and runs the React frontend served by Nginx.
**Why it exists:** The frontend container serves the React app via Nginx. Port mapping `3000:80` exposes it on `localhost:3000` while Nginx inside listens on 80. It depends on the API server being available so the proxy can forward API requests.
**Inputs:** None (static files, no environment variables needed at runtime).
**Returns:** Frontend available on port 3000.
**Depends on:** TASK-179, TASK-184

---

### FILE: `k8s/api-deployment.yaml`
**Purpose of this file:** Kubernetes manifests for deploying the API server with replicas, health probes, and a stable Service endpoint.
**Week:** Week 7
**Layer:** DevOps

---

#### TASK-188: Deployment Spec
**Signature:** `apiVersion: apps/v1`, `kind: Deployment`, `metadata.name: goshield-api`, `spec.replicas: 2`, template with container spec (image, ports, envFrom, readinessProbe, livenessProbe, resources)
**File Location:** `k8s/api-deployment.yaml`
**Purpose:** Defines how Kubernetes should run and manage the API server pods.
**Why it exists:** A Deployment manages a set of identical pods, ensuring the desired number of replicas (2) are always running. If a pod crashes, Kubernetes automatically replaces it. The probes use the `/health` endpoint (TASK-012) — the readiness probe determines when a pod can receive traffic (after startup), and the liveness probe detects hung processes for automatic restart. Resource limits prevent one service from consuming all cluster resources.
**Inputs:** Container image reference, environment variable references to Secrets.
**Returns:** Running pods managed by Kubernetes.
**Concepts to learn before implementing:**
- Kubernetes Deployment → search "kubernetes deployment tutorial"
- Readiness vs liveness probes → search "kubernetes readiness liveness probe"
- Resource limits and requests → search "kubernetes resource limits requests"
- Replicas for high availability → search "kubernetes replicas scaling"
- Pod template spec → search "kubernetes pod template"
**Depends on:** TASK-177, TASK-012

---

#### TASK-189: Service Spec
**Signature:** `apiVersion: v1`, `kind: Service`, `metadata.name: goshield-api`, `spec.type: ClusterIP`, `ports: [{port: 8080}]`, `selector` matching deployment labels
**File Location:** `k8s/api-deployment.yaml`
**Purpose:** Creates a stable network endpoint for the API pods.
**Why it exists:** Pods are ephemeral — they can be created, destroyed, and rescheduled with new IP addresses. A Kubernetes Service gives a stable DNS name (`goshield-api:8080`) that automatically load-balances across all healthy pods. Other services and the Ingress controller use this Service name to communicate with the API, never individual pod IPs.
**Inputs:** Label selector matching deployment pods.
**Returns:** Stable ClusterIP endpoint.
**Concepts to learn before implementing:**
- Kubernetes Service → search "kubernetes service types"
- ClusterIP vs NodePort vs LoadBalancer → search "kubernetes service types explained"
- Label selectors → search "kubernetes label selector"
- DNS-based service discovery → search "kubernetes service discovery dns"
**Depends on:** TASK-188

---

### FILE: `k8s/classifier-deployment.yaml`
**Purpose of this file:** Kubernetes manifests for deploying the gRPC classifier service.
**Week:** Week 7
**Layer:** DevOps

---

#### TASK-190: Deployment Spec
**Signature:** Deployment for the classifier: container runs the classifier binary, port 50051, gRPC health check probe.
**File Location:** `k8s/classifier-deployment.yaml`
**Purpose:** Manages gRPC classifier pods with health checking and scaling.
**Why it exists:** The classifier is stateless and CPU-bound (string pattern matching), making it ideal for horizontal scaling. Each replica handles classification requests independently. The gRPC health check uses the standard gRPC health checking protocol instead of HTTP probes.
**Inputs:** Container image, minimal environment.
**Returns:** Running classifier pods.
**Depends on:** TASK-177

---

#### TASK-191: Service Spec
**Signature:** ClusterIP Service on port 50051 with gRPC protocol.
**File Location:** `k8s/classifier-deployment.yaml`
**Purpose:** Provides a stable internal endpoint for gRPC communication.
**Why it exists:** The consumer and API server call the classifier via gRPC. The Service provides a stable DNS endpoint (`goshield-classifier:50051`) for these internal RPC calls, load-balancing across classifier replicas.
**Inputs:** Label selector.
**Returns:** Stable gRPC endpoint.
**Depends on:** TASK-190

---

### FILE: `k8s/consumer-deployment.yaml`
**Purpose of this file:** Kubernetes manifests for deploying the Kafka consumer.
**Week:** Week 7
**Layer:** DevOps

---

#### TASK-192: Deployment Spec
**Signature:** Deployment with `replicas: 1`, environment variables for Kafka/Postgres/Claude API.
**File Location:** `k8s/consumer-deployment.yaml`
**Purpose:** Manages the Kafka consumer pod(s).
**Why it exists:** The consumer starts with 1 replica because Kafka consumer groups handle parallelism — you can scale up to the number of topic partitions. If the topic has 3 partitions, you can run 3 consumer replicas, each processing one partition. Environment variables inject connection details from Kubernetes Secrets (TASK-194).
**Inputs:** Container image, environment from Secrets.
**Returns:** Running consumer pod(s).
**Concepts to learn before implementing:**
- Kafka consumer scaling with partitions → search "kafka consumer group partitions scaling"
- Environment variables from Secrets → search "kubernetes env from secret"
**Depends on:** TASK-177

---

#### TASK-193: Service Spec
**Signature:** ClusterIP Service for the consumer (primarily for monitoring/discovery).
**File Location:** `k8s/consumer-deployment.yaml`
**Purpose:** Provides a Service endpoint for the consumer.
**Why it exists:** Even though the consumer doesn't serve inbound HTTP/gRPC traffic (it only reads from Kafka and writes to PostgreSQL), a Service is useful for monitoring tools to discover and health-check the consumer pods. Some setups also use it for internal DNS resolution.
**Inputs:** Label selector.
**Returns:** Service endpoint.
**Depends on:** TASK-192

---

### FILE: `k8s/secrets.yaml`
**Purpose of this file:** Kubernetes Secrets manifest storing all sensitive configuration values used by the GoShield services.
**Week:** Week 7
**Layer:** DevOps

---

#### TASK-194: Secret Spec
**Signature:** `apiVersion: v1`, `kind: Secret`, `metadata.name: goshield-secrets`, `type: Opaque`, `data:` with base64-encoded values for DATABASE_URL, JWT_SECRET, REDIS_URL, KAFKA_BROKERS, CLAUDE_API_KEY
**File Location:** `k8s/secrets.yaml`
**Purpose:** Stores sensitive configuration values separately from deployment manifests.
**Why it exists:** Secrets keep passwords, API keys, and connection strings out of Deployment manifests and version control. Values are base64-encoded (not encrypted — base64 is encoding, not encryption). For real security, use tools like Sealed Secrets or HashiCorp Vault. Pods reference secrets via `envFrom: secretRef` in their container spec.
**Inputs:** Base64-encoded secret values.
**Returns:** A Kubernetes Secret object.
**Concepts to learn before implementing:**
- Kubernetes Secrets → search "kubernetes secrets tutorial"
- Base64 encoding → search "base64 encoding explained"
- Secret management best practices → search "kubernetes secret management"
- Sealed Secrets for GitOps → search "kubernetes sealed secrets"
- `kubectl create secret` command → search "kubectl create secret generic"
**Depends on:** None

---

### FILE: `.github/workflows/ci-backend.yml`
**Purpose of this file:** GitHub Actions CI/CD pipeline for the Go backend — runs tests on every push/PR and builds the Docker image on main branch merges.
**Week:** Week 7
**Layer:** DevOps

---

#### TASK-195: test Job
**Signature:** `job: test`, `runs-on: ubuntu-latest`, steps: checkout → setup-go → `go mod download` → `go test ./...`
**File Location:** `.github/workflows/ci-backend.yml`
**Purpose:** Runs all Go tests on every push and pull request.
**Why it exists:** Automated testing is the first gate in CI/CD. Every push and PR triggers the test suite. If any test fails, the pipeline stops and the PR cannot be merged. This prevents broken code from reaching the main branch. The `go test ./...` command runs tests in all packages recursively.
**Inputs:** Go source code from the repository.
**Returns:** Pass/fail status.
**Concepts to learn before implementing:**
- GitHub Actions basics → search "github actions tutorial"
- Go testing in CI → search "golang test ci github actions"
- Go module caching in CI → search "github actions cache go modules"
- Workflow triggers (push, pull_request) → search "github actions workflow triggers"
**Depends on:** TASK-199 (tests must exist to run)

---

#### TASK-196: build-and-push Job
**Signature:** `job: build-and-push`, `needs: test`, steps: checkout → docker/login-action → docker/build-push-action (tag with commit SHA + latest)
**File Location:** `.github/workflows/ci-backend.yml`
**Purpose:** Builds the Docker image and pushes it to a container registry after tests pass.
**Why it exists:** After tests pass, the Docker image is built and pushed to a container registry (Docker Hub or GitHub Container Registry). Tagging with the commit SHA ensures every build is uniquely identifiable — you can always roll back to a specific commit's image. The `needs: test` dependency ensures images are only published for code that passes all tests. This is the CD (Continuous Deployment) part of CI/CD.
**Inputs:** Dockerfile, registry credentials (stored as GitHub Secrets).
**Returns:** Docker image pushed to registry.
**Concepts to learn before implementing:**
- Docker image tagging strategies → search "docker image tag semver sha"
- GitHub Actions docker/build-push-action → search "github actions docker build push"
- Container registries → search "docker container registry"
- GitHub Secrets for CI → search "github actions secrets"
- CI/CD pipeline stages → search "ci cd pipeline stages"
**Depends on:** TASK-195, TASK-177

---

### FILE: `.github/workflows/ci-frontend.yml`
**Purpose of this file:** GitHub Actions CI/CD pipeline for the React frontend — lints code and verifies the production build succeeds.
**Week:** Week 7
**Layer:** DevOps

---

#### TASK-197: lint Job
**Signature:** `job: lint`, `runs-on: ubuntu-latest`, steps: checkout → setup-node → `npm ci` → `npm run lint`
**File Location:** `.github/workflows/ci-frontend.yml`
**Purpose:** Runs ESLint and TypeScript type checking on every push and PR.
**Why it exists:** Linting catches code quality issues (unused variables, inconsistent formatting), potential bugs (unreachable code, missing return statements), and type errors (accessing undefined properties) before they reach production. Running lint in CI ensures the entire team maintains consistent code quality even if individual developers forget to lint locally.
**Inputs:** Frontend source code.
**Returns:** Pass/fail status.
**Concepts to learn before implementing:**
- ESLint in CI → search "eslint ci pipeline"
- TypeScript type checking (`tsc --noEmit`) → search "typescript tsc noEmit"
- npm ci for reproducible installs → search "npm ci ci pipeline"
**Depends on:** TASK-169

---

#### TASK-198: build Job
**Signature:** `job: build`, `needs: lint`, steps: checkout → setup-node → `npm ci` → `npm run build`
**File Location:** `.github/workflows/ci-frontend.yml`
**Purpose:** Verifies that the production build compiles successfully.
**Why it exists:** The build job catches issues that linting alone misses — like dynamic import errors, missing assets, or environment variable problems that only surface during production bundling. If the build fails, the pipeline stops. Optionally, this job can also build and push the frontend Docker image.
**Inputs:** Frontend source code.
**Returns:** Pass/fail status (+ optionally Docker image).
**Depends on:** TASK-197, TASK-179

---
---

## 🧪 WEEK 8 — Testing & Demo Data

### FILE: `backend/tests/auth_test.go`
**Purpose of this file:** Unit tests for the JWT authentication service — verifies that token generation and validation work correctly and securely.
**Week:** Week 8
**Layer:** Backend

---

#### TASK-199: TestGenerateAndValidate
**Signature:** `func TestGenerateAndValidate(t *testing.T)`
**File Location:** `backend/tests/auth_test.go`
**Purpose:** Tests the happy path — generates a JWT token and then validates it, verifying the claims match.
**Why it exists:** This is the foundational auth test. It proves that a generated token can be validated and the embedded user ID and role are correctly stored and retrieved. If this test fails, authentication is fundamentally broken — no user can log in or access protected endpoints.
**Test steps:**
1. Create an `auth.Service` with a test secret and 1-hour expiry
2. Call `Generate` with a test UUID and role "analyst"
3. Call `Validate` with the returned token string
4. Assert `claims.UserID` matches the test UUID
5. Assert `claims.Role` equals "analyst"
6. Assert no error was returned
**Concepts to learn before implementing:**
- Go testing package → search "golang testing tutorial"
- Table-driven tests → search "golang table driven tests"
- `t.Run` for subtests → search "golang t.Run subtests"
- Test assertions → search "golang test assert equal"
- `testing.T` methods (Error, Fatal, Errorf) → search "golang testing.T methods"
**Depends on:** TASK-004, TASK-005, TASK-006

---

#### TASK-200: TestInvalidToken
**Signature:** `func TestInvalidToken(t *testing.T)`
**File Location:** `backend/tests/auth_test.go`
**Purpose:** Tests that Validate correctly rejects a garbled/tampered token string.
**Why it exists:** Security test — ensures the system rejects invalid tokens. An attacker might try to forge a token or modify its payload (changing their role from "analyst" to "admin"). This test verifies that any such attempt results in a validation error and nil claims.
**Test steps:**
1. Create an `auth.Service` with a test secret
2. Call `Validate` with the string "not.a.valid.token"
3. Assert an error IS returned
4. Assert claims is nil
**Depends on:** TASK-004, TASK-006

---

#### TASK-201: TestWrongSecret
**Signature:** `func TestWrongSecret(t *testing.T)`
**File Location:** `backend/tests/auth_test.go`
**Purpose:** Tests that a token generated with one secret fails validation with a different secret.
**Why it exists:** This verifies the HMAC signature verification — the most critical security property of JWT. If someone obtains a valid token from System A (with secret "key-a"), it must NOT be valid on System B (with secret "key-b"). Without this property, a leaked token from any system would compromise all systems.
**Test steps:**
1. Create `service1` with secret "secret-a"
2. Generate a token with service1
3. Create `service2` with secret "secret-b"
4. Call `Validate` on service2 with the token from service1
5. Assert an error IS returned
**Depends on:** TASK-004, TASK-005, TASK-006

---

### FILE: `backend/tests/detector_test.go`
**Purpose of this file:** Tests for the rule-based threat detection engine — verifies correct classification of various attack patterns and benign messages.
**Week:** Week 8
**Layer:** Backend

---

#### TASK-202: TestClassify
**Signature:** `func TestClassify(t *testing.T)` with table-driven test cases
**File Location:** `backend/tests/detector_test.go`
**Purpose:** Tests that the Classify function correctly identifies threats from log messages using a comprehensive set of test cases.
**Why it exists:** The detector is a critical security component. False negatives (missing real threats) mean attacks go undetected. False positives (flagging safe logs) create alert fatigue. Table-driven tests efficiently cover many scenarios with minimal code duplication.
**Table rows (each explained as a sub-case):**
- **"SQL injection detected"** — input: `"SELECT * FROM users WHERE 1=1"`, expected level: `"critical"`, expected rule: `"SQL Injection"` — tests the most common attack pattern
- **"XSS attack detected"** — input: `"<script>alert('xss')</script>"`, expected level: `"high"`, expected rule: `"XSS Attack"` — tests cross-site scripting detection
- **"Normal log (no threat)"** — input: `"User logged in successfully"`, expected level: `"low"`, expected rule: `""` — tests that benign messages aren't flagged
- **"Case insensitive match"** — input: `"drop table users"` (lowercase), expected level: `"critical"` — verifies case-insensitive pattern matching works
- **"Path traversal"** — input: `"GET /../../etc/passwd"`, expected level: `"high"`, expected rule: `"Path Traversal"` — tests directory traversal detection
- **"Brute force"** — input: `"failed login attempt from 192.168.1.1"`, expected level: `"medium"`, expected rule: `"Brute Force"` — tests authentication attack detection
**Concepts to learn before implementing:**
- Table-driven tests in Go → search "golang table driven tests pattern"
- Security testing → search "unit testing security functions"
- Test coverage → search "golang test coverage"
- `t.Run` with test case names → search "golang subtests table driven"
**Depends on:** TASK-048

---

### FILE: `backend/tests/log_handler_test.go`
**Purpose of this file:** Integration-style tests for the log HTTP handlers — verifies the full HTTP request/response cycle without a running server.
**Week:** Week 8
**Layer:** Backend

---

#### TASK-203: TestIngestHandler
**Signature:** `func TestIngestHandler(t *testing.T)`
**File Location:** `backend/tests/log_handler_test.go`
**Purpose:** Tests the POST /api/logs ingest endpoint by sending a mock HTTP request and verifying the response.
**Why it exists:** Handler tests verify the complete HTTP flow — JSON parsing, threat classification, database insertion, alert creation, and JSON response formatting — without needing a running server or real database. The `httptest` package provides mock request/response objects that simulate real HTTP traffic.
**Test steps:**
1. Create mock repositories (with stub methods) and a detector function
2. Create a `LogHandler` with the mocks
3. Create `httptest.NewRequest("POST", "/api/logs", jsonBody)` with a valid IngestRequest JSON body
4. Create `httptest.NewRecorder()` as the response writer
5. Call `handler.Ingest(recorder, request)`
6. Assert response status code is 201 Created
7. Decode the response body and verify the returned SecurityLog has the correct fields
**Concepts to learn before implementing:**
- httptest package → search "golang httptest tutorial"
- Mock dependencies in tests → search "golang mock interface testing"
- Testing HTTP handlers → search "golang test http handler"
- httptest.NewRecorder → search "golang httptest NewRecorder"
- httptest.NewRequest → search "golang httptest NewRequest"
**Depends on:** TASK-051, TASK-052

---

#### TASK-204: TestListHandler
**Signature:** `func TestListHandler(t *testing.T)`
**File Location:** `backend/tests/log_handler_test.go`
**Purpose:** Tests the GET /api/logs endpoint with query parameters for filtering and pagination.
**Why it exists:** Verifies that pagination parameters (`page`, `limit`) and filter parameters (`level`, `source`) are correctly extracted from the URL query string and passed to the repository. Also tests the response format including the total count and pagination metadata.
**Test steps:**
1. Create a mock LogRepo that returns sample data and a known total count
2. Create a LogHandler with the mock
3. Create `httptest.NewRequest("GET", "/api/logs?level=high&page=2&limit=10", nil)`
4. Create `httptest.NewRecorder()`
5. Call `handler.List(recorder, request)`
6. Assert response status is 200 OK
7. Decode response and verify `data` array, `total`, `page`, and `limit` fields
**Depends on:** TASK-051, TASK-053

---

### FILE: `backend/cmd/seed/main.go`
**Purpose of this file:** Demo data seeder — generates and sends realistic sample security logs to populate the system with test data for demonstrations and development.
**Week:** Week 8
**Layer:** Backend

---

#### TASK-205: sampleLogs Slice
**Signature:** `var sampleLogs = []model.IngestRequest{ ... }`
**File Location:** `backend/cmd/seed/main.go`
**Purpose:** Defines a collection of realistic sample security log messages covering all threat levels and sources.
**Why it exists:** A demo system needs realistic data to showcase all features — different threat levels (so stat cards show varied numbers), different sources (so the pie chart has multiple slices), and different attack types (so the rule engine demonstrates its detection capabilities). Example entries:
- SQL injection attempt from firewall (critical)
- Failed login from auth_service (medium)
- Port scan from IDS (medium)
- Normal health check from endpoint (low)
- XSS attempt from firewall (high)
- Path traversal from IDS (high)
- Successful auth from auth_service (low)
**Inputs:** N/A (package-level variable).
**Returns:** N/A (slice of IngestRequest objects).
**Concepts to learn before implementing:**
- Test data design → search "test data best practices"
- Realistic demo data → search "seed data database"
- Slice literals in Go → search "golang slice literal"
**Depends on:** TASK-033

---

#### TASK-206: sendLog Function
**Signature:** `func sendLog(client *http.Client, baseURL string, log model.IngestRequest) error`
**File Location:** `backend/cmd/seed/main.go`
**Purpose:** Sends a single log to the ingest API endpoint via HTTP POST.
**Why it exists:** Encapsulates the HTTP call logic — serializes the IngestRequest to JSON, POSTs it to `/api/logs`, and checks the response status. Having this as a separate function keeps the main loop clean and makes error handling per-log possible (log the error and continue with the next, rather than crashing).
**Inputs:**
- `client` (*http.Client) — configured HTTP client.
- `baseURL` (string) — API server URL, e.g., `"http://localhost:8080"`.
- `log` (model.IngestRequest) — the log to send.
**Returns:** `error` — returned if serialization, HTTP request, or response status check fails.
**Concepts to learn before implementing:**
- HTTP client POST in Go → search "golang http.Client POST"
- json.Marshal for request bodies → search "golang json marshal http body"
- bytes.NewReader for request body → search "golang bytes NewReader"
- Checking HTTP response status → search "golang http response status check"
**Depends on:** TASK-033

---

#### TASK-207: main Loop
**Signature:** `func main()` — iterates over sampleLogs and sends each one with a delay.
**File Location:** `backend/cmd/seed/main.go`
**Purpose:** Orchestrates the demo data population by sending all sample logs with realistic timing.
**Why it exists:** The seeder's main function sends logs with a configurable delay (e.g., 500ms) between each to simulate realistic ingestion timing. This exercises the full pipeline: API → Kafka → Consumer → Detector → AI → Database → WebSocket. The delay also prevents overwhelming the rate limiter (TASK-050). Logs are sent sequentially so you can watch them appear one-by-one on the dashboard.
**Steps:**
1. Parse base URL from command-line args or environment variable
2. Create an HTTP client
3. Loop over sampleLogs slice
4. For each log, call sendLog
5. Log success/failure for each
6. Sleep between sends (configurable delay)
**Concepts to learn before implementing:**
- time.Sleep for delays → search "golang time.Sleep"
- Command-line arguments (os.Args) → search "golang os.Args"
- Seeder/migration scripts → search "database seeder script"
- Sequential vs parallel sending → search "golang sequential processing"
**Depends on:** TASK-205, TASK-206

---
---

## Quick Reference Index

| Task # | Function/Component Name | File | Week | Layer |
|--------|------------------------|------|------|-------|
| TASK-001 | Load | `backend/internal/config/config.go` | Week 1 | Backend |
| TASK-002 | Role Constants | `backend/internal/model/user.go` | Week 1 | Backend |
| TASK-003 | User Struct | `backend/internal/model/user.go` | Week 1 | Backend |
| TASK-004 | NewService | `backend/internal/auth/jwt.go` | Week 1 | Backend |
| TASK-005 | Generate | `backend/internal/auth/jwt.go` | Week 1 | Backend |
| TASK-006 | Validate | `backend/internal/auth/jwt.go` | Week 1 | Backend |
| TASK-007 | JWTMiddleware | `backend/internal/auth/middleware.go` | Week 1 | Backend |
| TASK-008 | NewUserRepo | `backend/internal/repository/user_repo.go` | Week 1 | Backend |
| TASK-009 | Create | `backend/internal/repository/user_repo.go` | Week 1 | Backend |
| TASK-010 | FindByEmail | `backend/internal/repository/user_repo.go` | Week 1 | Backend |
| TASK-011 | FindByID | `backend/internal/repository/user_repo.go` | Week 1 | Backend |
| TASK-012 | Health | `backend/internal/handler/health.go` | Week 1 | Backend |
| TASK-013 | NewAuthHandler | `backend/internal/handler/auth_handler.go` | Week 1 | Backend |
| TASK-014 | Register | `backend/internal/handler/auth_handler.go` | Week 1 | Backend |
| TASK-015 | Login | `backend/internal/handler/auth_handler.go` | Week 1 | Backend |
| TASK-016 | Me | `backend/internal/handler/auth_handler.go` | Week 1 | Backend |
| TASK-017 | main — Load Config | `backend/cmd/server/main.go` | Week 1 | Backend |
| TASK-018 | main — Connect to PostgreSQL | `backend/cmd/server/main.go` | Week 1 | Backend |
| TASK-019 | main — Wire Repositories | `backend/cmd/server/main.go` | Week 1 | Backend |
| TASK-020 | main — Wire Handlers | `backend/cmd/server/main.go` | Week 1 | Backend |
| TASK-021 | main — Register Routes | `backend/cmd/server/main.go` | Week 1 | Backend |
| TASK-022 | main — Start HTTP Server | `backend/cmd/server/main.go` | Week 1 | Backend |
| TASK-023 | CREATE TYPE user_role | `migrations/001_create_users.sql` | Week 1 | Backend |
| TASK-024 | CREATE TABLE users | `migrations/001_create_users.sql` | Week 1 | Backend |
| TASK-025 | CREATE INDEX on users | `migrations/001_create_users.sql` | Week 1 | Backend |
| TASK-026 | CREATE TABLE security_logs | `migrations/002_create_logs.sql` | Week 2 | Backend |
| TASK-027 | CREATE INDEX on security_logs | `migrations/002_create_logs.sql` | Week 2 | Backend |
| TASK-028 | CREATE TABLE alerts | `migrations/003_create_alerts.sql` | Week 2 | Backend |
| TASK-029 | CREATE INDEX on alerts | `migrations/003_create_alerts.sql` | Week 2 | Backend |
| TASK-030 | ThreatLevel Constants | `backend/internal/model/log.go` | Week 2 | Backend |
| TASK-031 | LogSource Constants | `backend/internal/model/log.go` | Week 2 | Backend |
| TASK-032 | SecurityLog Struct | `backend/internal/model/log.go` | Week 2 | Backend |
| TASK-033 | IngestRequest Struct | `backend/internal/model/log.go` | Week 2 | Backend |
| TASK-034 | Alert Struct | `backend/internal/model/alert.go` | Week 2 | Backend |
| TASK-035 | NewLogRepo | `backend/internal/repository/log_repo.go` | Week 2 | Backend |
| TASK-036 | Insert | `backend/internal/repository/log_repo.go` | Week 2 | Backend |
| TASK-037 | InsertWithSummary | `backend/internal/repository/log_repo.go` | Week 2 | Backend |
| TASK-038 | List | `backend/internal/repository/log_repo.go` | Week 2 | Backend |
| TASK-039 | CountByLevel | `backend/internal/repository/log_repo.go` | Week 2 | Backend |
| TASK-040 | FindByID | `backend/internal/repository/log_repo.go` | Week 2 | Backend |
| TASK-041 | NewAlertRepo | `backend/internal/repository/alert_repo.go` | Week 2 | Backend |
| TASK-042 | Create | `backend/internal/repository/alert_repo.go` | Week 2 | Backend |
| TASK-043 | List | `backend/internal/repository/alert_repo.go` | Week 2 | Backend |
| TASK-044 | Resolve | `backend/internal/repository/alert_repo.go` | Week 2 | Backend |
| TASK-045 | FindByLogID | `backend/internal/repository/alert_repo.go` | Week 2 | Backend |
| TASK-046 | Rule Struct | `backend/internal/detector/rules.go` | Week 2 | Backend |
| TASK-047 | Default Rules Slice | `backend/internal/detector/rules.go` | Week 2 | Backend |
| TASK-048 | Classify | `backend/internal/detector/rules.go` | Week 2 | Backend |
| TASK-049 | New | `backend/internal/ratelimit/redis.go` | Week 2 | Backend |
| TASK-050 | Allow | `backend/internal/ratelimit/redis.go` | Week 2 | Backend |
| TASK-051 | NewLogHandler | `backend/internal/handler/log_handler.go` | Week 2 | Backend |
| TASK-052 | Ingest | `backend/internal/handler/log_handler.go` | Week 2 | Backend |
| TASK-053 | List | `backend/internal/handler/log_handler.go` | Week 2 | Backend |
| TASK-054 | GetByID | `backend/internal/handler/log_handler.go` | Week 2 | Backend |
| TASK-055 | NewAlertHandler | `backend/internal/handler/alert_handler.go` | Week 2 | Backend |
| TASK-056 | List | `backend/internal/handler/alert_handler.go` | Week 2 | Backend |
| TASK-057 | Resolve | `backend/internal/handler/alert_handler.go` | Week 2 | Backend |
| TASK-058 | NewProducer | `backend/internal/kafka/producer.go` | Week 3 | Backend |
| TASK-059 | Send | `backend/internal/kafka/producer.go` | Week 3 | Backend |
| TASK-060 | Close | `backend/internal/kafka/producer.go` | Week 3 | Backend |
| TASK-061 | NewConsumer | `backend/internal/kafka/consumer.go` | Week 3 | Backend |
| TASK-062 | Start | `backend/internal/kafka/consumer.go` | Week 3 | Backend |
| TASK-063 | Close | `backend/internal/kafka/consumer.go` | Week 3 | Backend |
| TASK-064 | main | `backend/cmd/consumer/main.go` | Week 3 | Backend |
| TASK-065 | ClassifyRequest Message | `backend/proto/classifier.proto` | Week 4 | Backend |
| TASK-066 | ClassifyResponse Message | `backend/proto/classifier.proto` | Week 4 | Backend |
| TASK-067 | ClassifierService Service | `backend/proto/classifier.proto` | Week 4 | Backend |
| TASK-068 | Classify RPC | `backend/proto/classifier.proto` | Week 4 | Backend |
| TASK-069 | Server Struct | `backend/classifier/server.go` | Week 4 | Backend |
| TASK-070 | Classify Method | `backend/classifier/server.go` | Week 4 | Backend |
| TASK-071 | main | `backend/cmd/classifier/main.go` | Week 4 | Backend |
| TASK-072 | Client Struct | `backend/internal/ai/claude.go` | Week 5 | Backend |
| TASK-073 | NewClient | `backend/internal/ai/claude.go` | Week 5 | Backend |
| TASK-074 | request Struct | `backend/internal/ai/claude.go` | Week 5 | Backend |
| TASK-075 | response Struct | `backend/internal/ai/claude.go` | Week 5 | Backend |
| TASK-076 | AnalyzeLog | `backend/internal/ai/claude.go` | Week 5 | Backend |
| TASK-077 | ShouldAnalyze | `backend/internal/ai/analyzer.go` | Week 5 | Backend |
| TASK-078 | ParseResponse | `backend/internal/ai/analyzer.go` | Week 5 | Backend |
| TASK-079 | Updated Start | `backend/internal/kafka/consumer.go` | Week 5 | Backend |
| TASK-080 | Hub Struct | `backend/internal/ws/hub.go` | Week 6 | Backend |
| TASK-081 | NewHub | `backend/internal/ws/hub.go` | Week 6 | Backend |
| TASK-082 | run | `backend/internal/ws/hub.go` | Week 6 | Backend |
| TASK-083 | Register | `backend/internal/ws/hub.go` | Week 6 | Backend |
| TASK-084 | Unregister | `backend/internal/ws/hub.go` | Week 6 | Backend |
| TASK-085 | Broadcast | `backend/internal/ws/hub.go` | Week 6 | Backend |
| TASK-086 | WebSocket Upgrade Handler | `backend/cmd/server/main.go` | Week 6 | Backend |
| TASK-087 | Role Type | `frontend/src/types/user.ts` | Week 6 | Frontend |
| TASK-088 | User Interface | `frontend/src/types/user.ts` | Week 6 | Frontend |
| TASK-089 | ThreatLevel Type | `frontend/src/types/log.ts` | Week 6 | Frontend |
| TASK-090 | LogSource Type | `frontend/src/types/log.ts` | Week 6 | Frontend |
| TASK-091 | SecurityLog Interface | `frontend/src/types/log.ts` | Week 6 | Frontend |
| TASK-092 | Alert Interface | `frontend/src/types/alert.ts` | Week 6 | Frontend |
| TASK-093 | axiosInstance Config | `frontend/src/api/client.ts` | Week 6 | Frontend |
| TASK-094 | Request Interceptor | `frontend/src/api/client.ts` | Week 6 | Frontend |
| TASK-095 | Response Interceptor | `frontend/src/api/client.ts` | Week 6 | Frontend |
| TASK-096 | login | `frontend/src/api/auth.api.ts` | Week 6 | Frontend |
| TASK-097 | register | `frontend/src/api/auth.api.ts` | Week 6 | Frontend |
| TASK-098 | me | `frontend/src/api/auth.api.ts` | Week 6 | Frontend |
| TASK-099 | getLogs | `frontend/src/api/logs.api.ts` | Week 6 | Frontend |
| TASK-100 | ingestLog | `frontend/src/api/logs.api.ts` | Week 6 | Frontend |
| TASK-101 | getLogById | `frontend/src/api/logs.api.ts` | Week 6 | Frontend |
| TASK-102 | getAlerts | `frontend/src/api/alerts.api.ts` | Week 6 | Frontend |
| TASK-103 | resolveAlert | `frontend/src/api/alerts.api.ts` | Week 6 | Frontend |
| TASK-104 | configureStore | `frontend/src/store/index.ts` | Week 6 | Frontend |
| TASK-105 | RootState Type | `frontend/src/store/index.ts` | Week 6 | Frontend |
| TASK-106 | AppDispatch Type | `frontend/src/store/index.ts` | Week 6 | Frontend |
| TASK-107 | initialState | `frontend/src/store/authSlice.ts` | Week 6 | Frontend |
| TASK-108 | setCredentials Action | `frontend/src/store/authSlice.ts` | Week 6 | Frontend |
| TASK-109 | logout Action | `frontend/src/store/authSlice.ts` | Week 6 | Frontend |
| TASK-110 | selectUser Selector | `frontend/src/store/authSlice.ts` | Week 6 | Frontend |
| TASK-111 | selectToken Selector | `frontend/src/store/authSlice.ts` | Week 6 | Frontend |
| TASK-112 | initialState | `frontend/src/store/alertSlice.ts` | Week 6 | Frontend |
| TASK-113 | addAlert Action | `frontend/src/store/alertSlice.ts` | Week 6 | Frontend |
| TASK-114 | clearAlerts Action | `frontend/src/store/alertSlice.ts` | Week 6 | Frontend |
| TASK-115 | selectAlerts Selector | `frontend/src/store/alertSlice.ts` | Week 6 | Frontend |
| TASK-116 | initialState | `frontend/src/store/logSlice.ts` | Week 6 | Frontend |
| TASK-117 | setLogs Action | `frontend/src/store/logSlice.ts` | Week 6 | Frontend |
| TASK-118 | setFilters Action | `frontend/src/store/logSlice.ts` | Week 6 | Frontend |
| TASK-119 | setPage Action | `frontend/src/store/logSlice.ts` | Week 6 | Frontend |
| TASK-120 | selectLogs Selector | `frontend/src/store/logSlice.ts` | Week 6 | Frontend |
| TASK-121 | selectFilters Selector | `frontend/src/store/logSlice.ts` | Week 6 | Frontend |
| TASK-122 | useAuth — login | `frontend/src/hooks/useAuth.ts` | Week 6 | Frontend |
| TASK-123 | useAuth — logout | `frontend/src/hooks/useAuth.ts` | Week 6 | Frontend |
| TASK-124 | useAuth — isAuthenticated | `frontend/src/hooks/useAuth.ts` | Week 6 | Frontend |
| TASK-125 | useAlerts — WebSocket setup | `frontend/src/hooks/useAlerts.ts` | Week 6 | Frontend |
| TASK-126 | useAlerts — onmessage | `frontend/src/hooks/useAlerts.ts` | Week 6 | Frontend |
| TASK-127 | useAlerts — cleanup | `frontend/src/hooks/useAlerts.ts` | Week 6 | Frontend |
| TASK-128 | useLogs — fetch | `frontend/src/hooks/useLogs.ts` | Week 6 | Frontend |
| TASK-129 | useLogs — pagination | `frontend/src/hooks/useLogs.ts` | Week 6 | Frontend |
| TASK-130 | useLogs — filters | `frontend/src/hooks/useLogs.ts` | Week 6 | Frontend |
| TASK-131 | useStats | `frontend/src/hooks/useStats.ts` | Week 6 | Frontend |
| TASK-132 | getThreatColor | `frontend/src/utils/threatColors.ts` | Week 6 | Frontend |
| TASK-133 | getThreatBgColor | `frontend/src/utils/threatColors.ts` | Week 6 | Frontend |
| TASK-134 | formatRelative | `frontend/src/utils/formatDate.ts` | Week 6 | Frontend |
| TASK-135 | formatAbsolute | `frontend/src/utils/formatDate.ts` | Week 6 | Frontend |
| TASK-136 | Badge | `frontend/src/components/ui/Badge.tsx` | Week 6 | Frontend |
| TASK-137 | Button | `frontend/src/components/ui/Button.tsx` | Week 6 | Frontend |
| TASK-138 | Input | `frontend/src/components/ui/Input.tsx` | Week 6 | Frontend |
| TASK-139 | Modal | `frontend/src/components/ui/Modal.tsx` | Week 6 | Frontend |
| TASK-140 | Spinner | `frontend/src/components/ui/Spinner.tsx` | Week 6 | Frontend |
| TASK-141 | StatCard | `frontend/src/components/ui/StatCard.tsx` | Week 6 | Frontend |
| TASK-142 | AppShell | `frontend/src/components/layout/AppShell.tsx` | Week 6 | Frontend |
| TASK-143 | Sidebar — navLinks | `frontend/src/components/layout/Sidebar.tsx` | Week 6 | Frontend |
| TASK-144 | Sidebar — logout | `frontend/src/components/layout/Sidebar.tsx` | Week 6 | Frontend |
| TASK-145 | TopBar | `frontend/src/components/layout/TopBar.tsx` | Week 6 | Frontend |
| TASK-146 | ThreatLevelBar — transform | `frontend/src/components/charts/ThreatLevelBar.tsx` | Week 6 | Frontend |
| TASK-147 | ThreatLevelBar — render | `frontend/src/components/charts/ThreatLevelBar.tsx` | Week 6 | Frontend |
| TASK-148 | LogTimelineLine | `frontend/src/components/charts/LogTimelineLine.tsx` | Week 6 | Frontend |
| TASK-149 | SourcePie | `frontend/src/components/charts/SourcePie.tsx` | Week 6 | Frontend |
| TASK-150 | AlertCard — toggle | `frontend/src/components/alerts/AlertCard.tsx` | Week 6 | Frontend |
| TASK-151 | AlertCard — render | `frontend/src/components/alerts/AlertCard.tsx` | Week 6 | Frontend |
| TASK-152 | AlertFeed | `frontend/src/components/alerts/AlertFeed.tsx` | Week 6 | Frontend |
| TASK-153 | LogFilters — handler | `frontend/src/components/logs/LogFilters.tsx` | Week 6 | Frontend |
| TASK-154 | LogFilters — state | `frontend/src/components/logs/LogFilters.tsx` | Week 6 | Frontend |
| TASK-155 | LogRow — toggle | `frontend/src/components/logs/LogRow.tsx` | Week 6 | Frontend |
| TASK-156 | LogRow — render | `frontend/src/components/logs/LogRow.tsx` | Week 6 | Frontend |
| TASK-157 | LogTable — sort | `frontend/src/components/logs/LogTable.tsx` | Week 6 | Frontend |
| TASK-158 | LogTable — render | `frontend/src/components/logs/LogTable.tsx` | Week 6 | Frontend |
| TASK-159 | LoginPage — state | `frontend/src/pages/LoginPage.tsx` | Week 6 | Frontend |
| TASK-160 | LoginPage — submit | `frontend/src/pages/LoginPage.tsx` | Week 6 | Frontend |
| TASK-161 | RegisterPage — state | `frontend/src/pages/RegisterPage.tsx` | Week 6 | Frontend |
| TASK-162 | RegisterPage — submit | `frontend/src/pages/RegisterPage.tsx` | Week 6 | Frontend |
| TASK-163 | DashboardPage — useEffect | `frontend/src/pages/DashboardPage.tsx` | Week 6 | Frontend |
| TASK-164 | DashboardPage — render | `frontend/src/pages/DashboardPage.tsx` | Week 6 | Frontend |
| TASK-165 | LogsPage — pagination | `frontend/src/pages/LogsPage.tsx` | Week 6 | Frontend |
| TASK-166 | LogsPage — filters | `frontend/src/pages/LogsPage.tsx` | Week 6 | Frontend |
| TASK-167 | AlertsPage — resolve | `frontend/src/pages/AlertsPage.tsx` | Week 6 | Frontend |
| TASK-168 | SettingsPage — update | `frontend/src/pages/SettingsPage.tsx` | Week 6 | Frontend |
| TASK-169 | App — router | `frontend/src/App.tsx` | Week 6 | Frontend |
| TASK-170 | App — ProtectedRoute | `frontend/src/App.tsx` | Week 6 | Frontend |
| TASK-171 | App — routes | `frontend/src/App.tsx` | Week 6 | Frontend |
| TASK-172 | main — Provider | `frontend/src/main.tsx` | Week 6 | Frontend |
| TASK-173 | main — StrictMode | `frontend/src/main.tsx` | Week 6 | Frontend |
| TASK-174 | main — render | `frontend/src/main.tsx` | Week 6 | Frontend |
| TASK-175 | defineConfig — proxy | `frontend/vite.config.ts` | Week 6 | Frontend |
| TASK-176 | Stage 1 — Builder | `Dockerfile.backend` | Week 7 | DevOps |
| TASK-177 | Stage 2 — Runner | `Dockerfile.backend` | Week 7 | DevOps |
| TASK-178 | Stage 1 — Builder | `Dockerfile.frontend` | Week 7 | DevOps |
| TASK-179 | Stage 2 — Nginx Runner | `Dockerfile.frontend` | Week 7 | DevOps |
| TASK-180 | postgres Service | `docker-compose.yml` | Week 7 | DevOps |
| TASK-181 | redis Service | `docker-compose.yml` | Week 7 | DevOps |
| TASK-182 | zookeeper Service | `docker-compose.yml` | Week 7 | DevOps |
| TASK-183 | kafka Service | `docker-compose.yml` | Week 7 | DevOps |
| TASK-184 | backend-api Service | `docker-compose.yml` | Week 7 | DevOps |
| TASK-185 | backend-consumer Service | `docker-compose.yml` | Week 7 | DevOps |
| TASK-186 | backend-classifier Service | `docker-compose.yml` | Week 7 | DevOps |
| TASK-187 | frontend Service | `docker-compose.yml` | Week 7 | DevOps |
| TASK-188 | Deployment Spec | `k8s/api-deployment.yaml` | Week 7 | DevOps |
| TASK-189 | Service Spec | `k8s/api-deployment.yaml` | Week 7 | DevOps |
| TASK-190 | Deployment Spec | `k8s/classifier-deployment.yaml` | Week 7 | DevOps |
| TASK-191 | Service Spec | `k8s/classifier-deployment.yaml` | Week 7 | DevOps |
| TASK-192 | Deployment Spec | `k8s/consumer-deployment.yaml` | Week 7 | DevOps |
| TASK-193 | Service Spec | `k8s/consumer-deployment.yaml` | Week 7 | DevOps |
| TASK-194 | Secret Spec | `k8s/secrets.yaml` | Week 7 | DevOps |
| TASK-195 | test Job | `.github/workflows/ci-backend.yml` | Week 7 | DevOps |
| TASK-196 | build-and-push Job | `.github/workflows/ci-backend.yml` | Week 7 | DevOps |
| TASK-197 | lint Job | `.github/workflows/ci-frontend.yml` | Week 7 | DevOps |
| TASK-198 | build Job | `.github/workflows/ci-frontend.yml` | Week 7 | DevOps |
| TASK-199 | TestGenerateAndValidate | `backend/tests/auth_test.go` | Week 8 | Backend |
| TASK-200 | TestInvalidToken | `backend/tests/auth_test.go` | Week 8 | Backend |
| TASK-201 | TestWrongSecret | `backend/tests/auth_test.go` | Week 8 | Backend |
| TASK-202 | TestClassify | `backend/tests/detector_test.go` | Week 8 | Backend |
| TASK-203 | TestIngestHandler | `backend/tests/log_handler_test.go` | Week 8 | Backend |
| TASK-204 | TestListHandler | `backend/tests/log_handler_test.go` | Week 8 | Backend |
| TASK-205 | sampleLogs Slice | `backend/cmd/seed/main.go` | Week 8 | Backend |
| TASK-206 | sendLog Function | `backend/cmd/seed/main.go` | Week 8 | Backend |
| TASK-207 | main Loop | `backend/cmd/seed/main.go` | Week 8 | Backend |
