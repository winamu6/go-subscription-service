package main

import (
    "log"

    "github.com/winnamu6/go-subscription-service/docs"

    "github.com/winnamu6/go-subscription-service/internal/config"
    "github.com/winnamu6/go-subscription-service/internal/db"
    "github.com/winnamu6/go-subscription-service/internal/model"
    "github.com/winnamu6/go-subscription-service/internal/repository/read_repository"
    "github.com/winnamu6/go-subscription-service/internal/repository/write_repository"
    "github.com/winnamu6/go-subscription-service/internal/router"
    "github.com/winnamu6/go-subscription-service/internal/service"

    swaggerFiles "github.com/swaggo/files"     
    ginSwagger "github.com/swaggo/gin-swagger" 

    "gorm.io/gorm"
    "github.com/gin-gonic/gin"
)

func main() {
    cfg := config.Load()

    docs.SwaggerInfo.Title = "Subscription Service API"
    docs.SwaggerInfo.Description = "API documentation for Subscription Service"
    docs.SwaggerInfo.Version = "1.0"
    docs.SwaggerInfo.Host = "localhost:" + cfg.AppPort
    docs.SwaggerInfo.BasePath = "/"

    database := db.Connect(cfg)
    runMigrations(database)

    readRepo := read_repository.NewSubscriptionReadRepo(database)
    writeRepo := write_repository.NewSubscriptionWriteRepo(database)

    readSvc := service.NewSubscriptionQueryService(readRepo)
    writeSvc := service.NewSubscriptionCommandService(writeRepo, readSvc)

    r := router.NewRouter(readSvc, writeSvc)

    r.GET("/", func(c *gin.Context) {
        c.Redirect(302, "/swagger/index.html")
    })

    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

    log.Printf("Starting server on port %s...", cfg.AppPort)
    if err := r.Run(":" + cfg.AppPort); err != nil {
        log.Fatalf("Server failed to start: %v", err)
    }
}

func runMigrations(db *gorm.DB) {
	log.Println("Running migrations...")
	if err := db.AutoMigrate(&model.Subscription{}); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	log.Println("Migrations completed")
}