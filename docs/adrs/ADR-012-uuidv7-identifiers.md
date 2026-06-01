# ADR-012: UUIDv7 Identifiers

## Status

Accepted

## Decision

Primary entities use UUIDv7 identifiers.

Entities include:

- Users
- Projects
- API keys
- Events

## Alternatives Considered

- Auto-increment IDs
- UUIDv4

## Rationale

UUIDv7 provides globally unique identifiers that are time sortable and have better indexing characteristics than fully random UUIDs.
