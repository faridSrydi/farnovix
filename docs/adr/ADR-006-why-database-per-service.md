# ADR-006: Why Database-per-Service

## Status
Accepted

## Context
In microservices architecture, a shared database creates tight coupling at the data layer, breaks domain boundaries, and makes independent schema migrations or scaling impossible.

## Decision
We enforce a strict **database-per-service** pattern. Each microservice logically and physically owns its dedicated database:
- `user_db` (User Service)
- `event_db` (Event Service)
- `booking_db` (Booking Service)
- `payment_db` (Payment Service)
- `ticket_db` (Ticket Service)
- `notification_db` (Notification Service)

## Rationale
1. **True Fault Isolation**: A database bottleneck or locking contention in `booking_db` cannot freeze `user_db` or `event_db`.
2. **Encapsulated Schema Evolution**: Services can evolve their schemas and migrations without fear of breaking foreign dependencies.
3. **Boundary Enforcement**: Services MUST communicate via clean REST APIs or RabbitMQ domain events, preventing hidden dependencies.
4. **Research Validity**: To realistically test RQ2 (service failure blast radius), services must not share database connections or transactional states.

## Consequences
- No cross-database SQL joins.
- Denormalization or API aggregation needed when composite views are presented.
- Distributed data consistency managed via saga or asynchronous events.
