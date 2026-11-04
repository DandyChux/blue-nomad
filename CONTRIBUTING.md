# Contributing to Blue Nomad

Thank you for your interest in contributing to Blue Nomad! This guide covers everything you need to know to get set up, develop locally, and submit contributions.

---

## Table of Contents

- [Code of Conduct & Proprietary Notice](#code-of-conduct--proprietary-notice)
- [Project Architecture](#project-architecture)
- [Prerequisites](#prerequisites)
- [Getting Started](#getting-started)
    - [1. Clone the Repository](#1-clone-the-repository)
    - [2. Environment Configuration](#2-environment-configuration)
    - [3. Start the Database](#3-start-the-database)
    - [4. Database Migrations](#4-database-migrations)
    - [5. Run Development Servers](#5-run-development-servers)
- [Development Workflow](#development-workflow)
    - [Branch Naming](#branch-naming)
    - [Backend Guidelines (Go)](#backend-guidelines-go)
    - [Frontend Guidelines (SvelteKit + Bun)](#frontend-guidelines-sveltekit--bun)
    - [Testing & Quality Checks](#testing--quality-checks)
- [Commit Message Guidelines](#commit-message-guidelines)
- [Submitting a Pull Request](#submitting-a-pull-request)

---

## Code of Conduct & Proprietary Notice

This repository contains proprietary code and assets developed by Blue Nomad Labs LLC. Unauthorized copying, distribution, or public disclosure is strictly prohibited. All contributors must have authorized access and adhere to professional, respectful communication.

---

## Project Architecture

```
blue-nomad/
├── cmd/server/          # Backend entry point (main.go)
├── internal/            # Core Go application logic (handlers, middleware, services)
├── migrations/          # Goose SQL schema migrations
├── ui/                  # SvelteKit + Tailwind CSS + Bun frontend
├── air.toml             # Live-reload configuration for Go backend
├── docker-compose.yml   # Local database (PostgreSQL) services
└── Makefile             # Common dev and build tasks
```

---

## Prerequisites

Ensure the following tools are installed on your machine:

- **Go**: `1.25+`
- **Bun**: `1.x` (preferred frontend package manager and runtime)
- **Docker & Docker Compose**: for local services
- **Air**: for Go hot-reloading
    ```sh
    go install github.com/air-verse/air@latest
    ```
- **Goose**: for database migrations
    ```sh
    go install github.com/pressly/goose/v3/cmd/goose@latest
    ```

---

## Getting Started

### 1. Clone the Repository

```sh
git clone git@github.com:dandychux/blue-nomad.git
cd blue-nomad
```

### 2. Environment Configuration

Create your `.env` (or `.env.local`) file based on project defaults:

```sh
cp .env.example .env # If available, or configure the required DB and app variables
```

Make sure your database connection string matches your local setup (e.g., `DATABASE_URL=postgres://user:password@localhost:5432/dbname?sslmode=disable`).

### 3. Start the Database

Run PostgreSQL via Docker Compose:

```sh
docker compose up -d
```

### 4. Database Migrations

Apply pending migrations using Goose:

```sh
make migrate-up DATABASE_URL="postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
```

To create a new migration:

```sh
goose -dir migrations create <migration_name> sql
```

### 5. Run Development Servers

Run the backend and frontend in separate terminals:

- **Backend (with live-reload):**

    ```sh
    make start-server
    # or simply: air
    ```

- **Frontend:**
    ```sh
    cd ui
    bun install
    bun run dev
    ```

---

## Development Workflow

### Branch Naming

Create feature branches off `main` using descriptive prefixes:

- `feat/feature-name`
- `fix/bug-description`
- `refactor/component-name`
- `chore/task-name`

### Backend Guidelines (Go)

- Follow standard Go formatting conventions (`gofmt` / `goimports`).
- Keep business logic in `internal/services/` and HTTP request parsing/response logic in `internal/handlers/`.
- Always handle errors explicitly; avoid panics in application code.
- Clean up dependencies before committing:
    ```sh
    make deps
    # or: go mod tidy
    ```

### Frontend Guidelines (SvelteKit + Bun)

- Use **Svelte 5 runes** (`$state`, `$derived`, `$props`, etc.) for reactive state.
- Keep components modular and reusable under `ui/src/lib/components/`.
- Style using Tailwind CSS v4 utility classes.
- Ensure TypeScript types are accurate and strict.

### Testing & Quality Checks

Run test suites and type checks before submitting a PR:

```sh
# Backend tests
make test
# or: go test -v ./...

# Frontend type and syntax checks
cd ui && bun run check

# Full production build verification
make build
```

---

## Commit Message Guidelines

We strictly follow the **[Conventional Commits](https://www.conventionalcommits.org/)** specification.

### Format

```
<type>(<optional scope>): <description>

[optional body]

[optional footer]
```

### Types

- `feat`: A new feature
- `fix`: A bug fix
- `docs`: Documentation-only changes
- `style`: Formatting, whitespace changes (no logic changes)
- `refactor`: Code changes that neither fix a bug nor add a feature
- `perf`: Performance improvements
- `test`: Adding or correcting tests
- `chore`: Changes to build process, tooling, or dependencies

### Subject Line Rules

- Use the imperative mood (e.g., `add` instead of `added` or `adds`).
- Do not capitalize the first letter.
- Do not end with a period (`.`).
- Keep the entire subject line under 50 characters.

### Examples

```
feat(auth): add email verification handler
fix(db): correct connection pool timeout setting
docs(readme): update local setup instructions
refactor(ui): migrate nav component to svelte 5 runes
```

---

## Submitting a Pull Request

1. **Keep it focused**: Limit each PR to a single logical change or feature.
2. **Verify locally**:
    - Backend tests pass (`make test`).
    - Frontend passes type-checks (`cd ui && bun run check`).
    - The application builds cleanly (`make build`).
3. **Open PR**:
    - Provide a concise description of the changes and motivation.
    - Reference any related issues or tickets.
    - Wait for review and address any feedback promptly.
