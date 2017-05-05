package protocol

type Receiver interface {
	Data() chan []byte
	Start() error
	Stop()
}
