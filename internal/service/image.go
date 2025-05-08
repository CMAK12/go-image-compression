package service

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"

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
	}
)

func newImageService(imageRepository repository.ImageRepository) ImageService {
	return &imageService{
		imageRepository: imageRepository,
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

	contentType := fileHeader.Header.Get("Content-Type")
	fileName := fileHeader.Filename

	image := model.NewImage(file, fileHeader.Size, consts.FullImageBucket, fileName, contentType)

	err = s.imageRepository.Create(ctx, image)
	if err != nil {
		return errx.NewInternal().WithDescriptionAndCause(codepath, err)
	}

	return nil
}
