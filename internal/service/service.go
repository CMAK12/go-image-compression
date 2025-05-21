package service

import (
	"go-image-compression/internal/repository"
	"go-image-compression/pkg/broker"
	"go-image-compression/pkg/resizer"
)

type (
	Services struct {
		ImageService ImageService
	}
)

func NewService(repositories repository.Repositories, event broker.Broker, compressor resizer.Compressor) Services {
	return Services{
		ImageService: newImageService(repositories.ImageRepository, event, compressor),
	}
}
