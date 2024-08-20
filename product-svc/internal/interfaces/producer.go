package interfaces

type Producer interface {
	SendMessage(key string, message []byte) error
}
