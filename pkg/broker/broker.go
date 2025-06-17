package broker

type (
	Message struct {
		Topic string
		Data  []byte
	}

	PublishResult struct {
		Stream  string
		Seq     uint64
		Success bool
	}

	Subscription interface {
		Unsubscribe() error
		Drain() error
	}
)

type HandlerFunc func(msg *Message) error

type Broker interface {
	Publish(topic string, data []byte) (*PublishResult, error)
	Subscribe(subject string, handler HandlerFunc) (Subscription, error)
	Close()
}
