package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"subscription-service/internal/service"
)

type SubscriptionReadHandler struct {
	queryService service.SubscriptionQueryService
}

func NewSubscriptionReadHandler(queryService service.SubscriptionQueryService) *SubscriptionReadHandler {
	return &SubscriptionReadHandler{queryService: queryService}
}

func (h *SubscriptionReadHandler) GetAll(c *gin.Context) {
	ctx := c.Request.Context()

	subs, err := h.queryService.GetAll(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, subs)
}

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