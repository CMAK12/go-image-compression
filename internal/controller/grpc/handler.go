package grpc_handler

import (
	"context"
	"mime/multipart"

	"go-image-compression/internal/model"
	"go-image-compression/internal/service"
	"go-image-compression/pkg/converter"
	"go-image-compression/pkg/pb"

	"go.uber.org/zap"
)

type (
	ImageService interface {
		Get(ctx context.Context, filter model.ListImageFilter) (multipart.File, error)
		Create(ctx context.Context, body []byte) error
	}

	CompressionHandler struct {
		pb.UnimplementedImageServiceServer
		imageService ImageService
		logger       *zap.Logger
		converter    converter.Converter
	}
)

func NewCompressionHandler(
	svc service.Services,
	logger *zap.Logger,
	converter converter.Converter,
) *CompressionHandler {
	return &CompressionHandler{
		imageService: &svc.ImageService,
		logger:       logger,
		converter:    converter,
	}
}
