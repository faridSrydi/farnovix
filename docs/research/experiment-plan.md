# Research Experiment Plan

## Working Title
> **Evaluasi Ketahanan dan Skalabilitas Arsitektur Microservices pada Sistem Pemesanan Tiket Event melalui Eksperimen Kegagalan Layanan dan Beban Konkuren**

## Research Questions & Experimental Mapping

| ID | Research Question | Primary Scenario | Key Measured Metrics |
|----|-------------------|------------------|----------------------|
| **RQ1** | How does increasing concurrent load affect response time, throughput, and error rate? | Scenario A, Scenario B | p50, p95, p99 Latency, RPS Throughput, HTTP Error Rate |
| **RQ2** | How does failure of one microservice affect the availability and functionality of other services? | Scenario C, Scenario D | Blast radius, Error cascade rate, Service availability % |
| **RQ3** | How effective is distributed locking in preventing double booking during concurrent seat reservations? | Scenario G | Lock acquisition rate, Conflict rate (409), Duplicate booking count (Invariant: 0) |
| **RQ4** | How long does the system require to recover after a service failure? | Scenario E, Scenario F | Mean Time to Recovery (MTTR), Queue lag drain duration |

---

## Experiment Matrix

### Scenario A: Normal Baseline Load
- **Objective**: Establish baseline performance characteristics under healthy conditions.
- **Traffic**: 10 to 50 concurrent virtual users (VUs) with linear ramp-up.
- **Endpoints**: `/api/v1/events`, `/api/v1/events/:id`, `/api/v1/seats/event/:id`.

### Scenario B: High Concurrent Scalability Load (RQ1)
- **Objective**: Test system elasticity, saturation point, and latency degradation curves.
- **Traffic**: 100, 250, 500 concurrent VUs.
- **Duration**: 5 minutes per stage.
- **Tool**: k6 with Prometheus remote write.

### Scenario C: Notification Service Failure (RQ2)
- **Objective**: Verify non-critical asynchronous service failure isolation.
- **Fault Injection**: `docker stop ticketing-notification-service` or Kubernetes pod deletion with replica 0.
- **Expected Outcome**:
  - Event browsing, booking creation, and payment continue at 100% success rate.
  - Domain messages accumulate in RabbitMQ queue `notification.queue` without message loss.
  - No customer-facing error propagation.

### Scenario D: Payment Service Failure (RQ2)
- **Objective**: Verify graceful failure and isolation of critical transactional dependency.
- **Fault Injection**: Stop Payment Service during checkout workflow.
- **Expected Outcome**:
  - Event discovery and seat browsing unaffected.
  - Payment requests return structured 503/504 errors without corrupting booking states.
  - Seats remain safely in `PAYMENT_PENDING` until expiration timer releases lock.

### Scenario E: Booking Service Replica Failure & Self-Healing (RQ4)
- **Objective**: Evaluate horizontal replica failover and recovery time in Kubernetes.
- **Deployment**: 3 replicas behind ClusterIP Service.
- **Fault Injection**: Kill active pod handling traffic.
- **Expected Outcome**:
  - Load balancer redirects subsequent traffic to remaining 2 replicas.
  - Replica reconstituted within MTTR threshold (< 15 seconds).

### Scenario F: RabbitMQ Consumer Downtime & Catch-Up (RQ4)
- **Objective**: Verify message durability and consumer queue lag recovery.
- **Fault Injection**: Halt ticket and notification consumers for 120 seconds while continuing bookings.
- **Recovery Measurement**: Measure timestamp delta from consumer restart until queue depth reaches 0.

### Scenario G: Concurrent Contention on Identical Seat (RQ3)
- **Objective**: Prove zero-double-booking invariant under extreme contention.
- **Traffic**: 100 concurrent requests targeting `Event X`, `Seat A12` simultaneously.
- **Expected Invariant**:
  - Exactly **1** request receives `201 Created` with lock.
  - Exactly **99** requests receive `409 Conflict`.
  - Database verification: exactly 1 booking record exists for seat A12.
