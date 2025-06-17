package service

import (
	"go-image-compression/internal/config"
	"go-image-compression/internal/repository"
	"go-image-compression/pkg/broker"
	"go-image-compression/pkg/resizer"
)

type Services struct {
	ImageService ImageService
}

func NewService(
	repositories repository.Repositories,
	broker broker.Broker,
	compressor resizer.Compressor,
	cfg *config.Config,
) Services {
	return Services{
		ImageService: newImageService(&repositories.ImageRepository, broker, compressor, cfg.Topic),
	}
}
