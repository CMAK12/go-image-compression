package grpc

import (
	"context"
	"io"

	"go-image-compression/internal/model"
	"go-image-compression/internal/service"
	pb "go-image-compression/pkg/proto"

	"github.com/nordew/go-errx"
	"google.golang.org/protobuf/types/known/emptypb"
)

type CompressionHandler struct {
	pb.UnimplementedImageServiceServer
	svc service.ImageService
}

func NewCompressionHandler(svc service.ImageService) *CompressionHandler {
	return &CompressionHandler{
		svc: svc,
	}
}

func (h *CompressionHandler) GetImage(ctx context.Context, req *pb.GetImageRequest) (*pb.GetImageResponse, error) {
	image, err := h.svc.Get(ctx, model.ListImageFilter{
		ID: req.GetId(),
	})
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(image)
	if err != nil {
		return nil, errx.NewInternal().WithDescriptionAndCause("grpc.GetImage: ", err)
	}

	return &pb.GetImageResponse{
		Data: data,
	}, nil
}

func (h *CompressionHandler) UploadImage(ctx context.Context, req *pb.UploadImageRequest) (*emptypb.Empty, error) {
	if req.GetData() == nil || len(req.GetData()) == 0 {
		return nil, errx.NewInternal().WithDescription("no data provided in request")
	}

	if err := h.svc.Create(ctx, nil); err != nil {
		return nil, err
	}

	return nil, nil
}
