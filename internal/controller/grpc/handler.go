package grpc_handler

import (
	"context"
	"io"

	"go-image-compression/internal/model"
	"go-image-compression/internal/service"
	pb "go-image-compression/pkg/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type CompressionHandler struct {
	pb.UnimplementedImageServiceServer
	svc service.Services
}

func NewCompressionHandler(svc service.Services) *CompressionHandler {
	return &CompressionHandler{
		svc: svc,
	}
}

func (h *CompressionHandler) GetImage(ctx context.Context, req *pb.GetImageRequest) (*pb.GetImageResponse, error) {
	image, err := h.svc.ImageService.Get(ctx, model.ListImageFilter{
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

	if err := h.svc.ImageService.Create(ctx, nil); err != nil {
		return nil, status.Errorf(codes.Internal, "grpc.UploadImage: %v", err)
	}

	return &emptypb.Empty{}, nil
}
