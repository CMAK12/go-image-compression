package app

import (
	"context"
	"fmt"
	"net"

	"go-image-compression/internal/config"
	consumer "go-image-compression/internal/controller/amqp"
	grpc_handler "go-image-compression/internal/controller/grpc"
	"go-image-compression/internal/controller/http"
	"go-image-compression/internal/repository"
	"go-image-compression/internal/service"
	"go-image-compression/pkg/broker"
	"go-image-compression/pkg/converter"
	"go-image-compression/pkg/db"
	"go-image-compression/pkg/pb"
	"go-image-compression/pkg/resizer"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func MustRun() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.MustLoad()
	if err != nil {
		return fmt.Errorf("server.MustRun: %w", err)
	}

	logger, err := zap.NewDevelopment()
	if err != nil {
		return fmt.Errorf("server.MustRun: %w", err)
	}
	defer logger.Sync()

	client, err := db.NewMinioStorage(cfg.Minio)
	if err != nil {
		return fmt.Errorf("server.MustRun.MinIO: %w", err)
	}

	if err = migrate(ctx, client, logger); err != nil {
		return fmt.Errorf("server.MustRun: %w", err)
	}

	nc, err := broker.NewNatsClient(cfg.NATS, cfg.Topic)
	if err != nil {
		return fmt.Errorf("server.MustRun: %w", err)
	}
	defer nc.Close()

	resizer := resizer.NewResizer()
	converter := converter.NewConverter()

	repositories := repository.NewRepository(client)
	service := service.NewService(repositories, nc, resizer, cfg)

	if err = consumer.Start(nc, service, cfg, logger); err != nil {
		return fmt.Errorf("server.MustRun: %w", err)
	}

	go func() {
		if err = startHTTPServer(service, cfg, logger); err != nil {
			logger.Error("failed to start HTTP server", zap.Error(err))
		}
	}()

	if err = startGRPCServer(service, cfg, logger, converter); err != nil {
		return fmt.Errorf("server.MustRun: %w", err)
	}

	return nil
}

func startHTTPServer(service service.Services, cfg *config.Config, logger *zap.Logger) error {
	app := fiber.New()

	handler := http.NewHandler(service, logger)
	handler.SetupRoutes(app)

	if err := app.Listen(cfg.HTTP.Port); err != nil {
		return fmt.Errorf("server.startHTTPServer: %w", err)
	}

	return nil
}

func startGRPCServer(service service.Services, cfg *config.Config, logger *zap.Logger, converter converter.Converter) error {
	lis, err := net.Listen("tcp", cfg.GRPC.Port)
	if err != nil {
		return fmt.Errorf("server.startGRPCServer: %w", err)
	}

	grpcHandler := grpc_handler.NewCompressionHandler(service, logger, converter)
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_handler.UnaryLoggerInterceptor(logger)),
	)
	pb.RegisterImageServiceServer(grpcServer, grpcHandler)

	logger.Info("Starting gRPC server on " + lis.Addr().String())
	grpcServer.Serve(lis)

	return nil
}

func migrate(ctx context.Context, client db.Storage, logger *zap.Logger) error {
	exists, err := client.BucketExists(ctx, "images")
	if err != nil {
		return fmt.Errorf("server.migrate: %w", err)
	}
	if !exists {
		err = client.CreateBucket(ctx, "images")
		if err != nil {
			return fmt.Errorf("server.migrate: %w", err)
		}
		logger.Info("Bucket \"images\" created successfully")

		return nil
	}
	logger.Info("Bucket images already exists")

	return nil
}
