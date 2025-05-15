package worker

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/png"
	"io"
	"log"
	"os"

	"go-image-compression/internal/consts"
	"go-image-compression/internal/model"
	"go-image-compression/internal/repository"

	"github.com/nats-io/nats.go"
	"github.com/nfnt/resize"
	"github.com/nordew/go-errx"
)

const codepath = "worker/image.go"

type ImageWorker struct {
	nats         *nats.Conn
	subscription *nats.Subscription
	repo         repository.ImageRepository
}

func NewImageWorker(natsURL string, repositories repository.Repositories) (ImageWorker, error) {
	nats, err := nats.Connect(natsURL)
	if err != nil {
		return ImageWorker{}, errx.NewInternal().WithDescriptionAndCause(codepath, err)
	}

	return ImageWorker{
		nats: nats,
		repo: repositories.ImageRepository,
	}, nil
}

func (w *ImageWorker) Start() error {
	var err error
	w.subscription, err = w.nats.Subscribe(
		consts.SubjectImageCompress,
		w.handleMessage,
	)
	if err != nil {
		log.Printf("%s: %v", codepath, err)
	}

	return err
}

func (w *ImageWorker) handleMessage(msg *nats.Msg) {
	imageID := string(msg.Data)

	recievedImage, err := w.repo.Get(context.Background(), model.ListImageFilter{
		ID: imageID,
	})
	if err != nil {
		log.Printf("%s: %v", codepath, err)
	}

	decodedImage, _, err := image.Decode(recievedImage)
	if err != nil {
		log.Printf("%s: %v", codepath, err)
		return
	}

	bounds := decodedImage.Bounds()
	imageWidth := bounds.Dx()
	imageHeight := bounds.Dy()
	imageID = imageID[:len(imageID)-4]

	for i := 0.75; i != 0; i -= 0.25 {
		go w.uploadImage(decodedImage, i, imageID, imageWidth, imageHeight)
	}
}

func (w *ImageWorker) uploadImage(
	recievedImage image.Image,
	compressPercent float64,
	imageID string,
	imageWidth,
	imageHeight int,
) {
	newFilename := fmt.Sprintf("%s_%d", imageID, int(compressPercent*100))

	widthToCompress := float64(imageWidth) * compressPercent
	heightToCompress := float64(imageHeight) * compressPercent

	resizedImage, err := compressImage(recievedImage, widthToCompress, heightToCompress, newFilename)
	if err != nil {
		log.Printf("%s: %v", codepath, err)
		return
	}
	defer resizedImage.Close()

	stat, err := resizedImage.Stat()
	if err != nil {
		log.Printf("%s: %v", codepath, err)
		return
	}

	modelImage := model.Image{
		ID:          newFilename,
		File:        resizedImage,
		FileSize:    stat.Size(),
		Bucket:      consts.BucketName,
		ContentType: "image/png",
	}

	if err := w.repo.Create(context.Background(), modelImage); err != nil {
		log.Printf("%s: %v", codepath, err)
	}
}

func compressImage(img image.Image, width, height float64, filename string) (*os.File, error) {
	resizedImage := resize.Resize(uint(width), uint(height), img, resize.Lanczos3)

	var buf bytes.Buffer
	if err := png.Encode(&buf, resizedImage); err != nil {
		return nil, fmt.Errorf("failed to encode image: %w", err)
	}

	tmpFile, err := os.CreateTemp("", filename+"-*.png")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %w", err)
	}

	if _, err := io.Copy(tmpFile, &buf); err != nil {
		tmpFile.Close()
		return nil, fmt.Errorf("failed to write to temp file: %w", err)
	}

	if _, err := tmpFile.Seek(0, io.SeekStart); err != nil {
		tmpFile.Close()
		return nil, fmt.Errorf("failed to seek temp file: %w", err)
	}

	return tmpFile, nil
}
