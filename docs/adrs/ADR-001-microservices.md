# ADR-001: Microservices Architecture

## Status

Accepted

## Decision

Sawako will use a microservices architecture.

Initial services:

- API Gateway
- Auth Service
- Event Service

## Alternatives Considered

| Alternative | Description |
| --- | --- |
| Monolith | One application with one database. |
| Modular monolith | One application with separate internal modules. |

## Rationale

Sawako's goals include learning distributed systems, service boundaries, independent deployment, and independent scaling. A monolith would be simpler but would not provide the same learning opportunities.

## Tradeoffs

| Benefits | Costs |
| --- | --- |
| Clear ownership boundaries. | More operational complexity. |
| Independent scaling. | Network communication between services. |
| Real-world architecture practice. | Multiple deployments and service coordination. |
