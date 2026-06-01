# ADR-007: JWT Authentication

## Status

Accepted

## Decision

Authentication uses:

- Access tokens
- Refresh tokens

## Alternative Considered

Server-side sessions.

## Rationale

JWT authentication works naturally with APIs, microservices, and mobile clients.

## Security Measures

- Access tokens are short-lived.
- Refresh tokens are stored as hashes only.
