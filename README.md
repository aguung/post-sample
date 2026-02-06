# Post Service API

A RESTful API service for managing posts and user profiles, built with Go (Gin), GORM, and PostgreSQL.

## Features

- **Authentication**: JWT-based auth with Auto-Refresh mechanism.
- **User Management**: Sign up, Sign in, Profile management.
- **Post Management**: CRUD operations for posts.
- **Database**: PostgreSQL with GORM ORM.
- **Migrations**: Managed by Atlas.
- **Docker Support**: Ready for containerization.

## Prerequisites

- Go 1.22 or higher
- PostgreSQL
- Docker (optional, for containerized run)
- [Atlas CLI](https://atlasgo.io/getting-started) (for migrations)

## Setup

1.  **Clone the repository:**

    ```bash
    git clone <repository-url>
    cd post
    ```

2.  **Environment Variables:**

    Copy the example environment file and configure it:

    ```bash
    cp .env.example .env
    ```

    Update `.env` with your database credentials and JWT secret.

3.  **Install Dependencies:**

    ```bash
    go mod download
    ```

## Database Migrations

This project uses **Atlas** for managing database migrations.

1.  **Apply Migrations:**

    ```bash
    atlas migrate apply --url "postgres://user:pass@localhost:5432/dbname?sslmode=disable"
    ```

    _Alternatively, you can use the provided Docker entrypoint which handles migrations automatically if configured._

## Running the Application

### Local Development

Run the API server:

```bash
go run main.go api
```

The server will start on port `8080` (or as configured in `.env`).

### Docker

Build and run using Docker:

```bash
# Build the image
docker build -t post-api .

# Run the container (Linux)
docker run --rm -p 8080:8080 --env-file .env --add-host=host.docker.internal:host-gateway -e DB_HOST=host.docker.internal post-api
```

> **Note for Linux Users**: The `--add-host` flag is necessary to allow the container to connect to localhost services (like PostgreSQL) on the host machine.

## Caching Strategy

The application uses an **In-Memory LRU (Least Recently Used) Cache** to optimize performance for read-heavy endpoints (e.g., retrieving all posts).

- **Implementation**: `hashicorp/golang-lru/v2`
- **Strategy**: Cache Aside
- **Invalidation**: Automatic on Create/Update/Delete operations specific to the entity.

## Authentication Mechanism

The API uses **JWT (JSON Web Token)** for authentication.

### Tokens

- **Access Token**: Short-lived token (default 15 minutes) used to access protected endpoints.
- **Refresh Token**: Long-lived token (default 7 days) used to obtain a new access token when the current one expires.

### Auto-Refresh Flow

The application implements an **Auto-Refresh** mechanism via middleware to ensure a seamless user experience.

1.  **Client Request**:
    Client sends a request with valid headers:
    - `Authorization: Bearer <access_token>`
    - `X-Refresh-Token: <refresh_token>`

2.  **Middleware Check**:
    - If the `Access Token` is valid, the request proceeds.
    - If the `Access Token` is **expired**:
      - The middleware validates the `Refresh Token` provided in the header.
        - If valid, a **New Access Token** is generated.
        - The new token is sent back in the response header `X-New-Token`.
        - The request is processed successfully using the fresh credentials.
        - **If `Refresh Token` is also expired/invalid**:
          - Server returns `401 Unauthorized` with message "Invalid or expired refresh token".
          - **Client Action**: Redirect user to Login page (Manual Re-login required).

3.  **Client Handling**:
    Client apps should check every response for the `X-New-Token` header. If present, replace the stored Access Token with the new value.

    ```javascript
    // Example Client Logic
    const response = await fetch('/api/resource', { ... });

    const newToken = response.headers.get('X-New-Token');
    if (newToken) {
        // Save new token
        localStorage.setItem('access_token', newToken);
    }
    ```

### Endpoints

- `POST /api/auth/signup`: Register a new user.
- `POST /api/auth/signin`: Login to receive Access and Refresh tokens.
