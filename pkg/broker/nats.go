package broker

import (
	"go-image-compression/internal/config"

	"github.com/nats-io/nats.go"
	"github.com/nordew/go-errx"
)

type (
	natsClient struct {
		natsClient *nats.Conn
	}
)

func NewNatsClient(cfg config.NATSConfig) (Broker, error) {
	nc, err := nats.Connect(cfg.URL)
	if err != nil {
		return nil, errx.NewInternal().WithDescriptionAndCause("nats:connect", err)
	}

	return &natsClient{
		natsClient: nc,
	}, nil
}

func (n *natsClient) Publish(subject string, data []byte) error {
	err := n.natsClient.Publish(subject, data)
	if err != nil {
		return err
	}

	return nil
}

func (n *natsClient) Subscribe(subject string, handler HandlerFunc) error {
	_, err := n.natsClient.Subscribe(subject, func(msg *nats.Msg) {
		handler(&Message{
			Topic: msg.Subject,
			Data:  msg.Data,
		})
	})
	if err != nil {
		return err
	}

	return nil
}

func (n *natsClient) Close() {
	if n.natsClient != nil {
		n.natsClient.Close()
	}
}
