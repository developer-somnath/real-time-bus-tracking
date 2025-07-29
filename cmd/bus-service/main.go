package main

import (
	"log"
	"real-time-bus-tracking/cmd/bus-service/handlers"
	"real-time-bus-tracking/cmd/bus-service/services"
	"real-time-bus-tracking/internal/models"
	"real-time-bus-tracking/pkg/logger"
)

func main() {
	logger.Init()

	if err := models.InitDB(); err != nil {
		log.Fatalf("Failed to initialize MySQL: %v", err)
	}
	if err := models.InitRedis(); err != nil {
		log.Fatalf("Failed to initialize Redis: %v", err)
	}

	busService := services.NewBusService(models.DB, models.Redis)
	go handlers.StartGRPCServer(":50051", busService)
	handlers.StartKafkaConsumer()

	select {}
}
