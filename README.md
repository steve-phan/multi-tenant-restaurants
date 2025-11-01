# Multi-Tenant Restaurant Backend on AWS

A comprehensive multi-tenant SaaS platform for restaurant management built with Go, Gin, PostgreSQL, and AWS.

## Architecture Overview

### Tech Stack
- **Backend**: Go (Golang) 1.22+
- **Web Framework**: Gin v1.9+
- **Database**: PostgreSQL 15+ (via Amazon RDS)
- **ORM**: GORM v1.25+ (with PostgreSQL driver)
- **Cloud**: AWS (RDS, S3, IAM, STS, CloudWatch)
- **Tenancy Model**: Bridge Model with PostgreSQL Row Level Security (RLS) and S3 Access Points

### Core Architecture Components
- **Database**: Amazon RDS PostgreSQL with RLS for tenant isolation
- **Storage**: Amazon S3 with Access Points for tenant-specific access
- **Compute**: AWS Fargate/ECS (containerized Go application)
- **Authentication**: JWT-based authentication
- **Monitoring**: AWS CloudWatch for logs and metrics

## Project Structure

```
restaurant-backend/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── config/
│   ├── handlers/
│   ├── middleware/
│   ├── models/
│   ├── repositories/
│   └── services/
├── pkg/
│   └── utils/
├── migrations/
├── docker/
└── go.mod
```

## Getting Started

### Prerequisites
- Go 1.23 or higher
- PostgreSQL 15+ (or Amazon RDS)
- AWS CLI configured (for S3 features)
- Docker (optional, for containerization)
- Make (optional, for Makefile commands)

### Quick Start

1. **Install dependencies**:
   ```bash
   make setup
   # or
   go mod download
   ```

2. **Set up environment variables**:
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

3. **Run database migrations**:
   ```bash
   make migrate
   # or
   go run cmd/server/main.go --migrate
   ```

4. **Start the server**:
   ```bash
   make run
   # or
   go run cmd/server/main.go
   ```

   For development mode with debug logging:
   ```bash
   make run-dev
   ```

### Makefile Commands

```bash
make help          # Show all available commands
make setup         # Install dependencies
make build         # Build the application
make run           # Run the application
make run-dev       # Run in development mode
make migrate       # Run database migrations
make test          # Run tests
make test-coverage # Run tests with coverage
make clean         # Clean build artifacts
make fmt           # Format code
make docker-build  # Build Docker image
make docker-run    # Run Docker container
```

### Manual Setup

1. **Initialize Go module** (if not already done):
   ```bash
   go mod init restaurant-backend
   go mod tidy
   ```

2. **Set up environment variables**:
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

3. **Run database migrations**:
   ```bash
   go run cmd/server/main.go --migrate
   ```

4. **Start the server**:
   ```bash
   go run cmd/server/main.go
   ```

## Development

### Running Tests
```bash
go test ./...
```

### Building
```bash
go build -o bin/server cmd/server/main.go
```

### Docker Build
```bash
docker build -t restaurant-backend:latest .
```

## Deployment

See [deployment guide](./docs/deployment.md) for detailed AWS deployment instructions.

## Documentation

- [Architecture Blueprint](./A%20Comprehensive%20Architectural%20Blueprint%20for%20a%20Multi-Tenancy%20Restaurant%20Backend%20on%20AWS.md)
- [Implementation Plan](./plan.md)
