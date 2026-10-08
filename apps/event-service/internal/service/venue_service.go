package service

import (
	"errors"
	"fmt"

	"github.com/event-ticketing/event-service/internal/model"
	"github.com/event-ticketing/event-service/internal/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var ErrVenueNotFound = errors.New("venue not found")

type VenueService interface {
	CreateVenue(req model.CreateVenueRequest) (*model.Venue, error)
	GetVenue(id uuid.UUID) (*model.Venue, error)
	ListVenues(page, perPage int) ([]model.Venue, int64, error)
	UpdateVenue(id uuid.UUID, req model.UpdateVenueRequest) (*model.Venue, error)
	DeleteVenue(id uuid.UUID) error
}

type venueService struct {
	repo repository.VenueRepository
}

func NewVenueService(repo repository.VenueRepository) VenueService {
	return &venueService{repo: repo}
}

func (s *venueService) CreateVenue(req model.CreateVenueRequest) (*model.Venue, error) {
	venue := &model.Venue{
		ID:       uuid.New(),
		Name:     req.Name,
		Address:  req.Address,
		City:     req.City,
		Capacity: req.Capacity,
	}
	if err := s.repo.Create(venue); err != nil {
		return nil, err
	}
	return venue, nil
}

func (s *venueService) GetVenue(id uuid.UUID) (*model.Venue, error) {
	venue, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrVenueNotFound
		}
		return nil, err
	}
	return venue, nil
}

func (s *venueService) ListVenues(page, perPage int) ([]model.Venue, int64, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}
	return s.repo.FindAll(page, perPage)
}

func (s *venueService) UpdateVenue(id uuid.UUID, req model.UpdateVenueRequest) (*model.Venue, error) {
	venue, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrVenueNotFound
		}
		return nil, err
	}

	if req.Name != nil {
		venue.Name = *req.Name
	}
	if req.Address != nil {
		venue.Address = *req.Address
	}
	if req.City != nil {
		venue.City = *req.City
	}
	if req.Capacity != nil {
		if *req.Capacity < 1 {
			return nil, fmt.Errorf("%w: capacity must be positive", ErrInvalidInput)
		}
		venue.Capacity = *req.Capacity
	}

	if err := s.repo.Update(venue); err != nil {
		return nil, err
	}
	return venue, nil
}

func (s *venueService) DeleteVenue(id uuid.UUID) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrVenueNotFound
		}
		return err
	}
	return s.repo.Delete(id)
}
