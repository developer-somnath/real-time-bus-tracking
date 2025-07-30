package routes

import (
	"real-time-bus-tracking/cmd/api-gateway/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.Use(middlewares.AuthMiddleware())
	r.Use(middlewares.LoggingMiddleware())

	// // Bus routes
	// r.POST("/buses", handlers.CreateBus)
	// r.GET("/buses/:id", handlers.GetBus)

	// // Route routes
	// r.POST("/routes", handlers.CreateRoute)
	// r.GET("/routes/:id/buses", handlers.GetBusesByRoute)

	// // Trip routes
	// r.POST("/trips", handlers.CreateTrip)
}
