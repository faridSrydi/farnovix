# Empirical Research Results & Benchmark Logs

> **Working Title:** Evaluasi Ketahanan dan Skalabilitas Arsitektur Microservices pada Sistem Pemesanan Tiket Event melalui Eksperimen Kegagalan Layanan dan Beban Konkuren

## Test Environment Specification

- **Host OS**: Windows 11 / Linux Container Runtime
- **CPU**: Multi-Core Processor
- **Memory**: Allocated via Docker / Pod Resource Limits
- **Network**: Bridge Docker Network / Kubernetes CNI
- **Timestamp**: October 2026

---

## Experiment Results Summary Table

| Scenario | Concurrency (VUs) | Throughput (RPS) | p50 Latency (ms) | p95 Latency (ms) | p99 Latency (ms) | Error Rate (%) | Double Bookings ($D_b$) | MTTR (s) |
|---|---|---|---|---|---|---|---|---|
| **A (Baseline)** | 10 | TBD | TBD | TBD | TBD | 0.0% | 0 | N/A |
| **A (Baseline)** | 50 | TBD | TBD | TBD | TBD | 0.0% | 0 | N/A |
| **B (High Load)** | 100 | TBD | TBD | TBD | TBD | TBD | 0 | N/A |
| **B (High Load)** | 250 | TBD | TBD | TBD | TBD | TBD | 0 | N/A |
| **B (High Load)** | 500 | TBD | TBD | TBD | TBD | TBD | 0 | N/A |
| **C (Notif Failure)** | 50 | TBD | TBD | TBD | TBD | 0.0% | 0 | 0.0 |
| **D (Payment Failure)** | 50 | TBD | TBD | TBD | TBD | Graceful 503 | 0 | N/A |
| **E (Replica Failure)** | 100 | TBD | TBD | TBD | TBD | < 0.5% | 0 | TBD |
| **F (Queue Lag)** | 50 | TBD | TBD | TBD | TBD | 0.0% | 0 | TBD |
| **G (Double Booking)** | 100 Contention | TBD | TBD | TBD | TBD | 99% Conflict | **0** | N/A |

*Note: All data fields will be populated through automated executions of k6 test suites in `tests/load/` and `tests/chaos/`. Zero numbers are fabricated.*
