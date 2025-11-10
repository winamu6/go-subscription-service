package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/winnamu6/go-subscription-service/internal/logger"
	"github.com/winnamu6/go-subscription-service/internal/service"
)

type SubscriptionReadHandler struct {
	queryService service.SubscriptionQueryService
}

func NewSubscriptionReadHandler(queryService service.SubscriptionQueryService) *SubscriptionReadHandler {
	return &SubscriptionReadHandler{queryService: queryService}
}

// GetAll godoc
func (h *SubscriptionReadHandler) GetAll(c *gin.Context) {
	log := logger.Get()
	log.Info("Handler: GetAll() called")

	ctx := c.Request.Context()
	subs, err := h.queryService.GetAll(ctx)
	if err != nil {
		log.Errorf("Failed to get all subscriptions: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Infof("Successfully retrieved %d subscriptions", len(subs))
	c.JSON(http.StatusOK, subs)
}

// GetByID godoc
func (h *SubscriptionReadHandler) GetByID(c *gin.Context) {
	log := logger.Get()
	idParam := c.Param("id")
	log.Infof("Handler: GetByID() called for ID=%s", idParam)

	ctx := c.Request.Context()
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		log.Warnf("Invalid subscription ID: %s", idParam)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	sub, err := h.queryService.GetByID(ctx, uint(id))
	if err != nil {
		log.Errorf("Failed to get subscription ID=%d: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if sub == nil {
		log.Warnf("Subscription ID=%d not found", id)
		c.JSON(http.StatusNotFound, gin.H{"error": "subscription not found"})
		return
	}

	log.Infof("Subscription ID=%d retrieved successfully", id)
	c.JSON(http.StatusOK, sub)
}

// GetByUserID godoc
func (h *SubscriptionReadHandler) GetByUserID(c *gin.Context) {
	log := logger.Get()
	userID := c.Param("user_id")
	log.Infof("Handler: GetByUserID() called for user_id=%s", userID)

	ctx := c.Request.Context()
	subs, err := h.queryService.GetByUserID(ctx, userID)
	if err != nil {
		log.Errorf("Failed to get subscriptions for user_id=%s: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Infof("Retrieved %d subscriptions for user_id=%s", len(subs), userID)
	c.JSON(http.StatusOK, subs)
}

// SumPriceByFilter godoc
func (h *SubscriptionReadHandler) SumPriceByFilter(c *gin.Context) {
	log := logger.Get()
	userID := c.Query("user_id")
	serviceName := c.Query("service_name")
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	log.Infof("Handler: SumPriceByFilter() called | user_id=%s service_name=%s start=%s end=%s",
		userID, serviceName, startDateStr, endDateStr)

	ctx := c.Request.Context()

	if startDateStr == "" || endDateStr == "" {
		log.Warn("Missing required parameters: start_date or end_date")
		c.JSON(http.StatusBadRequest, gin.H{"error": "start_date and end_date are required"})
		return
	}

	startDate, err := time.Parse(time.RFC3339, startDateStr)
	if err != nil {
		log.Warnf("Invalid start_date format: %s", startDateStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_date format"})
		return
	}

	endDate, err := time.Parse(time.RFC3339, endDateStr)
	if err != nil {
		log.Warnf("Invalid end_date format: %s", endDateStr)
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
		log.Errorf("Failed to calculate sum for filter: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Infof("SumPriceByFilter() result: %.2f", sum)
	c.JSON(http.StatusOK, gin.H{"total_price": sum})
}