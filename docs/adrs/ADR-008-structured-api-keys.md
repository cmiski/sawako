# ADR-008: Structured API Keys

## Status

Accepted

## Decision

Sawako uses structured API keys.

Format:

```text
swk_live_<key_id>.<secret>
```

Example:

```text
swk_live_kid_abcd1234.secret_xyz
```

## Alternative Considered

A single random token.

## Rationale

Structured API keys provide:

- Fast lookup
- Better observability
- Easier rotation

## Security Measures

Only secret hashes are stored. Raw API key secrets are never persisted.
