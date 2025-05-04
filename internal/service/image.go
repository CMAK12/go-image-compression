package service

import (
	"context"
	"go-image-compression/internal/model"
	"go-image-compression/internal/repository"
)

type (
	ImageService interface {
		List(ctx context.Context, filter model.ListImageFilter) ([]model.Image, error)
		Create(ctx context.Context, image model.Image) error
	}

	imageService struct {
		imageRepository repository.ImageRepository
	}
)

func newImageService(imageRepository repository.ImageRepository) ImageService {
	return &imageService{
		imageRepository: imageRepository,
	}
}

func (s *imageService) List(ctx context.Context, filter model.ListImageFilter) ([]model.Image, error) {
	return nil, nil
}

func (s *imageService) Create(ctx context.Context, image model.Image) error {
	return nil
}
