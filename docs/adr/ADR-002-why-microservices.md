# ADR-002: Why Microservices Architecture

## Status
Accepted

## Context
The Event Ticketing Platform serves as both a functional system and a research subject for evaluating resilience, scalability, and fault isolation. The architecture must support independent failure of components and concurrent load testing.

## Decision
We adopt a **microservices architecture** with seven independently deployable services.

## Rationale
1. **Research requirement**: The research questions (RQ1-RQ4) specifically require measurable fault isolation, independent scaling, and service-level failure experiments.
2. **Fault isolation**: A payment service failure must not crash the event browsing capability.
3. **Independent scaling**: The booking service can be scaled to multiple replicas independently.
4. **Database-per-service**: Each service owns its data, preventing tight coupling.
5. **Technology alignment**: Go microservices with Docker/Kubernetes enable realistic production deployment experiments.

## Consequences
- Higher operational complexity than a monolith.
- Requires inter-service communication patterns (REST + async messaging).
- Requires distributed tracing and correlation IDs for debugging.
- Network latency between services must be accounted for.
- Data consistency requires careful design (eventual consistency via events).

## Services
1. API Gateway
2. User Service
3. Event Service
4. Booking Service
5. Payment Service
6. Ticket Service
7. Notification Service
