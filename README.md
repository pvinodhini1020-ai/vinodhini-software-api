# Vinodhini Software API

A production-ready, scalable REST API built with Golang, Gin framework, and Clean Architecture pattern.

## Features

- ✅ Clean Architecture Pattern
- ✅ Role-Based Authentication (Admin, Employee, Client)
- ✅ JWT Token Authentication
- ✅ PostgreSQL with GORM ORM
- ✅ Password Hashing (bcrypt)
- ✅ Request Validation
- ✅ Pagination, Filtering & Search
- ✅ Centralized Error Handling
- ✅ Logging Middleware
- ✅ Rate Limiting
- ✅ CORS Configuration
- ✅ Security Headers
- ✅ Graceful Shutdown
- ✅ Docker Support
- ✅ Swagger Documentation

## Project Structure

```
.
├── cmd/                    # Application entry points
├── config/                 # Configuration files
├── internal/
│   ├── controllers/       # HTTP handlers
│   ├── services/          # Business logic
│   ├── repositories/      # Data access layer
│   ├── models/            # Domain models & DTOs
│   ├── middleware/        # HTTP middleware
│   └── routes/            # Route definitions
├── pkg/
│   └── utils/             # Utility functions
├── docs/                  # Swagger documentation
├── tests/                 # Test files
├── Dockerfile
├── docker-compose.yml
└── .env.example
```

## Prerequisites

- Go 1.21+
- PostgreSQL 15+
- Docker & Docker Compose (optional)

## Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd vinodhini-software-api
```

2. Copy environment file:
```bash
copy .env.example .env
```

3. Update `.env` with your configuration

4. Install dependencies:
```bash
go mod download
```

5. Run the application:
```bash
go run cmd/main.go
```

## Docker Deployment

```bash
docker-compose up -d
```

## API Endpoints

### Authentication
- `POST /api/auth/register` - Register new user
- `POST /api/auth/login` - Login user

### Users (Protected)
- `GET /api/users` - List users (Admin only)
- `GET /api/users/:id` - Get user by ID
- `PUT /api/users/:id` - Update user
- `DELETE /api/users/:id` - Delete user (Admin only)

### Projects (Protected)
- `POST /api/projects` - Create project (Admin only)
- `GET /api/projects` - List projects
- `GET /api/projects/:id` - Get project by ID
- `PUT /api/projects/:id` - Update project (Admin/Employee)
- `DELETE /api/projects/:id` - Delete project (Admin only)
- `POST /api/projects/:id/assign` - Assign employees (Admin only)
- `GET /api/projects/:project_id/messages` - Get project messages

### Service Requests (Protected)
- `POST /api/service-requests` - Create request (Client only)
- `GET /api/service-requests` - List requests
- `GET /api/service-requests/:id` - Get request by ID
- `PUT /api/service-requests/:id` - Update request (Admin/Employee)
- `DELETE /api/service-requests/:id` - Delete request (Admin only)

### Messages (Protected)
- `POST /api/messages` - Create message
- `GET /api/messages/:id` - Get message by ID
- `DELETE /api/messages/:id` - Delete message

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| SERVER_PORT | Server port | 8080 |
| DB_NAME | Database name | software_developer |
| JWT_SECRET | JWT secret key | your-secret-key |
| JWT_EXPIRY | JWT expiration | 24h |

## Testing

```bash
go test ./tests/...
```

## License

MIT
