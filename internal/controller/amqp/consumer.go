package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"go-image-compression/internal/model"
	"go-image-compression/internal/service"
	"go-image-compression/pkg/broker"
)

type (
	imageService interface {
		CompressImage(ctx context.Context, payload model.Payload) error
	}

	Consumer struct {
		broker   broker.Broker
		services imageService
	}
)

func Start(broker broker.Broker, services service.Services) error {
	c := Consumer{
		broker:   broker,
		services: services.ImageService,
	}

	err := c.broker.Subscribe("image.created", c.createImage)
	if err != nil {
		return err
	}

	return nil
}

func (c *Consumer) createImage(msg *broker.Message) error {
	var payload model.Payload

	if err := json.Unmarshal(msg.Data, &payload); err != nil {
		log.Printf("consumer.createImage: %v", err)
		return fmt.Errorf("consumer.createImage: %w", err)
	}

	if err := c.services.CompressImage(context.Background(), payload); err != nil {
		log.Println("consumer.createImage: ", err)
		return fmt.Errorf("consumer.createImage: %w", err)
	}

	return nil
}
