# Event Service API Reference

Base Path: `/api/v1`

All responses follow the standardized envelope schema:
- Success: `{ "data": <payload>, "meta": <pagination_info> }`
- Error: `{ "error": { "code": "<ERROR_CODE>", "message": "<description>", "request_id": "<uuid>" } }`

---

## 1. Events

### `POST /events`
Creates a new event.

**Request Body:**
```json
{
  "title": "Jakarta Music Festival 2026",
  "description": "Annual open air music festival",
  "venue_id": "c1f7a4a2-94b2-4d22-b5e1-c0817b35f082",
  "start_time": "2026-12-10T18:00:00Z",
  "end_time": "2026-12-10T23:00:00Z",
  "status": "PUBLISHED",
  "image_url": "https://example.com/poster.jpg"
}
```

**Response `201 Created`:**
```json
{
  "data": {
    "id": "e4b52c08-0112-4fcf-840f-7f551b9d4021",
    "title": "Jakarta Music Festival 2026",
    "description": "Annual open air music festival",
    "venue_id": "c1f7a4a2-94b2-4d22-b5e1-c0817b35f082",
    "venue": {
      "id": "c1f7a4a2-94b2-4d22-b5e1-c0817b35f082",
      "name": "Gelora Bung Karno",
      "capacity": 77000
    },
    "start_time": "2026-12-10T18:00:00Z",
    "end_time": "2026-12-10T23:00:00Z",
    "status": "PUBLISHED",
    "created_at": "2026-10-08T22:00:00Z"
  }
}
```

### `GET /events`
Lists events with filtering and pagination.

**Query Parameters:**
- `page` (default: 1)
- `per_page` (default: 20)
- `status` (e.g. `PUBLISHED`)
- `venue_id` (UUID filter)
- `search` (text search in title & description)

**Response `200 OK`:**
```json
{
  "data": [ ... ],
  "meta": {
    "page": 1,
    "per_page": 20,
    "total": 45,
    "total_pages": 3
  }
}
```

### `GET /events/:id`
Retrieves single event details including venue metadata.

### `PUT /events/:id`
Updates event fields.

### `DELETE /events/:id`
Soft deletes an event.

---

## 2. Venues

### `POST /venues`
Creates a new venue.
- `name`: string (required)
- `address`: string
- `city`: string
- `capacity`: integer > 0 (required)

### `GET /venues`
Lists venues.

### `GET /venues/:id`
Retrieves venue details.

---

## 3. Ticket Types

### `POST /ticket-types`
- `event_id`: UUID (required)
- `name`: string (e.g. `VIP`, `Regular`)
- `price`: float >= 0
- `quota`: integer > 0

### `GET /ticket-types/event/:event_id`
Retrieves all ticket categories configured for an event.

---

## 4. Seats

### `POST /seats`
Batch generates seats for an event tier:
```json
{
  "event_id": "e4b52c08-0112-4fcf-840f-7f551b9d4021",
  "ticket_type_id": "a93bc71d-84ff-4ea0-a612-4f90117b3333",
  "section": "VIP North",
  "row_prefix": "A",
  "count": 50
}
```

### `GET /seats/event/:event_id`
Returns real-time seat states (`AVAILABLE`, `LOCKED`, `SOLD`) for the interactive seat map.
