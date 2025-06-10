package app

import (
	"context"
	"fmt"
	"log"
	"net"

	"go-image-compression/internal/config"
	consumer "go-image-compression/internal/controller/amqp"
	grpc_handler "go-image-compression/internal/controller/grpc"
	"go-image-compression/internal/controller/http"
	"go-image-compression/internal/repository"
	"go-image-compression/internal/service"
	"go-image-compression/pkg/broker"
	"go-image-compression/pkg/db"
	pb "go-image-compression/pkg/proto"
	"go-image-compression/pkg/resizer"

	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc"
)

func MustRun() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.MustLoad()
	if err != nil {
		return fmt.Errorf("server.MustRun: %w", err)
	}

	client, err := db.NewMinioStorage(cfg.Minio)
	if err != nil {
		return fmt.Errorf("server.MustRun: %w", err)
	}

	if err = migrate(ctx, client); err != nil {
		return fmt.Errorf("server.MustRun: %w", err)
	}

	nc, err := broker.NewNatsClient(cfg.NATS)
	if err != nil {
		return fmt.Errorf("server.MustRun: %w", err)
	}
	defer nc.Close()

	resizer := resizer.NewResizer()

	repositories := repository.NewRepository(client)
	service := service.NewService(repositories, nc, resizer)

	if err = consumer.Start(nc, service); err != nil {
		return fmt.Errorf("server.MustRun: %w", err)
	}

	go func() {
		if err = startHTTPServer(service, cfg); err != nil {
			log.Println(err.Error())
		}
	}()

	if err = startGRPCServer(service, cfg); err != nil {
		return fmt.Errorf("server.MustRun: %w", err)
	}

	return nil
}

func startHTTPServer(service service.Services, cfg *config.Config) error {
	app := fiber.New()

	handler := http.NewHandler(service)
	handler.SetupRoutes(app)

	if err := app.Listen(cfg.HTTP.Port); err != nil {
		return fmt.Errorf("server.startHTTPServer: %w", err)
	}

	return nil
}

func startGRPCServer(service service.Services, cfg *config.Config) error {
	lis, err := net.Listen("tcp", cfg.GRPC.Port)
	if err != nil {
		return fmt.Errorf("server.startGRPCServer: %w", err)
	}

	grpcServer := grpc.NewServer()
	grpcHandler := grpc_handler.NewCompressionHandler(service)
	pb.RegisterImageServiceServer(grpcServer, grpcHandler)

	log.Printf("Starting gRPC server on %s", lis.Addr().String())
	grpcServer.Serve(lis)

	return nil
}

func migrate(ctx context.Context, client db.Storage) error {
	exists, err := client.BucketExists(ctx, "images")
	if err != nil {
		return fmt.Errorf("server.migrate: %w", err)
	}
	if !exists {
		err = client.CreateBucket(ctx, "images")
		if err != nil {
			return fmt.Errorf("server.migrate: %w", err)
		}
		log.Println("Bucket \"images\" created successfully")

		return nil
	}
	log.Println("Bucket images already exists")

	return nil
}
