# Repository Blueprint

## Table of Contents

- [Purpose](#purpose)
- [Repository Philosophy](#repository-philosophy)
- [Root Structure](#root-structure)
- [Root README](#root-readme)
- [API Gateway Structure](#api-gateway-structure)
- [Services Structure](#services-structure)
- [Node.js Service Structure Standard](#nodejs-service-structure-standard)
- [Documentation Structure](#documentation-structure)
- [Infrastructure Structure](#infrastructure-structure)
- [Shared Code Policy](#shared-code-policy)
- [Scripts](#scripts)
- [Environment Strategy](#environment-strategy)
- [Naming Conventions](#naming-conventions)
- [Documentation Rules](#documentation-rules)
- [Maturity Target](#maturity-target)
- [Related Documentation](#related-documentation)

## Purpose

This document defines the repository layout and organizational conventions for Sawako. It is the structural operating model for the project.

For engineering standards, see [Sawako Engineering Standards and Development Principles](engineering-standards-and-development-principles.md).

## Repository Philosophy

The repository should:

- Scale from 3 services to 10 or more services.
- Be easy to navigate.
- Keep documentation close to implementation decisions.
- Support independently owned services.
- Remain understandable after a year of development.

A repository is not just a collection of folders. It is the operating system of the project.

## Root Structure

```text
sawako/
|-- gateway/
|-- services/
|-- docs/
|-- infrastructure/
|-- shared/
|-- scripts/
|-- .gitignore
|-- README.md
`-- LICENSE
```

The current repository uses `docs/` as the canonical documentation directory.

## Root README

The root `README.md` is the public face of Sawako. A recruiter or contributor should understand the project in two minutes.

It should contain:

- What Sawako is
- Architecture overview
- Services
- Technology stack
- How to run the project
- Roadmap
- Documentation links

## API Gateway Structure

The `gateway/` directory contains the Go API Gateway.

```text
gateway/
|-- cmd/
|-- internal/
|-- configs/
|-- tests/
|-- go.mod
`-- README.md
```

### `cmd/`

Contains application entry points only.

```text
cmd/
`-- server/
    `-- main.go
```

Startup and dependency wiring belong here. Domain logic does not.

### `internal/`

Contains gateway implementation.

```text
internal/
|-- middleware/
|-- routing/
|-- auth/
|-- ratelimit/
|-- logging/
`-- config/
```

## Services Structure

All business services live under `services/`.

Initial services:

```text
services/
|-- auth/
`-- events/
```

Expected mature services:

```text
services/
|-- auth/
|-- events/
|-- organizations/
|-- notifications/
|-- analytics/
`-- search/
```

## Node.js Service Structure Standard

Every Node.js service follows the same structure.

Example:

```text
services/auth/
|-- src/
|-- tests/
|-- migrations/
|-- package.json
|-- .env.example
`-- README.md
```

### `src/`

```text
src/
|-- controllers/
|-- services/
|-- repositories/
|-- routes/
|-- middleware/
|-- validators/
|-- models/
|-- config/
|-- utils/
`-- app.js
```

### Layer Responsibilities

| Directory | Responsibility |
| --- | --- |
| `controllers/` | Receive requests, call services, return responses. Controllers should not contain business logic. |
| `services/` | Own business logic, use-case orchestration, and policy decisions. |
| `repositories/` | Own database access only. Repositories should not contain business logic. |
| `routes/` | Map HTTP paths to controllers. |
| `middleware/` | Service-specific middleware such as error handling and request logging. |
| `validators/` | Request validation schemas such as `RegisterSchema`, `LoginSchema`, and `CreateProjectSchema`. |
| `models/` | Domain or persistence models where useful. |
| `config/` | Environment and dependency configuration such as database and JWT settings. |
| `utils/` | Small service-local utilities. Do not use this as a dumping ground. |

### `migrations/`

Database changes are implemented through migrations.

Examples:

```text
001_create_users.sql
002_create_projects.sql
003_create_api_keys.sql
```

Production schemas must not be modified manually. Every schema change goes through a migration.

### `tests/`

Service-level tests are organized by test type.

```text
tests/
|-- unit/
|-- integration/
`-- fixtures/
```

## Documentation Structure

The `docs/` directory is the knowledge base for Sawako.

Recommended structure at maturity:

```text
docs/
|-- vision/
|-- architecture/
|-- adrs/
|-- api/
|-- database/
|-- diagrams/
`-- roadmap/
```

| Directory | Contents |
| --- | --- |
| `vision/` | Vision blueprint and scope definition. |
| `architecture/` | Service boundaries, system context, and engineering principles. |
| `adrs/` | One ADR per major architecture decision. |
| `api/` | API specifications such as `auth-api.md` and `event-api.md`. |
| `database/` | Database diagrams and explanations such as `auth-db.md` and `event-db.md`. |
| `diagrams/` | Architecture diagrams such as `context-diagram.png`, `service-map.png`, and `event-flow.png`. |
| `roadmap/` | Build roadmap and milestone tracking. |

Every major decision should be recorded as an individual ADR in [Architectural Decision Records](../adrs/README.md).

## Infrastructure Structure

Initial infrastructure:

```text
infrastructure/
|-- postgres/
`-- docker/
```

Expected mature infrastructure:

```text
infrastructure/
|-- postgres/
|-- rabbitmq/
|-- redis/
|-- elasticsearch/
`-- monitoring/
```

## Shared Code Policy

The `shared/` directory is only for truly shared code and contracts.

```text
shared/
|-- contracts/
|-- constants/
`-- types/
```

Do not dump random utilities into `shared/`. Shared code should remain minimal because excessive sharing creates coupling between services.

## Scripts

The `scripts/` directory contains developer automation only.

Examples:

```text
start-dev.sh
reset-db.sh
seed-data.sh
```

Scripts should make common development tasks repeatable, but business logic must remain inside services.

## Environment Strategy

Every service provides an `.env.example` file.

Example:

```text
PORT=
DATABASE_URL=
JWT_SECRET=
```

The real `.env` file is environment-specific and must never be committed.

## Naming Conventions

### Services

Use descriptive service names:

```text
auth-service
event-service
notification-service
```

### Tables

Use plural table names:

```text
users
projects
api_keys
events
```

### Columns

Use `snake_case` column names:

```text
created_at
updated_at
project_id
```

### API Endpoints

Use RESTful endpoint naming:

```text
POST /auth/login
POST /projects
GET /events
```

## Documentation Rules

Every new service must include:

- `README.md`
- Architecture notes
- Database notes
- API notes

No undocumented services should be merged into the repository.

## Maturity Target

At maturity, the repository should resemble:

```text
sawako/
|-- gateway/
|-- services/
|   |-- auth/
|   |-- events/
|   |-- organizations/
|   |-- notifications/
|   |-- analytics/
|   `-- search/
|-- docs/
|-- infrastructure/
|   |-- postgres/
|   |-- rabbitmq/
|   |-- redis/
|   `-- elasticsearch/
|-- shared/
`-- scripts/
```

## Related Documentation

- [Vision and Architecture Blueprint](../vision/vision-and-architecture-blueprint.md)
- [System Context and Service Boundary](system-context-and-service-boundary.md)
- [Sawako Engineering Standards and Development Principles](engineering-standards-and-development-principles.md)
- [ADRs](../adrs/README.md)
