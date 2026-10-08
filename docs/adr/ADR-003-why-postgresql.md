# ADR-003: Why PostgreSQL

## Status
Accepted

## Context
Each microservice requires a persistent relational data store. We need ACID compliance for booking and payment transactions where data consistency is critical.

## Decision
Use **PostgreSQL 16** as the relational database for all services.

## Rationale
1. **ACID compliance**: Critical for booking transactions where double-booking prevention depends on database-level constraints.
2. **UUID support**: Native `gen_random_uuid()` for generating public identifiers.
3. **JSONB**: Flexible storage for event metadata without schema migration.
4. **Proven at scale**: Battle-tested in production systems handling millions of transactions.
5. **Ecosystem**: Excellent Go driver support (pgx via GORM).
6. **Database-per-service**: PostgreSQL supports multiple databases in a single instance for development, with logical isolation per service.

## Consequences
- Each service owns its database schema.
- Cross-service queries are prohibited; services communicate via APIs or events.
- Migration management required per service.
