# API Gateway

Platform entry point for Sawako. Routes external traffic to downstream services and applies cross-cutting request controls.

## Environment

Copy `.env.example` to `.env` and adjust values for your local environment.

| Variable | Description |
| --- | --- |
| `PORT` | HTTP listen port. Defaults to `8080`. |
| `AUTH_SERVICE_URL` | Base URL for the Auth Service. Defaults to `http://localhost:8081`. |
| `EVENT_SERVICE_URL` | Base URL for the Event Service. Defaults to `http://localhost:8082`. |
| `JWT_SECRET` | Secret used to validate access tokens. Must match the Auth Service value. |

## Run

From the repository root:

```bash
go run ./gateway/cmd/server
```

## Routing

| Path prefix | Auth | Downstream service |
| --- | --- | --- |
| `/auth/*` | Public | Auth Service |
| `/events/*` | Public | Event Service |
| `/projects/*` | JWT required | Auth Service |

Examples:

- `POST /auth/register` → Auth Service (public)
- `POST /auth/login` → Auth Service (public)
- `POST /projects` → Auth Service (requires `Authorization: Bearer <access_token>`)
- `POST /events` → Event Service (when implemented)

The gateway validates JWT access tokens for protected routes, then forwards the request with an `X-User-ID` header derived from the token `sub` claim.

The gateway forwards the original request path and propagates `X-Request-ID` to downstream services.

## Health

`GET /health` returns gateway health. Downstream service health is checked at each service's own `/health` endpoint.
