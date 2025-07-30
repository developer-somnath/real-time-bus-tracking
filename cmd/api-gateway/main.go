package main

import (
	"real-time-bus-tracking/cmd/api-gateway/routes"
	"real-time-bus-tracking/internal/models"
	"real-time-bus-tracking/pkg/logger"

	"github.com/gin-gonic/gin"
)

func main() {
	log := logger.Init("api-gateway")
	defer log.Close()

	// Initialize database (for migrations, if RUN_MIGRATE=true)
	if err := models.InitDB(); err != nil {
		log.Fatalf("Failed to initialize MySQL: %v", err)
	}

	// Initialize Gin router
	r := gin.Default()
	routes.SetupRoutes(r)

	// Run server
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
