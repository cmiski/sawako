# Vision & Architecture Blueprint

## Table of Contents

- [Purpose](#purpose)
- [Product Vision](#product-vision)
- [Mission](#mission)
- [Core Architecture Principles](#core-architecture-principles)
- [Core Domain Model](#core-domain-model)
- [Technology Stack](#technology-stack)
- [Repository Strategy](#repository-strategy)
- [System Evolution Roadmap](#system-evolution-roadmap)
- [Scaling Strategy](#scaling-strategy)
- [Advanced Concepts Timeline](#advanced-concepts-timeline)
- [End-State Architecture](#end-state-architecture)
- [Related Documentation](#related-documentation)

## Purpose

This document defines the long-term product and architecture vision for Sawako. It is intended for backend engineers, technical leads, future contributors, and reviewers who need to understand what the platform is, why it exists, and how it is expected to evolve.

For implementation sequencing, see [Service-by-Service and Milestone-Based Roadmap](../roadmap/service-by-service-and-milestone-based-roadmap.md). For feature boundaries, see [Scope Definition](scope-definition.md).

## Product Vision

Sawako is an event-driven backend platform that enables applications to collect, process, search, analyze, and react to events at scale.

Applications integrated with Sawako send events such as:

```text
user_registered
user_logged_in
task_created
deployment_success
payment_failed
server_error
```

Sawako receives these events and provides the platform capabilities required to make them useful:

- Event ingestion
- Event storage
- Event querying
- Search
- Analytics
- Real-time streaming
- Notification delivery
- Monitoring and observability

## Mission

Sawako has two complementary goals.

### Product Goal

Build a developer platform capable of handling event-driven workloads similar to modern SaaS infrastructure products.

### Learning Goal

Use the platform as a structured path for mastering:

- Backend engineering
- System design
- Distributed systems
- Database design
- Messaging systems
- Search infrastructure
- Observability
- DevOps and release engineering

## Core Architecture Principles

### Events Are First-Class Citizens

Events are the central domain object in Sawako. Analytics, search, notifications, monitoring, and real-time systems all derive value from event data.

### Services Own Their Data

Each service owns its database and schema. No service directly reads from or writes to another service's database.

```text
Auth Service  -> Auth Database
Event Service -> Event Database
```

This decision is recorded in [ADR-002: Database Per Service](../adrs/ADR-002-database-per-service.md).

### Events Are Immutable

Once an event is stored, it is never modified. Corrections are represented by new events rather than updates to existing event records.

### APIs Before UI

Sawako is backend-first. The platform must remain useful without a dashboard, which means APIs, authentication, event ingestion, and retrieval semantics take priority over user interface work.

### Evolution Over Complexity

Sawako introduces complexity only when a real platform need appears. Kafka, Kubernetes, and similar technologies are not added for resume value. Every new technology must justify the operational complexity it introduces.

## Core Domain Model

| Entity | Description | Lifecycle |
| --- | --- | --- |
| User | A developer using Sawako. Users register, log in, create projects, manage API keys, and view events. | V1 |
| Project | An application connected to Sawako. Examples include a portfolio site, AI assistant, Discord bot, or task manager. | V1 |
| API Key | A credential used by external applications to send events to Sawako. | V1 |
| Event | A record describing something that happened in an integrated application. | V1 |
| Organization | A collection of users collaborating on projects. | V2 |

Example event:

```json
{
  "event_name": "user_registered",
  "category": "auth",
  "metadata": {
    "country": "India"
  }
}
```

Example structured API key:

```text
swk_live_kid_xxxx.secret_xxxx
```

## Technology Stack

| Area | Technology | Purpose |
| --- | --- | --- |
| API Gateway | Go | Routing, authentication, rate limiting, and request logging. |
| Business services | Go | Shared language with the API Gateway, strong typing, and a consistent monorepo module layout. |
| Primary database | PostgreSQL | Relationships, transactions, strong consistency, and JSONB metadata. |
| Cache and performance layer | Redis | Future caching, hot data storage, and distributed rate limiting. |
| Search infrastructure | Elasticsearch | Future full-text search and indexing. |
| Messaging | RabbitMQ | Future asynchronous processing, retries, and dead-letter queues. |
| Infrastructure | Docker and GitHub Actions | Local development, automation, CI/CD, and deployment workflows. |

## Repository Strategy

Sawako uses a monorepo so the platform can evolve with unified documentation, consistent conventions, and simple local development.

```text
sawako/
|-- gateway/
|-- services/
|-- docs/
|-- infrastructure/
|-- scripts/
`-- shared/
```

See [Repository Blueprint](../architecture/repository-blueprint.md) for detailed repository layout and service structure rules.

## System Evolution Roadmap

| Version | Theme | Primary Outcome |
| --- | --- | --- |
| V1 | Foundation | Functional event platform with API Gateway, Auth Service, and Event Service. |
| V2 | Organizations and RBAC | Team ownership, roles, permissions, and multi-tenancy. |
| V3 | Real-time infrastructure | WebSockets, live event streaming, and connection management. |
| V4 | Event-driven architecture | RabbitMQ, workers, retries, dead-letter queues, and asynchronous processing. |
| V5 | Search platform | Elasticsearch-based event search and indexing. |
| V6 | Analytics platform | Aggregations, metrics, dashboards, and trends. |
| V7 | Scalability layer | Redis caching, hot data storage, and performance optimization. |
| V8 | Observability | Structured logging, metrics, tracing, and health checks. |
| V9 | Production platform | Dockerization, CI/CD, automated testing, and deployment pipelines. |

## Scaling Strategy

| Stage | Approximate Volume | Architecture Direction |
| --- | --- | --- |
| Stage 1 | 100 events/day | PostgreSQL with simple schemas and straightforward queries. |
| Stage 2 | 10,000 events/day | Indexes and query optimization. |
| Stage 3 | 1,000,000 events/day | RabbitMQ and background workers. |
| Stage 4 | 100,000,000 events/day | Partitioned tables and dedicated search infrastructure. |
| Stage 5 | Billions of events | Distributed storage strategies and horizontal scaling. |

## Advanced Concepts Timeline

| Concept | Version |
| --- | --- |
| JWT authentication | V1 |
| Refresh tokens | V1 |
| Structured API keys | V1 |
| PostgreSQL JSONB | V1 |
| UUIDv7 identifiers | V1 |
| RBAC | V2 |
| Multi-tenancy | V2 |
| WebSockets | V3 |
| Pub/Sub | V3 |
| RabbitMQ | V4 |
| Retry queues | V4 |
| Dead-letter queues | V4 |
| Elasticsearch | V5 |
| Analytics pipelines | V6 |
| Redis | V7 |
| Distributed caching | V7 |
| Observability | V8 |
| CI/CD | V9 |

## End-State Architecture

At maturity, Sawako evolves from a three-service platform into a distributed event infrastructure system.

```text
Developer Applications
        |
        v
   API Gateway (Go)
        |
        +--> Auth Service
        +--> Event Service
        +--> Organization Service
        |
        v
   RabbitMQ Event Bus
        |
        +--> Notification Service
        +--> Search Service
        +--> Analytics Service
        +--> Future Services
        |
        v
   Redis / Elasticsearch / PostgreSQL
        |
        v
   Monitoring and Observability Stack
```

## Related Documentation

- [Scope Definition](scope-definition.md)
- [System Context and Service Boundary](../architecture/system-context-and-service-boundary.md)
- [Service-by-Service and Milestone-Based Roadmap](../roadmap/service-by-service-and-milestone-based-roadmap.md)
- [ADRs](../adrs/README.md)
