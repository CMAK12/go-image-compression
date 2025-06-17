package broker

import (
	"fmt"
	"go-image-compression/internal/config"
	"time"

	"github.com/nats-io/nats.go"
)

type (
	natsClient struct {
		natsClient *nats.Conn
		js         nats.JetStreamContext
	}
)

func NewNatsClient(natsConfig config.NATSConfig, topicConfig config.Topic) (Broker, error) {
	nc, err := nats.Connect(
		natsConfig.URL,
		nats.Timeout(5*time.Second),
		nats.MaxReconnects(10),
		nats.RetryOnFailedConnect(true),
	)
	if err != nil {
		return nil, fmt.Errorf("pkg.broker.nats: unable to connect to NATS server: %w", err)
	}

	js, err := nc.JetStream()
	if err != nil {
		return nil, fmt.Errorf("pkg.broker.nats: unable to get JetStream context: %w", err)
	}

	if _, err := js.AddStream(&nats.StreamConfig{
		Name:     topicConfig.ImageStream,
		Subjects: []string{topicConfig.ImageCreated},
		Storage:  nats.FileStorage,
	}); err != nil {
		return nil, fmt.Errorf("pkg.broker.nats: unable to create stream: %w", err)
	}

	return &natsClient{
		natsClient: nc,
		js:         js,
	}, nil
}

func (n *natsClient) Publish(subject string, data []byte) (*PublishResult, error) {
	ack, err := n.js.Publish(subject, data)
	if err != nil {
		return nil, err
	}

	return &PublishResult{
		Stream:  ack.Stream,
		Seq:     ack.Sequence,
		Success: ack != nil,
	}, nil
}

func (n *natsClient) Subscribe(subject string, handler HandlerFunc) (Subscription, error) {
	sub, err := n.js.Subscribe(subject, func(msg *nats.Msg) {
		err := handler(&Message{
			Topic: msg.Subject,
			Data:  msg.Data,
		})
		if err != nil {
			_ = msg.Nak()
			return
		}
		_ = msg.Ack()
	}, nats.ManualAck())

	if err != nil {
		return nil, err
	}

	return sub, nil
}

func (n *natsClient) Close() {
	if n.natsClient != nil {
		n.natsClient.Close()
	}
}
