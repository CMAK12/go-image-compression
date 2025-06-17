package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"go-image-compression/internal/config"
	"go-image-compression/internal/model"
	"go-image-compression/pkg/broker"
	"go-image-compression/pkg/resizer"

	"github.com/google/uuid"
)

type (
	ImageRepository interface {
		Get(ctx context.Context, filter model.ListImageFilter) (multipart.File, error)
		Create(ctx context.Context, img io.Reader, size int64, imageID, contentType string) error
	}

	ImageService struct {
		imageRepository ImageRepository
		broker          broker.Broker
		compressor      resizer.Compressor
		topics          config.Topic
	}
)

func newImageService(
	imageRepository ImageRepository,
	broker broker.Broker,
	compressor resizer.Compressor,
	topics config.Topic,
) ImageService {
	return ImageService{
		imageRepository: imageRepository,
		broker:          broker,
		compressor:      compressor,
		topics:          topics,
	}
}

func (s *ImageService) Get(ctx context.Context, filter model.ListImageFilter) (multipart.File, error) {
	obj, err := s.imageRepository.Get(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("service.image.GET: %w", err)
	}

	return obj, nil
}

func (s *ImageService) Create(ctx context.Context, body []byte) error {
	if len(body) == 0 {
		return fmt.Errorf("service.image.Create: body is empty")
	}

	reader := bytes.NewReader(body)

	imageID := fmt.Sprintf("%s_100", uuid.NewString())
	mimeType := http.DetectContentType(body)

	if err := s.imageRepository.Create(ctx, reader, reader.Size(), imageID, mimeType); err != nil {
		return fmt.Errorf("service.image.Create: %w", err)
	}

	payload := model.NewPayload(imageID, mimeType)
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("service.image.Create: %w", err)
	}

	if ack, err := s.broker.Publish(s.topics.ImageCreated, payloadBytes); err != nil || !ack.Success {
		return fmt.Errorf("service.image.Create: %w, ack: %t", err, ack.Success)
	}

	return nil
}

func (s *ImageService) CompressImage(ctx context.Context, payload model.Payload) error {
	if payload.ImageID == "" || payload.MIMEType == "" {
		return fmt.Errorf("service.image.CompressImage: imageID or MIME type is empty")
	}

	file, err := s.imageRepository.Get(ctx, model.ListImageFilter{
		ID: payload.ImageID,
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
		imageID := s.compressor.BuildImageID(payload.ImageID, quality)

		compressedImage, err := s.compressor.Compress(image, quality)
		if err != nil {
			return fmt.Errorf("service.image.CompressImage: %w", err)
		}

		reader, size, err := s.compressor.EncodeImage(compressedImage, format)
		if err != nil {
			return fmt.Errorf("service.image.CompressImage: %w", err)
		}

		if err := s.imageRepository.Create(ctx, reader, size, imageID, payload.MIMEType); err != nil {
			return fmt.Errorf("service.image.CompressImage: %w", err)
		}
	}

	return nil
}
