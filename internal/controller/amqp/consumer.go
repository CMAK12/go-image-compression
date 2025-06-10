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

func Start(broker broker.Broker, services service.Services) error {
	c := Consumer{
		broker:   broker,
		services: services,
	}

	err := c.broker.Subscribe(
		"image.created",
		c.createImage,
	)
	if err != nil {
		return err
	}

	return nil
}

func (c *Consumer) createImage(msg *broker.Message) error {
	imageID := string(msg.Data)

	if err := c.services.ImageService.CompressImage(context.Background(), imageID); err != nil {
		log.Println("consumer.createImage: ", err)
		return fmt.Errorf("consumer.createImage: %w", err)
	}

	return nil
}
