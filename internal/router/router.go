package router

import (
	"github.com/gin-gonic/gin"
	"github.com/winnamu6/go-subscription-service/internal/handler"
	"github.com/winnamu6/go-subscription-service/internal/service"
)

func NewRouter(readSvc service.SubscriptionQueryService, writeSvc service.SubscriptionCommandService) *gin.Engine {
	r := gin.Default()

	readHandler := handler.NewSubscriptionReadHandler(readSvc)
	writeHandler := handler.NewSubscriptionWriteHandler(writeSvc)

	subscriptions := r.Group("/subscriptions")
	{
		subscriptions.GET("", readHandler.GetAll)
		subscriptions.GET("/:id", readHandler.GetByID)
		subscriptions.GET("/user/:user_id", readHandler.GetByUserID)
		subscriptions.GET("/sum", readHandler.SumPriceByFilter)

		subscriptions.POST("", writeHandler.Create)
		subscriptions.PUT("/:id", writeHandler.Update)
		subscriptions.DELETE("/:id", writeHandler.Delete)
	}

	return r
}
