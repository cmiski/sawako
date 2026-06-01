# ADR-004: API Gateway Technology

## Status

Accepted

## Decision

The API Gateway will be written in Go.

## Alternative Considered

Node.js.

## Rationale

The API Gateway is responsible for:

- Routing
- Authentication
- Rate limiting
- Logging

Go is a strong fit for concurrency, networking, and performance. It also introduces production-grade Go development into the project.
