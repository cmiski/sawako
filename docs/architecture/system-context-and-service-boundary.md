# System Context & Service Boundary

## Table of Contents

- [Purpose](#purpose)
- [Golden Rule](#golden-rule)
- [V1 System Context](#v1-system-context)
- [Service Boundaries](#service-boundaries)
- [Service Ownership Matrix](#service-ownership-matrix)
- [Database Ownership Rules](#database-ownership-rules)
- [Communication Rules](#communication-rules)
- [Core Data Flows](#core-data-flows)
- [Future Context Map](#future-context-map)
- [Future Event Pipeline](#future-event-pipeline)
- [Architectural Invariants](#architectural-invariants)
- [Related Documentation](#related-documentation)

## Purpose

This document defines service ownership, responsibilities, forbidden responsibilities, communication rules, and data ownership boundaries for Sawako.

Use this document when adding new endpoints, tables, services, or integration paths. If a proposed change violates these boundaries, it requires an architectural decision record.

## Golden Rule

Every service must be able to answer:

> What is my job?

If the answer becomes "everything," the service is badly designed.

## V1 System Context

Sawako V1 intentionally contains only three services:

```text
Client or Application
        |
        v
API Gateway (Go)
        |
        +--> Auth Service
        |       |
        |       v
        |   Auth Database
        |
        +--> Event Service
                |
                v
            Event Database
```

The limited V1 surface area keeps the platform understandable while still teaching real service boundaries.

## Service Boundaries

### API Gateway

The API Gateway is the platform entry point. All external traffic enters through the gateway.

#### Responsibilities

| Responsibility | Description |
| --- | --- |
| Request routing | Route `/auth/*` traffic to the Auth Service and `/events/*` traffic to the Event Service. |
| JWT validation | Verify access tokens and claims for protected user-facing APIs. |
| API key validation | Validate application API keys for event ingestion. |
| Request logging | Record request metadata, response status, latency, and service routing decisions. |
| Rate limiting | Protect downstream services from abuse and accidental overload. |

#### Owns

The API Gateway owns no business data and has no database.

#### Must Never Do

The API Gateway must not perform:

- User registration
- User login
- Project creation
- Event storage
- Business logic

The API Gateway is a traffic controller, not a domain service.

### Auth Service

The Auth Service owns identity and ownership. It answers:

- Who are you?
- What do you own?

#### Owns

| Domain Object | Database Table |
| --- | --- |
| Users | `users` |
| Projects | `projects` |
| API keys | `api_keys` |
| Refresh tokens | `refresh_tokens` |

#### Responsibilities

| Capability | Example Endpoint |
| --- | --- |
| User registration | `POST /auth/register` |
| User login | `POST /auth/login` |
| Refresh tokens | `POST /auth/refresh` |
| Logout | `POST /auth/logout` |
| Project management | `POST /projects`, `GET /projects`, `PATCH /projects/:id`, `DELETE /projects/:id` |
| API key management | `POST /projects/:id/api-keys`, `DELETE /api-keys/:id` |

#### Must Never Do

The Auth Service must not:

- Store events
- Search events
- Send notifications
- Process analytics

Those responsibilities belong to other services.

### Event Service

The Event Service stores and retrieves events. It should do that job and nothing else.

#### Owns

| Domain Object | Database Table |
| --- | --- |
| Events | `events` |

#### Responsibilities

| Capability | Description |
| --- | --- |
| Event ingestion | Accept and validate `POST /events` requests. |
| Event retrieval | Serve `GET /events` requests. |
| Event filtering | Filter by date range, category, and project. |
| Pagination | Support large result sets safely. |

#### Must Never Do

The Event Service must not:

- Authenticate users
- Create projects
- Generate API keys
- Send notifications

Those responsibilities belong to the Auth Service, API Gateway, or future downstream services.

## Service Ownership Matrix

| Feature | API Gateway | Auth Service | Event Service |
| --- | --- | --- | --- |
| Login | No | Yes | No |
| JWT validation | Yes | No | No |
| API key validation | Yes | No | No |
| Create project | No | Yes | No |
| Generate API key | No | Yes | No |
| Store event | No | No | Yes |
| Query event | No | No | Yes |

This matrix is a guardrail against boundary drift.

## Database Ownership Rules

### Auth Database

Owned exclusively by the Auth Service.

Contains:

- `users`
- `projects`
- `api_keys`
- `refresh_tokens`

### Event Database

Owned exclusively by the Event Service.

Contains:

- `events`

### Forbidden Pattern

Services must never read or write another service's database.

```text
Event Service -> Auth Database    # forbidden
Auth Service  -> Event Database   # forbidden
```

Cross-service communication must happen through APIs or messaging.

## Communication Rules

### V1

V1 uses synchronous HTTP REST communication through the API Gateway.

```text
API Gateway -> Auth Service
API Gateway -> Event Service
```

### V2 and Later

Future versions introduce asynchronous communication through RabbitMQ when the platform needs event distribution, background workers, retries, and dead-letter queues.

## Core Data Flows

### User Login

```text
User
  |
  v
API Gateway
  |
  v
Auth Service
  |
  v
Auth Database
  |
  v
JWT returned
```

### Event Ingestion

```text
Application
  |
  v
API Gateway
  |
  v
API key validation
  |
  v
Event Service
  |
  v
Event Database
```

## Future Context Map

As Sawako grows, the API Gateway routes traffic to additional services:

```text
API Gateway
  |-- Auth Service
  |-- Event Service
  |-- Organization Service
  |-- Notification Service
  |-- Search Service
  |-- Analytics Service
  `-- Admin Service
```

Each service owns its API, logic, and database. There are no exceptions.

## Future Event Pipeline

In V4 and later, Sawako evolves into an event-driven platform:

```text
Application
    |
    v
API Gateway
    |
    v
Event Service
    |
    v
PostgreSQL
    |
    v
RabbitMQ
    |
    +--> Notification Service
    +--> Analytics Service
    `--> Search Service
```

This is where Sawako evolves from a microservice project into a distributed system.

## Architectural Invariants

These rules should not be broken:

1. Events are immutable.
2. Services own their data.
3. The API Gateway contains no business logic.
4. API key secrets are stored as hashes only.
5. Refresh tokens are stored as hashes only.
6. Communication happens through APIs or messaging, never shared databases.
7. Complexity is introduced only when justified.

## Related Documentation

- [Vision and Architecture Blueprint](../vision/vision-and-architecture-blueprint.md)
- [Scope Definition](../vision/scope-definition.md)
- [ADRs](../adrs/README.md)
