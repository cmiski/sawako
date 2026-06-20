# Sawako Engineering Standards & Development Principles

## Table of Contents

- [Purpose](#purpose)
- [Why This Document Matters](#why-this-document-matters)
- [Engineering Principles](#engineering-principles)
- [Technology Introduction Standard](#technology-introduction-standard)
- [Platform Test](#platform-test)
- [Production Thinking](#production-thinking)
- [Engineering Philosophy](#engineering-philosophy)
- [Related Documentation](#related-documentation)

## Purpose

This document defines how Sawako is built. It does not define the platform architecture or feature scope; those are covered in [Vision and Architecture Blueprint](../vision/vision-and-architecture-blueprint.md) and [Scope Definition](../vision/scope-definition.md).

These standards apply to every service from the first implementation through future platform maturity.

## Why This Document Matters

Projects become difficult to maintain when each service develops its own style:

```text
Auth Service         -> one coding style
Event Service        -> another coding style
Notification Service -> another coding style
```

Sawako should remain understandable after months of development. Consistent engineering standards are how that happens.

## Engineering Principles

### 1. Prefer Clean Architecture Over Quick Hacks

Services must use clear layers rather than placing database access directly behind routes.

Incorrect:

```text
Route -> Database
```

Expected:

```text
Route -> Handler -> Service -> Repository -> Database
```

Every service should follow this structure unless an ADR explicitly approves a different approach.

### 2. Business Logic Never Lives in Routes

Routes should remain thin. They connect HTTP paths to controllers and should not contain long business workflows.

Incorrect:

```go
r.Post("/login", func(w http.ResponseWriter, r *http.Request) {
    // 100 lines of authentication logic
})
```

Expected flow:

```text
Route -> Handler -> authentication.Service.Login()
```

### 3. Each Layer Has One Responsibility

| Layer | Responsibilities | Not Responsible For |
| --- | --- | --- |
| Handler | Request handling, response shaping, validation error mapping. | Business rules and database access. |
| Service | Business logic and use-case orchestration. | Raw SQL, HTTP routing, or response formatting. |
| Repository | Database operations. | Business rules and request handling. |

Example service-level questions:

- Can this user create a project?
- Can this API key be rotated?
- Is this event valid for this project?

### 4. Avoid Magic Values

Important values should be named constants or configuration values.

Incorrect:

```go
if tokenExpiresIn > 900 {
    // ...
}
```

Expected:

```go
const accessTokenExpirySeconds = 900
```

### 5. Use Environment-Driven Configuration

Secrets and environment-specific values must not be hard-coded.

Incorrect:

```go
dbPassword := "mypassword"
```

Expected:

```text
DATABASE_URL=
JWT_SECRET=
```

Each service must provide an `.env.example` file. Real `.env` files must not be committed.

### 6. Handle Errors Explicitly

Empty catch blocks and generic error responses hide production issues.

Expected error categories include:

- Validation errors
- Authentication errors
- Authorization errors
- Database errors
- Dependency errors
- Unexpected internal errors

Different errors should produce different responses and logs.

### 7. Use Structured Logging

Avoid unstructured logs such as:

```go
log.Println("Something happened")
```

Prefer structured logs:

```json
{
  "service": "auth",
  "action": "login",
  "userId": "123"
}
```

Structured logs become critical when Sawako reaches multiple services and asynchronous workers.

### 8. Design APIs Before Implementing Them

Before writing endpoint code, define:

- Request shape
- Response shape
- Status codes
- Error responses
- Authentication requirements

API design should lead implementation, not follow it.

### 9. Every Service Gets Its Own README

Each service must include a `README.md` explaining:

- Purpose
- Responsibilities
- Local setup
- Environment variables
- API surface
- Database ownership
- Tests

Examples:

```text
gateway/README.md
services/auth/README.md
services/events/README.md
```

### 10. Documentation Is Part of the Product

Sawako must maintain documentation continuously, not after implementation is complete.

Required documentation categories:

- Architecture documentation
- ADRs
- Database documentation
- API documentation
- Service-level README files

### 11. Build for Learning, Not Buzzwords

Sawako should not add Kafka, Kubernetes, GraphQL, or similar technologies because they look impressive. A technology is introduced only when the existing architecture has a real problem it can solve.

This mirrors how mature engineering organizations manage complexity.

### 12. Every New Technology Must Answer Three Questions

Before introducing new technology, document:

| Question | Required Answer |
| --- | --- |
| Problem | What problem does it solve? |
| Alternative | What other solutions were considered? |
| Tradeoff | What complexity does it introduce? |

Example:

| Question | Redis Example |
| --- | --- |
| Problem | Repeated database queries on hot paths. |
| Alternative | PostgreSQL-only implementation. |
| Tradeoff | Cache invalidation complexity. |

If these questions cannot be answered, the technology is rejected.

### 13. Every Feature Must Pass the Platform Test

Ask:

> Does this help Sawako become a better event platform?

Examples:

| Feature | Result | Reason |
| --- | --- | --- |
| Search Service | Pass | Improves event discovery. |
| RabbitMQ | Pass | Enables asynchronous event processing. |
| Social feed | Fail | Not part of the event platform domain. |
| Chat system | Fail | Not part of the event platform domain. |

### 14. Think About Production From Day One

Even with one user, Sawako should be designed with production behavior in mind:

- What happens at one million users?
- What happens when a service fails?
- What happens when events grow quickly?
- What happens when logs must explain an incident?

The goal is not premature scale. The goal is engineering judgment.

### 15. Introduce Complexity Incrementally

Sawako should not jump from a basic server directly to Kafka and Kubernetes.

Expected progression:

```text
Simple -> Working -> Reliable -> Scalable
```

Every stage must work before the next stage begins.

## Technology Introduction Standard

Any proposal for a new language, framework, infrastructure component, or major library should be documented in an ADR if it changes architecture, ownership, deployment, or operations.

See [Architectural Decision Records](../adrs/README.md) for the existing decision record format.

## Platform Test

The platform test is the primary scope filter:

> Does this make Sawako a better event platform?

If the answer is no, the feature belongs outside the project.

## Production Thinking

Production-grade thinking means each service should consider:

- Reliability
- Observability
- Security
- Data ownership
- Failure modes
- Operational complexity
- Local developer experience

## Engineering Philosophy

Build the simplest thing that solves today's problem while leaving room for tomorrow's growth.

This philosophy separates well-designed systems from overengineered portfolio projects.

## Related Documentation

- [Scope Definition](../vision/scope-definition.md)
- [Repository Blueprint](repository-blueprint.md)
- [Service-by-Service and Milestone-Based Roadmap](../roadmap/service-by-service-and-milestone-based-roadmap.md)
- [ADRs](../adrs/README.md)
