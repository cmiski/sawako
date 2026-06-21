# Auth Service

Identity and ownership service for Sawako. Owns users, credentials, refresh tokens, projects, and API keys.

## Database

The Auth Service uses the `auth_db` PostgreSQL database created by `infrastructure/postgres`.

## Environment

Copy `.env.example` to `.env` and adjust values for your local environment.

| Variable | Description |
| --- | --- |
| `PORT` | HTTP listen port. Defaults to `8081`. |
| `DATABASE_URL` | PostgreSQL connection string for `auth_db`. |
| `JWT_SECRET` | Secret used to sign access tokens. |
| `ACCESS_TOKEN_TTL_SECONDS` | Access token lifetime in seconds. Defaults to `900`. |

## Run

Apply migrations, then start the service from the repository root:

```bash
go run ./services/auth/cmd/migrate
go run ./services/auth/cmd/server
```

## API

| Method | Path | Description |
| --- | --- | --- |
| `GET` | `/health` | Service health check. |
| `POST` | `/auth/register` | Register a user with email and password. |
| `POST` | `/auth/login` | Authenticate and receive access and refresh tokens. |
| `POST` | `/auth/refresh` | Rotate a refresh token and receive new tokens. |
| `POST` | `/auth/logout` | Revoke a refresh token and end the session. |

### Register

```json
POST /auth/register
{
  "email": "user@example.com",
  "password": "securepassword"
}
```

Returns `201 Created` on success.

### Login

```json
POST /auth/login
{
  "email": "user@example.com",
  "password": "securepassword"
}
```

Returns:

```json
{
  "access_token": "...",
  "refresh_token": "..."
}
```

### Refresh

```json
POST /auth/refresh
{
  "refresh_token": "..."
}
```

Returns a new access token and refresh token. The presented refresh token is revoked as part of rotation.

### Logout

```json
POST /auth/logout
{
  "refresh_token": "..."
}
```

Returns `204 No Content`. Revokes the presented refresh token if it is still active. Repeated logout requests are treated as successful.

## Migrations

Start PostgreSQL from the repository root:

```bash
docker compose -f infrastructure/postgres/docker-compose.yaml up -d
```

Migrations live in `services/auth/migrations/` and are tracked in the `schema_migrations` table.

## PostgreSQL Repositories

Concrete repository implementations live in `internal/postgres/`:

- `UserRepository`
- `CredentialRepository`
- `RefreshTokenRepository`
- `TransactionManager`

These implement the domain repository interfaces used by the authentication workflow.
