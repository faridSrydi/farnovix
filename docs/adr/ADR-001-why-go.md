# ADR-001: Why Go

## Status
Accepted

## Context
We need to select a primary programming language for the backend microservices of the Event Ticketing Platform. The system requires high concurrency support, efficient resource utilization, and fast response times for seat locking and booking operations.

## Decision
We chose **Go (Golang)** as the backend language for all microservices.

## Rationale
1. **Concurrency model**: Go's goroutines and channels provide lightweight concurrency primitives ideal for handling many concurrent booking requests.
2. **Performance**: Compiled to native binaries with low memory footprint—critical for running multiple microservices on limited hardware.
3. **Standard library**: Rich HTTP, JSON, and database support out of the box.
4. **Docker-friendly**: Single static binary builds result in minimal container images.
5. **Ecosystem**: Gin (HTTP framework), GORM (ORM), and robust testing tools.
6. **Simplicity**: Explicit error handling and minimal abstractions reduce debugging complexity in distributed systems.

## Consequences
- Team must be proficient in Go.
- No generics-heavy patterns (Go generics are still maturing).
- Excellent fit for the research focus on concurrency and resilience.
