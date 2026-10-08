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

type VenueHandler struct {
	svc service.VenueService
}

func NewVenueHandler(svc service.VenueService) *VenueHandler {
	return &VenueHandler{svc: svc}
}

func (h *VenueHandler) RegisterRoutes(rg *gin.RouterGroup) {
	venues := rg.Group("/venues")
	{
		venues.POST("", h.Create)
		venues.GET("", h.List)
		venues.GET("/:id", h.Get)
		venues.PUT("/:id", h.Update)
		venues.DELETE("/:id", h.Delete)
	}
}

func (h *VenueHandler) Create(c *gin.Context) {
	var req model.CreateVenueRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.APIError{
			Error: model.ErrorDetail{Code: "VALIDATION_ERROR", Message: err.Error()},
		})
		return
	}

	venue, err := h.svc.CreateVenue(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.APIError{
			Error: model.ErrorDetail{Code: "INTERNAL_ERROR", Message: "failed to create venue"},
		})
		return
	}

	c.JSON(http.StatusCreated, model.APIResponse{Data: venue})
}

func (h *VenueHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.APIError{
			Error: model.ErrorDetail{Code: "INVALID_ID", Message: "invalid venue ID"},
		})
		return
	}

	venue, err := h.svc.GetVenue(id)
	if err != nil {
		if errors.Is(err, service.ErrVenueNotFound) {
			c.JSON(http.StatusNotFound, model.APIError{
				Error: model.ErrorDetail{Code: "VENUE_NOT_FOUND", Message: "venue not found"},
			})
			return
		}
		c.JSON(http.StatusInternalServerError, model.APIError{
			Error: model.ErrorDetail{Code: "INTERNAL_ERROR", Message: "internal server error"},
		})
		return
	}

	c.JSON(http.StatusOK, model.APIResponse{Data: venue})
}

func (h *VenueHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))

	venues, total, err := h.svc.ListVenues(page, perPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.APIError{
			Error: model.ErrorDetail{Code: "INTERNAL_ERROR", Message: "internal server error"},
		})
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(perPage)))

	c.JSON(http.StatusOK, model.APIResponse{
		Data: venues,
		Meta: &model.Meta{
			Page:       page,
			PerPage:    perPage,
			Total:      total,
			TotalPages: totalPages,
		},
	})
}

func (h *VenueHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.APIError{
			Error: model.ErrorDetail{Code: "INVALID_ID", Message: "invalid venue ID"},
		})
		return
	}

	var req model.UpdateVenueRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.APIError{
			Error: model.ErrorDetail{Code: "VALIDATION_ERROR", Message: err.Error()},
		})
		return
	}

	venue, err := h.svc.UpdateVenue(id, req)
	if err != nil {
		if errors.Is(err, service.ErrVenueNotFound) {
			c.JSON(http.StatusNotFound, model.APIError{
				Error: model.ErrorDetail{Code: "VENUE_NOT_FOUND", Message: "venue not found"},
			})
			return
		}
		c.JSON(http.StatusInternalServerError, model.APIError{
			Error: model.ErrorDetail{Code: "INTERNAL_ERROR", Message: "internal server error"},
		})
		return
	}

	c.JSON(http.StatusOK, model.APIResponse{Data: venue})
}

func (h *VenueHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.APIError{
			Error: model.ErrorDetail{Code: "INVALID_ID", Message: "invalid venue ID"},
		})
		return
	}

	if err := h.svc.DeleteVenue(id); err != nil {
		if errors.Is(err, service.ErrVenueNotFound) {
			c.JSON(http.StatusNotFound, model.APIError{
				Error: model.ErrorDetail{Code: "VENUE_NOT_FOUND", Message: "venue not found"},
			})
			return
		}
		c.JSON(http.StatusInternalServerError, model.APIError{
			Error: model.ErrorDetail{Code: "INTERNAL_ERROR", Message: "internal server error"},
		})
		return
	}

	c.JSON(http.StatusOK, model.APIResponse{Data: gin.H{"message": "venue deleted"}})
}
