package handler

import (
	"net/http"
	"online-cinema-analytics/internal/application/dto"
	"online-cinema-analytics/internal/application/event_service"

	"github.com/gin-gonic/gin"
)

type EventHandler struct {
	svc *eventservice.EventService
}

func NewEventHandler(svc *eventservice.EventService) *EventHandler {
	return &EventHandler{
		svc: svc,
	}
}

func (h *EventHandler) Create(c *gin.Context) {
	var req dto.CreateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	resp, err := h.svc.Publish(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}
