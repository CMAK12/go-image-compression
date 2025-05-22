package broker

type Message struct {
	Topic string
	Data  []byte
}

type HandlerFunc func(msg *Message) error

type Broker interface {
	Publish(topic string, data []byte) error
	Subscribe(topic string, handler HandlerFunc) error
	Close()
}
