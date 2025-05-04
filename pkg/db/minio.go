package db

import (
	"fmt"
	"log"

	"go-image-compression/internal/config"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func MustConnectMinio(cfg config.MinioConfig) (*minio.Client, error) {
	minioClient, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		Secure: cfg.SSLMode,
	})
	if err != nil {
		return nil, fmt.Errorf("pkg/db/minio.go: %v", err)
	}

	log.Printf("Connected to MinIO: %s", cfg.Endpoint)

	return minioClient, err
}
