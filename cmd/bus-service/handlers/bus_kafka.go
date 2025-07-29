package handlers

import (
	"context"
	"log"
	"real-time-bus-tracking/internal/events/kafka"
	"real-time-bus-tracking/internal/events/types"
)

func StartKafkaConsumer() {
	consumer := kafka.NewConsumer("driver.location.updated", "bus-service-group")
	ctx := context.Background()
	consumer.ConsumeDriverLocationUpdated(ctx, func(event types.DriverLocationUpdatedEvent) error {
		log.Printf("Processing DriverLocationUpdated: BusID=%d, Lat=%f", event.BusID, event.Latitude)
		return nil
	})
}
