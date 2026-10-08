# Research Metrics & Evaluation Formulas

## Metrics Taxonomy

| Metric Symbol | Metric Name | Unit | Source | Formula / Description |
|---|---|---|---|---|
| **$T_{p50}$** | Median Latency | ms | Prometheus / k6 | 50th percentile response duration |
| **$T_{p95}$** | 95th Percentile Latency | ms | Prometheus / k6 | 95th percentile response duration |
| **$T_{p99}$** | 99th Percentile Latency | ms | Prometheus / k6 | 99th percentile response duration |
| **$X$** | Throughput | req/sec | Prometheus | $\frac{\text{Total Completed Requests}}{\text{Duration (seconds)}}$ |
| **$E$** | Error Rate | % | k6 / Gateway Logs | $\frac{\text{Total } 5xx \text{ Responses}}{\text{Total Requests}} \times 100\%$ |
| **$L_c$** | Lock Conflict Rate | % | Redis / Booking Svc | $\frac{\text{Total } 409 \text{ Responses}}{\text{Total Lock Attempts}} \times 100\%$ |
| **$D_b$** | Double Booking Count | integer | PostgreSQL Invariant | $\sum \mathbb{I}(\text{Bookings per seat} > 1)$ — **Target: 0** |
| **$MTTR$** | Mean Time To Recovery | seconds | Kubernetes / K6 | $T_{\text{healthy}} - T_{\text{failure\_injected}}$ |
| **$Q_{lag}$** | Queue Drain Lag | messages | RabbitMQ API | $\text{Depth of } \text{queue}_{\text{notif}} \text{ over time}$ |

---

## Data Collection Procedure

1. **Automated Run Scripts**: All experiments are orchestrated by scripts in `tests/load/` and `tests/chaos/`.
2. **Raw Output Storage**: k6 outputs JSON and summary CSV files into `tests/results/raw/`.
3. **Reproducibility Guarantee**: Every test captures:
   - Git commit hash
   - Hardware specifications (CPU cores, RAM)
   - Concurrency parameter matrix
   - Raw logs from gateway, booking-service, and postgres
