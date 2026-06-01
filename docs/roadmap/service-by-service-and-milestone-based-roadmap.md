# Service-by-Service & Milestone-Based Roadmap

## Table of Contents

- [Purpose](#purpose)
- [Execution Rule](#execution-rule)
- [Phase 0: Foundation](#phase-0-foundation)
- [Phase 1: API Gateway](#phase-1-api-gateway)
- [Phase 2: Auth Service](#phase-2-auth-service)
- [Phase 3: Event Service](#phase-3-event-service)
- [Phase 4: Organizations and RBAC](#phase-4-organizations-and-rbac)
- [Phase 5: Real-Time Infrastructure](#phase-5-real-time-infrastructure)
- [Phase 6: RabbitMQ and Event Bus](#phase-6-rabbitmq-and-event-bus)
- [Phase 7: Notification Service](#phase-7-notification-service)
- [Phase 8: Search Service](#phase-8-search-service)
- [Phase 9: Analytics Service](#phase-9-analytics-service)
- [Phase 10: Redis and Performance](#phase-10-redis-and-performance)
- [Phase 11: Observability](#phase-11-observability)
- [Phase 12: DevOps and Production](#phase-12-devops-and-production)
- [Final Architecture](#final-architecture)
- [Related Documentation](#related-documentation)

## Purpose

This document defines the implementation roadmap for Sawako. It is intentionally milestone-based so each phase produces a working system and teaches a specific engineering capability.

For feature scope, see [Scope Definition](../vision/scope-definition.md). For architecture boundaries, see [System Context and Service Boundary](../architecture/system-context-and-service-boundary.md).

## Execution Rule

Do not move to the next milestone until the current milestone is implemented and understood.

Sawako should evolve through this sequence:

```text
Simple -> Working -> Reliable -> Scalable
```

## Phase 0: Foundation

### Goal

Create the engineering foundation before writing business logic.

| Milestone | Name | Deliverables | Learning Outcome |
| --- | --- | --- | --- |
| 0.1 | Project initialization | Monorepo, folder structure, Git setup, documentation structure. | Project organization and engineering workflow. |
| 0.2 | Development environment | Go setup, Node.js setup, PostgreSQL setup, Docker setup. | Local development environment design. |
| 0.3 | Architecture documentation | Blueprint, ADRs, context diagram, database diagrams. | Architecture thinking and decision recording. |

## Phase 1: API Gateway

### Goal

Build the first production-style service and learn API gateway fundamentals.

| Milestone | Name | Deliverables | Learning Outcome |
| --- | --- | --- | --- |
| 1.1 | Basic Go HTTP server | Gateway process starts and serves requests. | Go project structure, routing, middleware. |
| 1.2 | Request routing | `/auth/*` and `/events/*` routing works. | Reverse proxy concepts and request forwarding. |
| 1.3 | Request logging | Structured request logs. | Middleware and observability basics. |
| 1.4 | JWT verification | Gateway validates access tokens. | Token validation, claims, authorization. |
| 1.5 | Rate limiting | Rate limiting protects downstream services. | Token bucket and sliding-window strategies. |

### Result

A production-style API Gateway that routes traffic and enforces platform-level request controls.

## Phase 2: Auth Service

### Goal

Build the identity and ownership platform.

| Milestone | Name | Deliverables | Learning Outcome |
| --- | --- | --- | --- |
| 2.1 | Express architecture | Controllers, services, repositories, routes, validators. | Layered service architecture. |
| 2.2 | PostgreSQL integration | Connection pooling, queries, migrations. | Relational persistence in a service boundary. |
| 2.3 | User registration | Registration endpoint with validation and password hashing. | Input validation and credential storage. |
| 2.4 | Login | Login endpoint and authentication flow. | Identity verification. |
| 2.5 | JWT system | Access token creation and validation claims. | Expiration, claims, and stateless auth. |
| 2.6 | Refresh tokens | Refresh token persistence and rotation. | Session management and token lifecycle. |
| 2.7 | Project management | Project CRUD under authenticated ownership. | Authorization and ownership checks. |
| 2.8 | Structured API keys | Key IDs, secret hashing, usage tracking, rotation. | Application authentication and secret handling. |

### Result

A complete identity platform for users, projects, API keys, and refresh tokens.

## Phase 3: Event Service

### Goal

Build the core event platform.

| Milestone | Name | Deliverables | Learning Outcome |
| --- | --- | --- | --- |
| 3.1 | Event database design | Event schema, JSONB metadata, UUIDv7 IDs, indexes. | Flexible event modeling and indexing. |
| 3.2 | Event ingestion API | Validated `POST /events` endpoint. | API design and request validation. |
| 3.3 | API key authentication flow | Events accepted only from valid application credentials. | Application authentication. |
| 3.4 | Event retrieval | Filtering, sorting, and pagination. | Query design for product APIs. |
| 3.5 | Query optimization | Index review and query-plan analysis. | Performance tuning. |

### Result

A functional event platform capable of accepting, storing, and querying events.

## Phase 4: Organizations and RBAC

### Goal

Introduce multi-tenant architecture.

### Learning Outcomes

- Organizations
- Membership
- Roles
- Permissions
- Team ownership
- Access-control boundaries

### Result

Sawako begins to resemble an enterprise-grade platform with team-based ownership.

## Phase 5: Real-Time Infrastructure

### Goal

Support live event streaming.

### Learning Outcomes

- WebSockets
- Connection lifecycle management
- Presence
- Pub/Sub patterns

### Result

A live event feed for connected clients.

## Phase 6: RabbitMQ and Event Bus

### Goal

Introduce distributed-system fundamentals through asynchronous processing.

### Learning Outcomes

- Queues
- Exchanges
- Routing keys
- Consumers
- Retries
- Dead-letter queues

### Result

A true event-driven architecture with asynchronous consumers.

## Phase 7: Notification Service

### Goal

Consume events from RabbitMQ and notify external systems.

### Examples

- Email
- Slack
- Discord
- Webhooks

### Learning Outcomes

- Asynchronous workflows
- Delivery retries
- Failure handling

## Phase 8: Search Service

### Goal

Build search infrastructure for large event volumes.

### Learning Outcomes

- Elasticsearch
- Indexing
- Reindexing
- Autocomplete

### Result

Search across millions of events.

## Phase 9: Analytics Service

### Goal

Transform raw events into insights.

### Examples

- Events per hour
- Active users
- Error rates

### Learning Outcomes

- Aggregation
- Metrics
- Time-series thinking

## Phase 10: Redis and Performance

### Goal

Optimize high-traffic paths.

### Learning Outcomes

- Caching
- Distributed caching
- Hot data
- Rate-limit storage

## Phase 11: Observability

### Goal

Monitor Sawako itself.

### Learning Outcomes

- Logs
- Metrics
- Tracing
- Health checks

## Phase 12: DevOps and Production

### Goal

Make Sawako production-ready.

### Learning Outcomes

- Docker
- Docker Compose
- CI/CD
- GitHub Actions
- Deployment

## Final Architecture

```text
Client Applications
        |
        v
   API Gateway (Go)
        |
        +--> Auth Service
        +--> Event Service
        +--> Organization Service
        |
        v
     RabbitMQ
        |
        +--> Notification Service
        +--> Search Service
        +--> Analytics Service
        |
        v
 PostgreSQL / Redis / Elasticsearch
        |
        v
 Monitoring / CI/CD / Docker
```

## Related Documentation

- [Vision and Architecture Blueprint](../vision/vision-and-architecture-blueprint.md)
- [Scope Definition](../vision/scope-definition.md)
- [Repository Blueprint](../architecture/repository-blueprint.md)
- [Sawako Engineering Standards and Development Principles](../architecture/engineering-standards-and-development-principles.md)
