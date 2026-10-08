package repository

import (
	"github.com/event-ticketing/event-service/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SeatRepository interface {
	CreateBatch(seats []model.Seat) error
	FindByEventID(eventID uuid.UUID) ([]model.Seat, error)
	FindByID(id uuid.UUID) (*model.Seat, error)
	UpdateStatus(id uuid.UUID, status string) error
}

type seatRepository struct {
	db *gorm.DB
}

func NewSeatRepository(db *gorm.DB) SeatRepository {
	return &seatRepository{db: db}
}

func (r *seatRepository) CreateBatch(seats []model.Seat) error {
	return r.db.CreateInBatches(seats, 100).Error
}

func (r *seatRepository) FindByEventID(eventID uuid.UUID) ([]model.Seat, error) {
	var seats []model.Seat
	err := r.db.Where("event_id = ? AND deleted_at IS NULL", eventID).
		Order("section ASC, row ASC, number ASC").Find(&seats).Error
	return seats, err
}

func (r *seatRepository) FindByID(id uuid.UUID) (*model.Seat, error) {
	var seat model.Seat
	err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&seat).Error
	if err != nil {
		return nil, err
	}
	return &seat, nil
}

func (r *seatRepository) UpdateStatus(id uuid.UUID, status string) error {
	return r.db.Model(&model.Seat{}).Where("id = ?", id).Update("status", status).Error
}
