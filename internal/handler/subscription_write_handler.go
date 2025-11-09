package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/winnamu6/go-subscription-service/internal/model"
	"github.com/winnamu6/go-subscription-service/internal/service"
)

type SubscriptionWriteHandler struct {
	commandService service.SubscriptionCommandService
}

func NewSubscriptionWriteHandler(commandService service.SubscriptionCommandService) *SubscriptionWriteHandler {
	return &SubscriptionWriteHandler{commandService: commandService}
}

// Create godoc
// @Summary      Create a new subscription
// @Description  Creates a new subscription record
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        subscription  body      model.CreateSubscriptionRequest  true  "Subscription Data"
// @Success      201  {object}  model.Subscription
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /subscriptions [post]
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

// Update godoc
// @Summary      Update a subscription
// @Description  Updates an existing subscription by ID
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        id            path      int                              true  "Subscription ID"
// @Param        subscription  body      model.UpdateSubscriptionRequest  true  "Updated Subscription Data"
// @Success      200  {object}  model.Subscription
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /subscriptions/{id} [put]
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

// Delete godoc
// @Summary      Delete a subscription
// @Description  Deletes a subscription by its ID
// @Tags         subscriptions
// @Produce      json
// @Param        id   path      int  true  "Subscription ID"
// @Success      204  {string}  string  "No Content"
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /subscriptions/{id} [delete]
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
