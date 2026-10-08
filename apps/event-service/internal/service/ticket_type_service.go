package service

import (
	"errors"

	"github.com/event-ticketing/event-service/internal/model"
	"github.com/event-ticketing/event-service/internal/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var ErrTicketTypeNotFound = errors.New("ticket type not found")

type TicketTypeService interface {
	CreateTicketType(req model.CreateTicketTypeRequest) (*model.TicketType, error)
	GetTicketTypesByEvent(eventID uuid.UUID) ([]model.TicketType, error)
	DeleteTicketType(id uuid.UUID) error
}

type ticketTypeService struct {
	repo repository.TicketTypeRepository
}

func NewTicketTypeService(repo repository.TicketTypeRepository) TicketTypeService {
	return &ticketTypeService{repo: repo}
}

func (s *ticketTypeService) CreateTicketType(req model.CreateTicketTypeRequest) (*model.TicketType, error) {
	eventID, err := uuid.Parse(req.EventID)
	if err != nil {
		return nil, ErrInvalidInput
	}

	tt := &model.TicketType{
		ID:      uuid.New(),
		EventID: eventID,
		Name:    req.Name,
		Price:   req.Price,
		Quota:   req.Quota,
	}

	if err := s.repo.Create(tt); err != nil {
		return nil, err
	}
	return tt, nil
}

func (s *ticketTypeService) GetTicketTypesByEvent(eventID uuid.UUID) ([]model.TicketType, error) {
	return s.repo.FindByEventID(eventID)
}

func (s *ticketTypeService) DeleteTicketType(id uuid.UUID) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrTicketTypeNotFound
		}
		return err
	}
	return s.repo.Delete(id)
}
