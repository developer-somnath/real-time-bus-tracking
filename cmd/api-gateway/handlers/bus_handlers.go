package handlers

import (
	"net/http"
	"real-time-bus-tracking/internal/events/kafka"
	"real-time-bus-tracking/internal/events/types"

	"github.com/gin-gonic/gin"
)

func CreateBus(c *gin.Context) {
	var bus struct {
		LicensePlate      string `json:"license_plate" binding:"required"`
		WheelchairEnabled bool   `json:"wheelchair_enabled"`
	}
	if err := c.ShouldBindJSON(&bus); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	event := types.BusCreatedEvent{
		LicensePlate:      bus.LicensePlate,
		WheelchairEnabled: bus.WheelchairEnabled,
	}
	producer := kafka.NewProducer()
	if err := producer.PublishEvent(c, "bus.created", event); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to publish event"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "Bus created"})
}

func GetBus(c *gin.Context) {
	// Implement gRPC call to BusService
	c.JSON(http.StatusOK, gin.H{"id": c.Param("id")})
}
