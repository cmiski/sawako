# Scope Definition

## Table of Contents

- [Purpose](#purpose)
- [Core Identity](#core-identity)
- [Required V1 Features](#required-v1-features)
- [High-Priority Future Features](#high-priority-future-features)
- [Optional Features](#optional-features)
- [Explicit Non-Goals](#explicit-non-goals)
- [Stretch Goals](#stretch-goals)
- [Graduation User Journey](#graduation-user-journey)
- [Scope Governance](#scope-governance)
- [Related Documentation](#related-documentation)

## Purpose

This document defines what belongs in Sawako and what does not. Every future feature request should be evaluated against this scope before implementation.

For system architecture, see [Vision and Architecture Blueprint](vision-and-architecture-blueprint.md). For implementation order, see [Service-by-Service and Milestone-Based Roadmap](../roadmap/service-by-service-and-milestone-based-roadmap.md).

## Core Identity

Sawako is:

- An event platform
- Developer infrastructure
- An observability foundation

Sawako is not:

- A social network
- A chat application
- An e-commerce platform
- A CMS
- A project management tool

This distinction is important because Sawako should stay focused on event ingestion, event processing, and platform infrastructure.

## Required V1 Features

These capabilities are non-negotiable. Without them, Sawako is not yet a functional event platform.

| Category | Required Capabilities | Rationale |
| --- | --- | --- |
| Authentication | User registration, login, JWT authentication, refresh tokens, logout, password hashing. | Every future service depends on identity. |
| Projects | Create, update, delete, and list projects. | Projects are the ownership boundary for applications. |
| API Keys | Generate, revoke, rotate, and track last usage. | External applications authenticate through API keys. |
| Event Ingestion | Accept, validate, store, and query events. | Event ingestion is the core value proposition. |
| Event Retrieval | Filter by project, time range, and category with pagination. | Stored events must be usable and queryable. |

## High-Priority Future Features

These features are expected after V1, but they should not block the initial platform.

| Version | Capability | Included Concepts |
| --- | --- | --- |
| V2 | Organizations | Organizations, members, roles, permissions, and multi-tenancy. |
| V3 | Real-time streaming | WebSockets, live events, connection lifecycle, and presence. |
| V4 | Messaging infrastructure | RabbitMQ, workers, retry queues, and dead-letter queues. |
| V5 | Search | Elasticsearch, autocomplete, and full-text search. |
| V6 | Analytics | Metrics, aggregations, dashboards, and trends. |

## Optional Features

These features are useful but not essential. They should be considered only after the major platform versions are complete.

| Feature | Examples | Priority |
| --- | --- | --- |
| Notification Service | Email alerts, Slack alerts, Discord alerts, webhook alerts. | V6 or later. |
| Public SDKs | Node.js, Go, and Python SDKs. | Useful, but not required initially. |
| Billing | Free tier, usage tracking, plans. | Interesting, but not required for the learning goal. |
| Admin Dashboard | Internal management UI. | Useful, but backend work comes first. |
| Event Replay | Reprocessing historical events. | Advanced capability for later platform maturity. |

## Explicit Non-Goals

The following features are intentionally out of scope:

- Chat systems
- Social features such as likes, followers, comments, and posts
- Video processing
- File storage platform capabilities
- Generic AI chatbot functionality
- E-commerce features such as orders, carts, payments, and inventory

These domains do not support Sawako's purpose as an event platform.

## Stretch Goals

These capabilities would make Sawako significantly more advanced, but they should only be considered after V9.

| Capability | Description |
| --- | --- |
| Multi-region event processing | Distributed processing, replication, geo-routing, and regional failover. |
| Event replay engine | Replay events from a defined time window, such as the last seven days. |
| Audit Log Service | Make every administrative and user action traceable. |
| Event Schema Registry | Support event validation through versioned schemas. |
| Workflow automation engine | Enable rules such as `IF deployment_failed THEN send notification`. |

## Graduation User Journey

When Sawako is considered complete, the core user journey should work end to end:

```text
Developer creates account
        |
        v
Creates project
        |
        v
Generates API key
        |
        v
Application sends events
        |
        v
Sawako stores events
        |
        v
RabbitMQ distributes events
        |
        v
Search indexes events
        |
        v
Analytics processes events
        |
        v
Notifications trigger
        |
        v
Real-time dashboard updates
```

## Scope Governance

Every proposed feature must answer three questions:

1. Does this help Sawako become a better event platform?
2. Does this belong in the current phase of the roadmap?
3. Does this introduce complexity that is justified by a real platform need?

Features that fail this test should be deferred or rejected.

## Related Documentation

- [Vision and Architecture Blueprint](vision-and-architecture-blueprint.md)
- [Service-by-Service and Milestone-Based Roadmap](../roadmap/service-by-service-and-milestone-based-roadmap.md)
- [Sawako Engineering Standards and Development Principles](../architecture/engineering-standards-and-development-principles.md)
