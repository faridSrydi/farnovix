package service

import (
	"fmt"
	"strconv"

	"github.com/event-ticketing/event-service/internal/model"
	"github.com/event-ticketing/event-service/internal/repository"
	"github.com/google/uuid"
)

type SeatService interface {
	CreateSeats(req model.CreateSeatsRequest) ([]model.Seat, error)
	GetSeatsByEvent(eventID uuid.UUID) ([]model.Seat, error)
}

type seatService struct {
	seatRepo       repository.SeatRepository
	ticketTypeRepo repository.TicketTypeRepository
}

func NewSeatService(seatRepo repository.SeatRepository, ticketTypeRepo repository.TicketTypeRepository) SeatService {
	return &seatService{
		seatRepo:       seatRepo,
		ticketTypeRepo: ticketTypeRepo,
	}
}

func (s *seatService) CreateSeats(req model.CreateSeatsRequest) ([]model.Seat, error) {
	eventID, err := uuid.Parse(req.EventID)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid event_id", ErrInvalidInput)
	}

	ticketTypeID, err := uuid.Parse(req.TicketTypeID)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid ticket_type_id", ErrInvalidInput)
	}

	// Verify ticket type exists
	_, err = s.ticketTypeRepo.FindByID(ticketTypeID)
	if err != nil {
		return nil, fmt.Errorf("%w: ticket type not found", ErrInvalidInput)
	}

	seats := make([]model.Seat, req.Count)
	for i := 0; i < req.Count; i++ {
		num := strconv.Itoa(i + 1)
		seats[i] = model.Seat{
			ID:           uuid.New(),
			EventID:      eventID,
			TicketTypeID: ticketTypeID,
			Section:      req.Section,
			Row:          req.RowPrefix,
			Number:       num,
			Status:       "AVAILABLE",
			Label:        req.RowPrefix + num,
		}
	}

	if err := s.seatRepo.CreateBatch(seats); err != nil {
		return nil, err
	}

	return seats, nil
}

func (s *seatService) GetSeatsByEvent(eventID uuid.UUID) ([]model.Seat, error) {
	return s.seatRepo.FindByEventID(eventID)
}
