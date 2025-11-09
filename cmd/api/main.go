package main

import (
	"log"
	"subscription-service/internal/config"
	"subscription-service/internal/db"
	"subscription-service/internal/model"
	"subscription-service/internal/repository/read_repository"
	"subscription-service/internal/repository/write_repository"
	"subscription-service/internal/router"
	"subscription-service/internal/service"

	"gorm.io/gorm"
)

func main() {
	cfg := config.Load()

	database := db.Connect(cfg)

	runMigrations(database)

	readRepo := read_repository.NewSubscriptionReadRepo(database)
	writeRepo := write_repository.NewSubscriptionWriteRepo(database)

	readSvc := service.NewSubscriptionQueryService(readRepo)
	writeSvc := service.NewSubscriptionCommandService(writeRepo, readSvc)

	r := router.NewRouter(readSvc, writeSvc)

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