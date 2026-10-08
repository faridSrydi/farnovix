package handler

import (
	"net/http"

	"github.com/event-ticketing/event-service/internal/model"
	"github.com/event-ticketing/event-service/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TicketTypeHandler struct {
	svc service.TicketTypeService
}

func NewTicketTypeHandler(svc service.TicketTypeService) *TicketTypeHandler {
	return &TicketTypeHandler{svc: svc}
}

func (h *TicketTypeHandler) RegisterRoutes(rg *gin.RouterGroup) {
	tt := rg.Group("/ticket-types")
	{
		tt.POST("", h.Create)
		tt.GET("/event/:event_id", h.ListByEvent)
		tt.DELETE("/:id", h.Delete)
	}
}

func (h *TicketTypeHandler) Create(c *gin.Context) {
	var req model.CreateTicketTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.APIError{
			Error: model.ErrorDetail{Code: "VALIDATION_ERROR", Message: err.Error()},
		})
		return
	}

	tt, err := h.svc.CreateTicketType(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.APIError{
			Error: model.ErrorDetail{Code: "INTERNAL_ERROR", Message: "failed to create ticket type"},
		})
		return
	}

	c.JSON(http.StatusCreated, model.APIResponse{Data: tt})
}

func (h *TicketTypeHandler) ListByEvent(c *gin.Context) {
	eventID, err := uuid.Parse(c.Param("event_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.APIError{
			Error: model.ErrorDetail{Code: "INVALID_ID", Message: "invalid event ID"},
		})
		return
	}

	types, err := h.svc.GetTicketTypesByEvent(eventID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.APIError{
			Error: model.ErrorDetail{Code: "INTERNAL_ERROR", Message: "internal server error"},
		})
		return
	}

	c.JSON(http.StatusOK, model.APIResponse{Data: types})
}

func (h *TicketTypeHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.APIError{
			Error: model.ErrorDetail{Code: "INVALID_ID", Message: "invalid ticket type ID"},
		})
		return
	}

	if err := h.svc.DeleteTicketType(id); err != nil {
		c.JSON(http.StatusInternalServerError, model.APIError{
			Error: model.ErrorDetail{Code: "INTERNAL_ERROR", Message: "failed to delete ticket type"},
		})
		return
	}

	c.JSON(http.StatusOK, model.APIResponse{Data: gin.H{"message": "ticket type deleted"}})
}
