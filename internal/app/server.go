package app

import (
	"context"
	"fmt"
	"log"

	"go-image-compression/internal/config"
	consumer "go-image-compression/internal/controller/amqp"
	"go-image-compression/internal/controller/grpc"
	"go-image-compression/internal/controller/http"
	"go-image-compression/internal/repository"
	"go-image-compression/internal/service"
	"go-image-compression/pkg/broker"
	"go-image-compression/pkg/db"
	"go-image-compression/pkg/resizer"

	"github.com/gofiber/fiber/v2"
	"github.com/nordew/go-errx"
)

const codepath = "app/server.go"

func MustRun() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.MustLoad()
	if err != nil {
		return fmt.Errorf("%s: %w", codepath, err)
	}

	client, err := db.NewMinioStorage(cfg.Minio)
	if err != nil {
		return errx.NewInternal().WithDescriptionAndCause(codepath, err)
	}
	if err = migrate(ctx, client); err != nil {
		return fmt.Errorf("%s: %w", codepath, err)
	}

	nc, err := broker.NewNatsClient(cfg.NATS)
	if err != nil {
		return fmt.Errorf("%s: %w", codepath, err)
	}
	defer nc.Close()

	resizer := resizer.NewResizer()

	repositories := repository.NewRepository(client)
	services := service.NewService(repositories, nc, resizer)
	handler := http.NewHandler(services)

	consumer, err := consumer.NewConsumer(nc, services)
	if err != nil {
		return fmt.Errorf("%s: %w", codepath, err)
	}

	if err = consumer.Start(); err != nil {
		return fmt.Errorf("%s: %w", codepath, err)
	}

	go func() {
		if err := grpc.StartGRPCServer(services.ImageService); err != nil {
			log.Printf("%s: failed to start gRPC server: %v", codepath, err)
		}
	}()

	app := fiber.New()
	handler.SetupRoutes(app)

	if err := app.Listen(cfg.HTTP.Port); err != nil {
		return errx.NewInternal().WithDescriptionAndCause(codepath, err)
	}

	return nil
}

func migrate(ctx context.Context, client db.Storage) error {
	exists, err := client.BucketExists(ctx, "images")
	if err != nil {
		return errx.NewInternal().WithDescriptionAndCause(codepath, err)
	}
	if !exists {
		err = client.CreateBucket(ctx, "images")
		if err != nil {
			return errx.NewInternal().WithDescriptionAndCause(codepath, err)
		}
		log.Println("Bucket \"images\" created successfully")

		return nil
	}
	log.Println("Bucket images already exists")

	return nil
}
