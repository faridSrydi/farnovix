package handler

import (
	"errors"
	"math"
	"net/http"
	"strconv"

	"github.com/event-ticketing/event-service/internal/model"
	"github.com/event-ticketing/event-service/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type EventHandler struct {
	svc service.EventService
}

func NewEventHandler(svc service.EventService) *EventHandler {
	return &EventHandler{svc: svc}
}

func (h *EventHandler) RegisterRoutes(rg *gin.RouterGroup) {
	events := rg.Group("/events")
	{
		events.POST("", h.Create)
		events.GET("", h.List)
		events.GET("/:id", h.Get)
		events.PUT("/:id", h.Update)
		events.DELETE("/:id", h.Delete)
	}
}

// Create handles POST /events
func (h *EventHandler) Create(c *gin.Context) {
	var req model.CreateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.APIError{
			Error: model.ErrorDetail{
				Code:      "VALIDATION_ERROR",
				Message:   err.Error(),
				RequestID: c.GetString("request_id"),
			},
		})
		return
	}

	event, err := h.svc.CreateEvent(req)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, model.APIResponse{Data: event})
}

// Get handles GET /events/:id
func (h *EventHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.APIError{
			Error: model.ErrorDetail{
				Code:    "INVALID_ID",
				Message: "invalid event ID format",
			},
		})
		return
	}

	event, err := h.svc.GetEvent(id)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, model.APIResponse{Data: event})
}

// List handles GET /events
func (h *EventHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))

	filter := model.EventFilter{
		Status: c.Query("status"),
		Search: c.Query("search"),
		PaginationQuery: model.PaginationQuery{
			Page:    page,
			PerPage: perPage,
		},
	}

	if vid := c.Query("venue_id"); vid != "" {
		venueID, err := uuid.Parse(vid)
		if err == nil {
			filter.VenueID = &venueID
		}
	}

	events, total, err := h.svc.ListEvents(filter)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(perPage)))

	c.JSON(http.StatusOK, model.APIResponse{
		Data: events,
		Meta: &model.Meta{
			Page:       page,
			PerPage:    perPage,
			Total:      total,
			TotalPages: totalPages,
		},
	})
}

// Update handles PUT /events/:id
func (h *EventHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.APIError{
			Error: model.ErrorDetail{Code: "INVALID_ID", Message: "invalid event ID format"},
		})
		return
	}

	var req model.UpdateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.APIError{
			Error: model.ErrorDetail{Code: "VALIDATION_ERROR", Message: err.Error()},
		})
		return
	}

	event, err := h.svc.UpdateEvent(id, req)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, model.APIResponse{Data: event})
}

// Delete handles DELETE /events/:id
func (h *EventHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.APIError{
			Error: model.ErrorDetail{Code: "INVALID_ID", Message: "invalid event ID format"},
		})
		return
	}

	if err := h.svc.DeleteEvent(id); err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, model.APIResponse{Data: gin.H{"message": "event deleted"}})
}

func handleServiceError(c *gin.Context, err error) {
	requestID := c.GetString("request_id")

	if errors.Is(err, service.ErrEventNotFound) {
		c.JSON(http.StatusNotFound, model.APIError{
			Error: model.ErrorDetail{
				Code:      "EVENT_NOT_FOUND",
				Message:   "event not found",
				RequestID: requestID,
			},
		})
		return
	}

	if errors.Is(err, service.ErrInvalidInput) {
		c.JSON(http.StatusBadRequest, model.APIError{
			Error: model.ErrorDetail{
				Code:      "INVALID_INPUT",
				Message:   err.Error(),
				RequestID: requestID,
			},
		})
		return
	}

	c.JSON(http.StatusInternalServerError, model.APIError{
		Error: model.ErrorDetail{
			Code:      "INTERNAL_ERROR",
			Message:   "internal server error",
			RequestID: requestID,
		},
	})
}
