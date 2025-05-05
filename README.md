# Go/Gin Gonic Starter Template

A zero-configuration starter template for building RESTful APIs using [Gin Gonic](https://gin-gonic.com/) with:

- **JWT Authentication** for secure API access
- **Redis** for caching
- **GORM** for database interactions

## Features

- Pre-configured JWT middleware for authentication
- RBAC (Role-based access control) middleware
- Caching middleware using Redis
- GORM setup for database ORM
- RESTful route structure
- Dockerfile for containerization
- Environment variable management via `.env` file

## Prerequisites

- [Go](https://go.dev/) (>=1.20)
- [Redis](https://redis.io/) (>=6.0)
- A database supported by GORM (e.g., PostgreSQL, MySQL, SQLite)
- [Docker](https://www.docker.com/) (optional, for containerization)

## Getting Started

### Clone the Repository

```bash
git clone https://github.com/Hooannn/go-restful-starter.git
cd go-restful-starter
```

### Install Dependencies

```bash
go mod tidy
```

### Environment Variables

Create a `.env` file in the project root and configure the following variables:

```env
# Database
DATABASE_CONNECTION_STRING="host=db port=5432 user=postgres password=postgres dbname=EventPlatform sslmode=disable"

# Redis
REDIS_ADDRESS="redis:6379"
REDIS_USERNAME=""
REDIS_PASSWORD=""
REDIS_DB=0

# App settings
APP_NAME="EventPlatform"
APP_PORT=8080
APP_ENV="development"
RESET_PASSWORD_OTP_EXPIRE_MINUTES=5
DEFAULT_CACHE_EXPIRE_MINUTES=10

# Authentication
JWT_ACCESS_TOKEN_SECRET=""
JWT_REFRESH_TOKEN_SECRET=""
JWT_ACCESS_TOKEN_EXPIRE_HOURS=1
JWT_REFRESH_TOKEN_EXPIRE_HOURS=168

# Mail sender
EMAIL_SENDER=""
EMAIL_PASSWORD=""
SMTP_HOST=""
SMTP_PORT=587
```

### Run the Application

```bash
go run ./cmd
```

The server will start on `http://localhost:8080`.

## Project Structure

```plaintext
.
├── cmd                                      # Application entry point
│   └── event_platform.go
├── configs                                  # Application configs
│   └── config.go
├── internal
│   ├── constant                             # Success, error messages, redis prefix keys
│   │   └── constant.go
│   ├── redis
│   │   └── redis.go
│   ├── factory                              # Variables initialization entry point
│   │   └── factory.go
│   ├── handler                              # Controller logic
│   │   ├── auth_handler.go
│   │   └── user_handler.go
│   ├── entity                               # Models
│   │   ├── database.go
│   │   ├── permission.go
│   │   ├── role.go
│   │   └── user.go
│   ├── repository                           # Database access layer
│   │   └── user_repository.go
│   ├── middleware                           # Custom middleware (e.g., JWT)
│   │   ├── invalidate_cache.go
│   │   ├── with_cache.go
│   │   ├── with_jwt_auth_middleware.go
│   │   └── with_permissions.go
│   ├── services                             # Business logic (e.g., Auth)
│   │   ├── auth_service.go
│   │   └── user_service.go
│   ├── types
│   │   ├── auth.go
│   │   └── user.go
│   ├── util                                 # Helper functions
│   │   └── jwt.go
│   └── routes                               # API routes
│       ├── route.go
│       ├── auth_route.go
│       └── user_route.go
├── pkg                                      # Shared logic
│   └── api
│       ├── http_exception.go
│       └── http_response.go
├── docker-compose.yaml
├── docker-compose-prod.yaml
├── .env
├── Dockerfile
├── go.mod
└── go.sum
```

## Running with Docker

Run the Docker Compose stack:

```bash
docker compose -f docker-compose-prod.yml up
```

The server will start on `http://localhost:8080`.

## Dependencies

- [Gin Gonic](https://github.com/gin-gonic/gin)
- [GORM](https://github.com/go-gorm/gorm)
- [Redis](github.com/redis/go-redis/v9)
- [JWT](github.com/golang-jwt/jwt/v5)

## Todo

- Test implementation
