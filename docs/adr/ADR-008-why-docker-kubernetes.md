# ADR-008: Why Docker & Kubernetes for Deployment & Experimentation

## Status
Accepted

## Context
The research objective evaluates scalability, resilience, recovery time, and fault isolation under concurrent load and failure scenarios (RQ1–RQ4). Reproducibility across local and cluster environments is mandatory.

## Decision
We package every microservice as an immutable **OCI container image (Docker)** and orchestrate them via **Docker Compose** (for local development and initial baseline testing) and **Kubernetes** manifests (for replica scaling, pod failure recovery, and horizontal pod autoscaling).

## Rationale
1. **Reproducibility**: Eliminates "works on my machine" discrepancies across development, CI, and test environments.
2. **Failure Injection & Recovery (RQ2, RQ4)**: Kubernetes Pod self-healing (`restartPolicy`, liveness/readiness probes) enables precise measurement of Mean Time to Recovery (MTTR).
3. **Horizontal Pod Autoscaling (RQ1)**: Validates scalability under increasing concurrent load with multi-replica deployments (e.g. 3 replicas for booking-service).
4. **Environment Isolation**: Network policies and namespaces simulate network partitioning and fault blast radius.

## Consequences
- Requires container build pipelines for each microservice.
- Kubernetes manifest maintenance (Deployments, Services, ConfigMaps, Secrets, HPAs).
