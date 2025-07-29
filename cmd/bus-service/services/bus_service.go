package services

import (
	"context"
	"real-time-bus-tracking/cmd/bus-service/types"
	"real-time-bus-tracking/internal/events/kafka"
	"real-time-bus-tracking/internal/events/types"
	"real-time-bus-tracking/internal/models/schemas"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type BusService struct {
	db       *gorm.DB
	redis    *redis.Client
	producer *kafka.Producer
}

func NewBusService(db *gorm.DB, redis *redis.Client) *BusService {
	return &BusService{
		db:       db,
		redis:    redis,
		producer: kafka.NewProducer(),
	}
}

func (s *BusService) CreateBus(ctx context.Context, bus *types.Bus) error {
	dbBus := schemas.Bus{
		LicensePlate:      bus.LicensePlate,
		WheelchairEnabled: bus.WheelchairEnabled,
	}
	if err := s.db.Create(&dbBus).Error; err != nil {
		return err
	}
	bus.ID = dbBus.ID

	event := types.BusCreatedEvent{
		BusID:             bus.ID,
		LicensePlate:      bus.LicensePlate,
		WheelchairEnabled: bus.WheelchairEnabled,
		CreatedAt:         time.Now(),
	}
	return s.producer.PublishEvent(ctx, "bus.created", event)
}
