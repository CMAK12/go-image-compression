package consumer

import (
	"context"
	"fmt"
	"log"

	"go-image-compression/internal/service"
	"go-image-compression/pkg/broker"
)

type (
	Consumer struct {
		broker   broker.Broker
		services service.Services
	}
)

func NewConsumer(broker broker.Broker, services service.Services) (Consumer, error) {
	return Consumer{
		broker:   broker,
		services: services,
	}, nil
}

func (c *Consumer) Start() error {
	err := c.broker.Subscribe(
		"image.created",
		c.handleMessage,
	)
	if err != nil {
		return err
	}

	return nil
}

func (c *Consumer) handleMessage(msg *broker.Message) error {
	imageID := string(msg.Data)

	if err := c.services.ImageService.CompressImage(context.Background(), imageID); err != nil {
		log.Println("consumer.handleMessage: ", err)
		return fmt.Errorf("consumer.handleMessage: %w", err)
	}

	return nil
}
