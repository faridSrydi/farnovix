# Comprehensive API Test Script for Event Service

$baseUrl = "http://localhost:8082"
$ErrorActionPreference = "Stop"

Write-Output "=================================================="
Write-Output "TESTING EVENT SERVICE API ON $baseUrl"
Write-Output "=================================================="

# 1. Health check
Write-Output "`n[1] Testing GET /health ..."
$health = Invoke-RestMethod -Uri "$baseUrl/health" -Method Get
Write-Output "Response: $($health | ConvertTo-Json -Compress)"
if ($health.status -ne "up") { throw "Health check failed" }
Write-Output "PASS: Health check"

# 2. Readiness check
Write-Output "`n[2] Testing GET /ready ..."
$ready = Invoke-RestMethod -Uri "$baseUrl/ready" -Method Get
Write-Output "Response: $($ready | ConvertTo-Json -Compress)"
if ($ready.status -ne "ready") { throw "Readiness check failed" }
Write-Output "PASS: Readiness check (DB connected)"

# 3. Create Venue
Write-Output "`n[3] Testing POST /api/v1/venues ..."
$venuePayload = @{
    name = "Gelora Bung Karno Main Stadium"
    address = "Jl. Pintu Satu Senayan, Gelora"
    city = "Jakarta"
    capacity = 77000
} | ConvertTo-Json

$venueRes = Invoke-RestMethod -Uri "$baseUrl/api/v1/venues" -Method Post -Body $venuePayload -ContentType "application/json"
$venueId = $venueRes.data.id
Write-Output "Created Venue ID: $venueId"
Write-Output "Response: $($venueRes.data.name) in $($venueRes.data.city)"
if (-not $venueId) { throw "Venue creation failed" }
Write-Output "PASS: Create venue"

# 4. List Venues
Write-Output "`n[4] Testing GET /api/v1/venues ..."
$venuesList = Invoke-RestMethod -Uri "$baseUrl/api/v1/venues" -Method Get
Write-Output "Total Venues: $($venuesList.meta.total)"
if ($venuesList.meta.total -lt 1) { throw "List venues failed" }
Write-Output "PASS: List venues"

# 5. Create Event
Write-Output "`n[5] Testing POST /api/v1/events ..."
$eventPayload = @{
    title = "Coldplay Music of the Spheres World Tour"
    description = "Special stadium concert in Jakarta"
    venue_id = $venueId
    start_time = "2026-11-15T19:00:00Z"
    end_time = "2026-11-15T23:00:00Z"
    status = "PUBLISHED"
    image_url = "https://images.unsplash.com/photo-1540039155733-5bb30b53aa14"
} | ConvertTo-Json

$eventRes = Invoke-RestMethod -Uri "$baseUrl/api/v1/events" -Method Post -Body $eventPayload -ContentType "application/json"
$eventId = $eventRes.data.id
Write-Output "Created Event ID: $eventId"
Write-Output "Title: $($eventRes.data.title), Status: $($eventRes.data.status)"
if (-not $eventId) { throw "Event creation failed" }
Write-Output "PASS: Create event"

# 6. Get Event by ID
Write-Output "`n[6] Testing GET /api/v1/events/$eventId ..."
$getEvent = Invoke-RestMethod -Uri "$baseUrl/api/v1/events/$eventId" -Method Get
Write-Output "Retrieved Title: $($getEvent.data.title)"
Write-Output "Venue Name: $($getEvent.data.venue.name)"
if ($getEvent.data.venue.name -ne "Gelora Bung Karno Main Stadium") { throw "Get event preload failed" }
Write-Output "PASS: Get event by ID with venue preload"

# 7. List Events with search and pagination
Write-Output "`n[7] Testing GET /api/v1/events?search=Coldplay ..."
$searchRes = Invoke-RestMethod -Uri "$baseUrl/api/v1/events?search=Coldplay" -Method Get
Write-Output "Search Total: $($searchRes.meta.total)"
if ($searchRes.meta.total -lt 1) { throw "Event search failed" }
Write-Output "PASS: List events with search query"

# 8. Create Ticket Type
Write-Output "`n[8] Testing POST /api/v1/ticket-types ..."
$ttPayload = @{
    event_id = $eventId
    name = "VIP Cat 1"
    price = 3500000.00
    quota = 500
} | ConvertTo-Json

$ttRes = Invoke-RestMethod -Uri "$baseUrl/api/v1/ticket-types" -Method Post -Body $ttPayload -ContentType "application/json"
$ticketTypeId = $ttRes.data.id
Write-Output "Created Ticket Type ID: $ticketTypeId ($($ttRes.data.name) - Rp $($ttRes.data.price))"
if (-not $ticketTypeId) { throw "Ticket type creation failed" }
Write-Output "PASS: Create ticket type"

# 9. List Ticket Types by Event
Write-Output "`n[9] Testing GET /api/v1/ticket-types/event/$eventId ..."
$ttList = Invoke-RestMethod -Uri "$baseUrl/api/v1/ticket-types/event/$eventId" -Method Get
Write-Output "Found $($ttList.data.Count) ticket types"
if ($ttList.data.Count -lt 1) { throw "List ticket types failed" }
Write-Output "PASS: List ticket types by event"

# 10. Create Seats in Batch
Write-Output "`n[10] Testing POST /api/v1/seats (Batch Creation) ..."
$seatsPayload = @{
    event_id = $eventId
    ticket_type_id = $ticketTypeId
    section = "Cat 1A"
    row_prefix = "A"
    count = 10
} | ConvertTo-Json

$seatsRes = Invoke-RestMethod -Uri "$baseUrl/api/v1/seats" -Method Post -Body $seatsPayload -ContentType "application/json"
Write-Output "Generated $($seatsRes.meta.total) seats. Sample label: $($seatsRes.data[0].label)"
if ($seatsRes.meta.total -ne 10) { throw "Seat creation failed" }
Write-Output "PASS: Batch seat creation"

# 11. List Seats by Event
Write-Output "`n[11] Testing GET /api/v1/seats/event/$eventId ..."
$seatsList = Invoke-RestMethod -Uri "$baseUrl/api/v1/seats/event/$eventId" -Method Get
Write-Output "Total seats retrieved: $($seatsList.meta.total)"
if ($seatsList.meta.total -ne 10) { throw "List seats failed" }
Write-Output "PASS: List seats by event"

Write-Output "`n=================================================="
Write-Output "ALL 11 API TESTS PASSED SUCCESSFULLY!"
Write-Output "=================================================="
