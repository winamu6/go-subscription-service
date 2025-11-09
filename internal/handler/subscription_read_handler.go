package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/winnamu6/go-subscription-service/internal/service"
)

type SubscriptionReadHandler struct {
	queryService service.SubscriptionQueryService
}

func NewSubscriptionReadHandler(queryService service.SubscriptionQueryService) *SubscriptionReadHandler {
	return &SubscriptionReadHandler{queryService: queryService}
}

// GetAll godoc
// @Summary      Get all subscriptions
// @Description  Returns a list of all subscriptions in the system
// @Tags         subscriptions
// @Produce      json
// @Success      200  {array}   model.Subscription
// @Failure      500  {object}  map[string]string
// @Router       /subscriptions [get]
func (h *SubscriptionReadHandler) GetAll(c *gin.Context) {
	ctx := c.Request.Context()

	subs, err := h.queryService.GetAll(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, subs)
}

// GetByID godoc
// @Summary      Get subscription by ID
// @Description  Returns a subscription by its unique ID
// @Tags         subscriptions
// @Produce      json
// @Param        id   path      int  true  "Subscription ID"
// @Success      200  {object}  model.Subscription
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /subscriptions/{id} [get]
func (h *SubscriptionReadHandler) GetByID(c *gin.Context) {
	ctx := c.Request.Context()
	idParam := c.Param("id")

	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	sub, err := h.queryService.GetByID(ctx, uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if sub == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "subscription not found"})
		return
	}

	c.JSON(http.StatusOK, sub)
}

// GetByUserID godoc
// @Summary      Get subscriptions by User ID
// @Description  Returns all subscriptions associated with a given user
// @Tags         subscriptions
// @Produce      json
// @Param        user_id  path      string  true  "User ID"
// @Success      200  {array}   model.Subscription
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /subscriptions/user/{user_id} [get]
func (h *SubscriptionReadHandler) GetByUserID(c *gin.Context) {
	ctx := c.Request.Context()
	userID := c.Param("user_id")

	subs, err := h.queryService.GetByUserID(ctx, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, subs)
}

// SumPriceByFilter godoc
// @Summary      Get total subscription price by filters
// @Description  Calculates the total price of subscriptions filtered by user_id, service_name, start_date and end_date
// @Tags         subscriptions
// @Produce      json
// @Param        user_id      query     string  false  "User ID"
// @Param        service_name query     string  false  "Service Name"
// @Param        start_date   query     string  true   "Start Date (RFC3339 format)"
// @Param        end_date     query     string  true   "End Date (RFC3339 format)"
// @Success      200  {object}  map[string]float64
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /subscriptions/sum [get]
func (h *SubscriptionReadHandler) SumPriceByFilter(c *gin.Context) {
	ctx := c.Request.Context()
	userID := c.Query("user_id")
	serviceName := c.Query("service_name")
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if startDateStr == "" || endDateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start_date and end_date are required"})
		return
	}

	startDate, err := time.Parse(time.RFC3339, startDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_date format"})
		return
	}

	endDate, err := time.Parse(time.RFC3339, endDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_date format"})
		return
	}

	var userIDPtr, serviceNamePtr *string
	if userID != "" {
		userIDPtr = &userID
	}
	if serviceName != "" {
		serviceNamePtr = &serviceName
	}

	sum, err := h.queryService.SumPriceByFilter(ctx, userIDPtr, serviceNamePtr, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"total_price": sum})
}