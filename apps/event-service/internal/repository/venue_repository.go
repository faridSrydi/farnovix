package repository

import (
	"github.com/event-ticketing/event-service/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type VenueRepository interface {
	Create(venue *model.Venue) error
	FindByID(id uuid.UUID) (*model.Venue, error)
	FindAll(page, perPage int) ([]model.Venue, int64, error)
	Update(venue *model.Venue) error
	Delete(id uuid.UUID) error
}

type venueRepository struct {
	db *gorm.DB
}

func NewVenueRepository(db *gorm.DB) VenueRepository {
	return &venueRepository{db: db}
}

func (r *venueRepository) Create(venue *model.Venue) error {
	return r.db.Create(venue).Error
}

func (r *venueRepository) FindByID(id uuid.UUID) (*model.Venue, error) {
	var venue model.Venue
	err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&venue).Error
	if err != nil {
		return nil, err
	}
	return &venue, nil
}

func (r *venueRepository) FindAll(page, perPage int) ([]model.Venue, int64, error) {
	var venues []model.Venue
	var total int64

	query := r.db.Model(&model.Venue{}).Where("deleted_at IS NULL")

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("name ASC").Find(&venues).Error; err != nil {
		return nil, 0, err
	}

	return venues, total, nil
}

func (r *venueRepository) Update(venue *model.Venue) error {
	return r.db.Save(venue).Error
}

func (r *venueRepository) Delete(id uuid.UUID) error {
	return r.db.Where("id = ?", id).Delete(&model.Venue{}).Error
}
