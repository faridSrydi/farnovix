# ADR-005: Why RabbitMQ for Asynchronous Messaging

## Status
Accepted

## Context
Downstream processes such as ticket PDF generation, notification delivery, and expiration monitoring must not block the critical path of booking creation and payment processing. Furthermore, if a notification or ticket service experiences a temporary failure, events must not be lost and the primary customer journey must proceed without outage.

## Decision
We use **RabbitMQ** as the asynchronous message broker with persistent queues, topic exchanges, acknowledgments (ack/nack), and Dead-Letter Queues (DLQ).

## Rationale
1. **Fault Isolation (Research Question RQ2)**: If `notification-service` is down, booking and payment flows succeed unimpeded. Events wait in RabbitMQ until the consumer recovers.
2. **Durability & Reliability**: Durable exchanges and queues ensure messages persist across broker restarts.
3. **Decoupled Architecture**: Domain events (`booking.created`, `payment.completed`, `ticket.issued`, etc.) follow the Publish/Subscribe paradigm.
4. **Dead Letter Queues (DLQ)**: Poison pill messages or persistent consumer errors are safely routed to DLQs for inspection rather than blocking the main pipeline.

## Consequences
- Requires idempotent consumers to prevent duplicate side-effects from at-least-once delivery.
- RabbitMQ connection management, channel pooling, and heartbeat monitoring must be implemented in Go consumers/publishers.
