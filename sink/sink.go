package sink

type Sink interface {
	Write(data []byte) error
	Close() error
}
