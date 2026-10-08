package service_test

import (
	"errors"
	"testing"
	"time"

	"github.com/event-ticketing/event-service/internal/model"
	"github.com/event-ticketing/event-service/internal/service"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// MockEventRepository implements repository.EventRepository
type MockEventRepository struct {
	events map[uuid.UUID]*model.Event
}

func NewMockEventRepo() *MockEventRepository {
	return &MockEventRepository{events: make(map[uuid.UUID]*model.Event)}
}

func (m *MockEventRepository) Create(event *model.Event) error {
	m.events[event.ID] = event
	return nil
}

func (m *MockEventRepository) FindByID(id uuid.UUID) (*model.Event, error) {
	if ev, ok := m.events[id]; ok {
		return ev, nil
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *MockEventRepository) FindAll(filter model.EventFilter) ([]model.Event, int64, error) {
	var list []model.Event
	for _, ev := range m.events {
		if filter.Status != "" && ev.Status != filter.Status {
			continue
		}
		list = append(list, *ev)
	}
	return list, int64(len(list)), nil
}

func (m *MockEventRepository) Update(event *model.Event) error {
	m.events[event.ID] = event
	return nil
}

func (m *MockEventRepository) Delete(id uuid.UUID) error {
	delete(m.events, id)
	return nil
}

// MockVenueRepository implements repository.VenueRepository
type MockVenueRepository struct {
	venues map[uuid.UUID]*model.Venue
}

func NewMockVenueRepo() *MockVenueRepository {
	return &MockVenueRepository{venues: make(map[uuid.UUID]*model.Venue)}
}

func (m *MockVenueRepository) Create(v *model.Venue) error {
	m.venues[v.ID] = v
	return nil
}

func (m *MockVenueRepository) FindByID(id uuid.UUID) (*model.Venue, error) {
	if v, ok := m.venues[id]; ok {
		return v, nil
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *MockVenueRepository) FindAll(page, perPage int) ([]model.Venue, int64, error) {
	var list []model.Venue
	for _, v := range m.venues {
		list = append(list, *v)
	}
	return list, int64(len(list)), nil
}

func (m *MockVenueRepository) Update(v *model.Venue) error {
	m.venues[v.ID] = v
	return nil
}

func (m *MockVenueRepository) Delete(id uuid.UUID) error {
	delete(m.venues, id)
	return nil
}

func TestCreateEvent_Success(t *testing.T) {
	eventRepo := NewMockEventRepo()
	venueRepo := NewMockVenueRepo()

	venueID := uuid.New()
	venueRepo.Create(&model.Venue{
		ID:       venueID,
		Name:     "GBK Stadium",
		Capacity: 50000,
	})

	svc := service.NewEventService(eventRepo, venueRepo)

	req := model.CreateEventRequest{
		Title:       "Rock Festival 2026",
		Description: "Great rock festival",
		VenueID:     venueID.String(),
		StartTime:   time.Now().Add(24 * time.Hour).Format(time.RFC3339),
		EndTime:     time.Now().Add(30 * time.Hour).Format(time.RFC3339),
		Status:      "PUBLISHED",
	}

	ev, err := svc.CreateEvent(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if ev.Title != req.Title {
		t.Errorf("expected title %s, got %s", req.Title, ev.Title)
	}
	if ev.Status != "PUBLISHED" {
		t.Errorf("expected status PUBLISHED, got %s", ev.Status)
	}
}

func TestCreateEvent_InvalidVenue(t *testing.T) {
	eventRepo := NewMockEventRepo()
	venueRepo := NewMockVenueRepo()
	svc := service.NewEventService(eventRepo, venueRepo)

	nonExistentVenue := uuid.New().String()
	req := model.CreateEventRequest{
		Title:       "Orchestra",
		VenueID:     nonExistentVenue,
		StartTime:   time.Now().Add(24 * time.Hour).Format(time.RFC3339),
		EndTime:     time.Now().Add(28 * time.Hour).Format(time.RFC3339),
	}

	_, err := svc.CreateEvent(req)
	if err == nil {
		t.Fatalf("expected error for non-existent venue, got nil")
	}
}

func TestCreateEvent_InvalidTimes(t *testing.T) {
	eventRepo := NewMockEventRepo()
	venueRepo := NewMockVenueRepo()

	venueID := uuid.New()
	venueRepo.Create(&model.Venue{ID: venueID, Name: "Hall B", Capacity: 1000})

	svc := service.NewEventService(eventRepo, venueRepo)

	// End time before start time
	req := model.CreateEventRequest{
		Title:     "Concert",
		VenueID:   venueID.String(),
		StartTime: time.Now().Add(24 * time.Hour).Format(time.RFC3339),
		EndTime:   time.Now().Add(12 * time.Hour).Format(time.RFC3339),
	}

	_, err := svc.CreateEvent(req)
	if !errors.Is(err, service.ErrInvalidInput) {
		t.Errorf("expected ErrInvalidInput, got %v", err)
	}
}

func TestGetEvent_NotFound(t *testing.T) {
	eventRepo := NewMockEventRepo()
	venueRepo := NewMockVenueRepo()
	svc := service.NewEventService(eventRepo, venueRepo)

	_, err := svc.GetEvent(uuid.New())
	if !errors.Is(err, service.ErrEventNotFound) {
		t.Errorf("expected ErrEventNotFound, got %v", err)
	}
}
