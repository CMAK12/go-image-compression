package grpc_handler

import (
	"context"
	"io"
	"mime/multipart"

	"go-image-compression/internal/model"
	"go-image-compression/internal/service"
	pb "go-image-compression/pkg/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type (
	imageService interface {
		Get(ctx context.Context, filter model.ListImageFilter) (multipart.File, error)
		Create(ctx context.Context, fileHeader *multipart.FileHeader) error
	}

	CompressionHandler struct {
		pb.UnimplementedImageServiceServer
		imageService imageService
	}
)

func NewCompressionHandler(svc service.Services) *CompressionHandler {
	return &CompressionHandler{
		imageService: svc.ImageService,
	}
}

func (h *CompressionHandler) GetImage(ctx context.Context, req *pb.GetImageRequest) (*pb.GetImageResponse, error) {
	image, err := h.imageService.Get(ctx, model.ListImageFilter{
		ID: req.GetId(),
	})
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(image)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "grpc.GetImage: %v", err)
	}

	return &pb.GetImageResponse{
		Data: data,
	}, nil
}

func (h *CompressionHandler) UploadImage(ctx context.Context, req *pb.UploadImageRequest) (*emptypb.Empty, error) {
	if req.GetData() == nil || len(req.GetData()) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "grpc.UploadImage: data cannot be empty")
	}

	if err := h.imageService.Create(ctx, nil); err != nil {
		return nil, status.Errorf(codes.Internal, "grpc.UploadImage: %v", err)
	}

	return &emptypb.Empty{}, nil
}
