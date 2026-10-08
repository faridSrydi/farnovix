# ADR-004: Why Redis for Distributed Seat Locking

## Status
Accepted

## Context
In high-concurrency ticket sales (flash sales, high-demand concert events), thousands of users may attempt to reserve the same seat simultaneously. If concurrency control relies solely on database table/row locks, database connection pools quickly exhaust, throughput degrades severely, and deadlocks may occur.

## Decision
We use **Redis** with atomic `SET key value NX EX TTL` commands for distributed seat locking with a short time-to-live (e.g., 300 seconds / 5 minutes), supplemented by relational database transactional constraints as the ultimate source of truth.

## Rationale
1. **In-Memory Speed & Low Latency**: Redis processes locking commands in sub-millisecond time, shielding PostgreSQL from intense lock contention spikes.
2. **Atomic Mutual Exclusion**: `SET seat-lock:{event_id}:{seat_id} {booking_id} NX EX {ttl}` guarantees that exactly one client acquires the lock atomically.
3. **Automatic Expiration**: The TTL automatically releases abandoned locks without requiring manual cron sweeps or zombie lock retention.
4. **Defense in Depth**: Database constraints (`UNIQUE`, foreign keys, and status transitions) protect data integrity in case of Redis restarts or cache invalidations.

## Consequences
- Requires Redis as an operational infrastructure dependency.
- Releases and rollbacks must be handled idempotently.
- A lock acquisition failure returns `409 Conflict` immediately, keeping response times predictable.
