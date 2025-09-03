# API Gateway

The API Gateway for the Career Portal project, built with GoFiber.

## Setup

1. Install Go 1.25.
2. Set environment variables (see `.env.example`).
3. Run `go build -o api-gateway ./cmd && ./api-gateway`.

## Environment Variables

- `API_GATEWAY_USER_SERVICE_URL`: URL of User Service (default: `http://user-service:8081`)
- `API_GATEWAY_JWT_SECRET`: JWT secret key
- ...

## Endpoints

- `POST /api/login`: Authenticate user
- `GET /health`: Health check
- `GET /metrics`: Prometheus metrics
