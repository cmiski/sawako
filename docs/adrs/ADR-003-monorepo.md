# ADR-003: Monorepo Strategy

## Status

Accepted

## Decision

Sawako uses a monorepo.

```text
sawako/
|-- gateway/
|-- services/
|-- docs/
|-- infrastructure/
`-- shared/
```

## Alternative Considered

Multi-repository architecture.

## Rationale

Sawako is currently maintained by a single developer. A monorepo improves development speed, documentation consistency, learning experience, and portfolio presentation.
