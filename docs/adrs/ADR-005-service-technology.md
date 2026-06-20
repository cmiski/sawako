# ADR-005: Service Technology

## Status

Accepted (revised)

## Decision

Business services use Go.

The Auth Service is the first implemented business service and establishes the pattern for future services such as Event, Organization, and Notification.

## Alternatives Considered

- Node.js and Express
- NestJS
- Spring Boot

## Rationale

The API Gateway is already written in Go ([ADR-004](ADR-004-api-gateway-technology.md)). Implementing business services in Go keeps the monorepo on a single language, reuses shared packages under `shared/`, and avoids context-switching between runtimes during local development.

The original decision favored Node.js and Express for familiarity. After implementation began, Go proved a better fit for service boundaries, interface-driven design, and alignment with the gateway stack.

## Revision

| Date | Change |
| --- | --- |
| 2026-06 | Revised from Node.js and Express to Go to match the implemented Auth Service. |
