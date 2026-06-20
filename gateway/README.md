# API Gateway

Platform entry point for Sawako. Routes external traffic to downstream services and applies cross-cutting request controls.

## Environment

Copy `.env.example` to `.env` and adjust values for your local environment.

| Variable | Description |
| --- | --- |
| `PORT` | HTTP listen port. Defaults to `8080`. |
| `AUTH_SERVICE_URL` | Base URL for the Auth Service. Defaults to `http://localhost:8081`. |
| `EVENT_SERVICE_URL` | Base URL for the Event Service. Defaults to `http://localhost:8082`. |

## Run

From the repository root:

```bash
go run ./gateway/cmd/server
```

## Routing

| Path prefix | Downstream service |
| --- | --- |
| `/auth/*` | Auth Service |
| `/events/*` | Event Service |

Examples:

- `POST /auth/register` → Auth Service
- `POST /auth/login` → Auth Service
- `POST /events` → Event Service (when implemented)

The gateway forwards the original request path and propagates `X-Request-ID` to downstream services.

## Health

`GET /health` returns gateway health. Downstream service health is checked at each service's own `/health` endpoint.
