package handler

import (
	"net/http"

	"github.com/event-ticketing/event-service/internal/model"
	"github.com/event-ticketing/event-service/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SeatHandler struct {
	svc service.SeatService
}

func NewSeatHandler(svc service.SeatService) *SeatHandler {
	return &SeatHandler{svc: svc}
}

func (h *SeatHandler) RegisterRoutes(rg *gin.RouterGroup) {
	seats := rg.Group("/seats")
	{
		seats.POST("", h.Create)
		seats.GET("/event/:event_id", h.ListByEvent)
	}
}

func (h *SeatHandler) Create(c *gin.Context) {
	var req model.CreateSeatsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.APIError{
			Error: model.ErrorDetail{Code: "VALIDATION_ERROR", Message: err.Error()},
		})
		return
	}

	seats, err := h.svc.CreateSeats(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.APIError{
			Error: model.ErrorDetail{Code: "INTERNAL_ERROR", Message: err.Error()},
		})
		return
	}

	c.JSON(http.StatusCreated, model.APIResponse{
		Data: seats,
		Meta: &model.Meta{Total: int64(len(seats))},
	})
}

func (h *SeatHandler) ListByEvent(c *gin.Context) {
	eventID, err := uuid.Parse(c.Param("event_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.APIError{
			Error: model.ErrorDetail{Code: "INVALID_ID", Message: "invalid event ID"},
		})
		return
	}

	seats, err := h.svc.GetSeatsByEvent(eventID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.APIError{
			Error: model.ErrorDetail{Code: "INTERNAL_ERROR", Message: "internal server error"},
		})
		return
	}

	c.JSON(http.StatusOK, model.APIResponse{
		Data: seats,
		Meta: &model.Meta{Total: int64(len(seats))},
	})
}
