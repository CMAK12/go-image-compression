package service

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/png"
	"io"
	"mime/multipart"
	"os"

	"go-image-compression/internal/model"
	"go-image-compression/internal/repository"
	"go-image-compression/pkg/broker"
	"go-image-compression/pkg/resizer"

	"github.com/nordew/go-errx"
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
	image, err := buildImage(fileHeader)
	if err != nil {
		return fmt.Errorf("%s: %w", codepath, err)
	}

	if err = s.imageRepository.Create(ctx, image); err != nil {
		return fmt.Errorf("%s: %w", codepath, err)
	}

	if err = s.event.Publish("image.created", []byte(image.ID)); err != nil {
		return fmt.Errorf("%s: %w", codepath, err)
	}

	return nil
}

func (s *imageService) CompressImage(ctx context.Context, imageID string) error {
	file, err := s.imageRepository.Get(ctx, model.ListImageFilter{
		ID: imageID,
	})
	if err != nil {
		return fmt.Errorf("%s: %w", codepath, err)
	}
	defer file.Close()

	newImageName := imageID[:len(imageID)-4]

	for i := 0.75; i >= 0.25; i -= 0.25 {
		if err := s.compressAndStore(ctx, file, newImageName, i); err != nil {
			return fmt.Errorf("service.image.CompressImage: %w", err)
		}

		if _, err := file.Seek(0, io.SeekStart); err != nil {
			return errx.NewInternal().WithDescriptionAndCause("service.image.CompressImage: ", err)
		}
	}

	return nil
}

func buildImage(header *multipart.FileHeader) (model.Image, error) {
	file, err := header.Open()
	if err != nil {
		return model.Image{}, errx.NewInternal().WithDescriptionAndCause(codepath, err)
	}
	defer file.Close()

	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		return model.Image{}, errx.NewInternal().WithDescriptionAndCause(codepath, fmt.Errorf("missing content type"))
	}

	fileName := header.Filename
	if fileName == "" {
		return model.Image{}, errx.NewInternal().WithDescriptionAndCause(codepath, fmt.Errorf("missing file name"))
	}

	image := model.NewImage(file, header.Size, fileName, contentType)

	return image, nil
}

func (s *imageService) compressAndStore(ctx context.Context, file multipart.File, filename string, percent float64) error {
	decodedImage, _, err := image.Decode(file)
	if err != nil {
		return errx.NewInternal().WithDescriptionAndCause("service.image.compressAndStore: ", err)
	}

	resizedImage, err := s.compressor.Compress(decodedImage, percent)
	if err != nil {
		return errx.NewInternal().WithDescriptionAndCause("service.image.compressAndStore: ", err)
	}

	newFileName := s.compressor.BuildImageID(filename, percent)

	tmpFile, stat, err := parseFile(resizedImage, newFileName)
	if err != nil {
		return errx.NewInternal().WithDescriptionAndCause("service.image.compressAndStore: ", err)
	}
	defer tmpFile.Close()

	modelImage := model.NewImageWithID(tmpFile, stat.Size(), newFileName, "image/png")

	if err = s.imageRepository.Create(ctx, modelImage); err != nil {
		return errx.NewInternal().WithDescriptionAndCause("service.image.compressAndStore: ", err)
	}

	return nil
}

func parseFile(img image.Image, filename string) (*os.File, os.FileInfo, error) {
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, nil, errx.NewInternal().WithDescriptionAndCause("service.image.parseFile: ", err)
	}

	tmpFile, err := os.CreateTemp("", filename+"-*.png")
	if err != nil {
		return nil, nil, errx.NewInternal().WithDescriptionAndCause("service.image.parseFile: ", err)
	}

	if _, err := io.Copy(tmpFile, &buf); err != nil {
		tmpFile.Close()
		return nil, nil, errx.NewInternal().WithDescriptionAndCause("service.image.parseFile: ", err)
	}

	if _, err := tmpFile.Seek(0, io.SeekStart); err != nil {
		tmpFile.Close()
		return nil, nil, errx.NewInternal().WithDescriptionAndCause("service.image.parseFile: ", err)
	}

	stat, err := tmpFile.Stat()
	if err != nil {
		return nil, nil, errx.NewInternal().WithDescriptionAndCause("service.image.parseFile: ", err)
	}

	return tmpFile, stat, nil
}
