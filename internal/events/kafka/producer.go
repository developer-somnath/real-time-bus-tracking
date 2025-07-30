package kafka

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"real-time-bus-tracking/internal/utils"
	"time"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer() *Producer {
	brokers := []string{os.Getenv("KAFKA_BROKER_ADDR")}
	config := kafka.WriterConfig{
		Brokers:      brokers,
		BatchTimeout: utils.GetEnvDuration("KAFKA_BATCH_TIMEOUT", 20*time.Millisecond),
		MaxAttempts:  utils.GetEnvInt("KAFKA_MAX_RETRIES", 5),
	}
	writer := kafka.NewWriter(config)
	return &Producer{writer: writer}
}

func (p *Producer) PublishEvent(ctx context.Context, topic string, event interface{}) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}
	msg := kafka.Message{
		Topic: topic,
		Value: data,
	}
	for attempt := 1; attempt <= utils.GetEnvInt("KAFKA_MAX_RETRIES", 5); attempt++ {
		err = p.writer.WriteMessages(ctx, msg)
		if err == nil {
			log.Printf("Published event to topic %s: %v", topic, event)
			return nil
		}
		log.Printf("Kafka publish attempt %d/%d failed: %v", attempt, utils.GetEnvInt("KAFKA_MAX_RETRIES", 5), err)
		if attempt < utils.GetEnvInt("KAFKA_MAX_RETRIES", 5) {
			time.Sleep(utils.GetEnvDuration("KAFKA_RETRY_BACKOFF", 500*time.Millisecond))
		}
	}
	return err
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
