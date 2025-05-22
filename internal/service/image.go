package service

import (
	"context"
	"fmt"
	"mime/multipart"

	"go-image-compression/internal/model"
	"go-image-compression/internal/repository"
	"go-image-compression/pkg/broker"
	"go-image-compression/pkg/db"
	"go-image-compression/pkg/resizer"

	"github.com/google/uuid"
)

type (
	ImageService interface {
		Get(ctx context.Context, filter model.ListImageFilter) (multipart.File, error)
		Create(ctx context.Context, fileHeader *multipart.FileHeader) error
		CompressImage(ctx context.Context, imageID string) error
	}

	imageService struct {
		imageRepository repository.ImageRepository
		event           broker.Broker
		compressor      resizer.Compressor
	}
)

func newImageService(imageRepository repository.ImageRepository, event broker.Broker, compressor resizer.Compressor) ImageService {
	return &imageService{
		imageRepository: imageRepository,
		event:           event,
		compressor:      compressor,
	}
}

const codepath = "service/image.go"

func (s *imageService) Get(ctx context.Context, filter model.ListImageFilter) (multipart.File, error) {
	obj, err := s.imageRepository.Get(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", codepath, err)
	}

	return obj, nil
}

func (s *imageService) Create(ctx context.Context, fileHeader *multipart.FileHeader) error {
	image, size, fileName, contentType, err := db.GetFileStat(fileHeader)
	if err != nil {
		return fmt.Errorf("%s: %w", codepath, err)
	}

	imageID := fmt.Sprintf("%s_%s_100", uuid.NewString(), fileName)

	if err = s.imageRepository.Create(ctx, image, size, imageID, contentType); err != nil {
		return fmt.Errorf("%s: %w", codepath, err)
	}

	if err = s.event.Publish("image.created", []byte(imageID)); err != nil {
		return fmt.Errorf("%s: %w", codepath, err)
	}

	return nil
}

func (s *imageService) CompressImage(ctx context.Context, imageID string) error {
	file, err := s.imageRepository.Get(ctx, model.ListImageFilter{
		ID: imageID,
	})
	if err != nil {
		return fmt.Errorf("service.image.CompressImage: %w", err)
	}
	defer file.Close()

	qualities := []float64{0.75, 0.5, 0.25}

	image, format, err := s.compressor.GetImage(file)
	if err != nil {
		return fmt.Errorf("service.image.CompressImage: %w", err)
	}

	for _, quality := range qualities {
		imageID := s.compressor.BuildImageID(imageID, quality)

		compressedImage, err := s.compressor.Compress(image, quality)
		if err != nil {
			return fmt.Errorf("service.image.CompressImage: %w", err)
		}

		reader, size, err := s.compressor.EncodeImage(compressedImage, format)
		if err != nil {
			return fmt.Errorf("service.image.CompressImage: %w", err)
		}

		if err := s.imageRepository.Create(ctx, reader, size, imageID, format); err != nil {
			return fmt.Errorf("service.image.CompressImage: %w", err)
		}
	}

	return nil
}
