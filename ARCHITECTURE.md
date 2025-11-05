# Clean Architecture Design

## Overview

This IAM service follows Uncle Bob's Clean Architecture principles, ensuring:
- Independence of frameworks
- Testability
- Independence of UI
- Independence of database
- Independence of external agencies

## Layer Dependency Flow

```
┌─────────────────────────────────────────────────────────────┐
│                    External Systems                          │
│                  (HTTP, Database, etc.)                      │
└─────────────────────┬───────────────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────────────┐
│              Infrastructure Layer                            │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │   HTTP       │  │ Persistence  │  │   Config     │      │
│  │  Handlers    │  │ (Repositories)│  │              │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
└─────────────────────┬───────────────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────────────┐
│                 Use Case Layer                               │
│  ┌──────────────┐  ┌──────────────┐                         │
│  │ AuthUseCase  │  │ UserUseCase  │                         │
│  └──────────────┘  └──────────────┘                         │
└─────────────────────┬───────────────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────────────┐
│                   Domain Layer                               │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │   Entities   │  │  Repository  │  │ Business     │      │
│  │ (User, etc.) │  │  Interfaces  │  │    Rules     │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
└──────────────────────────────────────────────────────────────┘
```

**Rule**: Dependencies only flow inward. Inner layers know nothing about outer layers.

## SOLID Principles Application

### 1. Single Responsibility Principle (SRP)

Each component has one reason to change:

- **User Entity**: Represents user data and validation
- **AuthUseCase**: Authentication operations only
- **UserUseCase**: User management operations only
- **UserRepository**: User data persistence only
- **SessionRepository**: Session data persistence only

### 2. Open/Closed Principle (OCP)

The system is open for extension but closed for modification:

- New repository implementations can be added without modifying use cases
- New authentication methods can be added by extending AuthUseCase
- New handlers can be added without modifying existing ones

Example:
```go
// Add PostgreSQL without modifying use cases
type PostgresUserRepository struct { ... }

// Just implement the interface
func (r *PostgresUserRepository) Create(ctx context.Context, user *entity.User) error {
    // PostgreSQL implementation
}
```

### 3. Liskov Substitution Principle (LSP)

Any implementation of an interface can be substituted:

```go
// These are interchangeable
var userRepo repository.UserRepository

userRepo = persistence.NewInMemoryUserRepository()
// OR
userRepo = persistence.NewPostgresUserRepository()
// OR
userRepo = persistence.NewMongoUserRepository()

// Use case works with any implementation
authUseCase := usecase.NewAuthUseCase(userRepo, sessionRepo)
```

### 4. Interface Segregation Principle (ISP)

Interfaces are focused and minimal:

```go
// UserRepository only has methods related to user data access
type UserRepository interface {
    Create(ctx context.Context, user *entity.User) error
    GetByID(ctx context.Context, id string) (*entity.User, error)
    // ... only user-related methods
}

// SessionRepository only has methods related to session data access
type SessionRepository interface {
    Create(ctx context.Context, session *entity.Session) error
    GetByToken(ctx context.Context, token string) (*entity.Session, error)
    // ... only session-related methods
}
```

### 5. Dependency Inversion Principle (DIP)

High-level modules depend on abstractions:

```go
// Use case depends on interface, not implementation
type AuthUseCase struct {
    userRepo    repository.UserRepository    // Interface (abstraction)
    sessionRepo repository.SessionRepository // Interface (abstraction)
}

// Concrete implementations are injected from outside
authUseCase := usecase.NewAuthUseCase(userRepo, sessionRepo)
```

## Request Flow Example

### Registration Flow

```
1. HTTP Request → POST /api/auth/register
                 ↓
2. AuthHandler.Register()
                 ↓
3. AuthUseCase.Register()
                 ↓
4. entity.NewUser() (validation & password hashing)
                 ↓
5. UserRepository.Create()
                 ↓
6. InMemoryUserRepository.Create() (or any other implementation)
                 ↓
7. Response back through layers
```

### Login Flow

```
1. HTTP Request → POST /api/auth/login
                 ↓
2. AuthHandler.Login()
                 ↓
3. AuthUseCase.Login()
                 ↓
4. UserRepository.GetByEmail() or GetByUsername()
                 ↓
5. user.CheckPassword()
                 ↓
6. SessionRepository.Create()
                 ↓
7. Return session token
```

## Benefits of This Architecture

1. **Testability**: Each layer can be tested independently
   - Mock repositories for testing use cases
   - Mock use cases for testing handlers

2. **Flexibility**: Easy to change implementations
   - Switch from in-memory to PostgreSQL without changing business logic
   - Replace HTTP with gRPC without changing use cases

3. **Maintainability**: Clear separation of concerns
   - Business rules are in domain layer
   - Framework details are in infrastructure layer
   - Easy to locate and fix bugs

4. **Scalability**: Easy to extend
   - Add new features by adding new use cases
   - Add new data sources by implementing repository interfaces

5. **Independence**: Frameworks and tools are plugins
   - Not tied to any specific HTTP framework
   - Not tied to any specific database
   - Can delay decisions about external dependencies

## Adding New Features

### Example: Adding a New Authentication Method (OAuth)

1. Create new use case or extend AuthUseCase
```go
func (a *AuthUseCase) LoginWithOAuth(ctx context.Context, provider, token string) (*entity.Session, error) {
    // Implementation
}
```

2. Add new handler
```go
func (h *AuthHandler) OAuthLogin(w http.ResponseWriter, r *http.Request) {
    // Call use case
}
```

3. Register route
```go
mux.HandleFunc("/api/auth/oauth", router.authHandler.OAuthLogin)
```

No changes needed in domain layer or existing use cases!

### Example: Switching to PostgreSQL

1. Implement PostgresUserRepository
```go
type PostgresUserRepository struct {
    db *sql.DB
}

func (r *PostgresUserRepository) Create(ctx context.Context, user *entity.User) error {
    // PostgreSQL implementation
}
// ... implement other methods
```

2. Change dependency injection in main.go
```go
// Old
userRepo := persistence.NewInMemoryUserRepository()

// New
userRepo := persistence.NewPostgresUserRepository(db)
```

No changes needed in use cases or handlers!
