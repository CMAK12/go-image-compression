package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go-image-compression/internal/config"
	"go-image-compression/internal/model"
	"go-image-compression/internal/service"
	"go-image-compression/pkg/broker"

	"go.uber.org/zap"
)

type (
	ImageService interface {
		CompressImage(ctx context.Context, payload model.Payload) error
	}

	Consumer struct {
		broker   broker.Broker
		services ImageService
		topics   config.Topic
		logger   *zap.Logger
	}
)

func Start(broker broker.Broker, services service.Services, cfg *config.Config, logger *zap.Logger) error {
	c := Consumer{
		broker:   broker,
		services: &services.ImageService,
		topics:   cfg.Topic,
		logger:   logger,
	}

	_, err := c.broker.Subscribe(c.topics.ImageCreated, c.createImage)
	if err != nil {
		c.logger.Error("consumer.Start: ", zap.Error(err))
		return err
	}

	return nil
}

func (c *Consumer) createImage(msg *broker.Message) error {
	var payload model.Payload

	if err := json.Unmarshal(msg.Data, &payload); err != nil {
		c.logger.Error("consumer.createImage: %v", zap.Error(err))
		return fmt.Errorf("consumer.createImage: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := c.services.CompressImage(ctx, payload); err != nil {
		c.logger.Error("consumer.createImage: ", zap.Error(err))
		return fmt.Errorf("consumer.createImage: %w", err)
	}

	return nil
}
