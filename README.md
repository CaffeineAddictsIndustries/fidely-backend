# Fidely Backend

Go backend for the Fidely loyalty card platform.

## Tech Stack

- **Go 1.23**
- **Echo** — HTTP framework
- **pgx** — PostgreSQL driver
- **golang-migrate** — Database migrations

## Prerequisites

- Go 1.23+
- Docker Desktop (with WSL integration enabled)
- golang-migrate CLI

### Install golang-migrate

```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

## Getting Started

### 1. Start PostgreSQL

```bash
docker compose up -d
```

### 2. Run Migrations

```bash
~/go/bin/migrate -path migrations -database "postgres://fidely:fidely@host.docker.internal:5432/fidely?sslmode=disable" up
```

### 3. Start Server

```bash
DATABASE_URL="postgres://fidely:fidely@host.docker.internal:5432/fidely?sslmode=disable" go run cmd/api/main.go
```

Server runs on `http://localhost:8080`

### 4. Test Health Check

```bash
curl http://localhost:8080/health
```

## Database

- **Host:** `host.docker.internal` (from WSL) or `localhost` (from Windows)
- **Port:** 5432
- **User:** fidely
- **Password:** fidely
- **Database:** fidely

### Migration Commands

```bash
# Apply all migrations
~/go/bin/migrate -path migrations -database "postgres://fidely:fidely@host.docker.internal:5432/fidely?sslmode=disable" up

# Rollback last migration
~/go/bin/migrate -path migrations -database "postgres://fidely:fidely@host.docker.internal:5432/fidely?sslmode=disable" down 1

# Rollback all migrations
~/go/bin/migrate -path migrations -database "postgres://fidely:fidely@host.docker.internal:5432/fidely?sslmode=disable" down

# Check current version
~/go/bin/migrate -path migrations -database "postgres://fidely:fidely@host.docker.internal:5432/fidely?sslmode=disable" version
```

## Project Structure

```
fidely-backend/
├── cmd/api/
│   └── main.go                 # Entrypoint
├── internal/
│   ├── config/                 # Environment config
│   ├── db/                     # Database connection
│   ├── model/                  # Data models
│   ├── repository/             # Database queries (add as needed)
│   ├── service/                # Business logic (add as needed)
│   └── handler/                # HTTP handlers (add as needed)
├── migrations/                 # SQL migrations
├── docker-compose.yml
├── go.mod
└── go.sum
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `DATABASE_URL` | PostgreSQL connection string | (required) |
| `SERVER_PORT` | HTTP server port | 8080 |
