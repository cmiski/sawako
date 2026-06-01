# ADR-009: Event-Centric Architecture

## Status

Accepted

## Decision

Events are the central domain object in Sawako. Everything eventually derives from event data.

Future services that consume event data include:

- Search
- Analytics
- Notifications
- Monitoring

## Rationale

This creates a coherent platform vision and keeps future features aligned with the event-platform domain.
