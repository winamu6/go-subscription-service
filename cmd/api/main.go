import (
    "log"

    _ "subscription-service/docs"

    "subscription-service/internal/config"
    "subscription-service/internal/db"
    "subscription-service/internal/model"
    "subscription-service/internal/repository/read_repository"
    "subscription-service/internal/repository/write_repository"
    "subscription-service/internal/router"
    "subscription-service/internal/service"

    swaggerFiles "github.com/swaggo/files"     
    ginSwagger "github.com/swaggo/gin-swagger" 

    "gorm.io/gorm"
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

    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

    log.Printf("Starting server on port %s...", cfg.AppPort)
    if err := r.Run(":" + cfg.AppPort); err != nil {
        log.Fatalf("Server failed to start: %v", err)
    }
}