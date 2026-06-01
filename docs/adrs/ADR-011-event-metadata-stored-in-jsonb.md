# ADR-011: Event Metadata Stored in JSONB

## Status

Accepted

## Decision

Event metadata uses PostgreSQL JSONB.

Example:

```json
{
  "country": "India",
  "browser": "Chrome"
}
```

## Alternative Considered

Rigid metadata schemas.

## Rationale

Event payloads vary significantly. JSONB provides flexibility without requiring a schema migration for every metadata shape.
