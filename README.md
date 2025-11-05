# IAM Service

Identity and Access Management service built with Go, following Clean Architecture and SOLID principles.

## Architecture

This service implements Clean Architecture with the following layers:

### 1. Domain Layer (`internal/domain/`)
- **Entities**: Core business objects (User, Session)
- **Repository Interfaces**: Abstractions for data access (following Dependency Inversion Principle)

### 2. Use Case Layer (`internal/usecase/`)
- **AuthUseCase**: Handles authentication business logic (register, login, logout, validate session)
- **UserUseCase**: Handles user management business logic (CRUD operations)
- Each use case follows Single Responsibility Principle

### 3. Infrastructure Layer (`internal/infrastructure/`)
- **Persistence**: Repository implementations (currently in-memory, easily replaceable)
- **HTTP**: HTTP handlers and router (Interface Adapters)
- **Config**: Configuration management

### 4. Framework Layer (`cmd/api/`)
- Application entry point
- Dependency injection and wiring

## SOLID Principles

### Single Responsibility Principle (SRP)
- Each struct/interface has one reason to change
- `AuthUseCase` only handles authentication logic
- `UserUseCase` only handles user management logic
- Each repository handles only one entity type

### Open/Closed Principle (OCP)
- Repository interfaces allow extension without modification
- New repository implementations (e.g., PostgreSQL, MongoDB) can be added without changing use cases

### Liskov Substitution Principle (LSP)
- Any implementation of `UserRepository` can replace `InMemoryUserRepository`
- Any implementation of `SessionRepository` can replace `InMemorySessionRepository`

### Interface Segregation Principle (ISP)
- Repository interfaces are focused and minimal
- Clients only depend on methods they use

### Dependency Inversion Principle (DIP)
- Use cases depend on repository interfaces (abstractions), not concrete implementations
- High-level modules don't depend on low-level modules

## Project Structure

```
iam_service/
├── cmd/
│   └── api/              # Application entry point
│       └── main.go
├── internal/
│   ├── domain/           # Business logic layer
│   │   ├── entity/       # Core entities
│   │   └── repository/   # Repository interfaces
│   ├── usecase/          # Application business logic
│   ├── infrastructure/   # External concerns
│   │   ├── config/       # Configuration
│   │   ├── http/         # HTTP handlers & router
│   │   └── persistence/  # Repository implementations
└── pkg/
    └── logger/           # Shared utilities

```

## Getting Started

### Prerequisites
- Go 1.21 or higher

### Installation

1. Clone the repository:
```bash
git clone https://github.com/team-nino/iam_service.git
cd iam_service
```

2. Install dependencies:
```bash
make deps
```

3. Build the application:
```bash
make build
```

4. Run the application:
```bash
make run
```

The service will start on `http://localhost:8080`

## API Endpoints

### Health Check
- `GET /health` - Check service health

### Authentication
- `POST /api/auth/register` - Register a new user
  ```json
  {
    "email": "user@example.com",
    "username": "username",
    "password": "password"
  }
  ```

- `POST /api/auth/login` - Login
  ```json
  {
    "email_or_username": "user@example.com",
    "password": "password"
  }
  ```

- `POST /api/auth/logout` - Logout
  - Requires `Authorization` header with session token

### User Management
- `GET /api/users` - List users (with pagination)
  - Query params: `offset`, `limit`
- `GET /api/users/get?id={user_id}` - Get user by ID

## Development

### Available Make Commands

```bash
make help      # Show available commands
make build     # Build the application
make run       # Run the application
make test      # Run tests
make fmt       # Format code
make vet       # Run go vet
make lint      # Run linter (requires golangci-lint)
make clean     # Clean build artifacts
make deps      # Download dependencies
```

### Running Tests
```bash
make test
```

### Code Quality
```bash
make fmt    # Format code
make vet    # Static analysis
make lint   # Comprehensive linting
```

## Configuration

Configuration is loaded from environment variables:

- `SERVER_PORT` - Server port (default: 8080)
- `SERVER_HOST` - Server host (default: 0.0.0.0)
- `DB_DRIVER` - Database driver (default: memory)
- `DB_DSN` - Database connection string

## Future Enhancements

- [ ] Add PostgreSQL/MySQL repository implementations
- [ ] Implement JWT-based authentication
- [ ] Add role-based access control (RBAC)
- [ ] Add OAuth2/OpenID Connect support
- [ ] Add API rate limiting
- [ ] Add comprehensive unit and integration tests
- [ ] Add API documentation with Swagger/OpenAPI
- [ ] Add Docker support
- [ ] Add CI/CD pipeline
- [ ] Add monitoring and observability

## License

MIT License

