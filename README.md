# Louder - A Hexagonal Architecture Backend API

[![Go Version](https://img.shields.io/badge/Go-1.21-00ADD8.svg?style=flat-square)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg?style=flat-square)](https://opensource.org/licenses/MIT)
[![Architect](https://img.shields.io/badge/Architecture-Hexagonal-8B008B.svg?style=flat-square)]()

A demonstration of a professional-grade backend service built in Golang, designed around a **Hexagonal (Ports & Adapters) Architecture**. This project showcases best practices for creating decoupled, testable, and maintainable systems ready for real-world deployment.

The core application logic is completely independent of external frameworks and infrastructure, allowing for easy adaptation and long-term stability.

## Architectural Highlights

This project was built to demonstrate specific, high-level software design patterns:

1.  **Hexagonal (Ports & Adapters) Architecture:** The codebase is structured to isolate the core domain logic from external concerns.
    *   **Ports:** Explicit interfaces (`internal/core/ports`) define the contracts for interacting with the core application (e.g., `PostsRepository`).
    *   **Adapters:** Concrete implementations for these ports handle infrastructure concerns. The current adapters include:
        *   A **Gin handler** for the REST API.
        *   A **SQLite repository** for data persistence.
    *   **Benefit:** This decoupling allows for swapping infrastructure (e.g., from Gin to gRPC, or SQLite to PostgreSQL) without modifying the core business logic.

2.  **Performance-Optimized Primary Keys (UUIDv7):** Instead of standard auto-incrementing integers or random UUIDs, this project uses **UUIDv7**. This modern standard embeds a Unix timestamp, ensuring that new IDs are always chronologically ordered. This provides the uniqueness of a UUID while being highly efficient for database indexing, preventing fragmentation and improving write performance.

3.  **Dynamic Query Processing:** The repository layer contains a flexible query processor that safely builds dynamic SQL to handle client-side sorting and pagination. This is achieved using parameterized queries to prevent any risk of SQL injection.

## Tech Stack

*   **Language:** Golang
*   **Core Architecture:** Hexagonal (Ports & Adapters)
*   **API Adapter:** Gin Gonic
*   **Database Adapter:** SQLite 3 (for simple, file-based persistence)
*   **Containerization:** Docker & Docker Compose (with a production-ready PostgreSQL service defined)
*   **ID Generation:** Google UUID (v7)

## Getting Started

### Prerequisites

*   Go 1.21+
*   Git

### Running Locally (with SQLite)

The application is configured to run with a file-based SQLite database by default, requiring no external services.

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/while-maybe/louder_app.git
    cd louder_app
    ```

2.  **Run the server:**
    ```bash
    go run ./cmd/server
    ```

The API will be running on `http://localhost:8080` and will create a `louder.db` file in the root directory.

### Running with Docker (PostgreSQL)

The `docker-compose.yml` is configured for a more production-like environment using PostgreSQL. *(Note: The Go application will need a new PostgreSQL adapter to be written to use this setup).*

```bash
docker-compose up -d --build
