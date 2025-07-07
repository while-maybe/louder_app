# Louder - A Hexagonal Backend API in Pure Go

[![Go Version](https://img.shields.io/badge/Go-1.21-00ADD8.svg?style=flat-square)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg?style=flat-square)](https://opensource.org/licenses/MIT)
[![Architect](https://img.shields.io/badge/Architecture-Hexagonal-8B008B.svg?style=flat-square)]()

A professional-grade backend service written in Go, meticulously designed to showcase a deep understanding of software architecture and the power of Go's standard libraries. This project is a demonstration of a **Hexagonal (Ports & Adapters) Architecture** and features a **custom `net/http` router**, proving a mastery of core backend principles without reliance on external frameworks.

The core application logic, which manages a simple directory of people, is completely decoupled from external infrastructure, ensuring exceptional testability, maintainability, and long-term flexibility.

## Architectural Highlights

This project was built to demonstrate specific, high-level software design patterns using Go's built-in capabilities.

1.  **Custom HTTP Router (Multiplexer):** The API is powered by a custom-built router that implements Go's native `http.Handler` interface. This lightweight multiplexer maps request paths to handlers, proving a fundamental understanding of the `net/http` package without the "magic" of a third-party framework.

2.  **Embedded & Decoupled SQL Queries:** SQL queries are kept in separate `.sql` files and embedded directly into the application binary at compile time using Go's `embed` package. At startup, a custom loader reads these files into a map, completely decoupling SQL logic from the Go code. This results in a self-contained executable with no loose dependencies, simplifying deployment while maintaining clean separation of concerns.

3.  **Hexagonal (Ports & Adapters) Architecture:** The codebase strictly isolates the core domain logic from all external concerns.
    *   **Ports:** Abstract interfaces in `internal/core/ports` define the contracts for data persistence (e.g., `PeopleRepository`).
    *   **Adapters:** The `internal/server` package acts as the primary adapter, translating HTTP requests into calls to the core service. The `internal/db` package provides the concrete repository implementation.

4.  **Professional Database Layer:**
    *   **SQLX:** The project uses `github.com/jmoiron/sqlx` over the standard `database/sql` library to provide safer, more convenient query binding and result scanning directly into Go structs.
    *   **Database Migrations:** Schema changes are managed via raw SQL files in `internal/db/migrations` and applied using `golang-migrate/migrate`. This is the industry standard for versioning the database schema alongside application code.

5.  **Performance-Optimized Primary Keys (UUIDv7):** The project uses `github.com/gofrs/uuid` to generate **UUIDv7** primary keys. These modern, time-ordered identifiers provide the uniqueness of a UUID while being highly efficient for database indexing, preventing index fragmentation and improving write throughput.

## Summary of Demonstrated Best Practices

*   **Dependency Inversion:** The core application depends on abstract ports, not concrete infrastructure.
*   **Decoupled Architecture:** Business logic is cleanly separated from API and database concerns.
*   **Self-Contained Binaries:** SQL queries are embedded into the application, simplifying deployment.
*   **Database Schema Versioning:** Using `golang-migrate` with raw SQL files for repeatable deployments.
*   **Secure & Efficient Database Access:** Leveraging `sqlx` for safe query parameterization and struct scanning.
*   **`net/http` Mastery:** Proving fundamental knowledge by building a custom router.
*   **Optimized Data Modeling:** Using UUIDv7 for high-performance primary key indexing.
*   **Configuration Management:** Loading configuration from environment variables.
*   **Reproducible Environments:** A `docker-compose.yml` file is provided for a production-like Postgres environment.

## Tech Stack

*   **Language:** Golang
*   **Core HTTP:** `net/http` (with a custom router)
*   **File Embedding:** `embed` (Go Standard Library)
*   **Database Driver:** `github.com/jmoiron/sqlx`, `github.com/mattn/go-sqlite3`
*   **Database Migrations:** `github.com/golang-migrate/migrate/v4`
*   **ID Generation:** `github.com/gofrs/uuid`
*   **Containerization:** Docker & Docker Compose

## Getting Started

### Prerequisites

*   Go 1.21+
*   Git

### Running Locally (with SQLite)

The application runs with a file-based SQLite database by default, requiring no external services.

1.  **Clone the repository:** `git clone https://github.com/while-maybe/louder_app.git && cd louder_app`
2.  **Run the database migrations:** `go run ./cmd/migrate`
3.  **Run the server:** `go run ./cmd/server`

The API will be running on `http://localhost:8080`.

## API Endpoints

The base URL is `/api/v1`. The current implementation provides the core functionality for creating and retrieving people.

| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `GET` | `/people` | Retrieves a complete list of all people in the database. |
| `POST`| `/people` | Creates a new person. |

### Example cURL Requests

```bash
# Create a new person
curl -X POST http://localhost:8080/api/v1/people \
-H "Content-Type: application/json" \
-d '{"name": "Jane Doe"}'

# Get all people
curl http://localhost:8080/api/v1/people
```

## The Road Ahead 🗺️
This project is built on a solid foundation, ready for new features. Here is the vision for its evolution:
### 🧪 Quality & Robustness
Comprehensive Unit Testing: Implement a full test suite for the core services, using mocks for the repository port to ensure true isolation and validate all business logic.
### ✨ Feature Expansion
Dynamic Queries: Enhance the GET /people endpoint with robust server-side pagination and sorting, allowing clients to control the data they receive.
### 🐘 Production Readiness
PostgreSQL Adapter: Develop a new repository adapter to connect to the production-ready PostgreSQL database defined in docker-compose.yml.
JWT Authentication: Secure the API by implementing a security adapter with JWT for user authentication and authorization on protected routes.
### 🚀 API Evolution
Admin API with Gin: Launch a separate, feature-rich admin API on a new port using the Gin framework, demonstrating proficiency with modern Go frameworks for more complex use cases.
