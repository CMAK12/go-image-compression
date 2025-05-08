package app

import (
	"context"
	"fmt"
	"log"

	"go-image-compression/internal/config"
	"go-image-compression/internal/controller/v1/http"
	"go-image-compression/internal/repository"
	"go-image-compression/internal/service"
	"go-image-compression/pkg/db"

	"github.com/gofiber/fiber/v2"
	"github.com/minio/minio-go/v7"
	"github.com/nats-io/nats.go"
)

const codepath = "app/server.go"

func MustRun() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.MustLoad()
	if err != nil {
		return fmt.Errorf("%s: %v", codepath, err)
	}

	client, err := db.MustConnectMinio(cfg.Minio)
	if err != nil {
		return fmt.Errorf("%s: %v", codepath, err)
	}

	err = migrate(ctx, client)
	if err != nil {
		return fmt.Errorf("%s: %v", codepath, err)
	}

	natsClient, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		return fmt.Errorf("%s: %v", codepath, err)
	}
	defer natsClient.Close()

	repositories := repository.NewRepository(client)
	services := service.NewService(repositories)
	handler := http.NewHandler(services)

	app := fiber.New()
	handler.SetupRoutes(app)

	if err := app.Listen(cfg.HTTP.Port); err != nil {
		return fmt.Errorf("%s: %v", codepath, err)
	}

	return nil
}

func migrate(ctx context.Context, client *minio.Client) error {
	buckets := []string{
		"image-100",
		"image-75",
		"image-50",
		"image-25",
	}

	for _, bucket := range buckets {
		exists, err := client.BucketExists(ctx, bucket)
		if err != nil {
			return err
		}
		if !exists {
			err = client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
			if err != nil {
				return err
			}
			log.Printf("Bucket %s created successfully\n", bucket)
			continue
		}
		log.Printf("Bucket %s already exists\n", bucket)
	}

	return nil
}
