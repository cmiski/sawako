# ADR-002: Database Per Service

## Status

Accepted

## Decision

Every service owns its own database.

```text
Auth Service  -> Auth Database
Event Service -> Event Database
```

## Alternative Considered

Shared database.

## Rationale

Shared databases create tight coupling because services become dependent on each other's schemas.

## Benefits

- Service autonomy
- Independent scaling
- Technology flexibility

## Tradeoffs

- Cross-service joins are not available.
- Communication must happen through APIs or messaging.
