package repository

import (
	"github.com/event-ticketing/event-service/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EventRepository interface {
	Create(event *model.Event) error
	FindByID(id uuid.UUID) (*model.Event, error)
	FindAll(filter model.EventFilter) ([]model.Event, int64, error)
	Update(event *model.Event) error
	Delete(id uuid.UUID) error
}

type eventRepository struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) EventRepository {
	return &eventRepository{db: db}
}

func (r *eventRepository) Create(event *model.Event) error {
	return r.db.Create(event).Error
}

func (r *eventRepository) FindByID(id uuid.UUID) (*model.Event, error) {
	var event model.Event
	err := r.db.Preload("Venue").Where("id = ? AND deleted_at IS NULL", id).First(&event).Error
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func (r *eventRepository) FindAll(filter model.EventFilter) ([]model.Event, int64, error) {
	var events []model.Event
	var total int64

	query := r.db.Model(&model.Event{}).Where("deleted_at IS NULL")

	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.VenueID != nil {
		query = query.Where("venue_id = ?", *filter.VenueID)
	}
	if filter.Search != "" {
		search := "%" + filter.Search + "%"
		query = query.Where("title ILIKE ? OR description ILIKE ?", search, search)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (filter.Page - 1) * filter.PerPage
	if err := query.Preload("Venue").Offset(offset).Limit(filter.PerPage).
		Order("created_at DESC").Find(&events).Error; err != nil {
		return nil, 0, err
	}

	return events, total, nil
}

func (r *eventRepository) Update(event *model.Event) error {
	return r.db.Save(event).Error
}

func (r *eventRepository) Delete(id uuid.UUID) error {
	return r.db.Where("id = ?", id).Delete(&model.Event{}).Error
}
