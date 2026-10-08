# Service Boundaries & Domain Ownership

## Principles

1. **Strict Data Isolation**: No direct database access across service boundaries.
2. **Autonomous Modules**: Each service has its own `go.mod`, dependency tree, migrations, and Dockerfile.
3. **API-First Communication**: Synchronous queries go through REST APIs; side effects go through domain events.
4. **Idempotent State Transitions**: State machines are protected against duplicate requests and network retries.

## Services & Bounded Contexts

### 1. API Gateway
- **Domain**: Edge routing, policy enforcement.
- **Responsibilities**:
  - Reverse proxy routing for all `/api/v1/*` routes
  - JWT token validation and claims extraction
  - Rate limiting (token bucket / leaky bucket per IP and user)
  - Correlation ID injection (`X-Correlation-ID`)
  - Request logging and CORS headers

### 2. User Service
- **Domain**: Identity & Access Management (IAM).
- **Entities**: `User`, `Role`, `RefreshToken`
- **Responsibilities**:
  - Registration with bcrypt password hashing (cost factor 12)
  - Login issuing short-lived Access Tokens (15m) and Refresh Tokens (7d)
  - Profile management and role assignment (`CUSTOMER`, `ADMIN`)

### 3. Event Service
- **Domain**: Catalog & Venue Management.
- **Entities**: `Event`, `Venue`, `TicketType`, `Seat`
- **Responsibilities**:
  - Event lifecycle management (DRAFT, PUBLISHED, CANCELLED)
  - Venue definition (capacities, sections, rows)
  - Ticket tier configuration (VIP, Regular, quotas, prices)
  - Seat layout generation

### 4. Booking Service
- **Domain**: Reservation & Inventory Concurrency Control.
- **Entities**: `Booking`, `BookingItem`, `BookingStatusHistory`
- **State Machine**: `PENDING` -> `PAYMENT_PENDING` -> `CONFIRMED` | `CANCELLED` | `EXPIRED`
- **Responsibilities**:
  - Atomic distributed seat locking via Redis (`SET seat-lock:{event_id}:{seat_id} {booking_id} NX EX 300`)
  - Database transactional booking creation with seat availability verification
  - Automatic expiration handling (5-minute timer)
  - Prevention of double-booking under concurrent load (Research Question RQ3)

### 5. Payment Service
- **Domain**: Financial Transactions & Billing.
- **Entities**: `Payment`, `PaymentAttempt`, `Refund`
- **State Machine**: `PENDING` -> `PAID` | `FAILED` | `EXPIRED` | `REFUNDED`
- **Responsibilities**:
  - Idempotent payment intent creation (`Idempotency-Key`)
  - Mock payment gateway simulation (success/failure/delay toggle)
  - Webhook callback verification
  - Domain event publishing: `payment.completed`, `payment.failed`

### 6. Ticket Service
- **Domain**: Fulfillment & Access Credentials.
- **Entities**: `Ticket`, `TicketItem`
- **Responsibilities**:
  - Asynchronous consumption of `payment.completed`
  - Unique cryptographically secure ticket code generation
  - QR code rendering with verification tokens
  - Ticket validation API for venue scanning

### 7. Notification Service
- **Domain**: Customer Communication.
- **Entities**: `NotificationLog`, `DeliveryAttempt`
- **Responsibilities**:
  - Asynchronous consumption of domain events
  - Simulated email/SMS dispatch with retry and exponential backoff
  - Dead-letter routing on repeated delivery failure
  - Independent fault tolerance (Research Question RQ2)
