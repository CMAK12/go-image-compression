package grpc_handler

import (
	"context"

	"go-image-compression/internal/model"
	"go-image-compression/pkg/pb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *CompressionHandler) GetImage(ctx context.Context, req *pb.GetImageRequest) (*pb.GetImageResponse, error) {
	image, err := h.imageService.Get(ctx, model.ListImageFilter{
		ID: req.GetId(),
	})
	if err != nil {
		return nil, err
	}

	bytes, err := h.converter.ConvertToBytes(image)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "grpc.GetImage: %v", err)
	}

	return &pb.GetImageResponse{
		Data: bytes,
	}, nil
}

func (h *CompressionHandler) UploadImage(ctx context.Context, req *pb.UploadImageRequest) (*pb.Empty, error) {
	if err := h.imageService.Create(ctx, req.Data); err != nil {
		return nil, status.Errorf(codes.Internal, "grpc.UploadImage: %v", err)
	}

	return &pb.Empty{}, nil
}
