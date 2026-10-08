# ADR-007: Why Synchronous REST vs Asynchronous Messaging

## Status
Accepted

## Context
A microservices system requires both synchronous interactions (where an immediate response is necessary for user experience) and asynchronous interactions (where tasks can be deferred or processed in parallel).

## Decision
We adopt a **hybrid communication model**:
- **Synchronous REST (HTTP/JSON)**: Used for synchronous user queries and actions requiring immediate feedback (browsing events, querying seat maps, authentication, initial booking lock).
- **Asynchronous Messaging (RabbitMQ)**: Used for post-booking state propagation, payment status side-effects, ticket issuance, and notifications.

## Rationale
1. **User Experience**: Seat locking and booking creation require immediate client confirmation (`201 Created` or `409 Conflict`).
2. **Resilience & Throughput**: Decoupling non-blocking workflows (e.g., ticket QR rendering, email dispatch) ensures the critical path remains lean and responsive.
3. **Clear Contracts**: REST endpoints use standardized OpenAPI schemas, while asynchronous events adhere to explicit JSON schema versions.

## Consequences
- Requires idempotency handling on both REST endpoints (`Idempotency-Key` header) and message consumer handlers.
