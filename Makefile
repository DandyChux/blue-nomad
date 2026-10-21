# The .PHONY directive tells Make that these targets are not files
# This prevents conflicts if files with these names exist in the directory
# and ensures the targets always run when called
.PHONY: start-server start-db build clean test run migrate dev docker deps docs fmt lint

start-server:
	@echo "Starting server..."
	@air -c air.toml

start-db:
	@echo "Starting database..."
	@docker compose up -d redis

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@go mod tidy
	@go mod download

all: build

frontend:
	@echo "Building frontend..."
	cd ui && bun install && bun run build

build: frontend
	@echo "Building backend..."
	go build -o bin/blue-nomad ./cmd/server/main.go

clean:
	@echo "Cleaning..."
	rm -rf node_modules ui/build bin

# usage: make migrate-up DATABASE_URL=<your_database_url>
migrate-up:
	@goose -dir migrations postgres "$(DATABASE_URL)" up
