package model

import (
	"time"

	"github.com/google/uuid"
)

// Event represents an event entity.
type Event struct {
	ID          uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Title       string     `json:"title" gorm:"not null"`
	Description string     `json:"description"`
	VenueID     uuid.UUID  `json:"venue_id" gorm:"type:uuid;not null"`
	Venue       *Venue     `json:"venue,omitempty" gorm:"foreignKey:VenueID"`
	StartTime   time.Time  `json:"start_time" gorm:"not null"`
	EndTime     time.Time  `json:"end_time" gorm:"not null"`
	Status      string     `json:"status" gorm:"not null;default:'DRAFT'"` // DRAFT, PUBLISHED, CANCELLED
	ImageURL    string     `json:"image_url"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}

// Venue represents a venue where events take place.
type Venue struct {
	ID        uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name      string     `json:"name" gorm:"not null"`
	Address   string     `json:"address"`
	City      string     `json:"city"`
	Capacity  int        `json:"capacity" gorm:"not null"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}

// TicketType represents a category of ticket for an event.
type TicketType struct {
	ID        uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	EventID   uuid.UUID  `json:"event_id" gorm:"type:uuid;not null;index"`
	Name      string     `json:"name" gorm:"not null"`  // VIP, Regular, etc.
	Price     float64    `json:"price" gorm:"not null"`
	Quota     int        `json:"quota" gorm:"not null"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}

// Seat represents an individual seat within a venue for a specific event.
type Seat struct {
	ID           uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	EventID      uuid.UUID  `json:"event_id" gorm:"type:uuid;not null;index:idx_seat_event"`
	TicketTypeID uuid.UUID  `json:"ticket_type_id" gorm:"type:uuid;not null"`
	Section      string     `json:"section" gorm:"not null"`
	Row          string     `json:"row" gorm:"not null"`
	Number       string     `json:"number" gorm:"not null"`
	Status       string     `json:"status" gorm:"not null;default:'AVAILABLE'"` // AVAILABLE, LOCKED, SOLD
	Label        string     `json:"label"`                                       // e.g. "A12"
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}
