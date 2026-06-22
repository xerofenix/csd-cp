# User Service

Manages user authentication, profiles, and resumes for the Career Portal.

## Setup

1. Install Go 1.25 and PostgreSQL.
2. Set environment variables (see `.env.example`).
3. Run `go build -o user-service ./cmd && ./user-service`.

## Environment Variables

- `USER_SERVICE_DATABASE_URL`: PostgreSQL URL (default: `postgres://user:password@localhost:5432/career_db?sslmode=disable`)
- `USER_SERVICE_JWT_SECRET`: JWT secret key (must match API Gateway)
- `USER_SERVICE_PORT`: Server port (default: `8081`)
- `USER_SERVICE_UPLOAD_DIR`: Resume upload directory (default: `./uploads/resumes`)

## Endpoints

- `POST /register`: Create a user
- `POST /login`: Authenticate user
- `GET /users/:id`: Get user profile
- `PUT /users/:id`: Update user profile
- `POST /users/resume`: Upload resume (student only)
- `GET /health`: Health check
