.PHONY: up down build logs ps test clean

# Start development environment
up:
	docker compose up -d --build

# Stop all services
down:
	docker compose down

# Build all services
build:
	docker compose build

# View logs
logs:
	docker compose logs -f

# View running containers
ps:
	docker compose ps

# Clean everything including volumes
clean:
	docker compose down -v --remove-orphans

# Service-specific targets
event-service-logs:
	docker compose logs -f event-service

event-service-build:
	docker compose build event-service

# Database access
db-event:
	docker compose exec postgres psql -U ticketing -d event_db

# Health checks
health-event:
	curl -s http://localhost:8082/health | python -m json.tool

ready-event:
	curl -s http://localhost:8082/ready | python -m json.tool
