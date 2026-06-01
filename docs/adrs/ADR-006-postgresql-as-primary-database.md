# ADR-006: PostgreSQL as Primary Database

## Status

Accepted

## Decision

Sawako uses PostgreSQL as the primary database.

## Alternatives Considered

- MongoDB
- MySQL

## Rationale

Sawako requires relationships, transactions, JSONB, and strong consistency. PostgreSQL provides all of these capabilities.

## Future Considerations

Redis and Elasticsearch may be introduced later for caching, hot data, and search infrastructure.
