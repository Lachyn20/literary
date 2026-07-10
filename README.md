# Hemra Şirow — Literary Backend API

A Clean Architecture backend for a literary and biography website dedicated to Hemra Şirow, built with Go, PostgreSQL, and full-text search capabilities.

## Technology Stack

- **Language**: Go 1.25.0
- **Framework**: chi/v5 (HTTP routing)
- **Database**: PostgreSQL 15 with pgx driver
- **Authentication**: JWT (golang-jwt)
- **Password Hashing**: bcrypt
- **Validation**: go-playground/validator
- **File Storage**: Local disk adapter with multipart support
- **API Documentation**: Swagger/OpenAPI via swaggo
- **Containerization**: Docker + docker-compose

## Project Structure

```
.
├── cmd/api/               # Application entry point
├── internal/
│   ├── domain/            # Business entities and repository interfaces
│   ├── usecase/           # Application business logic
│   ├── infrastructure/    # External adapters (DB, storage, auth)
│   ├── presentation/http/ # HTTP handlers, DTOs, middleware, routing
│   └── di/                # Dependency injection container
├── migrations/            # Database schema and migrations
├── docs/                  # Generated Swagger documentation
├── uploads/               # Local file storage directory
├── Dockerfile             # Multi-stage Docker build
├── docker-compose.yml     # PostgreSQL + API services
└── .env                   # Environment configuration
```

## Quick Start (Docker)

### Prerequisites

- Docker and docker-compose installed
- Port 8080 (API) and 5432 (PostgreSQL) available

### Run with Docker

```bash
# Build and start services
docker-compose up --build

# First time: Apply database migrations
docker-compose exec api /app/server migrate

# API is now available at http://localhost:8081
# Swagger UI at http://localhost:8081/swagger/index.html
```

> Uploaded files persist on the host at `/srv/hemra-sirow/uploads` and are served directly by Nginx. Ensure that directory exists and has the correct permissions before running `docker compose up`.

### Environment Variables

The docker-compose.yml automatically sets:

- `DB_HOST=db` (PostgreSQL service)
- `DB_PORT=5432`
- `DB_USER=postgres`
- `DB_PASSWORD=postgres`
- `DB_NAME=hemra_siirow`
- `SERVER_PORT=8081`
- `UPLOAD_BASE_PATH=/app/uploads`

File uploads persist in the `uploads` Docker volume across container restarts.

## Local Development (without Docker)

### Prerequisites

- Go 1.25.0 or later
- PostgreSQL 15 running locally
- Environment variables configured in `.env`

### Setup

1. **Create `.env` file** (see [Configuration](#configuration))

2. **Start PostgreSQL** (or use an existing instance)

3. **Apply migrations**:

   ```bash
   psql -h localhost -U postgres -d hemra_siirow -f migrations/000001_create_schema.up.sql
   ```

4. **Install dependencies**:

   ```bash
   go mod download
   ```

5. **Run the app**:

   ```bash
   go run ./cmd/api
   ```

   Server listens on `http://localhost:8081`

## Configuration

Create a `.env` file in the project root:

```env
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=hemra_siirow

# Server
SERVER_PORT=8081
JWT_SECRET=your-secret-key-here-min-32-chars

# File Storage
UPLOAD_BASE_PATH=./uploads

# CORS
CORS_ALLOWED_ORIGIN=http://localhost:3000
```

## Database Migrations

Migrations are in `migrations/` directory using standard SQL format.

**Apply migration**:

```bash
psql "postgres://postgres:postgres@localhost:5432/hemra_siirow?sslmode=disable" \
  -f migrations/000001_create_schema.up.sql
```

**Rollback migration**:

```bash
psql "postgres://postgres:postgres@localhost:5432/hemra_siirow?sslmode=disable" \
  -f migrations/000001_create_schema.down.sql
```

### Features

- Full-text search on `works` table using PostgreSQL `tsvector` and GIN index
- Automatic search vector update via database trigger
- Support for all content types: books, broadcasts, films, photos, theatre, translations, biographies

## Swagger / OpenAPI Documentation

### View API Docs

Once the server is running, access:

```
http://localhost:8081/swagger/index.html
```

### Regenerate Swagger Docs

After modifying handlers, regenerate documentation:

```bash
# Install swag CLI if not present
go install github.com/swaggo/swag/cmd/swag@latest

# Generate docs
swag init -g cmd/api/main.go
```

Docs are generated in `docs/` directory (docs.go, swagger.json, swagger.yaml).

## API Endpoints

### Works (Full-text Search)

```bash
# List all works
curl http://localhost:8081/api/works

# Search works with pagination
curl "http://localhost:8081/api/works?search=keyword&page=1&limit=10"
```

### Books (with file uploads)

```bash
# Create book
curl -X POST http://localhost:8081/api/books \
  -F "title=Book Title" \
  -F "cover=@cover.jpg" \
  -F "pdf=@book.pdf"

# Get book by ID
curl http://localhost:8081/api/books/{id}

# List books
curl http://localhost:8081/api/books
```

### Other Resources

Similar endpoints available for:

- `/api/broadcasts` (with audio/video files)
- `/api/films` (with video files)
- `/api/photos` (with image files)
- `/api/theatre-productions`
- `/api/translations`
- `/api/biographies`
- `/api/personal-letters`

## Testing

### Run All Tests

```bash
go test ./...
```

### Run Tests with Coverage

```bash
go test -cover ./...
```

### Run Specific Test

```bash
go test -run TestListWorks ./internal/usecase/work
```

### Example Unit Tests

Unit test example for work usecase:

```bash
go test -v ./internal/usecase/work
```

## File Uploads

Supported file types and size limits:

- **Images/Scans**: JPEG, PNG (max 32 MB)
- **Audio**: MP3, WAV, OGG (max 64 MB)
- **Video**: MP4, WebM, MOV (max 128 MB)

Files are stored in `uploads/` directory. In Docker, files persist in the `uploads` named volume.

## Architecture

This project follows **Clean Architecture** principles:

1. **Domain Layer** (`internal/domain/`): Business logic, entities, and repository interfaces
2. **Use Case Layer** (`internal/usecase/`): Application-specific business rules
3. **Infrastructure Layer** (`internal/infrastructure/`): External service implementations (DB, storage, auth)
4. **Presentation Layer** (`internal/presentation/http/`): HTTP handlers, DTOs, and middleware
5. **Dependency Injection** (`internal/di/`): Wires dependencies and creates service instances

Benefits:

- Easy to test (interfaces enable mocking)
- Independent of frameworks and databases
- Clear separation of concerns
- Easy to replace implementations (e.g., switch storage backend)

## Development Workflow

1. **Define business logic** in `internal/domain/`
2. **Implement use cases** in `internal/usecase/`
3. **Create repository adapters** in `internal/infrastructure/`
4. **Add HTTP handlers** in `internal/presentation/http/`
5. **Wire dependencies** in `internal/di/`
6. **Write tests** alongside implementation
7. **Update Swagger docs** via handler comments

## Troubleshooting

### Database Connection Errors

- Verify PostgreSQL is running
- Check database credentials in `.env`
- Ensure `hemra_siirow` database exists

### File Upload Issues

- Verify `uploads/` directory exists and is writable
- Check file size limits in `internal/infrastructure/storage/local.go`
- For Docker: ensure `uploads` volume is mounted

### Swagger UI Not Loading

- Verify `docs/` directory exists
- Run `swag init -g cmd/api/main.go` to regenerate
- Check that handler methods have Swagger comments

## License

This project is created for the Hemra Şirow literary collection.
