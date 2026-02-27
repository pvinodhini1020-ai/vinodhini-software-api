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

# Backend Setup Guide - Vinodhini Software API

## Prerequisites
- Go 1.21 or higher
- MongoDB Atlas account
- Git

## Quick Start

### 1. Navigate to Backend Directory
```bash
cd vinodhini-software-api
```

### 2. Install Dependencies
```bash
go mod download
go mod tidy
```

### 3. Configure Environment Variables
Create a `.env` file in the root of `vinodhini-software-api`:

```env
# MongoDB Configuration
MONGODB_URI=mongodb+srv://software_developer:Vino%40102055@cluster0.xbvh07l.mongodb.net/

# Server Configuration
PORT=8080
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
JWT_EXPIRY=24h

# CORS Configuration
CORS_ORIGIN=*
```

### 4. Start the Server
```bash
go run cmd/main.go
```

The server will start on `http://localhost:8080`

## Verify Installation

### Check Server Health
```bash
curl https://vinodhini-software-api.onrender.com/api/health
```

Expected response:
```json
{"status":"ok"}
```

### Test Login Endpoint
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"vinodhini@admin.com","password":"admin123"}'
```

## Project Structure
```
vinodhini-software-api/
├── cmd/
│   └── main.go              # Application entry point
├── config/
│   ├── config.go            # Configuration loader
│   └── database.go          # Database connection
├── internal/
│   ├── controllers/         # HTTP request handlers
│   ├── services/            # Business logic
│   ├── repositories/        # Database operations
│   ├── models/              # Data structures
│   ├── middleware/          # HTTP middleware
│   └── routes/              # Route definitions
├── pkg/                     # Shared utilities
├── .env                     # Environment variables
├── go.mod                   # Go dependencies
└── go.sum                   # Dependency checksums
```

## Available API Endpoints

### Authentication
- `POST /api/auth/login` - User login
- `POST /api/auth/register` - User registration

### Users
- `GET /api/users` - List all users (Admin)
- `GET /api/users/dashboard/stats` - Dashboard statistics
- `PUT /api/users/:id` - Update user
- `DELETE /api/users/:id` - Delete user (Admin)

### Employees
- `GET /api/employees` - List employees (Admin)
- `POST /api/employees` - Create employee (Admin)
- `PUT /api/employees/:id` - Update employee (Admin)
- `DELETE /api/employees/:id` - Delete employee (Admin)

### Projects
- `GET /api/projects` - List projects
- `POST /api/projects` - Create project (Admin)
- `PUT /api/projects/:id` - Update project
- `DELETE /api/projects/:id` - Delete project (Admin)

### Service Requests
- `GET /api/service-requests` - List requests
- `POST /api/service-requests` - Create request (Client)
- `PUT /api/service-requests/:id/approve` - Approve (Admin)

### Messages
- `GET /api/messages` - List messages
- `POST /api/messages` - Send message

## Troubleshooting

### MongoDB Connection Issues
- Verify `MONGODB_URI` in `.env`
- Check MongoDB Atlas IP whitelist (add `0.0.0.0/0` for development)
- Ensure database user credentials are correct

### Port Already in Use
Change `PORT` in `.env` to another port (e.g., `8081`)

### CORS Errors
Update `CORS_ORIGIN` in `.env` to match your frontend URL

## Build for Production
```bash
go build -o vinodhini-api cmd/main.go
./vinodhini-api
```

## Run Tests
```bash
go test ./...
```

## Default Test Accounts
- **Admin**: `admin@vinodhini.com` / `admin123`
- **Employee**: `employee@vinodhini.com` / `employee123`
- **Client**: `client@vinodhini.com` / `client123`

