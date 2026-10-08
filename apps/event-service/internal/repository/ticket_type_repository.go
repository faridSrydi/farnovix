package repository

import (
	"github.com/event-ticketing/event-service/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TicketTypeRepository interface {
	Create(tt *model.TicketType) error
	FindByID(id uuid.UUID) (*model.TicketType, error)
	FindByEventID(eventID uuid.UUID) ([]model.TicketType, error)
	Update(tt *model.TicketType) error
	Delete(id uuid.UUID) error
}

type ticketTypeRepository struct {
	db *gorm.DB
}

func NewTicketTypeRepository(db *gorm.DB) TicketTypeRepository {
	return &ticketTypeRepository{db: db}
}

func (r *ticketTypeRepository) Create(tt *model.TicketType) error {
	return r.db.Create(tt).Error
}

func (r *ticketTypeRepository) FindByID(id uuid.UUID) (*model.TicketType, error) {
	var tt model.TicketType
	err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&tt).Error
	if err != nil {
		return nil, err
	}
	return &tt, nil
}

func (r *ticketTypeRepository) FindByEventID(eventID uuid.UUID) ([]model.TicketType, error) {
	var types []model.TicketType
	err := r.db.Where("event_id = ? AND deleted_at IS NULL", eventID).
		Order("price ASC").Find(&types).Error
	return types, err
}

func (r *ticketTypeRepository) Update(tt *model.TicketType) error {
	return r.db.Save(tt).Error
}

func (r *ticketTypeRepository) Delete(id uuid.UUID) error {
	return r.db.Where("id = ?", id).Delete(&model.TicketType{}).Error
}
