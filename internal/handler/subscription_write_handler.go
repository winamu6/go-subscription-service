package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"subscription-service/internal/model"
	"subscription-service/internal/service"
)

type SubscriptionWriteHandler struct {
	commandService service.SubscriptionCommandService
}

func NewSubscriptionWriteHandler(commandService service.SubscriptionCommandService) *SubscriptionWriteHandler {
	return &SubscriptionWriteHandler{commandService: commandService}
}

func (h *SubscriptionWriteHandler) Create(c *gin.Context) {
	ctx := c.Request.Context()
	var req model.CreateSubscriptionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sub, err := h.commandService.Create(ctx, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, sub)
}

func (h *SubscriptionWriteHandler) Update(c *gin.Context) {
	ctx := c.Request.Context()
	idParam := c.Param("id")

	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req model.UpdateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sub, err := h.commandService.Update(ctx, uint(id), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, sub)
}

func (h *SubscriptionWriteHandler) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	idParam := c.Param("id")

	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.commandService.Delete(ctx, uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}