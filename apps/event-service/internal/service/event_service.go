package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/event-ticketing/event-service/internal/model"
	"github.com/event-ticketing/event-service/internal/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrEventNotFound = errors.New("event not found")
	ErrInvalidInput  = errors.New("invalid input")
)

type EventService interface {
	CreateEvent(req model.CreateEventRequest) (*model.Event, error)
	GetEvent(id uuid.UUID) (*model.Event, error)
	ListEvents(filter model.EventFilter) ([]model.Event, int64, error)
	UpdateEvent(id uuid.UUID, req model.UpdateEventRequest) (*model.Event, error)
	DeleteEvent(id uuid.UUID) error
}

type eventService struct {
	eventRepo repository.EventRepository
	venueRepo repository.VenueRepository
}

func NewEventService(eventRepo repository.EventRepository, venueRepo repository.VenueRepository) EventService {
	return &eventService{
		eventRepo: eventRepo,
		venueRepo: venueRepo,
	}
}

func (s *eventService) CreateEvent(req model.CreateEventRequest) (*model.Event, error) {
	venueID, err := uuid.Parse(req.VenueID)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid venue_id", ErrInvalidInput)
	}

	// Verify venue exists
	_, err = s.venueRepo.FindByID(venueID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: venue not found", ErrInvalidInput)
		}
		return nil, err
	}

	startTime, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid start_time format, use RFC3339", ErrInvalidInput)
	}

	endTime, err := time.Parse(time.RFC3339, req.EndTime)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid end_time format, use RFC3339", ErrInvalidInput)
	}

	if endTime.Before(startTime) {
		return nil, fmt.Errorf("%w: end_time must be after start_time", ErrInvalidInput)
	}

	status := "DRAFT"
	if req.Status != "" {
		status = req.Status
	}

	event := &model.Event{
		ID:          uuid.New(),
		Title:       req.Title,
		Description: req.Description,
		VenueID:     venueID,
		StartTime:   startTime,
		EndTime:     endTime,
		Status:      status,
		ImageURL:    req.ImageURL,
	}

	if err := s.eventRepo.Create(event); err != nil {
		return nil, err
	}

	// Reload with venue
	return s.eventRepo.FindByID(event.ID)
}

func (s *eventService) GetEvent(id uuid.UUID) (*model.Event, error) {
	event, err := s.eventRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrEventNotFound
		}
		return nil, err
	}
	return event, nil
}

func (s *eventService) ListEvents(filter model.EventFilter) ([]model.Event, int64, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.PerPage < 1 || filter.PerPage > 100 {
		filter.PerPage = 20
	}
	return s.eventRepo.FindAll(filter)
}

func (s *eventService) UpdateEvent(id uuid.UUID, req model.UpdateEventRequest) (*model.Event, error) {
	event, err := s.eventRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrEventNotFound
		}
		return nil, err
	}

	if req.Title != nil {
		event.Title = *req.Title
	}
	if req.Description != nil {
		event.Description = *req.Description
	}
	if req.VenueID != nil {
		vid, err := uuid.Parse(*req.VenueID)
		if err != nil {
			return nil, fmt.Errorf("%w: invalid venue_id", ErrInvalidInput)
		}
		if _, err := s.venueRepo.FindByID(vid); err != nil {
			return nil, fmt.Errorf("%w: venue not found", ErrInvalidInput)
		}
		event.VenueID = vid
	}
	if req.StartTime != nil {
		t, err := time.Parse(time.RFC3339, *req.StartTime)
		if err != nil {
			return nil, fmt.Errorf("%w: invalid start_time", ErrInvalidInput)
		}
		event.StartTime = t
	}
	if req.EndTime != nil {
		t, err := time.Parse(time.RFC3339, *req.EndTime)
		if err != nil {
			return nil, fmt.Errorf("%w: invalid end_time", ErrInvalidInput)
		}
		event.EndTime = t
	}
	if req.Status != nil {
		event.Status = *req.Status
	}
	if req.ImageURL != nil {
		event.ImageURL = *req.ImageURL
	}

	if err := s.eventRepo.Update(event); err != nil {
		return nil, err
	}

	return s.eventRepo.FindByID(id)
}

func (s *eventService) DeleteEvent(id uuid.UUID) error {
	_, err := s.eventRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrEventNotFound
		}
		return err
	}
	return s.eventRepo.Delete(id)
}
