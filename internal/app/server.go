package app

import (
	"context"
	"fmt"
	"log"

	"go-image-compression/internal/broker"
	"go-image-compression/internal/config"
	"go-image-compression/internal/consts"
	"go-image-compression/internal/controller/v1/http"
	"go-image-compression/internal/repository"
	"go-image-compression/internal/service"
	"go-image-compression/internal/worker"
	"go-image-compression/pkg/db"

	"github.com/gofiber/fiber/v2"
	"github.com/minio/minio-go/v7"
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

	client, err := db.MustConnectMinio(cfg.Minio)
	if err != nil {
		return errx.NewInternal().WithDescriptionAndCause(codepath, err)
	}
	if err = migrate(ctx, client); err != nil {
		return fmt.Errorf("%s: %w", codepath, err)
	}

	broker, err := broker.NewImageProducer(cfg.NATS.URL)
	if err != nil {
		return fmt.Errorf("%s: %w", codepath, err)
	}
	defer broker.Close()

	repositories := repository.NewRepository(client)
	services := service.NewService(repositories, broker)
	handler := http.NewHandler(services)

	w, err := worker.NewImageWorker(cfg.NATS.URL, repositories)
	if err != nil {
		return fmt.Errorf("%s: %w", codepath, err)
	}

	if err = w.Start(); err != nil {
		return fmt.Errorf("%s: %w", codepath, err)
	}

	app := fiber.New()
	handler.SetupRoutes(app)

	if err := app.Listen(cfg.HTTP.Port); err != nil {
		return errx.NewInternal().WithDescriptionAndCause(codepath, err)
	}

	return nil
}

func migrate(ctx context.Context, client *minio.Client) error {
	exists, err := client.BucketExists(ctx, consts.BucketName)
	if err != nil {
		return errx.NewInternal().WithDescriptionAndCause(codepath, err)
	}
	if !exists {
		err = client.MakeBucket(ctx, consts.BucketName, minio.MakeBucketOptions{})
		if err != nil {
			return errx.NewInternal().WithDescriptionAndCause(codepath, err)
		}
		log.Printf("Bucket \"%s\" created successfully\n", consts.BucketName)

		return nil
	}
	log.Printf("Bucket %s already exists\n", consts.BucketName)

	return nil
}
