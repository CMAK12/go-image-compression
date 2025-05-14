package service

import (
	"go-image-compression/internal/broker"
	"go-image-compression/internal/repository"
)

type Services struct {
	ImageService ImageService
}

func NewService(repositories repository.Repositories, imageBroker broker.ImageProducer) Services {
	return Services{
		ImageService: newImageService(repositories.ImageRepository, imageBroker),
	}
}
