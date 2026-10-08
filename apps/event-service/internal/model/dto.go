package model

import "github.com/google/uuid"

// --- Request DTOs ---

type CreateEventRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	VenueID     string `json:"venue_id" binding:"required,uuid"`
	StartTime   string `json:"start_time" binding:"required"`
	EndTime     string `json:"end_time" binding:"required"`
	Status      string `json:"status"`
	ImageURL    string `json:"image_url"`
}

type UpdateEventRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	VenueID     *string `json:"venue_id"`
	StartTime   *string `json:"start_time"`
	EndTime     *string `json:"end_time"`
	Status      *string `json:"status"`
	ImageURL    *string `json:"image_url"`
}

type CreateVenueRequest struct {
	Name     string `json:"name" binding:"required"`
	Address  string `json:"address"`
	City     string `json:"city"`
	Capacity int    `json:"capacity" binding:"required,min=1"`
}

type UpdateVenueRequest struct {
	Name     *string `json:"name"`
	Address  *string `json:"address"`
	City     *string `json:"city"`
	Capacity *int    `json:"capacity"`
}

type CreateTicketTypeRequest struct {
	EventID string  `json:"event_id" binding:"required,uuid"`
	Name    string  `json:"name" binding:"required"`
	Price   float64 `json:"price" binding:"required,min=0"`
	Quota   int     `json:"quota" binding:"required,min=1"`
}

type CreateSeatsRequest struct {
	EventID      string `json:"event_id" binding:"required,uuid"`
	TicketTypeID string `json:"ticket_type_id" binding:"required,uuid"`
	Section      string `json:"section" binding:"required"`
	RowPrefix    string `json:"row_prefix" binding:"required"`
	Count        int    `json:"count" binding:"required,min=1,max=500"`
}

// --- Response DTOs ---

type APIResponse struct {
	Data interface{} `json:"data,omitempty"`
	Meta *Meta       `json:"meta,omitempty"`
}

type Meta struct {
	Page       int   `json:"page,omitempty"`
	PerPage    int   `json:"per_page,omitempty"`
	Total      int64 `json:"total,omitempty"`
	TotalPages int   `json:"total_pages,omitempty"`
}

type APIError struct {
	Error ErrorDetail `json:"error"`
}

type ErrorDetail struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	RequestID string `json:"request_id,omitempty"`
}

// PaginationQuery holds pagination params.
type PaginationQuery struct {
	Page    int
	PerPage int
}

// EventFilter for listing events.
type EventFilter struct {
	Status  string
	VenueID *uuid.UUID
	Search  string
	PaginationQuery
}
