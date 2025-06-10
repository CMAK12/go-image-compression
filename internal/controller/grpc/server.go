package grpc

import (
	"log"
	"net"

	"go-image-compression/internal/service"
	pb "go-image-compression/pkg/proto"

	"github.com/nordew/go-errx"
	"google.golang.org/grpc"
)

func StartGRPCServer(service service.ImageService) error {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		return errx.NewInternal().WithDescriptionAndCause("grpc.StartGRPCServer: ", err)
	}

	grpcServer := grpc.NewServer()
	handler := NewCompressionHandler(service)
	pb.RegisterImageServiceServer(grpcServer, handler)

	log.Printf("Starting gRPC server on %s", lis.Addr().String())
	return grpcServer.Serve(lis)
}
