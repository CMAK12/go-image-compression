package service

import (
	"go-image-compression/internal/repository"
)

type Services struct {
	ImageService ImageService
}

func NewService(repositories repository.Repositories) Services {
	return Services{
		ImageService: newImageService(repositories.ImageRepository),
	}
}
