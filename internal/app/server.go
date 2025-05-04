package app

import (
	"context"
	"fmt"
	"log"

	"go-image-compression/internal/config"
	"go-image-compression/internal/controller/v1/http"
	"go-image-compression/pkg/db"

	"github.com/gofiber/fiber/v2"
	"github.com/minio/minio-go/v7"
)

func MustRun() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.MustLoad()
	if err != nil {
		return err
	}

	client, err := db.MustConnectMinio(cfg.Minio)
	if err != nil {
		return fmt.Errorf("app/server.go: %v", err)
	}

	err = migrate(ctx, client)
	if err != nil {
		return fmt.Errorf("app/server.go: %v", err)
	}

	handler := http.NewHandler()

	app := fiber.New()
	handler.SetupRoutes(app)

	if err := app.Listen(cfg.HTTP.Port); err != nil {
		return err
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
