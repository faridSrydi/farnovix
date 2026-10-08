# System Context Diagram (C4 Model - Level 1)

## Overview

The Event Ticketing Platform is a high-concurrency microservices system for browsing events, reserving seats via distributed locking, and processing bookings and payments with asynchronous ticket issuance and notifications.

```mermaid
C4Context
    title System Context Diagram - Event Ticketing Platform

    Person(customer, "Customer", "Browses events, reserves seats, pays for tickets, receives e-tickets")
    Person(admin, "Administrator", "Manages events, venues, ticket quotas, and views platform metrics")

    System_Boundary(platform_boundary, "Event Ticketing Platform") {
        System(ticketing_system, "Event Ticketing Microservices", "Handles authentication, catalog, seat reservation, booking workflow, payments, tickets, and notifications")
    }

    System_Ext(mock_payment, "Sandbox Payment Provider", "Simulates credit card/bank transfer checkout, callbacks, and webhooks")
    System_Ext(email_service, "SMTP / Notification Gateway", "Simulates email/SMS delivery for tickets and receipts")

    Rel(customer, ticketing_system, "Browses events, reserves seats, checks out", "HTTPS / JSON")
    Rel(admin, ticketing_system, "Configures venues, events, quotas", "HTTPS / JSON")
    Rel(ticketing_system, mock_payment, "Initiates payments, receives webhooks", "HTTPS / REST")
    Rel(ticketing_system, email_service, "Dispatches booking confirmations & e-tickets", "Async RabbitMQ")
```

## Stakeholders & Users

1. **Customer**:
   - Authentication (Registration / Login / JWT)
   - Real-time event and seat browsing
   - High-concurrency seat locking (5-minute hold)
   - Booking and mock checkout
   - E-ticket retrieval with unique QR tokens

2. **Administrator**:
   - Venue management (capacities, layouts)
   - Event lifecycle (Draft, Published, Cancelled)
   - Ticket type pricing & quotas
   - System monitoring & metrics observation
