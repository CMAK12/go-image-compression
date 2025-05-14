package broker

import (
	"context"

	"go-image-compression/internal/consts"

	"github.com/nats-io/nats.go"
	"github.com/nordew/go-errx"
)

const codepath = "broker/image.go"

type ImageProducer struct {
	natsClient *nats.Conn
}

func NewImageProducer(natsURL string) (ImageProducer, error) {
	natsClient, err := nats.Connect(natsURL)
	if err != nil {
		return ImageProducer{}, errx.NewInternal().WithDescriptionAndCause(codepath, err)
	}

	return ImageProducer{
		natsClient: natsClient,
	}, nil
}

func (p *ImageProducer) Publish(ctx context.Context, imageID string) error {
	err := p.natsClient.Publish(consts.SubjectImageCompress, []byte(imageID))
	if err != nil {
		return errx.NewInternal().WithDescriptionAndCause(codepath, err)
	}

	return nil
}

func (p *ImageProducer) Close() {
	if p.natsClient != nil {
		p.natsClient.Close()
	}
}
