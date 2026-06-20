# Auth Service

Identity and ownership service for Sawako. Owns users, credentials, refresh tokens, projects, and API keys.

## Database

The Auth Service uses the `auth_db` PostgreSQL database created by `infrastructure/postgres`.

## Environment

Copy `.env.example` to `.env` and adjust values for your local environment.

| Variable | Description |
| --- | --- |
| `DATABASE_URL` | PostgreSQL connection string for `auth_db`. |

## Migrations

Start PostgreSQL from the repository root:

```bash
docker compose -f infrastructure/postgres/docker-compose.yaml up -d
```

Apply migrations:

```bash
go run ./services/auth/cmd/migrate
```

Migrations live in `services/auth/migrations/` and are tracked in the `schema_migrations` table.

## PostgreSQL Repositories

Concrete repository implementations live in `internal/postgres/`:

- `UserRepository`
- `CredentialRepository`
- `RefreshTokenRepository`
- `TransactionManager`

These implement the domain repository interfaces used by the authentication workflow.
