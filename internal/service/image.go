package service

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"

	"go-image-compression/internal/broker"
	"go-image-compression/internal/consts"
	"go-image-compression/internal/model"
	"go-image-compression/internal/repository"

	"github.com/nordew/go-errx"
)

type (
	ImageService interface {
		Get(ctx context.Context, filter model.ListImageFilter) (io.Reader, error)
		Create(ctx context.Context, fileHeader *multipart.FileHeader) error
	}

	imageService struct {
		imageRepository repository.ImageRepository
		imageProducer   broker.ImageProducer
	}
)

func newImageService(imageRepository repository.ImageRepository, imageProducer broker.ImageProducer) ImageService {
	return &imageService{
		imageRepository: imageRepository,
		imageProducer:   imageProducer,
	}
}

const codepath = "service/image.go"

func (s *imageService) Get(ctx context.Context, filter model.ListImageFilter) (io.Reader, error) {
	obj, err := s.imageRepository.Get(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", codepath, err)
	}

	return obj, nil
}

func (s *imageService) Create(ctx context.Context, fileHeader *multipart.FileHeader) error {
	file, err := fileHeader.Open()
	if err != nil {
		return errx.NewInternal().WithDescriptionAndCause(codepath, err)
	}
	defer file.Close()

	image, err := buildImage(fileHeader, file)
	if err != nil {
		return fmt.Errorf("%s: %w", codepath, err)
	}

	if err = s.imageRepository.Create(ctx, image); err != nil {
		return fmt.Errorf("%s: %w", codepath, err)
	}

	if err = s.imageProducer.Publish(ctx, image.ID); err != nil {
		return fmt.Errorf("%s: %w", codepath, err)
	}

	return nil
}

func buildImage(header *multipart.FileHeader, file multipart.File) (model.Image, error) {
	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		return model.Image{}, errx.NewInternal().WithDescriptionAndCause(codepath, fmt.Errorf("missing content type"))
	}

	fileName := header.Filename
	if fileName == "" {
		return model.Image{}, errx.NewInternal().WithDescriptionAndCause(codepath, fmt.Errorf("missing file name"))
	}

	image := model.NewImage(file, header.Size, consts.BucketName, fileName, contentType)
	return image, nil
}
